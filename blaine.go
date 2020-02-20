package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
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
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
	visit := new(Visit)
	visit.Time = time.Now()
	visit.IpAddress = GetIP(req)
	visit.UserAgent = req.Header.Get("User-Agent")

	visitJSON, err := json.Marshal(visit)
	if err != nil {
		fmt.Println(err)
	}

	redisClient.RPush(REDIS_KEY, visitJSON)

	http.Redirect(w, req, os.Getenv("TARGET_URL"), 301)
}

func main() {
	fmt.Println("Starting server...")
	http.HandleFunc("/", root)

	http.ListenAndServe(":8090", nil)
}
