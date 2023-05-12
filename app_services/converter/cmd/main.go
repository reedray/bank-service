package main

import (
	"fmt"
	"github.com/reedray/bank-service/app_services/converter/config"
	"github.com/reedray/bank-service/app_services/converter/internal/grpcHandler"
	"github.com/reedray/bank-service/app_services/converter/internal/usecase"
	"github.com/reedray/bank-service/app_services/converter/internal/usecase/repository"
	"github.com/reedray/bank-service/app_services/converter/internal/usecase/webapi"
	"github.com/reedray/bank-service/app_services/converter/pkg/api/api_pb"
	"github.com/reedray/bank-service/pkg/logger"
	"google.golang.org/grpc"
	"net"
)

var (
	configPath = "./config/config.yml"
)

func main() {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Println(err)
		//TODO: replace by graceful shutdown
		return
	}
	log, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		fmt.Println(err)
		//TODO: replace by graceful shutdown
		return
	}
	log.Info("Logger created")

	redisRepository, err := repository.NewRedis(cfg)
	if err != nil {
		log.Error(err.Error())
		//TODO: replace by graceful shutdown
		return
	}
	log.Info("redis repository created")

	webAPI := webapi.New()
	log.Info("webAPI created")

	convertUseCase := usecase.New(redisRepository, webAPI)
	log.Info("convert UseCase created")

	handler := grpcHandler.New(convertUseCase)

	listener, err := net.Listen("tcp", cfg.Grpc.Addr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	server := grpc.NewServer()
	api_pb.RegisterConvertServiceServer(server, handler)
	log.Info("starting server")
	if err = server.Serve(listener); err != nil {
		log.Fatal(err.Error())
		return
	}
}
