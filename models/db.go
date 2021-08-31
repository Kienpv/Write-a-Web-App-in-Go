package models

import (
	"github.com/go-redis/redis"
)

var client *redis.Client


func Init()  {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost: 6379", 	// where the redis serve is. In this example, we're running it on the same machine 
									// as our web application on port 6379 which is the default port for Redis.
	})
}