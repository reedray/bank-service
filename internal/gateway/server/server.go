package server

import (
	"context"
	"encoding/json"
	"github.com/reedray/bank-service/api/pb/transact/gen_transact"
	"github.com/reedray/bank-service/config/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Server struct {
	http.Server
	Cfg *gateway.Config
}

func New(cfg *gateway.Config) *Server {
	return &Server{Cfg: cfg}
}

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("cmd/gateway/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	conn, err := grpc.Dial(s.Cfg.Transact, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := gen_transact.NewTransferServiceClient(conn)
	response, err := client.Login(context.Background(), &gen_transact.LoginRequest{
		Username: username,
		Password: password,
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("failed to log in"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "Authorization",
		Value:   response.Token,
		Expires: time.Now().Add(time.Hour),
	})

}
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	conn, err := grpc.Dial(s.Cfg.Transact, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := gen_transact.NewTransferServiceClient(conn)
	response, err := client.Register(context.Background(), &gen_transact.RegisterRequest{
		Username: username,
		Password: password,
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to Register"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "Authorization",
		Value:   response.Token,
		Expires: time.Now().Add(time.Hour),
	})

}

func (s *Server) Balance(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("Authorization")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn, err := grpc.Dial(s.Cfg.Transact, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := gen_transact.NewTransferServiceClient(conn)
	response, err := client.Balance(context.Background(), &gen_transact.BalanceRequest{Token: cookie.Value})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response.ErrorMsg.Err != "" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	bytes, err := json.Marshal(response.Total)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (s *Server) Deposit(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	currency := r.FormValue("currency")
	amount := r.FormValue("amount")

	conn, err := grpc.Dial(s.Cfg.Transact, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := gen_transact.NewTransferServiceClient(conn)
	response, err := client.Deposit(context.Background(), &gen_transact.DepositRequest{
		Token: cookie.Value,
		Total: &gen_transact.Money{
			Amount:       currency,
			CurrencyCode: amount,
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response.Err != "" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successful operation "))

}

func (s *Server) Withdraw(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	currency := r.FormValue("currency")
	amount := r.FormValue("amount")

	conn, err := grpc.Dial(s.Cfg.Transact, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := gen_transact.NewTransferServiceClient(conn)
	response, err := client.Withdraw(context.Background(), &gen_transact.WithdrawRequest{
		Token: cookie.Value,
		Total: &gen_transact.Money{
			Amount:       currency,
			CurrencyCode: amount,
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response.Err != "" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successful operation "))

}
func (s *Server) Transfer(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	currency := r.FormValue("currency")
	amount := r.FormValue("amount")
	idTo := r.FormValue("idTo")

	conn, err := grpc.Dial(s.Cfg.Transact, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := gen_transact.NewTransferServiceClient(conn)
	response, err := client.Transfer(context.Background(), &gen_transact.TransferRequest{
		IdTo:  idTo,
		Token: cookie.Value,
		Total: &gen_transact.Money{
			Amount:       currency,
			CurrencyCode: amount,
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response.Err != "" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successful operation "))

}

func (s *Server) Credit(w http.ResponseWriter, r *http.Request) {

}
