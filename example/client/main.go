package main

import (
	"context"
	"log"
	"strconv"
)

func main() {
	ctx := context.Background()
	client := NewClient(":6379")

	for i := range 10 {
		key, value := "key"+strconv.Itoa(i), "value"+strconv.Itoa(i)
		err := client.Set(ctx, key, value)
		if err != nil {
			panic(err)
		}
		log.Printf("send the key '%v' and value '%v'\n", key, value)
	}

	for i := range 10 {
		resp, err := client.Ping(ctx)
		if err != nil {
			panic(err)
		}
		log.Printf("ping '%d' response: '%v'\n", i, resp)
	}
}
