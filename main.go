package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// DB setup
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	// DB init
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// AWS setup & init
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("Couldn't load default AWS configuration, error: ", err)
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	server := NewAPIServer(":3000", store, s3Client)
	server.Run()
}
