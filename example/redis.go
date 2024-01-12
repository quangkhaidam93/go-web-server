package example

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type author struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func PlayWithRedis() {
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})

	// 0 -> unexpired key-value
	if err := redisClient.Set("name", "Deft", 0).Err(); err != nil {
		fmt.Println(err)
	}

	val, err := redisClient.Get("name").Result()

	if err != nil {
		fmt.Println(err)
	}

	jsonVal, err := json.Marshal(author{Name: "Deft", Age: 28})

	if err != nil {
		fmt.Println(err)
	}

	if err := redisClient.Set("id1234", jsonVal, 0).Err(); err != nil {
		fmt.Println(err)
	}

	val2, err := redisClient.Get("id1234").Result()

	if err != nil {
		fmt.Println(err)
	}

	var author author

	if err := json.Unmarshal([]byte(val2), &author); err != nil {
		fmt.Println(err)
	}

	fmt.Println(val)
	fmt.Printf("Author: %+v\n", author)
}
