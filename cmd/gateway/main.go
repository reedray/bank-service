package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/reedray/bank-service/config/gateway"
	"github.com/reedray/bank-service/internal/gateway/server"
	"log"
)

var (
	configPath = "./config/gateway/config.yml"
)

func main() {
	config, err := gateway.NewConfig(configPath)
	if err != nil {
		return
	}
	s := server.New(config)
	r := mux.NewRouter()
	r.HandleFunc("/", s.Home)
	r.HandleFunc("/login", s.Login)
	r.HandleFunc("/register", s.Register)
	r.HandleFunc("/deposit", s.Deposit)
	r.HandleFunc("/withdraw", s.Withdraw)
	r.HandleFunc("/balance", s.Balance)
	r.HandleFunc("/transfer", s.Transfer)
	s.Server.Addr = s.Cfg.Gateway
	s.Server.Handler = r
	fmt.Println("server started on", s.Cfg.Gateway)
	log.Fatal(s.ListenAndServe())
}
