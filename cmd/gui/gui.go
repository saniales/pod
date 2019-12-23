package gui

import (
	"github.com/p9c/pod/pkg/gio/app"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/pkg/log"
)

func WalletGUI(duo *duoui.DuoUI) (err error) {
	log.DEBUG("starting Wallet GUI")
	go func() {
		if err := duoui.DuoUImainLoop(duo); err != nil {
			log.FATAL(err.Error(), "- shutting down")
		}
	}()
	app.Main()
	log.DEBUG("GUI shut down")
	return
}
