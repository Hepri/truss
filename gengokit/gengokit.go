package gengokit

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	generatego "github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/pkg/errors"

	"github.com/TuneLab/truss/gengokit/clientarggen"
	"github.com/TuneLab/truss/gengokit/httptransport"
	"github.com/TuneLab/truss/svcdef"
)

type Renderable interface {
	Render(string, *Data) (io.Reader, error)
}

type Config struct {
	GoPackage   string
	PBPackage   string
	Version     string
	VersionDate string

	PreviousFiles map[string]io.Reader
}

// FuncMap contains a series of utility functions to be passed into
// templates and used within those templates.
var FuncMap = template.FuncMap{
	"ToLower": strings.ToLower,
	"GoName":  generatego.CamelCase,
}

// Data is passed to templates as the executing struct; its fields
// and methods are used to modify the template
type Data struct {
	// import path for the directory containing the definition .proto files
	ImportPath string
	// import path for .pb.go files containing service structs
	PBImportPath string
	// PackageName is the name of the package containing the service definition
	PackageName string
	// GRPC/Protobuff service, with all parameters and return values accessible
	Service        *svcdef.Service
	DataStoreTypes []*DataStoreType
	ClientArgs     *clientarggen.ClientServiceArgs
	// A helper struct for generating http transport functionality.
	HTTPHelper *httptransport.Helper
	FuncMap    template.FuncMap

	Version     string
	VersionDate string
}

func NewData(sd *svcdef.Svcdef, conf Config) (*Data, error) {
	return &Data{
		ImportPath:     conf.GoPackage,
		PBImportPath:   conf.PBPackage,
		PackageName:    sd.PkgName,
		Service:        sd.Service,
		DataStoreTypes: getDataStoreTypes(sd.Messages),
		ClientArgs:     clientarggen.New(sd.Service),
		HTTPHelper:     httptransport.NewHelper(sd.Service),
		FuncMap:        FuncMap,
		Version:        conf.Version,
		VersionDate:    conf.VersionDate,
	}, nil
}

// ApplyTemplate applies the passed template with the Data
func (e *Data) ApplyTemplate(templ string, templName string) (io.Reader, error) {
	return ApplyTemplate(templ, templName, e, e.FuncMap)
}

// ApplyTemplate is a helper methods that packages can call to render a
// template with any data and func map
func ApplyTemplate(templ string, templName string, data interface{}, funcMap template.FuncMap) (io.Reader, error) {
	codeTemplate, err := template.New(templName).Funcs(funcMap).Parse(templ)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create template")
	}

	outputBuffer := bytes.NewBuffer(nil)
	err = codeTemplate.Execute(outputBuffer, data)
	if err != nil {
		return nil, errors.Wrap(err, "template error")
	}

	return outputBuffer, nil
}

type DataStoreType struct {
	Name        string
	HasSchoolId bool
}

func getDataStoreTypes(arr []*svcdef.Message) []*DataStoreType {
	filtered := make([]*DataStoreType, 0)
	for _, m := range arr {
		isReq := strings.HasSuffix(m.Name, "Request")
		isResp := strings.HasSuffix(m.Name, "Response")
		isEmpty := m.Name == "Empty"
		if !(isReq || isResp || isEmpty) {
			var hasSchoolId bool
			for _, f := range m.Fields {
				if strings.ToLower(f.Name) == "schoolid" {
					hasSchoolId = true
					break
				}
			}

			filtered = append(filtered, &DataStoreType{
				Name:        m.Name,
				HasSchoolId: hasSchoolId,
			})
		}
	}
	return filtered
}
