package main

import (
	"context"
	"fmt"
	"github.com/reedray/bank-service/api/pb/transact/gen_transact"
	"github.com/reedray/bank-service/config/transact"
	"github.com/reedray/bank-service/internal/transact/storage"
	"github.com/reedray/bank-service/internal/transact/transport"
	"github.com/reedray/bank-service/internal/transact/transport/grpc_transport"
	"github.com/reedray/bank-service/internal/transact/usecase"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	configPath = "./config/transact/config.yml"
)

func main() {
	config, err := transact.NewConfig(configPath)
	if err != nil {
		return
	}

	background := context.Background()
	db, err := storage.NewConnectionPoolDB(background, config.ConnString)
	if err != nil {
		fmt.Println(err)
		return
	}
	ari := storage.NewAuthRepositotyImpl(db)
	cri := storage.NewCustomerRepository(db)
	tri := storage.NewTransacioRepository(db)

	a := usecase.NewAuthUseCase(ari, cri)
	b := usecase.NewBalance(cri)
	d := usecase.NewDeposit(cri)
	tx := usecase.NewTransaction(tri, db)
	tr := usecase.NewTransfer(cri)
	w := usecase.NewWithdraw(cri)

	handler := transport.NewTransferHandler(config, a, d, w, tr, b, tx)

	serverImpl := grpc_transport.NewServerHandler(handler)
	server := grpc.NewServer()
	listener, err := net.Listen("tcp", config.GRPC.Addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	gen_transact.RegisterTransferServiceServer(server, serverImpl)

	log.Println("Starting gRPC server...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
