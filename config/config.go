package config

import (
	"os"
	"strconv"
	"time"
)

var (
	Address             = get("MM_ADDR", "127.0.0.1:25565")
	ServerName          = get("MM_SRVNM", "marmalade")
	ServerMOTD          = get("MM_SRVMOTD", "placeholder MOTD, ask the server owner to set one!")
	BufferFlushInterval = time.Second / time.Duration(mustAtoi(get("MM_TICKRATE", "20"))) // value to be passed into the AFCBW constructor
)

func get(key, fallback string) string {
	val, found := os.LookupEnv(key)
	if found {
		return val
	} else {
		return fallback
	}
}

func mustAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}
