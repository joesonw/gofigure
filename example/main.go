package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/joesonw/gofigure"
	"github.com/joesonw/gofigure/feature"
)

//go:embed config/app.yaml
var defaultYaml []byte

//go:embed config/prod/app.yaml
var envYaml []byte

//go:embed config/storage/db.yaml
var defaultDBYaml []byte

//go:embed config/prod/storage/db.yaml
var envDBYaml []byte

func main() {
	loader := gofigure.New().WithFeatures(
		feature.Reference(),
		feature.Template(),
		feature.Include(os.DirFS("./config")),
	)
	die(loader.Load("app.yaml", defaultYaml))
	die(loader.Load("storage/db.yaml", defaultDBYaml))
	die(loader.Load("app.yaml", envYaml))
	die(loader.Load("storage/db.yaml", envDBYaml))
	var app struct {
		Env      string `yaml:"env"`
		Listen   string `yaml:"listen"`
		Database string `yaml:"database"`
		External string `yaml:"external"`
	}
	die(loader.Get(context.Background(), "app", &app))
	fmt.Println(app.Env)
	fmt.Println(app.Listen)
	fmt.Println(app.Database)
	fmt.Println(app.External)
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}
