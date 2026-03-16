package main

import (
	"esdep/internal/config"
	"esdep/internal/deployment"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "esdep")
		fmt.Fprintln(os.Stderr, "Deploys and watches configured repositories for updates.")
		fmt.Fprintln(os.Stderr, "")
		flag.PrintDefaults()
	}
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	interval := flag.Duration("interval", 2*time.Minute+30*time.Second, "how often to check for remote updates")
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
		log.Printf("Deploying %s...", entry.Path)
		err := deployment.Deploy(entry)
		if err != nil {
			log.Fatalf("Failed to deploy: %v", err)
		}
		log.Printf("Deployed %s", entry.Path)
	}

	log.Printf("Watching for updates every %s", *interval)
	deployment.RunUpdateChecker(cfg.DeployEntries, *interval)
}
