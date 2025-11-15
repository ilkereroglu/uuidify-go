package main

import (
	"context"
	"fmt"
	"log"

	uuidify "github.com/ilkereroglu/uuidify-go"
)

func main() {
	c, err := uuidify.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}

	ulid, err := c.ULID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ULID:", ulid)
}
