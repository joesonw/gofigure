package feature

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"

	"github.com/joesonw/gofigure"
)

var _ gofigure.Feature = (*TemplateFeature)(nil)

type TemplateFeature struct {
	funcs  template.FuncMap
	values map[string]any
}

func Template() *TemplateFeature {
	return &TemplateFeature{
		funcs:  template.FuncMap{},
		values: map[string]any{},
	}
}

func (f *TemplateFeature) Funcs(funcs template.FuncMap) *TemplateFeature {
	f.funcs = funcs
	return f
}

func (f *TemplateFeature) Valeus(values map[string]any) *TemplateFeature {
	f.values = values
	return f
}

func (*TemplateFeature) Name() string {
	return "!tpl"
}

func (f *TemplateFeature) Resolve(ctx context.Context, loader *gofigure.Loader, node *gofigure.Node) (*gofigure.Node, error) {
	if node.Kind() != yaml.ScalarNode {
		return nil, fmt.Errorf("!tpl only supports scalar node")
	}

	tpl, err := template.New(node.Keypath()).
		Funcs(f.funcs).
		Funcs(template.FuncMap{
			"config": func(path string) (string, error) {
				result, err := loader.GetNode(ctx, path)
				if err != nil {
					return "", err
				}
				if result == nil {
					return "", nil
				}

				b, err := yaml.Marshal(result)
				if err != nil {
					return "", err
				}
				return strings.TrimSpace(string(b)), nil
			},
			"env": os.Getenv,
		}).
		Parse(node.Value())
	if err != nil {
		return nil, fmt.Errorf("unable to parse template: %w", err)
	}

	var w bytes.Buffer
	if err := tpl.Execute(&w, f.values); err != nil {
		return nil, fmt.Errorf("unable to execute template: %w", err)
	}

	return gofigure.NewScalarNode(strings.TrimSpace(w.String())), nil
}
