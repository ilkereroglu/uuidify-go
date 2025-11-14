package main

import (
	"context"
	"fmt"
	"log"

	uuidify "github.com/ilkereroglu/uuidify-go"
)

func main() {
	c := uuidify.NewClient()

	uuids, err := c.UUIDBatch(context.Background(), "v7", 5)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UUID v7 batch:", uuids)
}
