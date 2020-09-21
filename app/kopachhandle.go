package app

import (
	"github.com/stalker-loki/pod/app/config"
	"os"

	"github.com/urfave/cli"

	"github.com/stalker-loki/pod/cmd/kopach"
	"github.com/stalker-loki/pod/pkg/chain/config/netparams"
	"github.com/stalker-loki/pod/pkg/chain/fork"

	"github.com/stalker-loki/pod/app/conte"
	"github.com/stalker-loki/pod/pkg/util/interrupt"
)

func KopachHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		Info("starting up kopach standalone miner for parallelcoin")
		config.Configure(cx, c.Command.Name, true)
		if cx.ActiveNet.Name == netparams.TestNet3Params.Name {
			fork.IsTestnet = true
		}
		quit := make(chan struct{})
		interrupt.AddHandler(func() {
			Debug("KopachHandle interrupt")
			close(quit)
			os.Exit(0)
		})
		err = kopach.KopachHandle(cx)(c)
		<-quit
		Debug("kopach main finished")
		return
	}
}
