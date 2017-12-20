package main

import (
	"flag"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/apex/log"
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

	renderReport := make(chan OpsStatus, 10)
	defer close(renderReport)

	go bootstrapHealthCheck(config.HealthCheckPort, renderReport)

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

	duration, err := config.getInterval()
	if err != nil {
		log.WithError(err).Error("Failed to parse interval")
		return
	}

	if *dryRun {
		render(config.OutTemplate, clientset, tmpl, renderReport)
	} else {
		for range time.NewTicker(duration).C {
			err = render(config.OutTemplate, clientset, tmpl, renderReport)
			if err != nil {
				log.WithError(err).Error("Failed to render template")
				continue //we don't bother to exec hooks since the rendering failed
			}
			execHooks(config, renderReport)
		}
	}
}

func bootstrapHealthCheck(port uint32, hookStatus <-chan OpsStatus) {
	lastReport := OpsStatus{isSuccess: true, error: nil}
	http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		lastReport.mux.Lock()
		select {
		case lastReport = <-hookStatus:
			createHealthResponse(lastReport, writer)
		default:
			createHealthResponse(lastReport, writer)
		}
		lastReport.mux.Unlock()
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.WithError(err).Error("Couldn't bootstrap http health server")
		return
	}
}

func createHealthResponse(lastReport OpsStatus, writer http.ResponseWriter) {
	if lastReport.isSuccess {
		fmt.Fprintf(writer, "Healthy !\n")
		writer.WriteHeader(http.StatusOK)
	} else {
		fmt.Fprintf(writer, "Unhealthy: %s !\n", lastReport.error)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

// OpsStatus holds information to track failures/success of render and execHooks functions
// this information gets bubbled up to the health check.
type OpsStatus struct {
	isSuccess bool
	error     error
	mux       sync.Mutex
}

func execHooks(config Config, renderReport chan<- OpsStatus) {
	log.Info("Running post hook")
	out, err := ExecHook(config.Hooks.PostRender)
	if err != nil {
		log.WithError(err).Error("Failed to run post hook")
		renderReport <- OpsStatus{isSuccess: false, error: err}
		return
	}
	log.Info("Output from post hook")
	fmt.Println(out)
	renderReport <- OpsStatus{isSuccess: true, error: nil}
}

func render(outPath string, clientset *kubernetes.Clientset, tmpl *template.Template, renderReport chan<- OpsStatus) error {
	irules, err := ScrapeIngresses(clientset, "")
	cxt := ICxt{IngRules: ToIngressifyRule(irules)}
	err = RenderTemplate(tmpl, outPath, cxt)
	if err != nil {
		renderReport <- OpsStatus{isSuccess: false, error: err}
		return err
	}
	return nil
}
