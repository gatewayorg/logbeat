package main

import (
	"os"

	"github.com/Ankr-network/kit/mlog"
	"github.com/gatewayorg/logbeat/share"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var log = mlog.Logger("main")

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     share.LOG_DIR,
			Usage:    "log dir path",
			Required: true,
		},
	}
	svr := cli.NewApp()
	svr.Action = mainServe
	svr.Flags = flags
	err := svr.Run(os.Args)
	if err != nil {
		log.Fatal("Serverice Crash", zap.Error(err))
	}
}

func mainServe(c *cli.Context) error {

	_, err := os.Stat(c.String(share.LOG_DIR))
	if err != nil {
		log.Error("dir not exist", zap.Error(err))
		return err
	}

	// app.StartProcess(c.String(share.LOG_DIR))

	return nil

}
