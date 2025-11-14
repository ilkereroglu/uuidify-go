package main

import (
	"context"
	"fmt"
	"log"

	uuidify "github.com/ilkereroglu/uuidify-go"
)

func main() {
	c := uuidify.NewClient()

	ulid, err := c.ULID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ULID:", ulid)
}
