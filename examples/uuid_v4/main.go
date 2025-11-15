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

	uuid, err := c.UUIDv4(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UUID v4:", uuid)
}
