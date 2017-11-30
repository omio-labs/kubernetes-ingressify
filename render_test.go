package main

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"html/template"
	"io/ioutil"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestRenderNginxTemplate(t *testing.T) {
	client := fake.NewSimpleClientset(&testRules)
	irules, err := ScrapeIngresses(client, "")
	if err != nil {
		panic(err)
	}

	config := ReadConfig("./examples/config.yaml")

	fmap := template.FuncMap{
		"GroupByHost": GroupByHost,
		"GroupByPath": GroupByPath,
	}

	tmpl, err := PrepareTemplate(config.InTemplate, BuildFuncMap(fmap, sprig.FuncMap()))
	if err != nil {
		panic(err)
	}

	cxt := ICxt{IngRules: irules}
	err = RenderTemplate(tmpl, config.OutTemplate, cxt)

	nginxActual, err := ioutil.ReadFile("/tmp/nginx.actual")
	if err != nil {
		panic(err)
	}
	nginxExpected, err := ioutil.ReadFile("./examples/nginx.expected")
	if err != nil {
		panic(err)
	}

	if string(nginxActual) != string(nginxExpected) {
		t.Errorf("Template results differ")
		fmt.Println("Expected:")
		fmt.Printf("%s\n", nginxExpected)
		fmt.Println("----------------------")
		fmt.Println("Got: ")
		fmt.Printf("%s\n", nginxActual)
	}

}
