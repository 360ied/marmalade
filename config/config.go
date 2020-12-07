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
	WorldXSize          = mustAtoi(get("MM_WXSIZE", "256"))
	WorldYSize          = mustAtoi(get("MM_WYSIZE", "128"))
	WorldZSize          = mustAtoi(get("MM_WZSIZE", "256"))
	SpawnPointX         = mustAtoi(get("MM_WSPX", "1"))
	SpawnPointY         = mustAtoi(get("MM_WSPY", "1"))
	SpawnPointZ         = mustAtoi(get("MM_WSPZ", "1"))
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
