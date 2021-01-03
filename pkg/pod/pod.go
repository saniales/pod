package pod

import (
	"reflect"
	"sort"
	"sync"
	"time"
	
	"github.com/p9c/pod/app/appdata"
	log "github.com/p9c/pod/pkg/util/logi"
	
	"github.com/urfave/cli"
)

const AppName = "pod"

var (
	GitCommit string
	BuildTime string
	Tag       string
)

type Schema struct {
	Groups Groups `json:"groups"`
}
type Groups []Group

type Group struct {
	Legend string `json:"legend"`
	Fields `json:"fields"`
}

type Fields []Field

func (f Fields) Len() int {
	return len(f)
}

func (f Fields) Less(i, j int) bool {
	return f[i].Label < f[j].Label
}

func (f Fields) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type Field struct {
	Datatype    string   `json:"datatype"`
	Description string   `json:"help"`
	Featured    string   `json:"featured"`
	Group       string   `json:"group"`
	Hooks       string   `json:"hooks"`
	Label       string   `json:"label"`
	Model       string   `json:"model"`
	Options     []string `json:"options"`
	Restart     string   `json:"restart"`
	Slug        string   `json:"slug"`
	Type        string   `json:"type"`
	Widget      string   `json:"inputType"`
}

// GetConfigSchema returns a schema for a given config
func GetConfigSchema(cfg *Config) Schema {
	t := reflect.TypeOf(cfg)
	t = t.Elem()
	var levelOptions, network []string
	for _, i := range log.Levels {
		levelOptions = append(levelOptions, i)
	}
	network = []string{"mainnet", "testnet", "regtestnet", "simnet"}
	rawFields := make(map[string]Fields)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		var options []string
		switch {
		case field.Name == "LogLevel":
			options = levelOptions
		case field.Name == "Network":
			options = network
		}
		f := Field{
			Datatype:    field.Type.String(),
			Description: field.Tag.Get("description"),
			Featured:    field.Tag.Get("featured"),
			Group:       field.Tag.Get("group"),
			Hooks:       field.Tag.Get("hooks"),
			Label:       field.Tag.Get("label"),
			Model:       field.Tag.Get("json"),
			Options:     options,
			Slug:        field.Name,
			Type:        field.Tag.Get("type"),
			Widget:      field.Tag.Get("widget"),
		}
		if f.Group != "" {
			rawFields[f.Group] = append(rawFields[f.Group], f)
		}
	}
	for i := range rawFields {
		sort.Sort(rawFields[i])
	}
	var outGroups Groups
	var rf []string
	for i := range rawFields {
		rf = append(rf, i)
	}
	sort.Strings(rf)
	for i := range rf {
		rf[i], rf[len(rf)-1-i] = rf[len(rf)-1-i], rf[i]
	}
	for i := range rf {
		group := Group{
			Legend: rf[i],
			Fields: rawFields[rf[i]],
		}
		outGroups = append(Groups{group}, outGroups...)
	}
	return Schema{
		Groups: outGroups,
	}
}

// Config is
type Config struct {
	sync.Mutex
	AddCheckpoints         *cli.StringSlice `group:"debug" label:"AddCheckpoints" description:"add custom checkpoints" type:"" widget:"multi" json:"AddCheckpoints" hook:"restart"`
	AddPeers               *cli.StringSlice `group:"node" label:"Add Peers" description:"manually adds addresses to try to connect to" type:"address" widget:"multi" json:"AddPeers" hook:"addpeer"`
	AddrIndex              *bool            `group:"node" label:"Addr Index" description:"maintain a full address-based transaction index which makes the searchrawtransactions RPC available" type:"" widget:"toggle"  json:"AddrIndex" hook:"dropaddrindex"`
	AutoPorts              *bool            `group:"" label:"AutomaticPorts" description:"RPC and controller ports are randomized, use with controller for automatic peer discovery" type:"" widget:"toggle" json:"AutoPorts" hook:"restart"`
	BanDuration            *time.Duration   `group:"debug" label:"Ban Duration" description:"how long a ban of a misbehaving peer lasts" type:"" widget:"time" json:"BanDuration" hook:"restart"`
	BanThreshold           *int             `group:"debug" label:"Ban Threshold" description:"ban score that triggers a ban (default 100)" type:"" widget:"integer" json:"BanThreshold" hook:"restart"`
	BlockMaxSize           *int             `group:"mining" label:"Block Max Size" description:"maximum block size in bytes to be used when creating a block" type:"" widget:"integer" json:"BlockMaxSize" hook:"restart"`
	BlockMaxWeight         *int             `group:"mining" label:"Block Max Weight" description:"maximum block weight to be used when creating a block" type:"" widget:"integer" json:"BlockMaxWeight" hook:"restart"`
	BlockMinSize           *int             `group:"mining" label:"Block Min Size" description:"minimum block size in bytes to be used when creating a block" type:"" widget:"integer" json:"BlockMinSize" hook:"restart"`
	BlockMinWeight         *int             `group:"mining" label:"Block Min Weight" description:"minimum block weight to be used when creating a block" type:"" widget:"integer" json:"BlockMinWeight" hook:"restart"`
	BlockPrioritySize      *int             `group:"mining" label:"Block Priority Size" description:"size in bytes for high-priority/low-fee transactions when creating a block" type:"" widget:"integer" json:"BlockPrioritySize" hook:"restart"`
	BlocksOnly             *bool            `group:"node" label:"Blocks Only" description:"do not accept transactions from remote peers" type:"" widget:"toggle" json:"BlocksOnly" hook:"restart"`
	CAFile                 *string          `group:"tls" label:"Certificate Authority File" description:"certificate authority file for TLS certificate validation" type:"path" widget:"string" json:"CAFile" hook:"restart"`
	ConfigFile             *string          `group:"" label:"Configuration File" description:"location of configuration file, cannot actually be changed" type:"path" widget:"string" json:"ConfigFile" hook:"restart"`
	ConnectPeers           *cli.StringSlice `group:"node" label:"Connect Peers" description:"connect ONLY to these addresses (disables inbound connections)" type:"address" widget:"multi" json:"ConnectPeers" hook:"restart"`
	Controller             *string          `group:"node" label:"Controller Listener" description:"address to bind miner controller to" type:"address" widget:"string" json:"Controller" hook:"controller"`
	CPUProfile             *string          `group:"debug" label:"CPU Profile" description:"write cpu profile to this file" type:"path" widget:"string" json:"CPUProfile" hook:"restart"`
	DataDir                *string          `group:"" label:"Data Directory" description:"root folder where application data is stored" type:"path" widget:"string" json:"DataDir" hook:"restart"`
	DbType                 *string          `group:"" label:"Database Type" description:"type of database storage engine to use (only one right now)" type:"" widget:"string" json:"DbType" hook:"restart"`
	DisableBanning         *bool            `group:"debug" label:"Disable Banning" description:"disables banning of misbehaving peers" type:"" widget:"toggle" json:"DisableBanning" hook:"restart"`
	DisableCheckpoints     *bool            `group:"debug" label:"Disable Checkpoints" description:"disables all checkpoints" type:"" widget:"toggle" json:"DisableCheckpoints" hook:"restart"`
	DisableDNSSeed         *bool            `group:"node" label:"Disable DNS Seed" description:"disable seeding of addresses to peers" type:"" widget:"toggle" json:"DisableDNSSeed" hook:"restart"`
	DisableListen          *bool            `group:"node" label:"Disable Listen" description:"disables inbound connections for the peer to peer network" type:"" widget:"toggle" json:"DisableListen" hook:"restart"`
	DisableRPC             *bool            `group:"" label:"Disable RPC" description:"disable rpc servers" type:"" widget:"toggle" json:"DisableRPC" hook:"restart"`
	ExternalIPs            *cli.StringSlice `group:"node" label:"External IP Addresses" description:"extra addresses to tell peers they can connect to" type:"address" widget:"multi" json:"ExternalIPs" hook:"restart"`
	FreeTxRelayLimit       *float64         `group:"policy" label:"Free Tx Relay Limit" description:"limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute" type:"" widget:"float" json:"FreeTxRelayLimit" hook:"restart"`
	Generate               *bool            `group:"mining" label:"Generate Blocks" description:"turn on Kopach CPU miner" type:"" widget:"toggle" json:"Generate" hook:"generate"`
	GenThreads             *int             `group:"mining" label:"Gen Threads" description:"number of threads to mine with" type:"" widget:"integer" json:"GenThreads" hook:"genthreads"`
	Language               *string          `group:"config" label:"Language" description:"user interface language i18 localization" type:"" widget:"string" json:"Language" hook:"language"`
	LimitPass              *string          `group:"rpc" label:"Limit Pass" description:"limited user password" type:"" widget:"password" json:"LimitPass" hook:"restart"`
	LimitUser              *string          `group:"rpc" label:"Limit User" description:"limited user name" type:"" widget:"string" json:"LimitUser" hook:"restart"`
	Listeners              *cli.StringSlice `group:"node" label:"Listeners" description:"list of addresses to bind the node listener to" type:"address" widget:"multi" json:"Listeners" hook:"restart"`
	LogDir                 *string          `group:"config" label:"Log Dir" description:"folder where log files are written" type:"path" widget:"string" json:"LogDir" hook:"restart"`
	LogLevel               *string          `group:"config" label:"Log Level" description:"maximum log level to output\n(fatal error check warning info debug trace - what is selected includes all items to the left of the one in that list)" type:"" widget:"radio" json:"LogLevel" hook:"loglevel"`
	MaxOrphanTxs           *int             `group:"policy" label:"Max Orphan Txs" description:"max number of orphan transactions to keep in memory" type:"" widget:"integer" json:"MaxOrphanTxs" hook:"restart"`
	MaxPeers               *int             `group:"node" label:"Max Peers" description:"maximum number of peers to hold connections with" type:"" widget:"integer" json:"MaxPeers" hook:"restart"`
	MinerPass              *string          `group:"mining" label:"Miner Pass" description:"password that encrypts the connection to the mining controller" type:"" widget:"password" json:"MinerPass" hook:"restart"`
	MiningAddrs            *cli.StringSlice `group:"" label:"Mining Addrs" description:"addresses to pay block rewards to (TODO, make this auto)" type:"base58" widget:"multi" json:"MiningAddrs" hook:"miningaddr"`
	MinRelayTxFee          *float64         `group:"policy" label:"Min Relay Tx Fee" description:"the minimum transaction fee in DUO/kB to be considered a non-zero fee" type:"" widget:"float" json:"MinRelayTxFee" hook:"restart"`
	Network                *string          `group:"node" label:"Network" description:"connect to this network: mainnet, testnet)" type:"" widget:"radio" json:"Network" hook:"restart"`
	NoCFilters             *bool            `group:"node" label:"No CFilters" description:"disable committed filtering (CF) support" type:"" widget:"toggle" json:"NoCFilters" hook:"restart"`
	NodeOff                *bool            `group:"debug" label:"Node Off" description:"turn off the node backend" type:"" widget:"toggle" json:"NodeOff" hook:"node"`
	NoInitialLoad          *bool            `group:"" label:"No initial load" description:"do not load a wallet at startup" type:"" widget:"toggle" json:"NoInitialLoad" hook:"restart"`
	NoPeerBloomFilters     *bool            `group:"node" label:"No Peer Bloom Filters" description:"disable bloom filtering support" type:"" widget:"toggle" json:"NoPeerBloomFilters" hook:"restart"`
	NoRelayPriority        *bool            `group:"policy" label:"No Relay Priority" description:"do not require free or low-fee transactions to have high priority for relaying" type:"" widget:"toggle" json:"NoRelayPriority" hook:"restart"`
	OneTimeTLSKey          *bool            `group:"wallet" label:"One Time TLS Key" description:"generate a new TLS certificate pair at startup, but only write the certificate to disk" type:"" widget:"toggle" json:"OneTimeTLSKey" hook:"restart"`
	Onion                  *bool            `group:"proxy" label:"Onion" description:"enable tor proxy" type:"" widget:"toggle" json:"Onion" hook:"restart"`
	OnionProxy             *string          `group:"proxy" label:"Onion Proxy" description:"address of tor proxy you want to connect to" type:"address" widget:"string" json:"OnionProxy" hook:"restart"`
	OnionProxyPass         *string          `group:"proxy" label:"Onion Proxy Pass" description:"password for tor proxy" type:"" widget:"password" json:"OnionProxyPass" hook:"restart"`
	OnionProxyUser         *string          `group:"proxy" label:"Onion Proxy User" description:"tor proxy username" type:"" widget:"string" json:"OnionProxyUser" hook:"restart"`
	Password               *string          `group:"rpc" label:"Password" description:"password for client RPC connections" type:"" widget:"password" json:"Password" hook:"restart"`
	PipeLog                *bool            `group:"" label:"Pipe Logger" description:"enable pipe based loggerIPC" type:"" widget:"toggle" json:"PipeLog" hook:""`
	Profile                *string          `group:"debug" label:"Profile" description:"http profiling on given port (1024-40000)" type:"url" widget:"string" json:"Profile" hook:"restart"`
	Proxy                  *string          `group:"proxy" label:"Proxy" description:"address of proxy to connect to for outbound connections" type:"url" widget:"string" json:"Proxy" hook:"restart"`
	ProxyPass              *string          `group:"proxy" label:"Proxy Pass" description:"proxy password, if required" type:"" widget:"password" json:"ProxyPass" hook:"restart"`
	ProxyUser              *string          `group:"proxy" label:"ProxyUser" description:"proxy username, if required" type:"" widget:"string" json:"ProxyUser" hook:"restart"`
	RejectNonStd           *bool            `group:"node" label:"Reject Non Std" description:"reject non-standard transactions regardless of the default settings for the active network" type:"" widget:"toggle" json:"RejectNonStd" hook:"restart"`
	RelayNonStd            *bool            `group:"node" label:"Relay Non Std" description:"relay non-standard transactions regardless of the default settings for the active network" type:"" widget:"toggle" json:"RelayNonStd" hook:"restart"`
	RPCCert                *string          `group:"rpc" label:"RPC Cert" description:"location of RPC TLS certificate" type:"path" widget:"string" json:"RPCCert" hook:"restart"`
	RPCConnect             *string          `group:"wallet" label:"RPC Connect" description:"full node RPC for wallet" type:"address" widget:"string" json:"RPCConnect" hook:"restart"`
	RPCKey                 *string          `group:"rpc" label:"RPC Key" description:"location of rpc TLS key" type:"path" widget:"string" json:"RPCKey" hook:"restart"`
	RPCListeners           *cli.StringSlice `group:"rpc" label:"RPC Listeners" description:"addresses to listen for RPC connections" type:"address" widget:"multi" json:"RPCListeners" hook:"restart"`
	RPCMaxClients          *int             `group:"rpc" label:"Maximum RPC Clients" description:"maximum number of clients for regular RPC" type:"" widget:"integer" json:"RPCMaxClients" hook:"restart"`
	RPCMaxConcurrentReqs   *int             `group:"rpc" label:"Maximum RPC Concurrent Reqs" description:"maximum number of requests to process concurrently" type:"" widget:"integer" json:"RPCMaxConcurrentReqs" hook:"restart"`
	RPCMaxWebsockets       *int             `group:"rpc" label:"Maximum RPC Websockets" description:"maximum number of websocket clients to allow" type:"" widget:"integer" json:"RPCMaxWebsockets" hook:"restart"`
	RPCQuirks              *bool            `group:"rpc" label:"RPC Quirks" description:"enable bugs that replicate bitcoin core RPC's JSON" type:"" widget:"toggle" json:"RPCQuirks" hook:"restart"`
	ServerPass             *string          `group:"rpc" label:"Server Pass" description:"password for server connections" type:"" widget:"password" json:"ServerPass" hook:"restart"`
	ServerTLS              *bool            `group:"wallet" label:"Server TLS" description:"enable TLS for the wallet connection to node RPC server" type:"" widget:"toggle" json:"ServerTLS" hook:"restart"`
	ServerUser             *string          `group:"rpc" label:"Server User" description:"username for chain server connections" type:"" widget:"string" json:"ServerUser" hook:"restart"`
	SigCacheMaxSize        *int             `group:"node" label:"Sig Cache Max Size" description:"the maximum number of entries in the signature verification cache" type:"" widget:"integer" json:"SigCacheMaxSize" hook:"restart"`
	Solo                   *bool            `group:"mining" label:"Solo Generate" description:"mine even if not connected to a network" type:"" widget:"toggle" json:"Solo" hook:"restart"`
	TLS                    *bool            `group:"tls" label:"TLS" description:"enable TLS for RPC connections" type:"" widget:"toggle" json:"TLS" hook:"restart"`
	TLSSkipVerify          *bool            `group:"tls" label:"TLS Skip Verify" description:"skip TLS certificate verification (ignore CA errors)" type:"" widget:"toggle" json:"TLSSkipVerify" hook:"restart"`
	TorIsolation           *bool            `group:"proxy" label:"Tor Isolation" description:"makes a separate proxy connection for each connection" type:"" widget:"toggle" json:"TorIsolation" hook:"restart"`
	TrickleInterval        *time.Duration   `group:"policy" label:"Trickle Interval" description:"minimum time between attempts to send new inventory to a connected peer" type:"" widget:"time" json:"TrickleInterval" hook:"restart"`
	TxIndex                *bool            `group:"node" label:"Tx Index" description:"maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC" type:"" widget:"toggle" json:"TxIndex" hook:"droptxindex"`
	UPNP                   *bool            `group:"node" label:"UPNP" description:"enable UPNP for NAT traversal" type:"" widget:"toggle" json:"UPNP" hook:"restart"`
	UserAgentComments      *cli.StringSlice `group:"" label:"User Agent Comments" description:"comment to add to the user agent -- See BIP 14 for more information" type:"" widget:"multi" json:"UserAgentComments" hook:"restart"`
	Username               *string          `group:"rpc" label:"Username" description:"password for client RPC connections" type:"" widget:"string" json:"Username" hook:"restart"`
	Wallet                 *bool            `group:"debug" label:"Connect to Wallet" description:"set ctl to connect to wallet instead of chain server" type:"" widget:"toggle" json:"Wallet"`
	WalletFile             *string          `group:"config" label:"Wallet File" description:"wallet database file" type:"path" widget:"string" featured:"true" json:"WalletFile" hook:"restart"`
	WalletOff              *bool            `group:"debug" label:"Wallet Off" description:"turn off the wallet backend" type:"" widget:"toggle" json:"WalletOff" hook:"wallet"`
	WalletPass             *string          `group:"" label:"Wallet Pass" description:"password encrypting public data in wallet - hash is stored so give on command line" type:"" widget:"password" json:"WalletPass" hook:"restart"`
	WalletRPCListeners     *cli.StringSlice `group:"wallet" label:"Legacy RPC Listeners" description:"addresses for wallet RPC server to listen on" type:"address" widget:"multi" json:"WalletRPCListeners" hook:"restart"`
	WalletRPCMaxClients    *int             `group:"wallet" label:"Legacy RPC Max Clients" description:"maximum number of RPC clients allowed for wallet RPC" type:"" widget:"integer" json:"WalletRPCMaxClients" hook:"restart"`
	WalletRPCMaxWebsockets *int             `group:"wallet" label:"Legacy RPC Max Websockets" description:"maximum number of websocket clients allowed for wallet RPC" type:"" widget:"integer" json:"WalletRPCMaxWebsockets" hook:"restart"`
	WalletServer           *string          `group:"wallet" label:"Wallet Server" description:"node address to connect wallet server to" type:"address" widget:"string" json:"WalletServer" hook:"restart"`
	Whitelists             *cli.StringSlice `group:"debug" label:"Whitelists" description:"peers that you don't want to ever ban" type:"address" widget:"multi" json:"Whitelists" hook:"restart"`
	LAN                    *bool            `group:"debug" label:"LAN" description:"run without any connection to nodes on the internet (does not apply on mainnet)" type:"" widget:"toggle" json:"LAN" hook:"restart"`
	DarkTheme              *bool            `group:"config" label:"Dark Theme" description:"sets dark theme for GUI" type:"" widget:"toggle" json:"DarkTheme" hook:"restart"`
	RunAsService           *bool            `group:"" label:"Run As Service" description:"shuts down on lock timeout" type:"" widget:"toggle" json:"" hook:"restart"`
	CAPI                   *bool            `group:"" label:"Enable cAPI" description:"disable cAPI rpc" type:"" widget:"toggle" json:"CAPI" hook:"restart"`
}

func EmptyConfig() (c *Config, conf map[string]interface{}) {
	datadir := appdata.Dir(AppName, false)
	c = &Config{
		AddCheckpoints:         newStringSlice(),
		AddPeers:               newStringSlice(),
		AddrIndex:              newbool(),
		AutoPorts:              newbool(),
		BanDuration:            newDuration(),
		BanThreshold:           newint(),
		BlockMaxSize:           newint(),
		BlockMaxWeight:         newint(),
		BlockMinSize:           newint(),
		BlockMinWeight:         newint(),
		BlockPrioritySize:      newint(),
		BlocksOnly:             newbool(),
		CAFile:                 newstring(),
		CAPI:                   newbool(),
		ConfigFile:             newstring(),
		ConnectPeers:           newStringSlice(),
		Controller:             newstring(),
		CPUProfile:             newstring(),
		DarkTheme:              newbool(),
		DataDir:                &datadir,
		DbType:                 newstring(),
		DisableBanning:         newbool(),
		DisableCheckpoints:     newbool(),
		DisableDNSSeed:         newbool(),
		DisableListen:          newbool(),
		DisableRPC:             newbool(),
		ExternalIPs:            newStringSlice(),
		FreeTxRelayLimit:       newfloat64(),
		Generate:               newbool(),
		GenThreads:             newint(),
		LAN:                    newbool(),
		Language:               newstring(),
		LimitPass:              newstring(),
		LimitUser:              newstring(),
		Listeners:              newStringSlice(),
		LogDir:                 newstring(),
		LogLevel:               newstring(),
		MaxOrphanTxs:           newint(),
		MaxPeers:               newint(),
		MinerPass:              newstring(),
		MiningAddrs:            newStringSlice(),
		MinRelayTxFee:          newfloat64(),
		Network:                newstring(),
		NoCFilters:             newbool(),
		NodeOff:                newbool(),
		NoInitialLoad:          newbool(),
		NoPeerBloomFilters:     newbool(),
		NoRelayPriority:        newbool(),
		OneTimeTLSKey:          newbool(),
		Onion:                  newbool(),
		OnionProxy:             newstring(),
		OnionProxyPass:         newstring(),
		OnionProxyUser:         newstring(),
		Password:               newstring(),
		PipeLog:                newbool(),
		Profile:                newstring(),
		Proxy:                  newstring(),
		ProxyPass:              newstring(),
		ProxyUser:              newstring(),
		RejectNonStd:           newbool(),
		RelayNonStd:            newbool(),
		RPCCert:                newstring(),
		RPCConnect:             newstring(),
		RPCKey:                 newstring(),
		RPCListeners:           newStringSlice(),
		RPCMaxClients:          newint(),
		RPCMaxConcurrentReqs:   newint(),
		RPCMaxWebsockets:       newint(),
		RPCQuirks:              newbool(),
		RunAsService:           newbool(),
		ServerPass:             newstring(),
		ServerTLS:              newbool(),
		ServerUser:             newstring(),
		SigCacheMaxSize:        newint(),
		Solo:                   newbool(),
		TLS:                    newbool(),
		TLSSkipVerify:          newbool(),
		TorIsolation:           newbool(),
		TrickleInterval:        newDuration(),
		TxIndex:                newbool(),
		UPNP:                   newbool(),
		UserAgentComments:      newStringSlice(),
		Username:               newstring(),
		Wallet:                 newbool(),
		WalletFile:             newstring(),
		WalletOff:              newbool(),
		WalletPass:             newstring(),
		WalletRPCListeners:     newStringSlice(),
		WalletRPCMaxClients:    newint(),
		WalletRPCMaxWebsockets: newint(),
		WalletServer:           newstring(),
		Whitelists:             newStringSlice(),
	}
	conf = map[string]interface{}{
		"AddCheckpoints":         c.AddCheckpoints,
		"AddPeers":               c.AddPeers,
		"AddrIndex":              c.AddrIndex,
		"AutoPorts":              c.AutoPorts,
		"BanDuration":            c.BanDuration,
		"BanThreshold":           c.BanThreshold,
		"BlockMaxSize":           c.BlockMaxSize,
		"BlockMaxWeight":         c.BlockMaxWeight,
		"BlockMinSize":           c.BlockMinSize,
		"BlockMinWeight":         c.BlockMinWeight,
		"BlockPrioritySize":      c.BlockPrioritySize,
		"BlocksOnly":             c.BlocksOnly,
		"CAFile":                 c.CAFile,
		"CAPI":                   c.CAPI,
		"ConfigFile":             c.ConfigFile,
		"ConnectPeers":           c.ConnectPeers,
		"Controller":             c.Controller,
		"CPUProfile":             c.CPUProfile,
		"DarkTheme":              c.DarkTheme,
		"DataDir":                c.DataDir,
		"DbType":                 c.DbType,
		"DisableBanning":         c.DisableBanning,
		"DisableCheckpoints":     c.DisableCheckpoints,
		"DisableDNSSeed":         c.DisableDNSSeed,
		"DisableListen":          c.DisableListen,
		"DisableRPC":             c.DisableRPC,
		"ExternalIPs":            c.ExternalIPs,
		"FreeTxRelayLimit":       c.FreeTxRelayLimit,
		"Generate":               c.Generate,
		"GenThreads":             c.GenThreads,
		"LAN":                    c.LAN,
		"Language":               c.Language,
		"LimitPass":              c.LimitPass,
		"LimitUser":              c.LimitUser,
		"Listeners":              c.Listeners,
		"LogDir":                 c.LogDir,
		"LogLevel":               c.LogLevel,
		"MaxOrphanTxs":           c.MaxOrphanTxs,
		"MaxPeers":               c.MaxPeers,
		"MinerPass":              c.MinerPass,
		"MiningAddrs":            c.MiningAddrs,
		"MinRelayTxFee":          c.MinRelayTxFee,
		"Network":                c.Network,
		"NoCFilters":             c.NoCFilters,
		"NodeOff":                c.NodeOff,
		"NoInitialLoad":          c.NoInitialLoad,
		"NoPeerBloomFilters":     c.NoPeerBloomFilters,
		"NoRelayPriority":        c.NoRelayPriority,
		"OneTimeTLSKey":          c.OneTimeTLSKey,
		"Onion":                  c.Onion,
		"OnionProxy":             c.OnionProxy,
		"OnionProxyPass":         c.OnionProxyPass,
		"OnionProxyUser":         c.OnionProxyUser,
		"Password":               c.Password,
		"PipeLog":                c.PipeLog,
		"Profile":                c.Profile,
		"Proxy":                  c.Proxy,
		"ProxyPass":              c.ProxyPass,
		"ProxyUser":              c.ProxyUser,
		"RejectNonStd":           c.RejectNonStd,
		"RelayNonStd":            c.RelayNonStd,
		"RPCCert":                c.RPCCert,
		"RPCConnect":             c.RPCConnect,
		"RPCKey":                 c.RPCKey,
		"RPCListeners":           c.RPCListeners,
		"RPCMaxClients":          c.RPCMaxClients,
		"RPCMaxConcurrentReqs":   c.RPCMaxConcurrentReqs,
		"RPCMaxWebsockets":       c.RPCMaxWebsockets,
		"RPCQuirks":              c.RPCQuirks,
		"RunAsService":           c.RunAsService,
		"ServerPass":             c.ServerPass,
		"ServerTLS":              c.ServerTLS,
		"ServerUser":             c.ServerUser,
		"SigCacheMaxSize":        c.SigCacheMaxSize,
		"Solo":                   c.Solo,
		"TLS":                    c.TLS,
		"TLSSkipVerify":          c.TLSSkipVerify,
		"TorIsolation":           c.TorIsolation,
		"TrickleInterval":        c.TrickleInterval,
		"TxIndex":                c.TxIndex,
		"UPNP":                   c.UPNP,
		"UserAgentComments":      c.UserAgentComments,
		"Username":               c.Username,
		"Wallet":                 c.Wallet,
		"WalletFile":             c.WalletFile,
		"WalletOff":              c.WalletOff,
		"WalletPass":             c.WalletPass,
		"WalletRPCListeners":     c.WalletRPCListeners,
		"WalletRPCMaxClients":    c.WalletRPCMaxClients,
		"WalletRPCMaxWebsockets": c.WalletRPCMaxWebsockets,
		"WalletServer":           c.WalletServer,
		"Whitelists":             c.Whitelists,
	}
	return
}

func newbool() *bool {
	o := false
	return &o
}
func newStringSlice() *cli.StringSlice {
	o := cli.StringSlice{}
	return &o
}
func newfloat64() *float64 {
	o := 1.0
	return &o
}
func newint() *int {
	o := 1
	return &o
}
func newstring() *string {
	o := ""
	return &o
}
func newDuration() *time.Duration {
	o := time.Second - time.Second
	return &o
}
