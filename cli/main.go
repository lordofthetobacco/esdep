package main

import (
	"esdep/internal/config"
	"esdep/internal/deployment"
	"flag"
	"log"
	"time"
)

func main() {
	flag.Usage = func() {
		log.Printf("Usage of %s:\n", "esdep")
		flag.PrintDefaults()
		log.Println(`
Deploys and watches configured repositories for updates.

Options:
  -config <file>   Path to config file (default: config.yaml)
  -h, --help       Show this help message`)
	}
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	help := flag.Bool("help", false, "show help")
	h := flag.Bool("h", false, "show help")
	flag.Parse()

	if *help || *h {
		flag.Usage()
		return
	}

	cfg, err := config.LoadConfigFile(*cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	for _, entry := range cfg.DeployEntries {
		err := deployment.Deploy(entry)
		if err != nil {
			log.Fatalf("Failed to deploy: %v", err)
		}
	}

	deployment.RunUpdateChecker(cfg.DeployEntries, 2*time.Minute+30*time.Second)
}
