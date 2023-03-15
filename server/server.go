package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	HTTP_TIMEOUT = 200 * time.Millisecond
)

type USDBRL struct {
	USDBRL cotacao
}
type cotacao struct {
	Code      string `json:"code"`
	Codein    string `json:"codeIn"`
	Name      string `json:"name"`
	High      string `json:"highValue"`
	Low       string `json:"lowValue"`
	VarBid    string `json:"varBid"`
	PctChange string `json:"pctChange"`
	BID       string `json:"bid"`
	ASK       string `json:"ask"`
	Timestamp string `json:",omitempty"`
	CreatedAt string `json:"createdAt"`
}

func main() {
	http.HandleFunc("/cotacao", getDolarRateHandler)
	http.ListenAndServe(":8080", nil)
}

func getDolarRateHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		http.Error(w, "route not found", http.StatusNotFound)
	}

	ctx := r.Context()
	cotacao, err := getDolarRate(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Cotação:", cotacao)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cotacao.USDBRL.BID))

	select {
	case <-time.After(HTTP_TIMEOUT):
		log.Println("Request com tempo esgotado")
		// w.WriteHeader(http.StatusRequestTimeout)
		// w.Write([]byte("Request com tempo esgotado\n"))

	case <-ctx.Done():
		log.Println("Request cancelado pelo cliente")
		// http.Error(w, "Request cancelado pelo cliente", http.StatusBadRequest)
	}
}


func getDolarRate(ctx context.Context) (*USDBRL, error) {
	// "jump of the cat"
	ctx, cancel := context.WithTimeout(ctx, HTTP_TIMEOUT)
	defer cancel()

	r, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	rs, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	body, err := ioutil.ReadAll(rs.Body)
	log.Println("Body:", body)
	if err != nil {
		return nil, err
	}

	var cotacao USDBRL
	_ = json.Unmarshal(body, &cotacao)
	return &cotacao, nil
}
