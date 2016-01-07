package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

const (
	redisDefaultDDb     = 3
	redisDefaultAddress = ":6379"
)

var rconn redis.Conn

func getList() []string {
	rep, err := redis.Strings(rconn.Do("SMEMBERS", "numbers"))
	if err == redis.ErrNil || len(rep) == 0 {
		return []string{}
	} else if err != nil {
		panic("Failed to SMEMBERS numbers.")
	}

	fmt.Printf("WHAT? %+v", rep)

	return rep
}

func addNumber(number string) error {
	nd, err := redis.Int64(rconn.Do("SADD", "numbers", number))
	if err != nil {
		panic("Failed to add cellphone number " + number + err.Error())
	}
	fmt.Printf("Added #%s %d.\n", number, nd)
	return nil
}

func startRedis() redis.Conn {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = redisDefaultAddress
	}

	dbIndex, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		dbIndex = redisDefaultDDb
	}

	client, err := redis.Dial("tcp", addr, redis.DialDatabase(dbIndex))
	if err != nil {
		panic("Failed to start redis.\n" + err.Error())
	}

	log.Printf("Redis started addr=%s dbIndex=%d.\n", addr, dbIndex)

	return client
}

func closeRedis(conn redis.Conn) {
	conn.Close()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	fmt.Printf("=======================\n\n\n\n===============\n==================\n===================PUTA UQE ME PARIU!!\n")
	log.Printf("ADDR: %s\n", os.Getenv("REDIS_ADDR"))
	rconn = startRedis()
	defer closeRedis(rconn)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			log.Printf("Received task %s scheduled at %s\n",
				r.Header.Get("X-Aws-Sqsd-Taskname"),
				r.Header.Get("X-Aws-Sqsd-Scheduled-At"))
		}

		if r.Method == "GET" {
			number := r.FormValue("num")
			if number != "" {
				addNumber(number)
				getList()
			}
		}
	})

	log.Printf("Listening on port %s\n\n", port)
	http.ListenAndServe(":"+port, nil)
}
