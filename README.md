# GoFigure - a simple yet powerful and extensible configuration library for Golang

## Example
`config/app.yaml`
```yaml
env: dev
port: 8080
host: localhost
listen: !tpl {{.host}}:{{.port}}
```

`config/prod/app.yaml`
```yaml
env: prod
port: 80
```

`main.go`
```go
var defaultYaml []byte // config/app.yaml
var envYaml []byte // config/prod/app.yaml

loader := gofigure.NewLoader(feature.All()) // or you can manually pick individual features and with options
_ = loader.Load("app.yaml", defaultYaml)
_ = loader.Load("app.yaml", envYaml)
var app struct {
	Env string yaml "env"
	Listen string yaml "listen"
}
_ = loader.Get(context.Background(), "app", &app)
fmt.Println(app.Listen) // localhost:80
fmt.Println(app.env) // prod
```

## Introduction

GoFigure is a tool to allow maximum flexibility in configuration loading and parsing. It is designed to be simple to use, yet powerful and extensible. It comes with default features like include other files, render a go template and reference other values, etc.

You can easily extend GoFigure with your own features with ease, please check [feature](./feature) for examples.