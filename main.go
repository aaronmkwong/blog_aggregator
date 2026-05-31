package main

import (
	"fmt"
	"log"

	"github.com/aaronmkwong/blog_aggregator/internal/config"
)

func main() {
	// 1. Read config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Set current user
	err = cfg.SetUser("aaron")
	if err != nil {
		log.Fatal(err)
	}

	// 3. Read config again
	updatedCfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// print contents
	fmt.Println(updatedCfg)
}