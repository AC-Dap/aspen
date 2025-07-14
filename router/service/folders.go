package service

import "github.com/rs/zerolog/log"

// All service folders are created in this folder
var globalFolder string = ""

func SetGlobalFolder(folder string) {
	globalFolder = folder
}

func getServiceFolder(serviceId string) string {
	if globalFolder == "" {
		log.Fatal().Msg("Global service folder not initialized.")
	}

	return globalFolder + "/" + serviceId
}
