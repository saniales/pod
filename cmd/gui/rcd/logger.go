package rcd

import (
	log "github.com/p9c/pod/pkg/logi"
)

var (
	// MaxLogLength is a var so it can be changed dynamically
	MaxLogLength = 16384
)

func (r *RcVar) DuoUIloggerController() {
	L.LogChan = make(chan log.Entry)
	r.Log.LogChan = L.LogChan
	log.L.SetLevel(*r.cx.Config.LogLevel, true, "pod")
	go func() {
	out:
		for {
			select {
			case n := <-L.LogChan:
				le, ok := r.Log.LogMessages.Load().([]log.Entry)
				if ok {
					le = append(le, n)
					// Once length exceeds MaxLogLength we trim off the start to keep it the same size
					ll := len(le)
					if ll > MaxLogLength {
						le = le[ll-MaxLogLength:]
					}
					r.Log.LogMessages.Store(le)
				} else {
					r.Log.LogMessages.Store([]log.Entry{n})
				}
			case <-r.Log.StopLogger:
				defer func() {
					r.Log.StopLogger = make(chan struct{})
				}()
				r.Log.LogMessages.Store([]log.Entry{})
				L.LogChan = nil
				break out
			}
		}
	}()
}
