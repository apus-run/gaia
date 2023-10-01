module github.com/apus-run/gaia/plugins/transport/fasthttp

go 1.19

require (
	github.com/apus-run/gaia v0.9.0
	github.com/apus-run/sea-kit/log v0.0.0-20230930061415-4e76dbc0e8a9
	github.com/fasthttp/router v1.4.20
	github.com/valyala/fasthttp v1.49.0
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/klauspost/compress v1.16.6 // indirect
	github.com/savsgio/gotils v0.0.0-20230208104028-c358bd845dee // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
)

replace github.com/apus-run/gaia => ../../../
