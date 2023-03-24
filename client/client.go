package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const CLIENT_TIMEOUT = 500 * time.Millisecond

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), CLIENT_TIMEOUT)
	defer cancel()

	exchangeRate, err := getDolarRate(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("exchangeRate:", exchangeRate)

	err = saveRateToText(exchangeRate)
	if err != nil {
		panic(err)
	}
}


func getDolarRate(ctx context.Context) (float32, error) {
	r, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return 0.0, err
	}
	
	rs, err := http.DefaultClient.Do(r)
	if err != nil {
		return 0.0, err 
	}
	defer rs.Body.Close()
	
	dolarRate, err := io.ReadAll(rs.Body)
	if err != nil {
		return 0.0, err
	}

	value, _ := strconv.ParseFloat(string(dolarRate[:]), 32)
	exchangeRate := float32(value)
	return exchangeRate, nil
}


func saveRateToText(exchangeRate float32) error {
	rate := strconv.FormatFloat(float64(exchangeRate), 'f', 5, 32) 
	currentRate := "Dólar: " + rate 

	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}

	io.Copy(file, strings.NewReader(currentRate))
	return nil
}
