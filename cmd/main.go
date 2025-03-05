package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

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
	_ = godotenv.Load("../.env")

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		DB:   0,
	})

	jsonPath := "../dict.json"
	dict, _ := jsonDecode(jsonPath)

	redisConn(rdb, dict)

}

func jsonDecode(jsonPath string) ([]Dict, error) {
	file, err := os.Open(jsonPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var dict []Dict
	// decode
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&dict); err != nil {
		log.Fatal(err)
	}
	return dict, nil
}

func redisConn(rdb *redis.Client, dict []Dict) {
	ctx := context.Background()

	for _, v := range dict {
		marshal, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}
		key := fmt.Sprintf("word:%s", v.Word)

		if err := rdb.SetEx(ctx, key, marshal, 15*time.Second).Err(); err != nil {
			log.Fatal(err)
		}
	}
}
