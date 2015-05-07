package main

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

func main() {
	var redisConn redis.Conn
	var err error
	redisConn, err = redis.Dial("tcp", ":6379")
	if err != nil {
		log.Println(err)
		return
	}
	r, err := redisConn.Do("SET", "200", "197")
	if err != nil {
		log.Println(err)
	}
	log.Printf("r is [%T] with value %v\n", r, r)
	redisConn.Close()
	log.Printf("redisConn: %v\n", redisConn)

	var redisConn2 redis.Conn
	var err2 error
	redisConn2, err2 = redis.Dial("tcp", ":6379")
	if err2 != nil {
		log.Println(err2)
	}

	r2, err2 := redis.String(redisConn2.Do("GET", "100"))
	if err2 != nil {
		log.Println(err2)
	}
	log.Printf("r2 is type [%T] and value %v\n", r2, r2)
	log.Printf("")
}
