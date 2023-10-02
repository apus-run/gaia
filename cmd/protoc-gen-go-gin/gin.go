package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	contextPackage     = protogen.GoImportPath("context")
	ginPackage         = protogen.GoImportPath("github.com/gin-gonic/gin")
	metadataPackage    = protogen.GoImportPath("google.golang.org/grpc/metadata")
	ginxPackage        = protogen.GoImportPath("github.com/apus-run/gaia/pkg/ginx")
	errCodePackage     = protogen.GoImportPath("github.com/apus-run/gaia/pkg/errcode")
	deprecationComment = "// Deprecated: Do not use."
)

var methodSets = make(map[string]int)

// generateFile generates a _gin.pb.go file.
func generateFile(gen *protogen.Plugin, file *protogen.File, omitempty bool, omitemptyPrefix string) *protogen.GeneratedFile {
	if len(file.Services) == 0 || (omitempty && !hasHTTPRule(file.Services)) {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_gin.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated protoc-gen-go-gin. DO NOT EDIT.")
	g.P(fmt.Sprintf("// protoc-gen-go-gin %s", version))
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the eagle package it is being compiled against.")
	g.P()
	g.P("// ", contextPackage.Ident(""))
	g.P("// ", metadataPackage.Ident(""))
	g.P("// ", ginPackage.Ident(""), ginxPackage.Ident(""), errCodePackage.Ident(""))
	g.P()

	generateFileContent(gen, file, g, omitempty, omitemptyPrefix)
	return g
}

// generateFileContent generates the gaia errors definitions, excluding the package statement.
func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, omitempty bool, omitemptyPrefix string) {
	if len(file.Services) == 0 {
		return
	}
	for _, service := range file.Services {
		genService(gen, file, g, service, omitempty, omitemptyPrefix)
	}
}

func genService(
	_ *protogen.Plugin,
	file *protogen.File,
	g *protogen.GeneratedFile,
	s *protogen.Service,
	omitempty bool,
	omitemptyPrefix string,
) {
	if s.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}
	// HTTP Server.
	sd := &service{
		Name:     s.GoName,
		FullName: string(s.Desc.FullName()),
		FilePath: file.Desc.Path(),
	}

	for _, method := range s.Methods {
		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			continue
		}
		// 存在 http rule 配置
		rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
		if rule != nil && ok {
			for _, bind := range rule.AdditionalBindings {
				sd.Methods = append(sd.Methods, buildHTTPRule(g, s, method, bind, omitemptyPrefix))
			}
			sd.Methods = append(sd.Methods, buildHTTPRule(g, s, method, rule, omitemptyPrefix))
		} else if !omitempty {
			path := fmt.Sprintf("/%s/%s", s.Desc.FullName(), method.Desc.Name())
			sd.Methods = append(sd.Methods, buildMethodDesc(g, method, http.MethodPost, path))
		}
	}
	if len(sd.Methods) != 0 {
		g.P(sd.execute())
	}
}

func hasHTTPRule(services []*protogen.Service) bool {
	for _, service := range services {
		for _, method := range service.Methods {
			if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
				continue
			}
			rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
			if rule != nil && ok {
				return true
			}
		}
	}
	return false
}

func buildHTTPRule(
	g *protogen.GeneratedFile,
	service *protogen.Service,
	m *protogen.Method,
	rule *annotations.HttpRule,
	omitemptyPrefix string,
) *method {
	var (
		path   string
		method string
	)
	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		path = pattern.Get
		method = http.MethodGet
	case *annotations.HttpRule_Put:
		path = pattern.Put
		method = http.MethodPut
	case *annotations.HttpRule_Post:
		path = pattern.Post
		method = http.MethodPost
	case *annotations.HttpRule_Delete:
		path = pattern.Delete
		method = http.MethodDelete
	case *annotations.HttpRule_Patch:
		path = pattern.Patch
		method = http.MethodPatch
	case *annotations.HttpRule_Custom:
		path = pattern.Custom.Path
		method = pattern.Custom.Kind
	}
	if method == "" {
		method = http.MethodPost
	}
	if path == "" {
		path = fmt.Sprintf("%s/%s/%s", omitemptyPrefix, service.Desc.FullName(), m.Desc.Name())
	}
	md := buildMethodDesc(g, m, method, path)
	return md
}

func buildMethodDesc(_ *protogen.GeneratedFile, m *protogen.Method, httpMethod, path string) *method {
	defer func() { methodSets[m.GoName]++ }()
	md := &method{
		Name:    m.GoName,
		Num:     methodSets[m.GoName],
		Request: m.Input.GoIdent.GoName,
		Reply:   m.Output.GoIdent.GoName,
		Path:    path,
		Method:  httpMethod,
	}
	md.initPathParams()
	return md
}

var matchFirstCap = regexp.MustCompile("([A-Z])([A-Z][a-z])")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(input string) string {
	output := matchFirstCap.ReplaceAllString(input, "${1}_${2}")
	output = matchAllCap.ReplaceAllString(output, "${1}_${2}")
	output = strings.ReplaceAll(output, "-", "_")
	return strings.ToLower(output)
}
