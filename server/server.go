package main

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	HTTP_TIMEOUT = 200 * time.Millisecond
	DB_TIMEOUT   = 100 * time.Millisecond
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
	log.Println("Ouvindo requisições em :8080")
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

	err = saveDolarRate(ctx, cotacao.USDBRL.BID)
	if err != nil {
		log.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cotacao.USDBRL.BID))

	select {
	case <-time.After(HTTP_TIMEOUT):
		log.Println("Request executado com sucesso")

	case <-ctx.Done():
		log.Println("Request cancelado pelo cliente")
	}
}

func getDolarRate(ctx context.Context) (*USDBRL, error) {
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
	if err != nil {
		return nil, err
	}

	var cotacao USDBRL
	_ = json.Unmarshal(body, &cotacao)
	return &cotacao, nil
}

func saveDolarRate(ctx context.Context, rate string) error {
	ctx, cancel := context.WithTimeout(ctx, DB_TIMEOUT)
	defer cancel()

	// log.Println("Dólar:", rate)
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}

	query := `
	DROP TABLE IF EXISTS cotacao;
	CREATE TABLE cotacao(dolar TEXT);
	`
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	query = "INSERT INTO cotacao (dolar) VALUES (?)"
	prep, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = prep.ExecContext(ctx, rate)
	if err != nil {
		return err
	}

	return nil
}
