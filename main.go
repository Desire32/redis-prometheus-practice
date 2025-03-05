package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Dict struct {
	Word          string `json:"word"`
	Transcription string `json:"transcription"`
	Translation   string `json:"translation"`
	Description   string `json:"description"`
}

func main() {
	// env
	_ = godotenv.Load(".env")

	jsonPath := os.Getenv("JSON_PATH")

	dict, _ := jsonDecode(jsonPath)

	redisConn(dict)
}

func jsonDecode(jsonPath string) ([]Dict, error) {
	file, err := os.Open(jsonPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var dict []Dict
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&dict); err != nil {
		log.Fatal(err)
	}
	return dict, nil
}

func redisConn(dict []Dict) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		DB:   0,
	})

	for _, v := range dict {
		marshal, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}
		if err := rdb.RPush(ctx, os.Getenv("LIST_NAME"), marshal).Err(); err != nil {
			log.Fatal(err)
		}
	}
}
