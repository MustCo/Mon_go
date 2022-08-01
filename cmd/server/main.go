package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/Snakder/Mon_go/internal/server"
	"github.com/Snakder/Mon_go/internal/utils"
	"github.com/caarlos0/env/v6"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.json", "path to config json file")
}

func main() {
	flag.Parse()
	config := utils.NewConfig()
	ctx := context.TODO()
	if err := env.Parse(config); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", config.Address)
	server := server.New(config)
	log.Fatal(server.Start(ctx))

}
