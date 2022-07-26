package xgin

import (
	"net/http"

	"github.com/apus-run/gaia/pkg/errcode"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"google.golang.org/grpc/status"

	"github.com/apus-run/gaia/pkg/pagination"
	httpstatus "github.com/apus-run/gaia/transport/http/status"
)

// Context a wrapper of gin.Context
type Context struct {
	*gin.Context
}

// HandlerFunc defines the handler to wrap gin.Context
type HandlerFunc func(c *Context)

// Handle convert HandlerFunc to gin.HandlerFunc
func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			c,
		}
		h(ctx)
	}
}

const (
	// CodeOK means a successful response
	CodeOK = 0
	// CodeErr means a failure response
	CodeErr = 1
)

// Response defines HTTP JSON response
type Response struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	Details []string    `json:"details,omitempty"`
}

// ResponsePagination defines HTTP JSON response with extra pagination data
type ResponsePagination struct {
	Response
	Pagination pagination.Pagination `json:"pagination"`
}

// JSON returns JSON response
// e.x. {"code":<code>, "msg":<msg>, "data":<data>, "details":<details>}
func (c *Context) JSON(httpStatus int, resp Response) {
	c.Context.JSON(httpStatus, resp)
}

// JSONOK returns JSON response with successful business code and data
// e.x. {"code":0, "msg":"成功", "data":<data>}
func (c *Context) JSONOK(data interface{}) {
	j := new(Response)
	j.Code = CodeOK
	j.Msg = "ok"
	if data == nil {
		j.Data = gin.H{}
	} else {
		j.Data = data
	}
	c.Context.JSON(http.StatusOK, j)
	return
}

func (c *Context) Success(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	j := &Response{
		Code: errcode.Success.Code(),
		Msg:  errcode.Success.Msg(),
		Data: nil,
	}

	c.Context.JSON(http.StatusOK, j)
	return
}

// JSONE returns JSON response with failure business code ,msg and data
// e.x. {"code":<code>, "msg":<msg>, "data":<data>}
func (c *Context) JSONE(code int, msg string, data interface{}) {
	j := new(Response)
	j.Code = code
	j.Msg = msg
	switch d := data.(type) {
	case error:
		j.Data = d.Error()
	default:
		j.Data = data
	}

	c.Context.JSON(http.StatusOK, j)
	return
}

func (c *Context) Error(err error) {
	if err == nil {
		c.JSON(http.StatusOK, Response{
			Code: errcode.Success.Code(),
			Msg:  errcode.Success.Msg(),
			Data: gin.H{},
		})
		return
	}

	if v, ok := err.(*errcode.Error); ok {
		response := Response{
			Code:    v.Code(),
			Msg:     v.Msg(),
			Data:    gin.H{},
			Details: []string{},
		}
		details := v.Details()
		if len(details) > 0 {
			response.Details = details
		}
		c.JSON(errcode.ToHTTPStatusCode(v.Code()), response)
		return
	} else {
		// receive gRPC error
		if st, ok := status.FromError(err); ok {
			response := Response{
				Code:    int(st.Code()),
				Msg:     st.Message(),
				Data:    gin.H{},
				Details: []string{},
			}
			details := st.Details()
			if len(details) > 0 {
				for _, v := range details {
					response.Details = append(response.Details, cast.ToString(v))
				}
			}
			// https://httpstatus.in/
			// https://github.com/grpc-ecosystem/grpc-gateway/blob/master/runtime/errors.go#L15
			// https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
			c.JSON(httpstatus.FromGRPCCode(st.Code()), response)
			return
		}
	}
}

// JSONPage returns JSON response with pagination
// e.x. {"code":<code>, "msg":<msg>, "data":<data>, "pagination":<pagination>}
// <pagination> { "pageNumber":1, "pageSize":20, "total": 9 }
func (c *Context) JSONPage(data interface{}, p pagination.Pagination) {
	j := new(ResponsePagination)
	j.Code = CodeOK
	j.Data = data
	j.Pagination = p
	c.Context.JSON(http.StatusOK, j)
}

// Bind wraps gin context.Bind() with custom validator
func (c *Context) Bind(obj interface{}) (err error) {
	return validate(c.Context.Bind(obj))
}

// ShouldBind wraps gin context.ShouldBind() with custom validator
func (c *Context) ShouldBind(obj interface{}) (err error) {
	return validate(c.Context.ShouldBind(obj))
}

// RouteNotFound 未找到相关路由
func (c *Context) RouteNotFound() {
	c.String(http.StatusNotFound, "the route not found")
}
