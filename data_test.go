package main

import (
	"fmt"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
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

type IngressRule v1beta1.IngressRule
type Ingress v1beta1.Ingress

func buildIngressRule() IngressRule {
	return IngressRule{IngressRuleValue: v1beta1.IngressRuleValue{HTTP: &v1beta1.HTTPIngressRuleValue{Paths: []v1beta1.HTTPIngressPath{}}}}
}

func (ir IngressRule) build() v1beta1.IngressRule {
	return v1beta1.IngressRule(ir)
}

func buildIngress() Ingress {
	return Ingress{}
}

func (i Ingress) build() v1beta1.Ingress {
	return v1beta1.Ingress(i)
}

func (i Ingress) withName(name string) Ingress {
	i.Name = name
	return i
}

func (i Ingress) withNamespace(namespace string) Ingress {
	i.Namespace = namespace
	return i
}

func (i Ingress) withRule(ir v1beta1.IngressRule) Ingress {
	i.Spec.Rules = append(i.Spec.Rules, ir)
	return i
}

func (ir IngressRule) withHost(host string) IngressRule {
	ir.Host = host
	return ir
}

func (ir IngressRule) withPathBackend(path string, svcname string, svcport int32) IngressRule {
	x := v1beta1.HTTPIngressPath{Path: path, Backend: v1beta1.IngressBackend{ServiceName: svcname,
		ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: svcport}}}
	ir.HTTP.Paths = append(ir.HTTP.Paths, x)
	return ir
}

var testRules = v1beta1.IngressList{Items: []v1beta1.Ingress{
	buildIngress().
		withName("n1").
		withNamespace("ns1").
		withRule(buildIngressRule().
			withHost("h1").
			withPathBackend("/svc_n1", "svc1", 30400).
			withPathBackend("/svc_n1-2", "svc2", 30401).
			build()).
		withRule(buildIngressRule().
			withHost("h1-1").
			withPathBackend("", "svc1-1", 30402).
			build()).
		build(),
	buildIngress().
		withName("n2").
		withNamespace("ns2").
		withRule(buildIngressRule().
			withHost("").
			withPathBackend("/svc_2", "svc2", 30403).
			build()).
		withRule(buildIngressRule().
			withHost("h2").
			withPathBackend("/svc_n2", "svc2", 30404).
			build()).
		build(),
}}

//TODO missing checking if raw is present. Use deepcompare from reflect
func TestToIngressifyRule(t *testing.T) {
	testRuleCopy := testRules.DeepCopy()
	ingressifyRules := ToIngressifyRule(testRuleCopy)
	if numIngRules, numIngTestRules := len(ingressifyRules), sizeIngTest(*testRuleCopy); numIngRules != numIngTestRules {
		t.Errorf("IngressifyRules size, got %s, expected %s", numIngRules, numIngTestRules)
	}
	for _, ingTestRule := range testRules.Items {
		if ok, missingRule := checkIntegrity(ingTestRule, ingressifyRules); !ok {
			t.Errorf("Missing rule: %s", missingRule)
		}
	}
}

func TestGroupByHost(t *testing.T) {
	testRulesCopy := testRules.DeepCopy()
	ingressifyRules := ToIngressifyRule(testRulesCopy)
	byHost := GroupByHost(ingressifyRules)
	// All hosts are in the map
	for _, r := range ingressifyRules {
		if _, ok := byHost[r.Host]; !ok {
			var h string = r.Host
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
	testRuleCopy := testRules.DeepCopy()
	ingressifyRules := ToIngressifyRule(testRuleCopy)
	byPath := GroupByPath(ingressifyRules)
	// All paths are in the map
	for _, r := range ingressifyRules {
		if _, ok := byPath[r.Path]; !ok {
			var p string = r.Path
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
			if !isRulePresent(ing.Name, ing.Namespace, r.Host, p.Path, p.Backend, irs) {
				return false, fmt.Sprintf("Name: %s, Namespace: %s, Host: %s, Path: %s, ServicePort: %d, "+
					"ServiceName: %s", ing.Name, ing.Namespace, r.Host, p.Path, p.Backend.ServicePort.IntVal, p.Backend.ServiceName)
			}
		}
	}
	return true, ""
}

func isRulePresent(name string, namespace string, host string, path string, backend v1beta1.IngressBackend, irs []IngressifyRule) bool {
	for _, r := range irs {
		if r.Name == name && r.Namespace == namespace && r.Host == host && r.Path == path &&
			r.ServiceName == backend.ServiceName && r.ServicePort == backend.ServicePort.IntVal {
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
				count += 1
				continue
			}
			count += len(tr.IngressRuleValue.HTTP.Paths)
		}
	}
	return count
}
