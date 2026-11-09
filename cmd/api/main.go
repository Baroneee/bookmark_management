package main

import (
	"github.com/Baroneee/EBVN_Assignment_1/internal/api"
	"github.com/Baroneee/EBVN_Assignment_1/pkg/redis"
)

// @title Bookmark Management API
// @version 1.0.0
// @description This is a simple bookmark management API.
// BasePath
// main starts the bookmark management API.
func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(cfg, redisClient)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
