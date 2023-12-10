package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/botscubes/bql/internal/app"
	"github.com/botscubes/bql/pkg/logger"
)

func main() {
	log, err := logger.NewLogger(logger.Config{
		Type: "dev",
	})
	if err != nil {
		fmt.Printf("Create logger: %v\n", err)
		return
	}

	defer func() {
		if runtime.GOOS != "windows" {
			if err := log.Sync(); err != nil {
				log.Errorw("failed log sync", "error", err)
			}
		}
	}()

	if len(os.Args) != 2 {
		log.Info(`example usage: ./main code.txt`)
		return
	}

	app.Start(log, os.Args[1])

}
