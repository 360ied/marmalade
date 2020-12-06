package config

import "os"

var (
	Address = get("MM_ADDR", "127.0.0.1:25565")
)

func get(key, fallback string) string {
	val, found := os.LookupEnv(key)
	if found {
		return val
	} else {
		return fallback
	}
}
