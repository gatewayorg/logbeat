package main

import (
	"os"
	"time"

	"github.com/Ankr-network/kit/mlog"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var log = mlog.Logger("main")

func main() {
	svr := cli.NewApp()
	svr.Action = mainServe
	err := svr.Run(os.Args)
	if err != nil {
		log.Fatal("Serverice Crash", zap.Error(err))
	}
}

func mainServe(c *cli.Context) error {
	for {
		time.Sleep(5 * time.Second)
		log.Info("start")
	}

}
