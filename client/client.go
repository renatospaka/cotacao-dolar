package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Millisecond)
	defer cancel()

	go func() {
		exchangeRate, err := getDolarRate(ctx)
		if err != nil {
			panic(err)
		}
		log.Println("exchangeRate:", exchangeRate)
	}()

	select {
	case <-time.After(300 * time.Millisecond):
		log.Println("Status: ", http.StatusRequestTimeout)
		log.Println("Request com tempo esgotado")

	case <-ctx.Done():
		log.Println("Status: ", http.StatusBadRequest)
		log.Println("Request cancelado pelo cliente")
	}
}


func getDolarRate(ctx context.Context) (exchangeRate float32, err error) {
	r, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)
	if err != nil {
		return 
	}
	log.Println("PASSOU: http.NewRequestWithContext")
	
	rs, err := http.DefaultClient.Do(r)
	if err != nil {
		return 
	}
	defer rs.Body.Close()
	log.Println("PASSOU: http.DefaultClient.Do")
	
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		return
	}
	log.Println("PASSOU: io.ReadAll")
	log.Println("Cotação (body):", body)

	var cotacao any
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		return
	}
	log.Println("Cotação:", cotacao)
	return 
}


// func saveRateToText(exchangeRate float32) (err error) {
// 	return err
// }
