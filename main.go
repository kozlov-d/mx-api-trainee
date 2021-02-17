package main

import (
	"log"

	"github.com/kozlov-d/mx-api-trainee/common"
	"github.com/kozlov-d/mx-api-trainee/server"
)

func main() {
	conf := common.FromEnv()
	s := server.NewServer(conf)
	log.Fatal(s.Start())
}
