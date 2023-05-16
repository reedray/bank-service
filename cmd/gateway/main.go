package main

import (
	"context"
	"fmt"
	"github.com/reedray/bank-service/api/pb/converter"
	converter_client "github.com/reedray/bank-service/internal/converter/transport/grpc_transport"
)

func main() {

	client := converter_client.NewClient("localhost:8080")
	req := converter.Money{
		Amount:       "1",
		CurrencyCode: "EUR",
	}
	converted, err := client.Convert(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(converted.Amount, converted.CurrencyCode)
	client.Conn.Close()
}
