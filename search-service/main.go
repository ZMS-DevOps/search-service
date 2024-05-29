package main

import (
	"github.com/ZMS-DevOps/search-service/startup"
	cfg "github.com/ZMS-DevOps/search-service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
