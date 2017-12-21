package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func handlerBuilder() HealthHandler {
	var opsReport = make(chan *OpsStatus, 10)
	defaultReport := OpsStatus{isSuccess: true, timestamp: time.Now()}
	hhandler := HealthHandler{opsStatus: opsReport, cacheExpirationTime: REFRESH_INTERVAL, lastReport: &defaultReport}
	return hhandler
}

func TestBootstrapHealthCheck_should_return_initial_cache_when_no_report(t *testing.T) {
	hhandler := handlerBuilder()
	r, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	hhandler.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Errorf("wrong code returned")
	}
}

func TestBootstrapHealthCheck_should_return_500_when_cache_has_expired(t *testing.T) {
	hhandler := handlerBuilder()
	r, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	hhandler.ServeHTTP(w, r)
	time.Sleep(REFRESH_INTERVAL) //sleep interval so the cache expires
	r, _ = http.NewRequest("GET", "/health", nil)
	w = httptest.NewRecorder()
	hhandler.ServeHTTP(w, r)
	if w.Code != 500 {
		t.Errorf("Should return 500 when cache has expired, got: %d, expected %d", w.Code, 500)
	}
}

func TestBootstrapHealthCheck_should_not_return_cache_when_reporting_ops(t *testing.T) {
	hhandler := handlerBuilder()
	statusError := OpsStatus{isSuccess: false, timestamp: time.Now(), error: errors.New("This one failed")}
	hhandler.opsStatus <- &statusError
	r, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	hhandler.ServeHTTP(w, r)
	if w.Code != 500 {
		t.Errorf("Should return 500 from the published status error, got: %d, expected %d", w.Code, 500)
	}
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("Should return non empty body. Something went wrong: %s", err)
	}
	expectedBody := fmt.Sprintf("Unhealthy: %s !\n", statusError.error)
	if string(body) != expectedBody {
		t.Errorf("Body is wrong, got: %s, expected: %s", string(body), expectedBody)
	}
}

func TestBootstrapHealthCheck_should_return_last_report_as_cache_when_no_report(t *testing.T) {
	hhandler := handlerBuilder()
	statusError := OpsStatus{isSuccess: false, timestamp: time.Now(), error: errors.New("This one failed")}
	hhandler.opsStatus <- &statusError
	r, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	hhandler.ServeHTTP(w, r)
	if w.Code != 500 {
		t.Errorf("Should return 500 from the published status error, got: %d, expected %d", w.Code, 500)
	}
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("Should return non empty body. Something went wrong: %s", err)
	}
	expectedBody := fmt.Sprintf("Unhealthy: %s !\n", statusError.error)
	if string(body) != expectedBody {
		t.Errorf("Body is wrong, got: %s, expected: %s", string(body), expectedBody)
	}
	// call again without making report should give us the last response again
	r, _ = http.NewRequest("GET", "/health", nil)
	w = httptest.NewRecorder()
	hhandler.ServeHTTP(w, r)
	body, err = ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("Should return 500 when cache has expired, got: %d, expected %d", w.Code, 500)
	}
	expectedBody = fmt.Sprintf("Unhealthy: %s !\n", statusError.error)
	if string(body) != expectedBody {
		t.Errorf("Body is wrong, got: %s, expected: %s", string(body), expectedBody)
	}
}

func TestBootstrapHealthCheck_should_not_return_cache_after_report_and_cache_expiration(t *testing.T) {
	hhandler := handlerBuilder()
	statusError := OpsStatus{isSuccess: false, timestamp: time.Now(), error: errors.New("This one failed")}
	hhandler.opsStatus <- &statusError
	r, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	hhandler.ServeHTTP(w, r)
	if w.Code != 500 {
		t.Errorf("Should return 500 from the published status error, got: %d, expected %d", w.Code, 500)
	}
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("Should return non empty body. Something went wrong: %s", err)
	}
	expectedBody := fmt.Sprintf("Unhealthy: %s !\n", statusError.error)
	if string(body) != expectedBody {
		t.Errorf("Body is wrong, got: %s, expected: %s", string(body), expectedBody)
	}
	time.Sleep(REFRESH_INTERVAL)
	// call again without making report should give us the last response again
	r, _ = http.NewRequest("GET", "/health", nil)
	w = httptest.NewRecorder()
	hhandler.ServeHTTP(w, r)
	body, err = ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("Should return 500 when cache has expired, got: %d, expected %d", w.Code, 500)
	}
	expectedBody = fmt.Sprintf("Unhealthy: %s !\n", statusError.error)
	if string(body) == expectedBody {
		t.Errorf("Body is wrong, got: %s, expected: %s", string(body), expectedBody)
	}
}
