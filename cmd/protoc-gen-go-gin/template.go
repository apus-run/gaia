package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"
)

var httpTemplate = `

type {{ $.InterfaceName }} interface {
{{- range .MethodSet}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}
func Register{{ $.InterfaceName }}(r gin.IRouter, srv {{ $.InterfaceName }}) {
	s := &{{.Name}}{
		server: srv,
		router: r,
	}
	s.RegisterService()
}

type {{$.Name}} struct{
	server {{ $.InterfaceName }}
	router gin.IRouter
}

{{range .Methods}}
func (s *{{$.Name}}) {{ .HandlerName }} (ctx *xgin.Context) {
	var in {{.Request}}
{{if .HasPathParams }}
	if err := ctx.ShouldBindUri(&in); err != nil {
		e := errcode.ErrInvalidParam.WithDetails(err.Error())
		ctx.Error(e)
		return
	}
{{else if eq .Method "GET" }}
	if err := ctx.ShouldBindQuery(&in); err != nil {
        e := errcode.ErrInvalidParam.WithDetails(err.Error())
		ctx.Error(e)
		return
	}
{{else if eq .Method "POST" "PUT" "PATCH" "DELETE"}}
	if err := ctx.ShouldBindJSON(&in); err != nil {
        e := errcode.ErrInvalidParam.WithDetails(err.Error())
		ctx.Error(e)
		return
	}
{{else}}
	if err := ctx.ShouldBind(&in); err != nil {
        e := errcode.ErrInvalidParam.WithDetails(err.Error())
		ctx.Error(e)
		return
	}
{{end}}
	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.({{ $.InterfaceName }}).{{.Name}}(newCtx, &in)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Success(out)
}
{{end}}

func (s *{{$.Name}}) RegisterService() {
{{- range .Methods}}
		s.router.Handle("{{.Method}}", "{{.Path}}", xgin.Handle(s.{{ .HandlerName }}) )
{{- end}}
}
`

type service struct {
	Name     string // Greeter
	FullName string // helloworld.Greeter
	FilePath string // api/helloworld/helloworld.proto

	Methods   []*method
	MethodSet map[string]*method
}

func (s *service) execute() string {
	if s.MethodSet == nil {
		s.MethodSet = map[string]*method{}
		for _, m := range s.Methods {
			m := m
			s.MethodSet[m.Name] = m
		}
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(httpTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return buf.String()
}

// InterfaceName service interface name
func (s *service) InterfaceName() string {
	return s.Name + "HTTPServer"
}

type method struct {
	Name    string // SayHello
	Num     int    // ?????? rpc ???????????????????????? http ??????
	Request string // SayHelloReq
	Reply   string // SayHelloResp
	// http_rule
	Path         string // ??????
	Method       string // HTTP Method
	Body         string
	ResponseBody string
}

// HandlerName for gin handler name
func (m *method) HandlerName() string {
	return fmt.Sprintf("%s_%d", m.Name, m.Num)
}

// HasPathParams ????????????????????????
func (m *method) HasPathParams() bool {
	paths := strings.Split(m.Path, "/")
	for _, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			return true
		}
	}
	return false
}

// initPathParams ?????????????????? {xx} --> :xx
func (m *method) initPathParams() {
	paths := strings.Split(m.Path, "/")
	for i, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			paths[i] = ":" + p[1:len(p)-1]
		}
	}
	m.Path = strings.Join(paths, "/")
}
