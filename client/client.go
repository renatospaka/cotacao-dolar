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
		_, err := getDolarRate(ctx)
		if err != nil {
			panic(err)
		}
	}()

	select {
	case <-time.After(100 * time.Millisecond):
		log.Println("Status: ", http.StatusRequestTimeout)
		log.Println("Request com tempo esgotado")

	case <-ctx.Done():
		log.Println("Status: ", http.StatusBadRequest)
		log.Println("Request cancelado pelo cliente")
	}
}

func getDolarRate(ctx context.Context) (exchangeRate float32, err error) {
	exchangeRate = 0.0;
	err = nil

	r, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)
	if err != nil {
		return 0.0, err
	}

	rs, err := http.DefaultClient.Do(r)
	if err != nil {
		return 0.0, err
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		return
	}
	log.Println("Cotação:", body)

	var cotacao interface{}
	err = json.Unmarshal(body, &cotacao)
	return 
}

// func saveRateToText(exchangeRate float32) (err error) {
// 	return err
// }
