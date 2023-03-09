package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	_, err := getDolarRate(ctx)
	if err != nil {
		panic(err)
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

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return
	}
	log.Println("Cotação:", body)

	var cotacao interface{}
	err = json.Unmarshal(body, &cotacao)
	return 
}

func saveRateToText(exchangeRate float32) (err error) {
	return err
}