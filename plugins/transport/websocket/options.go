package websocket

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/apus-run/sea-kit/encoding"
	ws "github.com/gorilla/websocket"
)

type PayloadType uint8

const (
	PayloadTypeBinary = 0
	PayloadTypeText   = 1
)

type Server struct {
	*http.Server

	lis      net.Listener
	tlsConf  *tls.Config
	upgrader *ws.Upgrader

	network     string
	address     string
	path        string
	strictSlash bool

	timeout time.Duration

	err   error
	codec encoding.Codec

	messageHandlers MessageHandlerMap

	sessionMgr *SessionManager

	register   chan *Session
	unregister chan *Session

	payloadType PayloadType
}

// defaultServer return a default config server
func defaultServer() *Server {
	return &Server{
		network:     "tcp",
		address:     ":0",
		timeout:     1 * time.Second,
		strictSlash: true,
		path:        "/",

		messageHandlers: make(MessageHandlerMap),

		sessionMgr: NewSessionManager(),
		upgrader: &ws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},

		register:   make(chan *Session),
		unregister: make(chan *Session),

		payloadType: PayloadTypeBinary,
	}
}

type ServerOption func(o *Server)

func WithNetwork(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func WithPath(path string) ServerOption {
	return func(s *Server) {
		s.path = path
	}
}

func WithConnectHandle(h ConnectHandler) ServerOption {
	return func(s *Server) {
		s.sessionMgr.RegisterConnectHandler(h)
	}
}

func WithTLSConfig(c *tls.Config) ServerOption {
	return func(o *Server) {
		o.tlsConf = c
	}
}

func WithListener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

func WithCodec(c string) ServerOption {
	return func(s *Server) {
		s.codec = encoding.GetCodec(c)
	}
}

func WithChannelBufferSize(size int) ServerOption {
	return func(_ *Server) {
		channelBufSize = size
	}
}

func WithPayloadType(payloadType PayloadType) ServerOption {
	return func(s *Server) {
		s.payloadType = payloadType
	}
}

////////////////////////////////////////////////////////////////////////////////

type ClientOption func(o *Client)

func WithClientCodec(c string) ClientOption {
	return func(o *Client) {
		o.codec = encoding.GetCodec(c)
	}
}

func WithEndpoint(uri string) ClientOption {
	return func(o *Client) {
		o.url = uri
	}
}

func WithClientPayloadType(payloadType PayloadType) ClientOption {
	return func(c *Client) {
		c.payloadType = payloadType
	}
}
