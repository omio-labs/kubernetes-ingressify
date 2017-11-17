package main

import (
	"flag"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/apex/log"
	"html/template"
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

	configpath := flag.String("config", "", "path to the config file")
	dryrun := flag.Bool("dry-run", false, "Renders template without calling post hooks")
	loghook := flag.Bool("log-hook", false, "Logs output from hooks")
	flag.Parse()
	//TODO handle wrong flags, show usage, help command, dry-run, debug mode

	config := ReadConfig(*configpath)
	//TODO schema validation for yaml configuration ?

	fmap := template.FuncMap{
		"GroupByHost": GroupByHost,
		"GroupByPath": GroupByPath,
	}

	tmpl, err := PrepareTemplate(config.InTemplate, BuildFuncMap(fmap, sprig.FuncMap()))
	if err != nil {
		log.WithError(err).Error("Failed to prepare template")
		return
	}

	duration, err := config.GetInterval()
	if err != nil {
		log.WithError(err).Error("Failed to parse interval")
		return
	}

	for range time.NewTicker(duration).C {
		log.Info("Reloading configuration")
		irules, err := ScrapeIngresses(config.Kubeconfig, "")
		cxt := ICxt{IngRules: irules}
		err = RenderTemplate(tmpl, config.OutTemplate, cxt)
		if err != nil {
			log.WithError(err).Error("Failed to render template")
			return
		}
		if !*dryrun {
			log.Info("Running post hook")
			out, err := ExecHook(config.Hooks.PostRender)
			if err != nil {
				log.WithError(err).Error("Failed to run post hook")
				return
			}
			if *loghook {
				log.Info("Output from post hook")
				fmt.Println(out)
			}
		}
	}
}
