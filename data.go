package main

import (
	"hash/fnv"
	"reflect"
	"strings"

	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
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
				ir.Hash = hash(ing.Namespace + ing.Name + path.Backend.ServiceName + ir.Host + ir.Path)
				ir.IngressRaw = ing
				ifyrules = append(ifyrules, ir)
			}
		}
	}
	return ifyrules
}

// GroupByHost returns a map of IngressifyRule grouped by ir.Host
func GroupByHost(rules []IngressifyRule) map[string][]IngressifyRule {
	return groupByGeneric(rules, "Host")
}

// GroupByPath returns a map of IngressifyRule grouped by ir.Path
func GroupByPath(rules []IngressifyRule) map[string][]IngressifyRule {
	return groupByGeneric(rules, "Path")
}

// GroupByName returns a map of IngressifyRule grouped by ir.Name
func GroupByName(rules []IngressifyRule) map[string][]IngressifyRule {
	return groupByGeneric(rules, "ServiceName", "Namespace")
}

func groupByGeneric(rules []IngressifyRule, fields ...string) map[string][]IngressifyRule {
	m := make(map[string][]IngressifyRule)
	for _, rule := range rules {
		value := getFieldString(&rule, fields...)
		if m[value] != nil {
			m[value] = append(m[value], rule)
		} else {
			m[value] = []IngressifyRule{rule}
		}
	}
	return m
}

func getFieldString(ir *IngressifyRule, fields ...string) string {
	r := reflect.ValueOf(ir)
	var key string
	for _, field := range fields {
		f := reflect.Indirect(r).FieldByName(field)
		key = key + "-" + f.String()
	}
	return strings.TrimPrefix(key, "-")
}
