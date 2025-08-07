package main

import (
	"log"
	"log/slog"
	"net/http"

	"shortener.ru/internal/handler"
)

func init() {
	slog.LevelDebug
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", handler.GetShortUrl{})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}

}
