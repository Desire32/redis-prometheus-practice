package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// json format
type Dict struct {
	Word          string `json:"word"`
	Transcription string `json:"transcription"`
	Translation   string `json:"translation"`
	Description   string `json:"description"`
}

func main() {

	// env
	_ = godotenv.Load(".env")

	var dict []Dict

	fmt.Print(jsonDecode(dict, os.Getenv("JSON_PATH")))
}

// json encoder
func jsonDecode(dict []Dict, json_s string) []Dict {
	file_dict, _ := os.Open(json_s)

	defer file_dict.Close()

	encrypt_dict := json.NewDecoder(file_dict)
	if err := encrypt_dict.Decode(&dict); err != nil {
		log.Fatal(err)
	}
	return dict
}

// func redisConn() {

// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})
// }
