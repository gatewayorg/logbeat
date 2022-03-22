package main

import (
	"os"

	"github.com/Ankr-network/kit/mlog"
	"github.com/gatewayorg/logbeat/app"
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
		&cli.StringSliceFlag{
			Name:     share.MQ_ADDRESS,
			Required: true,
			Usage:    "mq address",
		},
		&cli.StringSliceFlag{
			Name:     share.FILTER_SETTINGS,
			Required: false,
			Usage:    "filter settings",
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

	log.Info("init", zap.Any("filter", c.StringSlice(share.FILTER_SETTINGS)))

	log.Info("init", zap.Any("pub", c.StringSlice(share.MQ_ADDRESS)))
	pubMetrics := app.NewPubMetrics(c.StringSlice(share.MQ_ADDRESS))
	log.Info("init", zap.Any("dir", c.String(share.LOG_DIR)))
	app.StartProcess(c.String(share.LOG_DIR), pubMetrics, c.StringSlice(share.FILTER_SETTINGS))

	// app.StartProcess(c.String(share.LOG_DIR))

	return nil

}
