package main

import (
	"k8s.io/api/extensions/v1beta1"
	"hash/fnv"
)

// IngressifyRule is a denormalization of the Ingresses rules coming from k8s
type IngressifyRule struct {
	Hash        uint32
	ServiceName string
	ServicePort int32
	Host        string
	Path        string
	Namespace   string
	Name        string
	IngressRaw  v1beta1.Ingress
}

// ICxt holds data used for rendering
type ICxt struct {
	IngRules []IngressifyRule
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// ToIngressifyRule converts from *v1beta1.IngressList (normalized) to IngressifyRule (denormalized)
func ToIngressifyRule(il *v1beta1.IngressList) []IngressifyRule {
	var ifyrules []IngressifyRule
	for _, ing := range il.Items {
		var ir IngressifyRule
		ir.Namespace = ing.Namespace
		ir.Name = ing.Name
		for _, rule := range ing.Spec.Rules {
			ir.Host = rule.Host
			for _, path := range rule.HTTP.Paths {
				ir.Path = path.Path
				ir.ServiceName = path.Backend.ServiceName
				ir.ServicePort = path.Backend.ServicePort.IntVal
				ir.Hash = hash(ing.Namespace + path.Backend.ServiceName + ir.Host + ir.Path)
				ir.IngressRaw = ing
				ifyrules = append(ifyrules, ir)
			}
		}
	}
	return ifyrules
}

// GroupByHost returns a map of IngressifyRule grouped by host
func GroupByHost(rules []IngressifyRule) map[string][]IngressifyRule {
	m := make(map[string][]IngressifyRule)
	for _, rule := range rules {
		if m[rule.Host] != nil {
			m[rule.Host] = append(m[rule.Host], rule)
		} else {
			m[rule.Host] = []IngressifyRule{rule}
		}
	}
	return m
}

// GroupByPath returns a map of IngressifyRule grouped by path
func GroupByPath(rules []IngressifyRule) map[string][]IngressifyRule {
	m := make(map[string][]IngressifyRule)
	for _, rule := range rules {
		if m[rule.Path] != nil {
			m[rule.Path] = append(m[rule.Path], rule)
		} else {
			m[rule.Path] = []IngressifyRule{rule}
		}
	}
	return m
}
