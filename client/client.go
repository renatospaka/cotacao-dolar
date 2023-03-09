package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	r, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)
	if err != nil {
		panic(err)
	}

	rs, err := http.DefaultClient.Do(r)
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, rs.Body)
}