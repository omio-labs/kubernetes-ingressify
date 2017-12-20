package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"reflect"
	"testing"
)

type Set struct {
	Set map[string]bool
}

func NewSet() Set {
	return Set{make(map[string]bool)}
}

func (f Set) Add(member string) {
	f.Set[member] = true
}

func (f Set) IsMember(x string) bool {
	_, prs := f.Set[x]
	return prs
}

func generateRules(from string) v1beta1.IngressList {
	test, err := ioutil.ReadFile(from)
	il := v1beta1.IngressList{}
	err = json.Unmarshal(test, &il)
	if err != nil {
		panic(err)
	}
	return il
}

func TestToIngressifyRule(t *testing.T) {
	testRules := generateRules("./examples/ingressList.json")
	ingressifyRules := ToIngressifyRule(&testRules)
	if numIngRules, numIngTestRules := len(ingressifyRules), sizeIngTest(testRules); numIngRules != numIngTestRules {
		t.Errorf("IngressifyRules size, got %s, expected %s", numIngRules, numIngTestRules)
	}
	for _, ingTestRule := range testRules.Items {
		if ok, missingRule := checkIntegrity(ingTestRule, ingressifyRules); !ok {
			t.Errorf("Missing rule: %s", missingRule)
		}
	}
}

func TestGroupByHost(t *testing.T) {
	testRules := generateRules("./examples/ingressList.json")
	ingressifyRules := ToIngressifyRule(&testRules)
	byHost := GroupByHost(ingressifyRules)
	// All hosts are in the map
	for _, r := range ingressifyRules {
		if _, ok := byHost[r.Host]; !ok {
			var h = r.Host
			if h == "" {
				h = "\"\""
			}
			t.Errorf("Host %s not found !", h)
		}
	}
	// There are exactly as many keys as hosts
	set := NewSet()
	for _, r := range ingressifyRules {
		set.Add(r.Host)
	}
	if len(set.Set) != len(byHost) {
		t.Errorf("Number of hosts is different than number of keys, got %d, expected %d", len(byHost), len(set.Set))
	}
	// All IngressifyRules are mapped
	for _, r := range ingressifyRules {
		mr, _ := byHost[r.Host]
		if !isIngressifyRulePresent(r, mr) {
			t.Errorf("Missing rule, Name: %s, Namespace: %s, Host: %s, Path: %s, ServicePort: %d, ServiceName: %s",
				r.Name, r.Namespace, r.Host, r.Path, r.ServicePort, r.ServiceName)
		}
	}
}

func TestGroupByPath(t *testing.T) {
	testRules := generateRules("./examples/ingressList.json")
	ingressifyRules := ToIngressifyRule(&testRules)
	byPath := GroupByPath(ingressifyRules)
	// All paths are in the map
	for _, r := range ingressifyRules {
		if _, ok := byPath[r.Path]; !ok {
			var p = r.Path
			if p == "" {
				p = "\"\""
			}
			t.Errorf("Path %s not found !", p)
		}
	}
	// There are exactly as many keys as paths
	set := NewSet()
	for _, r := range ingressifyRules {
		set.Add(r.Path)
	}
	if len(set.Set) != len(byPath) {
		t.Errorf("Number of paths is different than number of keys, got %d, expected %d", len(byPath), len(set.Set))
	}
	// All IngressifyRules are mapped
	for _, r := range ingressifyRules {
		mr, _ := byPath[r.Path]
		if !isIngressifyRulePresent(r, mr) {
			t.Errorf("Missing rule, Name: %s, Namespace: %s, Host: %s, Path: %s, ServicePort: %d, ServiceName: %s",
				r.Name, r.Namespace, r.Host, r.Path, r.ServicePort, r.ServiceName)
		}
	}
}

func isIngressifyRulePresent(ir IngressifyRule, irs []IngressifyRule) bool {
	for _, r := range irs {
		if ir.Namespace == r.Namespace && ir.Name == r.Name && ir.ServicePort == r.ServicePort &&
			ir.ServiceName == r.ServiceName && ir.Path == r.Path {
			return true
		}
	}
	return false
}

func checkIntegrity(ing v1beta1.Ingress, irs []IngressifyRule) (bool, string) {
	for _, r := range ing.Spec.Rules {
		for _, p := range r.IngressRuleValue.HTTP.Paths {
			if !isRulePresent(ing.Name, ing.Namespace, r.Host, p.Path, p.Backend, ing, irs) {
				return false, fmt.Sprintf("Name: %s, Namespace: %s, Host: %s, Path: %s, ServicePort: %d, "+
					"ServiceName: %s", ing.Name, ing.Namespace, r.Host, p.Path, p.Backend.ServicePort.IntVal, p.Backend.ServiceName)
			}
		}
	}
	return true, ""
}

func isRulePresent(name string, namespace string, host string, path string, backend v1beta1.IngressBackend, raw v1beta1.Ingress, irs []IngressifyRule) bool {
	for _, r := range irs {
		if r.Name == name && r.Namespace == namespace && r.Host == host && r.Path == path &&
			r.ServiceName == backend.ServiceName && r.ServicePort == backend.ServicePort.IntVal &&
			reflect.DeepEqual(r.IngressRaw, raw) {
			return true
		}
	}
	return false
}

func sizeIngTest(ingListTest v1beta1.IngressList) int {
	var count = 0
	for _, ing := range ingListTest.Items {
		for _, tr := range ing.Spec.Rules {
			if tr.Host != "" && len(tr.IngressRuleValue.HTTP.Paths) == 0 {
				count++
				continue
			}
			count += len(tr.IngressRuleValue.HTTP.Paths)
		}
	}
	return count
}
