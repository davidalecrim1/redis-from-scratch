package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	client := NewClient(":6379")

	err := client.Set(ctx, "key1", "value1")
	if err != nil {
		fmt.Printf("client set key1 value1 - err: %v\n", err)
	}

	err = client.Set(ctx, "key2", "value2")
	if err != nil {
		fmt.Printf("client set key2 value2 - err: %v\n", err)
	}

	for i := range 2 {
		resp, err := client.Ping(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Printf("ping %d response: %v\n", i, resp)
	}
}
