package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

// Embed HTML file

//go:embed static/index.html
var indexHtml string

// ---------- Redis ----------
var rdb *redis.Client
var ctx = context.Background()

func initRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379" // default for local dev
	}
	rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	// Simple ping to verify connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Unable to connect to Redis at %s: %v", addr, err)
	}
}

// ---------- Handlers ----------
func getCount(w http.ResponseWriter, r *http.Request) {
	val, err := rdb.Get(ctx, "clicks").Result()
	if err == redis.Nil {
		val = "0"
	} else if err != nil {
		http.Error(w, "Redis error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"count": val})
}

func incCount(w http.ResponseWriter, r *http.Request) {
	newVal, err := rdb.Incr(ctx, "clicks").Result()
	if err != nil {
		http.Error(w, "Redis error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"count": newVal})
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Host == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexHtml)
		return
	}
}

// ---------- Main ----------
func main() {
	initRedis()

	// API routes
	http.HandleFunc("/api/count", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getCount(w, r)
			return
		}
		if r.Method == http.MethodPost {
			incCount(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", defaultHandler)

	fmt.Printf("🚀 Server listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
