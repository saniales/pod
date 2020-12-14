package main

import (
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/pipe/consume"
	"os"
	"strings"
	"time"
)

func main() {
	var err error
	logi.L.SetLevel("debug", false, "pod")
	command := strings.Join(os.Args[1:], " ") // "./pod -D test0 -n testnet -l trace --solo --lan --pipelog node"
	quit := make(chan struct{})
	w := consume.Log(quit, consume.SimpleLog, consume.FilterNone, strings.Split(command, " ")...)
	Debug("starting")
	consume.Start(w)
	time.Sleep(time.Second * 5)
	Debug("killing")
	consume.Kill(w)
	if err = w.Wait(); Check(err) {
	}
}