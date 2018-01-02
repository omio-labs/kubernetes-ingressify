package main

import (
	"fmt"
	"net/http"
	"testing"
)

func init() {
	var opsReport = make(chan *OpsStatus, 10)
	go runHealthCheckServer(opsReport, REFRESHINTERVAL, PORT)
}

func BenchmarkBootstrapHealthCheck(b *testing.B) {
	for n := 0; n < b.N; n++ {
		http.Get(fmt.Sprintf("http://localhost:%d/health", PORT))
	}
}
