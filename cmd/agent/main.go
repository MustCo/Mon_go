package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Snakder/Mon_go/internal/agent"
	"github.com/Snakder/Mon_go/internal/utils"
	"github.com/caarlos0/env/v6"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiagent.json", "path to config file")
}

func main() {
	flag.Parse()
	config := utils.NewConfig()

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()
	json.Unmarshal(data, config)
	if err := env.Parse(config); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", config.Address)
	fmt.Printf("%+v\n", config.PollInterval)
	fmt.Printf("%+v\n", config.ReportInterval)
	agent := agent.New(config)
	err = agent.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
