package main

import (
	"log"
	config "url-shortener/internal"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	cfg, err := config.MustLoad()
	if err != nil {
		return err
	}

	println(cfg) //TODO: DELETE ME

	//TODO: logger

	//TODO: storage

	//TODO: router

	//TODO: start server

	return nil
}
