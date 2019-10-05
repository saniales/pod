package controller

import (
	"context"
	"math/rand"
	"time"

	"github.com/ugorji/go/codec"

	"github.com/p9c/pod/pkg/broadcast"
	chain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
)

type Blocks []*mining.BlockTemplate

// Run starts a controller instance
func Run(cx *conte.Xt) (cancel context.CancelFunc) {
	log.WARN("starting controller")
	var mh codec.MsgpackHandle
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	ciph := gcm.GetCipher(*cx.Config.MinerPass)
	outAddr, err := broadcast.New(*cx.Config.BroadcastAddress)
	if err != nil {
		log.ERROR(err)
		cancel()
		return
	}
	blockChan := make(chan Blocks)
	bytes := make([]byte, 0, broadcast.MaxDatagramSize)
	enc := codec.NewEncoderBytes(&bytes, &mh)
	go func() {
		for {
			// work dispatch loop
			select {
			case lb := <-blockChan:
				// send out block broadcast
				log.DEBUG("sending out block broadcast")
				// serialize initialBlocks
				//log.SPEW(lb)
				err := enc.Encode(lb)
				if err != nil {
					log.ERROR(err)
					break
				}
				log.SPEW(bytes)
				err = broadcast.Send(outAddr, bytes, ciph, broadcast.Template)
				if err != nil {
					log.ERROR(err)
				}
				// reset the bytes for next round
				bytes = bytes[:0]
				enc.ResetBytes(&bytes)
			case <-ctx.Done():
				// cancel has been called
				return
				// default:
			}
		}
	}()
	initialBlocks := Blocks{}
	// generate initial Blocks
	for algo := range fork.List[fork.GetCurrent(cx.RPCServer.Cfg.Chain.
		BestSnapshot().Height+1)].Algos {
		// Choose a payment address at random.
		rand.Seed(time.Now().UnixNano())
		payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(cx.
			StateCfg.ActiveMiningAddrs))]
		template, err := cx.RPCServer.Cfg.Generator.NewBlockTemplate(0,
			payToAddr, algo)
		if err != nil {
			log.ERROR("failed to create new block template:", err)
			continue
		}
		initialBlocks = append(initialBlocks, template)
	}
	blockChan <- initialBlocks
	// create subscriber for new block event
	cx.RPCServer.Cfg.Chain.Subscribe(func(n *chain.
	Notification) {
		switch n.Type {
		case chain.NTBlockConnected:
			log.WARN("new block found")
			blocks := Blocks{}
			// generate Blocks
			for algo := range fork.List[fork.GetCurrent(cx.RPCServer.Cfg.Chain.
				BestSnapshot().Height+1)].Algos {
				// Choose a payment address at random.
				rand.Seed(time.Now().UnixNano())
				payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(cx.
					StateCfg.ActiveMiningAddrs))]
				template, err := cx.RPCServer.Cfg.Generator.NewBlockTemplate(0,
					payToAddr, algo)
				if err != nil {
					log.ERROR("failed to create new block template:", err)
					continue
				}
				blocks = append(blocks, template)
			}
			blockChan <- blocks
		}
	})
	// goroutine loop checking for connection and sync status
	go func() {
		lastTxUpdate := cx.RPCServer.Cfg.Generator.GetTxSource().LastUpdated()
		time.Sleep(time.Second * 5)
		for {
			time.Sleep(time.Second)
			connCount := cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
			current := cx.RPCServer.Cfg.SyncMgr.IsCurrent()
			// if out of sync or disconnected,
			// once a second send out empty initialBlocks
			if (connCount < 1 && !*cx.Config.Solo) || !current {
				blockChan <- Blocks{}
			}
			// when new transactions are received the last updated timestamp
			// changes, when this happens a new dispatch needs to be made
			if lastTxUpdate != cx.RPCServer.Cfg.Generator.GetTxSource().LastUpdated() {
				log.WARN("new transactions added to block, " +
					"dispatching new block templates")
				blocks := Blocks{}
				// generate Blocks
				for algo := range fork.List[fork.GetCurrent(cx.RPCServer.Cfg.Chain.
					BestSnapshot().Height+1)].Algos {
					// Choose a payment address at random.
					rand.Seed(time.Now().UnixNano())
					payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(cx.
						StateCfg.ActiveMiningAddrs))]
					template, err := cx.RPCServer.Cfg.Generator.NewBlockTemplate(0,
						payToAddr, algo)
					if err != nil {
						log.ERROR("failed to create new block template:", err)
						continue
					}
					blocks = append(blocks, template)
				}
				blockChan <- blocks
			}
			select {
			case <-ctx.Done():
				break
			default:
			}
		}
	}()
	return
}
