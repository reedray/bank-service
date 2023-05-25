package usecase

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/reedray/bank-service/internal/transact"
	"github.com/reedray/bank-service/internal/transact/entity"
	"log"
	"time"
)

type AuthUseCaseImpl struct {
	transact.AuthRepository
	transact.CustomerRepository
}

func NewAuth(a transact.AuthRepository, c transact.CustomerRepository) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		AuthRepository:     a,
		CustomerRepository: c,
	}
}

func (a *AuthUseCaseImpl) Login(ctx context.Context, username, password, secret string) (string, error) {
	log.Println("finding user by credentials", username, password)
	customer, err := a.AuthRepository.FindByCredentials(ctx, username, password)
	if customer.CreatedAt.IsZero() || err != nil {
		return "", fmt.Errorf("no such user")
	}
	token, err := generateToken(*customer, secret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *AuthUseCaseImpl) Register(ctx context.Context, username, password, secret string) (string, error) {
	customerByCred, err := a.FindByCredentials(ctx, username, password)
	if customerByCred != nil && (!customerByCred.CreatedAt.IsZero() || err != nil) {
		return "", fmt.Errorf("user already exists")
	}
	customer := entity.Customer{
		ID:         uuid.New(),
		Username:   username,
		Password:   password,
		Role:       "customer",
		BalanceRaw: []byte(`{"BYN":0,"USD":0,"EUR":0}`),
		Status:     "active",
		CreatedAt:  time.Now(),
		DeletedAat: time.Time{},
	}
	token, err := generateToken(customer, secret)
	if err != nil {
		return "", err
	}
	err = a.CustomerRepository.Save(ctx, &customer)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *AuthUseCaseImpl) ValidateToken(tokenString, secret string) (bool, error) {
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

func generateToken(customer entity.Customer, secret string) (string, error) {
	log.Println("Creating token")
	claims := jwt.MapClaims{
		"exp":    time.Now().Add(time.Hour).Unix(),
		"iat":    time.Now(),
		"id":     customer.ID.String(),
		"role":   customer.Role,
		"status": customer.Status,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedString, nil

}
