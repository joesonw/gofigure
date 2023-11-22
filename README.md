# GoFigure - a simple yet powerful and extensible configuration library for Golang


[![Lint Status](https://img.shields.io/github/actions/workflow/status/joesonw/gofigure/golangci-lint.yml?branch=master&label=lint&style=flat-square&color=75C46B)](https://github.com/joesonw/gofigure/actions/workflows/golangci-lint.yml)
[![Test Status](https://img.shields.io/github/actions/workflow/status/joesonw/gofigure/gotest.yml?branch=master&label=test&style=flat-square&color=75C46B)](https://github.com/joesonw/gofigure/actions/workflows/gotest.yml)
[![Coverage Status](https://coveralls.io/repos/github/joesonw/gofigure/badge.svg?branch=master)](https://coveralls.io/github/joesonw/gofigure?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/joesonw/gofigure.svg)](https://pkg.go.dev/github.com/joesonw/gofigure)

## Example
`config/app.yaml`
```yaml
env: dev
port: 8080
host: localhost
listen: !tpl |
  {{ config "app.host" }}:{{ config "app.port" }}
db_host: !ref storage.db.host
database: !tpl |
  mysql://{{ config "storage.db.user" }}:{{ config "storage.db.password" }}@{{ config "storage.db.host" }}:{{ config "storage.db.port" }}
external: !include
  file:
    path: external/test.yaml
    parse: true
    key: value
```

`config/storage/db.yaml`
```yaml
host: localhost
port: 3306
user: root
```

`config/prod/app.yaml`
```yaml
env: prod
port: 80
```

`config/prod/storage/db.yaml`
```yaml
host: remote-address
password: supersecret
```

`config/external/test.yaml`
```yaml
value: hello world
```

`main.go`
```go
var defaultYaml []byte // config/app.yaml
var envYaml []byte // config/prod/app.yaml
var defaultDbYaml []byte // config/storage/db.yaml
var envDbYaml []byte // config/prod/storage/db.yaml

loader := gofigure.New().WithFeatures(
	feature.Reference(),
    feature.Template()/*.WithFuncs(template.Funcs{}).WithValeus(map[stirng]any{}) */,
	feature.Include()
)
_ = loader.Load("app.yaml", defaultYaml)
_ = loader.Load("storage/db.yaml", defaultDbYaml)
_ = loader.Load("app.yaml", envYaml)
_ = loader.Load("storage/db.yaml", envDbYaml)
var app struct {
    Env      string `yaml:"env"`
    Listen   string `yaml:"listen"`
    Database string `yaml:"database"`
    External string `yaml:"external"`
}
_ = loader.Get(context.Background(), "app", &app)
fmt.Println(app.Env) // prod
fmt.Println(app.Listen) // localhost:80
fmt.Println(app.Database) // mysql://root:supersecret@remote-address:3306
fmt.Println(app.External) // hello world 
```

## Introduction

GoFigure is a tool to allow maximum flexibility in configuration loading and parsing. It is designed to be simple to use, yet powerful and extensible. It comes with default features like include other files, render a go template and reference other values, etc.

You can easily extend GoFigure with your own features with ease, please check [feature](./feature) for examples.