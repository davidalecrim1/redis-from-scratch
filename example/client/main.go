package main

import (
	"context"
	"log"
	"strconv"
)

func main() {
	ctx := context.Background()
	client := NewClient(":6379")

	for i := range 2 {
		err := client.Set(ctx, "key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
		if err != nil {
			panic(err)
		}

		// time.Sleep(1 * time.Second)
	}

	for i := range 4 {
		resp, err := client.Ping(ctx)
		if err != nil {
			panic(err)
		}

		// time.Sleep(1 * time.Second)
		log.Printf("ping %d response: %v\n", i, resp)
	}
}
