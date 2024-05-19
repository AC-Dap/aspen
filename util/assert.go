package util

import "log"

func Assert(condition bool, message string) {
	if !condition {
		log.Panic(message)
	}
}
