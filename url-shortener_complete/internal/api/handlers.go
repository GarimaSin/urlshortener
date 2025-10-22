package api

import (
    "context"
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "golang.org/x/time/rate"
    "github.com/you/url-shortener/internal/id"
    "github.com/you/url-shortener/internal/store"
    "github.com/you/url-shortener/internal/util"
)

type shortenReq struct {
    URL string `json:"url"`
}

type shortenResp struct {
    ShortCode string `json:"short_code"`
    ShortURL  string `json:"short_url"`
}

type Handler struct {
    s        store.Store
    idgen    *id.Generator
    baseURL  string
    limiter  *rate.Limiter
}

func NewHandler(s store.Store, idgen *id.Generator, baseURL string) *Handler {
    return &Handler{s: s, idgen: idgen, baseURL: baseURL, limiter: rate.NewLimiter(rate.Limit(5000), 10000)}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
    if !h.limiter.Allow() {
        http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
        return
    }
    b, _ := io.ReadAll(r.Body)
    defer r.Body.Close()
    var req shortenReq
    if err := json.Unmarshal(b, &req); err != nil || req.URL == "" {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    // basic validation
    if len(req.URL) > 2048 {
        http.Error(w, "url too long", http.StatusBadRequest)
        return
    }

    idv, err := h.idgen.Next()
    if err != nil {
        http.Error(w, "internal", http.StatusInternalServerError)
        return
    }
    code := util.EncodeBase62(idv)
    if err := h.s.Put(code, req.URL); err != nil {
        http.Error(w, "store error", http.StatusInternalServerError)
        return
    }
    resp := shortenResp{ShortCode: code, ShortURL: h.baseURL + "/" + code}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
    short := chi.URLParam(r, "shortCode")
    if short == "" {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }
    long, err := h.s.Get(short)
    if err != nil {
        if errors.Is(err, store.ErrNotFound) {
            http.Error(w, "not found", http.StatusNotFound)
            return
        }
        http.Error(w, "internal", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, long, http.StatusFound)
}
