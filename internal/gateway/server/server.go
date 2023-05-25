package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/reedray/bank-service/api/pb/converter/gen_convert"
	"github.com/reedray/bank-service/api/pb/transact/gen_transact"
	"github.com/reedray/bank-service/config/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	http.Server
	Cfg *gateway.Config
}

var secret_hard_code = "some_secret"

func New(cfg *gateway.Config) *Server {
	return &Server{Cfg: cfg}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("cmd/gateway/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) Auth(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("cmd/gateway/auth.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		log.Println("Failed to get token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	valid, err := validateToken(cookie.Value, secret_hard_code)
	if err != nil || !valid {
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username, password)

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
		log.Println("Failed to get token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	currency := r.FormValue("currency")
	amount := r.FormValue("amount")
	currency = strings.ToUpper(currency)
	fmt.Println(currency, amount)

	conn, err := grpc.Dial(s.Cfg.Transact, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := gen_transact.NewTransferServiceClient(conn)
	_, err = client.Deposit(context.Background(), &gen_transact.DepositRequest{
		Token: cookie.Value,
		Total: &gen_transact.Money{
			Amount:       amount,
			CurrencyCode: currency,
		},
	})
	if err != nil {
		log.Println("gRPC deposit failed with error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("gRPC deposit finished successful")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successful operation "))

}

func (s *Server) Withdraw(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		log.Println("Failed to get token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	currency := r.FormValue("currency")
	amount := r.FormValue("amount")
	currency = strings.ToUpper(currency)
	log.Println(currency, amount)
	conn, err := grpc.Dial(s.Cfg.Transact, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := gen_transact.NewTransferServiceClient(conn)
	_, err = client.Withdraw(context.Background(), &gen_transact.WithdrawRequest{
		Token: cookie.Value,
		Total: &gen_transact.Money{
			Amount:       amount,
			CurrencyCode: currency,
		},
	})
	if err != nil {
		log.Println("gRPC withdraw failed with error", err.Error())
		if err.Error() == "not enough funds to withdraw" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("gRPC withdraw finished successful")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successful operation "))

}
func (s *Server) Transfer(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		log.Println("Failed to get token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	currency := r.FormValue("currency")
	amount := r.FormValue("amount")
	idTo := r.FormValue("recipientId")
	currency = strings.ToUpper(currency)

	conn, err := grpc.Dial(s.Cfg.Transact, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := gen_transact.NewTransferServiceClient(conn)
	_, err = client.Transfer(context.Background(), &gen_transact.TransferRequest{
		IdTo:  idTo,
		Token: cookie.Value,
		Total: &gen_transact.Money{
			Amount:       amount,
			CurrencyCode: currency,
		},
	})

	if err != nil {
		log.Println("gRPC transfer failed with error", err.Error())
		if err.Error() == "not enough funds to transfer" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("gRPC transfer finished successful")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successful operation "))

}

func (s *Server) Convert(w http.ResponseWriter, r *http.Request) {

	currency := r.FormValue("currency")
	amount := r.FormValue("amount")
	currency = strings.ToUpper(currency)
	fmt.Println(currency, amount)

	conn, err := grpc.Dial(s.Cfg.Converter, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := gen_convert.NewConvertServiceClient(conn)
	convert, err := client.Convert(context.Background(), &gen_convert.Money{
		Amount:       amount,
		CurrencyCode: currency,
	})

	if err != nil {
		log.Println("gRPC convert failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("gRPC convert finished successful")
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(convert)
	w.Write(bytes)

}

func validateToken(tokenString, secret string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is HMAC with the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for validation
		return []byte(secret), nil
	})

	// Check for parsing errors
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return false, fmt.Errorf("token is malformed")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return false, fmt.Errorf("token is expired or not active yet")
			} else {
				return false, fmt.Errorf("token validation failed: %v", err)
			}
		}
		return false, fmt.Errorf("token validation failed: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return false, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("invalid claims format")
	}

	status, ok := claims["status"].(string)
	if !ok {
		return false, fmt.Errorf("status claim is missing or has invalid typle")
	}
	if status == "banned" {
		return false, fmt.Errorf("customer is banned")
	}
	return true, nil
}
