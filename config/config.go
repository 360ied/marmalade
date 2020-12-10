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
	WorldPath           = get("MM_WPATH", "world.ucw")                                    // uncompressed classic world
	WorldScratchPath    = get("MM_WSPATH", WorldPath+"2")
	WorldTempPath       = get("MM_WTMPRNPATH", WorldPath+"_TMP")
	WorldSaveDelay      = time.Second * time.Duration(mustAtoi(get("MM_WSAVEDELAY", "30")))
	CommandPrefix       = get("MM_CMDPRFX", "/")
	WelcomeMessage      = get("MM_WELCOMEMSG",
		"marmalade is free and open source software licensed under the GNU Affero General Public License. "+
			"The full source code can be found at https://github.com/360ied/marmalade")
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
