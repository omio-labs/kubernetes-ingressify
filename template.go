package main

import (
	"github.com/apex/log"
	"html/template"
	"io/ioutil"
	"os"
)

func readTemplate(tmplpath string) ([]byte, error) {
	tmpl, err := ioutil.ReadFile(tmplpath)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// BuildFuncMap merges template.FuncMap's
func BuildFuncMap(funcs ...template.FuncMap) template.FuncMap {
	resmap := make(template.FuncMap)
	for _, fumap := range funcs {
		for k, v := range fumap {
			resmap[k] = v
		}
	}
	return resmap
}

// PrepareTemplate creates a template from `tmplpath` initialized with `withfuncs`
func PrepareTemplate(tmplpath string, withfuncs template.FuncMap) (*template.Template, error) {
	tmplstr, err := readTemplate(tmplpath)
	if err != nil {
		return nil, err
	}
	tmpl := template.Must(template.New("template").Funcs(withfuncs).Parse(string(tmplstr)))
	return tmpl, nil
}

// RenderTemplate renders the template and writes the output to `outpath`
func RenderTemplate(tmpl *template.Template, outpath string, cxt ICxt) error {
	output, err := os.Create(outpath)
	defer output.Close()
	if err != nil {
		log.WithError(err).Error("Failed to render template")
		return err
	}
	log.Info("Rendering template")
	err = tmpl.Execute(output, cxt)
	if err != nil {
		log.WithError(err).Error("Failed to render template")
		return err
	}
	log.Info("Template successfully rendered")
	return nil
}
