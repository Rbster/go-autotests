package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"shortener.ru/internal"
)

type GetShortUrl struct {
	service internal.ShortenerService
}

type GetFullUrl struct {
	service internal.ShortenerService
}

func (h *GetShortUrl) ServeHTTP(rs http.ResponseWriter, rq *http.Request) {
	if rq.Method != http.MethodPost {
		slog.Error("method not allowed", "error_code", http.StatusMethodNotAllowed)
		rs.Header().Add("Allow", http.MethodPost)
		http.Error(rs, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !strings.Contains(rq.Header.Get("Content-Type"), "text/plain") {
		slog.Error("unsupported media type", "error_code", http.StatusUnsupportedMediaType)
		http.Error(rs, "unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	var buff []byte

	if n, err := rq.Body.Read(buff); err != nil {
		slog.Error(err.Error(), "error_code", http.StatusInternalServerError)
		http.Error(rs, "server error occured", http.StatusInternalServerError)
		return
	} else if n <= 0 {
		defer rq.Body.Close()
		slog.Error("zero string received", "error_code", http.StatusBadRequest)
		http.Error(rs, "bad request", http.StatusBadRequest)
		return
	}
	defer rq.Body.Close()

	shortUrl, err := h.service.ToShort(string(buff), rq.URL.Host)
	if err != nil {
		slog.Error(err.Error(), "error_code", http.StatusInternalServerError)
		http.Error(rs, "server error occured", http.StatusInternalServerError)
		return
	}

	slog.Debug("succesfully delivered", "rs_body", shortUrl)

	rs.WriteHeader(http.StatusCreated)
	fmt.Fprint(rs, shortUrl)
}

func (h *GetFullUrl) ServeHTTP(rs http.ResponseWriter, rq *http.Request) {
	if rq.Method != http.MethodGet {
		slog.Error("method not allowed", "error_code", http.StatusMethodNotAllowed)
		rs.Header().Add("Allow", http.MethodGet)
		http.Error(rs, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	paths := strings.Split(rq.URL.Path, "/")
	if len(paths) != 1 {
		slog.Error("path has more or less then 1 argument", "error_code", http.StatusNotFound)
		http.Error(rs, "not found", http.StatusNotFound)
		return
	}
	if !strings.Contains(rq.Header.Get("Content-Type"), "text/plain") {
		slog.Error("unsupported media type", "error_code", http.StatusUnsupportedMediaType)
		http.Error(rs, "unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	shortUrl := paths[0]

	fullUrl, err := h.service.ToFull(shortUrl)
	if err != nil {
		slog.Error(err.Error(), "error_code", http.StatusInternalServerError)
		http.Error(rs, "server error occured", http.StatusInternalServerError)
		return
	}

	slog.Debug("succesfully delivered", "rs_body", shortUrl)
	rs.Header().Add("Location", fullUrl)
	rs.WriteHeader(http.StatusTemporaryRedirect)
	fmt.Fprint(rs, shortUrl)
}
