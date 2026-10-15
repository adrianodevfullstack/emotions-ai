package main

import (
	"log"

	"github.com/adrianodevfullstack/emotions-ai/internal/emotions"
)

func main() {
	if err := emotions.Run(); err != nil {
		log.Fatal(err)
	}
}
