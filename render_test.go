package main

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"html/template"
	"io/ioutil"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

var rules = generateRules()

func TestRenderNginxTemplate(t *testing.T) {
	client := fake.NewSimpleClientset(&rules)
	irules, err := ScrapeIngresses(client, "")
	if err != nil {
		panic(err)
	}

	config := ReadConfig("./examples/config_for_nginx.yaml")

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

func TestRenderHaproxyTemplate(t *testing.T) {
	client := fake.NewSimpleClientset(&rules)
	irules, err := ScrapeIngresses(client, "")
	if err != nil {
		panic(err)
	}

	config := ReadConfig("./examples/config_for_haproxy.yaml")

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

	haproxyActual, err := ioutil.ReadFile("/tmp/haproxy.actual")
	if err != nil {
		panic(err)
	}
	haproxyExpected, err := ioutil.ReadFile("./examples/haproxy.expected")
	if err != nil {
		panic(err)
	}

	if string(haproxyActual) != string(haproxyExpected) {
		t.Errorf("Template results differ")
		fmt.Println("Expected:")
		fmt.Printf("%s\n", haproxyExpected)
		fmt.Println("----------------------")
		fmt.Println("Got: ")
		fmt.Printf("%s\n", haproxyActual)
	}

}
