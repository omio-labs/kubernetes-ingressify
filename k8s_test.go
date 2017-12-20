package main

import (
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestScrapeIngressesForAllNamespaces(t *testing.T) {
	rules := generateRules("./examples/ingressList.json")
	client := fake.NewSimpleClientset(&rules)
	irules, err := ScrapeIngresses(client, "")
	if err != nil {
		t.Errorf("Something went wrong scraping ingress for rules: %s\n", err)
		return
	}
	if len(irules.Items) != 2 {
		t.Errorf("Didn't scrape all rules, got: %d, expected: %d ", irules.Size(), 2)
	}
}

func TestScrapeIngressesForAGivenNamespace(t *testing.T) {
	rules := generateRules("./examples/ingressList.json")
	client := fake.NewSimpleClientset(&rules)
	irules, err := ScrapeIngresses(client, "ns1")
	if err != nil {
		t.Errorf("Something went wrong scraping ingress for rules: %s\n", err)
		return
	}
	if len(irules.Items) != 1 {
		t.Errorf("Didn't scrape all rules, got: %d, expected: %d ", irules.Size(), 2)
	}
}
