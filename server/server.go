package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type USDBRL struct {
	USDBRL cotacao
}
type cotacao struct {
	Code      string  `json:"code"`
	Codein    string  `json:"code-in"`
	Name      string  `json:"nae"`
	High      string `json:"high-value"`
	Low       string `json:"low-value"`
	VarBid    string `json:"var-bid"`
	PctChange string `json:"pct-change"`
	BID       string `json:"bid"`
	ASK       string `json:"ask"`
	timestamp string  `json:"-"`
	CreatedAt string  `json:"created-at"`
}

type BID struct {
	BID 	float32 `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", getDolarRateHandler)
	http.ListenAndServe(":8080", nil)
}

func getDolarRateHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		http.Error(w, "route not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	defer log.Println("Request finalizado")

	cotacao, err := getDolarRate(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	time.Sleep(500 * time.Millisecond)
	
	w.WriteHeader(http.StatusOK)
	var bid BID
	value, _ := strconv.ParseFloat(cotacao.USDBRL.BID, 32)
	bid.BID = float32(value)
	log.Println("BID:", bid)
	json.NewEncoder(w).Encode(bid)

	// select {
	// case <-time.After(10 * time.Second):
	// 	w.WriteHeader(http.StatusRequestTimeout)
	// 	w.Write([]byte("Request com tempo esgotado\n"))

	// case <-ctx.Done():
	// 	log.Println("Request cancelado pelo cliente")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	// http.Error(w, "Request cancelado pelo cliente", http.StatusRequestTimeout)
	// }
}

func getDolarRate(ctx context.Context) (cotacao *USDBRL, err error) {
	r, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Println("Parece que não achou", err.Error())
		return
	}

	rs, err := http.DefaultClient.Do(r)
	if err != nil {
		return
	}
	defer rs.Body.Close()

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &cotacao)
	return
}
