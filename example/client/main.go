package main

import (
	"context"
	"log"
	"strconv"
)

func main() {
	ctx := context.Background()
	client := NewClient(":6379")
	defer client.Close()

	for i := range 10 {
		key, value := "key"+strconv.Itoa(i), "value"+strconv.Itoa(i)
		err := client.Set(ctx, key, value)
		if err != nil {
			panic(err)
		}

		log.Printf("send the key '%v' and value '%v'\n", key, value)

		resp, err := client.Get(ctx, key)
		if err != nil {
			panic(err)
		}

		log.Printf("received the value '%v' with getting the key '%v'", string(resp), key)
	}

	for i := range 10 {
		resp, err := client.Ping(ctx)
		if err != nil {
			panic(err)
		}
		log.Printf("ping '%d' response: '%v'\n", i, resp)
	}

	log.Println("Trying to get a key that doesn't exist")

	resp, err := client.Get(ctx, "key_invalid")
	if err != nil {
		panic(err)
	}

	log.Printf("received the value '%v' with getting the key '%v'", string(resp), "key_invalid")
}
