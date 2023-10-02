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
func (s *{{$.Name}}) {{ .HandlerName }}_HTTP_Handler (ctx *ginx.Context) {
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
		s.router.Handle("{{.Method}}", "{{.Path}}", ginx.Handle(s.{{ .HandlerName }}_HTTP_Handler) )
{{- end}}
}