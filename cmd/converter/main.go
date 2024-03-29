package main

import (
	"fmt"
	"github.com/reedray/bank-service/api/pb/converter/gen_convert"
	"github.com/reedray/bank-service/config/converter"
	"github.com/reedray/bank-service/internal/converter/storage"
	"github.com/reedray/bank-service/internal/converter/transport/grpc_transport"
	"github.com/reedray/bank-service/internal/converter/usecase"
	"github.com/reedray/bank-service/internal/converter/usecase/webapi"
	"github.com/reedray/bank-service/pkg/logger"
	"google.golang.org/grpc"
	"net"
)

var (
	configPath = "./config/converter/config.yml"
)

func main() {
	cfg, err := converter.NewConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	log, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Info("Logger created")

	redisRepository, err := storage.NewRedis(cfg)
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Info("redis storage created")

	webAPI := webapi.New()
	log.Info("webAPI created")

	convertUseCase := usecase.New(redisRepository, webAPI)
	log.Info("convert UseCase created")

	handler := grpc_transport.NewConvertController(convertUseCase)

	listener, err := net.Listen("tcp", cfg.Grpc.Addr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	server := grpc.NewServer()
	gen_convert.RegisterConvertServiceServer(server, handler)
	log.Info("starting server")
	if err = server.Serve(listener); err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Info("shutting down  server")
}
