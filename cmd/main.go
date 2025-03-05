package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
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

type Data struct {
	ctx context.Context
	rdb *redis.Client
}

func NewData() *Data {
	_ = godotenv.Load("../.env")

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		DB:   0,
	})

	ctx := context.Background()

	return &Data{
		ctx: ctx,
		rdb: rdb,
	}
}

func main() {
	data := NewData()

	jsonPath := "../dict.json"
	dict, err := data.jsonDecode(jsonPath)
	if err != nil {
		log.Fatal(err)
	}

	data.redisConn(dict)

}

func (ctx *Data) jsonDecode(jsonPath string) ([]Dict, error) {
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

func (ctx *Data) redisConn(dict []Dict) {

	for _, v := range dict {
		marshal, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}

		_ = ctx.rdb.LPush(ctx.ctx, os.Getenv("LIST_NAME"), marshal).Err()

	}

	length, err := ctx.rdb.LLen(ctx.ctx, os.Getenv("LIST_NAME")).Result()
	if err != nil {
		log.Fatal(err)
	}

	for {
		random := rand.Intn(int(length))

		element, err := ctx.rdb.LIndex(ctx.ctx, os.Getenv("LIST_NAME"), int64(random)).Result()
		if err != nil {
			log.Fatal(err)
		}
		_ = ctx.rdb.LRem(ctx.ctx, os.Getenv("LIST_NAME"), 1, element)

		time.Sleep(time.Second * 5)
	}
}
