package main

import (
	"flag"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/apex/log"
	"github.com/pkg/errors"
	"html/template"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {
	data, err := Asset("gen/version")
	if err != nil {
		log.WithError(err).Error("error")
	}

	version := strings.TrimSpace(string(data))
	log.Infof("kubernetes-ingressify version %s", version)

	configPath := flag.String("config", "", "path to the config file")
	dryRun := flag.Bool("dry-run", false, "Run once without hooks and exits")
	flag.Parse()

	config := ReadConfig(*configPath)

	opsStatus := make(chan *OpsStatus, 10)
	defer close(opsStatus)

	duration, err := config.getInterval()
	if err != nil {
		log.WithError(err).Error("Failed to parse interval")
		return
	}

	fmap := template.FuncMap{
		"GroupByHost": GroupByHost,
		"GroupByPath": GroupByPath,
	}

	clientset, err := GetKubeClient(config.Kubeconfig)
	if err != nil {
		log.WithError(err).Error("Failed to build k8s client")
		return
	}

	tmpl, err := PrepareTemplate(config.InTemplate, BuildFuncMap(fmap, sprig.FuncMap()))
	if err != nil {
		log.WithError(err).Error("Failed to prepare template")
		return
	}

	go func() {
		if *dryRun {
			render(config.OutTemplate, clientset, tmpl)
		} else {
			for range time.NewTicker(duration).C {
				err = render(config.OutTemplate, clientset, tmpl)
				if err != nil {
					log.WithError(err).Error("Failed to render template")
					opsStatus <- &OpsStatus{isSuccess: false, timestamp: time.Now(), error: err}
					continue //we don't bother to exec hooks since the rendering failed
				}
				err = execHooks(config, opsStatus)
				if err != nil {
					opsStatus <- &OpsStatus{isSuccess: false, timestamp: time.Now(), error: err}
					continue
				}
				opsStatus <- &OpsStatus{isSuccess: true, timestamp: time.Now()}
			}
		}
	}()
	log.WithError(runHealthCheckServer(opsStatus, duration, config.HealthCheckPort)).Error("Health server is down...")
}

func runHealthCheckServer(status chan *OpsStatus, duration time.Duration, port uint32) error {
	lastReport := OpsStatus{isSuccess: true, timestamp: time.Now()}
	healthHandler := healthHandler{opsStatus: status, cacheExpirationTime: duration, lastReport: &lastReport}
	http.HandleFunc("/health", healthHandler.ServeHTTP)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

type healthHandler struct {
	opsStatus           chan *OpsStatus
	cacheExpirationTime time.Duration
	lastReport          *OpsStatus
	sync.Mutex
}

func (hh healthHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	hh.Lock()
	select {
	case currentReport := <-hh.opsStatus:
		*hh.lastReport = *currentReport
		createHealthResponse(*currentReport, writer)
	default:
		if time.Now().Sub(hh.lastReport.timestamp) > hh.cacheExpirationTime {
			createHealthResponse(OpsStatus{
				isSuccess: false, error: errors.New("Seems that k8s-ingressify is stuck")}, writer)
		} else {
			createHealthResponse(*hh.lastReport, writer)
		}
	}
	hh.Unlock()
}

func createHealthResponse(lastReport OpsStatus, writer http.ResponseWriter) {
	if lastReport.isSuccess {
		writer.WriteHeader(http.StatusOK)
		fmt.Fprint(writer, "Healthy !\n")
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "Unhealthy: %s !\n", lastReport.error)
	}
}

// OpsStatus holds information to track failures/success of render and execHooks functions
// this information gets bubbled up to the health check.
type OpsStatus struct {
	isSuccess bool
	error     error
	timestamp time.Time
}

func execHooks(config Config, renderReport chan<- *OpsStatus) error {
	log.Info("Running post hook")
	out, err := ExecHook(config.Hooks.PostRender)
	if err != nil {
		log.WithError(err).Error("Failed to run post hook")
		return err
	}
	log.Info("Output from post hook")
	fmt.Println(out)
	return nil
}

func render(outPath string, clientset *kubernetes.Clientset, tmpl *template.Template) error {
	irules, err := ScrapeIngresses(clientset, "")
	cxt := ICxt{IngRules: ToIngressifyRule(irules)}
	err = RenderTemplate(tmpl, outPath, cxt)
	if err != nil {
		return err
	}
	return nil
}
