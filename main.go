package main

import (
	"log"

	_ "net/http/pprof"

	_ "github.com/kozlov-d/mx-api-trainee/docs"

	"github.com/kozlov-d/mx-api-trainee/common"
	"github.com/kozlov-d/mx-api-trainee/server"
)

// @title MerchantX Trainee API
// @license.name
// @host localhost
// @BasePath /
func main() {
	conf := common.FromEnv()
	s := server.NewServer(conf)
	log.Fatal(s.Start())
}
