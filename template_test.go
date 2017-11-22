package main

import (
	"fmt"
	"html/template"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestBuildFuncMap(t *testing.T) {
	funmap1 := template.FuncMap{
		"Join":    strings.Join,
		"Compare": strings.Compare,
	}
	funmap2 := template.FuncMap{
		"GroupByHost": GroupByHost,
		"GroupByPath": GroupByPath,
	}
	funmap3 := template.FuncMap{}
	funMap := BuildFuncMap(funmap1, funmap2)
	// Must contain all keys and values
	if ok, reason := checkFunctionsInMap(funMap, funmap1); !ok {
		t.Errorf(reason)
	}
	if ok, reason := checkFunctionsInMap(funMap, funmap2); !ok {
		t.Errorf(reason)
	}
	//The length of resulting map is the sum of the length of the input maps
	if len(funMap) != len(funmap1)+len(funmap2) {
		t.Errorf("The lenght of resulting map is not the sum of the length of the input maps, expected %d, got %d", len(funmap1)+len(funmap2), len(funMap))
	}
	//Building with empty returns empty map
	if funMap = BuildFuncMap(funmap3); len(funMap) != 0 {
		t.Errorf("Building with empty map returns empty map, expected 0, got %d", len(funMap))
	}
}

func checkFunctionsInMap(funMap template.FuncMap, testData template.FuncMap) (bool, string) {
	for k, v := range testData {
		if funVal, ok := funMap[k]; !ok || getFunctionPointer(funVal) != getFunctionPointer(v) {
			return false, fmt.Sprintf("Missing key %s or function %s", k, getFunctionName(v))
		}
	}
	return true, ""
}

func getFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func getFunctionPointer(f interface{}) uintptr {
	return reflect.ValueOf(f).Pointer()
}
