package main

import (
	"flag"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/apex/log"
	"html/template"
	"k8s.io/client-go/kubernetes"
	"strings"
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
	dryRun := flag.Bool("dry-run", false, "Renders template without calling post hooks")
	logHook := flag.Bool("log-hook", false, "Logs output from hooks")
	runOnce := flag.Bool("run-once", false, "Run once and exits")
	flag.Parse()

	config := ReadConfig(*configPath)

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

	if *runOnce {
		render(config, clientset, tmpl, dryRun, logHook)
	} else {
		for range time.NewTicker(duration).C {
			render(config, clientset, tmpl, dryRun, logHook)
		}
	}
}

func render(config Config, clientset *kubernetes.Clientset, tmpl *template.Template, dryRun *bool, logHook *bool) {
	irules, err := ScrapeIngresses(clientset, "")
	cxt := ICxt{IngRules: irules}
	err = RenderTemplate(tmpl, config.OutTemplate, cxt)
	if err != nil {
		log.WithError(err).Error("Failed to render template")
		return
	}
	if !*dryRun {
		log.Info("Running post hook")
		out, err := ExecHook(config.Hooks.PostRender)
		if err != nil {
			log.WithError(err).Error("Failed to run post hook")
			return
		}
		if *logHook {
			log.Info("Output from post hook")
			fmt.Println(out)
		}
	}
}
