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
	SpawnPointX         = 32 * uint16(mustAtoi(get("MM_WSP_X", "1")))
	SpawnPointY         = 32 * uint16(mustAtoi(get("MM_WSP_Y", "1")))
	SpawnPointZ         = 32 * uint16(mustAtoi(get("MM_WSP_Z", "1")))
	SpawnPointYaw       = uint8(mustAtoi(get("MM_WSP_YAW", "96"))) // southeast, looking away from 0, 0
	SpawnPointPitch     = uint8(mustAtoi(get("MM_WSP_PITCH", "0")))
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
