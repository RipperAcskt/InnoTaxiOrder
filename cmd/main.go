package main

import (
	"log"

	_ "net/http/pprof"

	"github.com/RipperAcskt/innotaxiorder/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("app run failed: %v", err)
	}
}
