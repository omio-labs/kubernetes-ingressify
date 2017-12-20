package main

import (
	"testing"
	"k8s.io/client-go/kubernetes/fake"
)

func TestScrapeIngressesForAllNamespaces(t *testing.T) {
	rules := generateRules("./examples/ingressList.json")
	client := fake.NewSimpleClientset(&rules)
	irules, err := ScrapeIngresses(client, "")
	if err != nil {
		t.Errorf("Something went wrong scraping ingress for rules: %s\n", err)
		return
	}
	if irules.Size() != 2 {
		t.Errorf("Didn't scrape all rules, got: %d, expected: %d ", irules.Size(), 2)
	}
}
