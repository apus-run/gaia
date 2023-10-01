module github.com/apus-run/gaia/plugins/transport/websocket

go 1.19

require (
	github.com/apus-run/gaia v1.9.0
	github.com/apus-run/sea-kit/encoding v0.0.0-20230930061415-4e76dbc0e8a9
	github.com/apus-run/sea-kit/log v0.0.0-20230930061415-4e76dbc0e8a9
	github.com/google/uuid v1.3.1
	github.com/gorilla/websocket v1.5.0
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/apus-run/gaia => ../../../
