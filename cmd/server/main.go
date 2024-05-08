package main

import "github.com/amho/go-redis/redis"

const port = 6379

func main() {
	rds := redis.New()
	rds.ListenAndServe(port)
}
