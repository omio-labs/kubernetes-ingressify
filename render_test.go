package main

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"html/template"
	"io/ioutil"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func runRenderFor(router string) (actual string, expected string) {
	rules := generateRules("./examples/ingressList.json")
	client := fake.NewSimpleClientset(&rules)
	irules, err := ScrapeIngresses(client, "")
	if err != nil {
		panic(err)
	}

	config := ReadConfig(fmt.Sprintf("./examples/config_for_%s.yaml", router))

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

	actualRes, err := ioutil.ReadFile(fmt.Sprintf("/tmp/%s.actual", router))
	if err != nil {
		panic(err)
	}
	expectedRes, err := ioutil.ReadFile(fmt.Sprintf("./examples/%s.expected", router))
	if err != nil {
		panic(err)
	}
	actual = string(actualRes)
	expected = string(expectedRes)
	return
}

func TestRenderNginxTemplate(t *testing.T) {
	nginxActual, nginxExpected := runRenderFor("nginx")

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
	haproxyActual, haproxyExpected := runRenderFor("haproxy")

	if string(haproxyActual) != string(haproxyExpected) {
		t.Errorf("Template results differ")
		fmt.Println("Expected:")
		fmt.Printf("%s\n", haproxyExpected)
		fmt.Println("----------------------")
		fmt.Println("Got: ")
		fmt.Printf("%s\n", haproxyActual)
	}

}
