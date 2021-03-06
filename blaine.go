package main

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"log"
	"net/http"
	"os"
	"time"
)

const REDIS_KEY = "visits"

type Visit struct {
	Time      time.Time `json:"time"`
	IpAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
}

// From https://golangcode.com/get-the-request-ip-addr/
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func root(w http.ResponseWriter, req *http.Request) {
	log.Println("Processing request")
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Println(err)
	}

	redisClient := redis.NewClient(opt)

	visit := new(Visit)
	visit.Time = time.Now()
	visit.IpAddress = GetIP(req)
	visit.UserAgent = req.Header.Get("User-Agent")

	visitJSON, err := json.Marshal(visit)
	if err != nil {
		log.Println(err)
	}

	err = redisClient.RPush(REDIS_KEY, visitJSON).Err()
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, req, os.Getenv("TARGET_URL"), 301)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.Println("Starting server...")
	http.HandleFunc("/", root)

	http.ListenAndServe(":"+port, nil)
}
