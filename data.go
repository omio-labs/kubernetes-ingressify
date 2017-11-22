package main

import (
	"k8s.io/api/extensions/v1beta1"
)

type IngressifyRule struct {
	ServiceName string
	ServicePort int32
	Host        string
	Path        string
	Namespace   string
	Name        string
	IngressRaw  v1beta1.Ingress
}

type ICxt struct {
	IngRules []IngressifyRule
}

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
				ir.IngressRaw = ing
				ifyrules = append(ifyrules, ir)
			}
		}
	}
	return ifyrules
}

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
