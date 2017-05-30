package main

import (
	"errors"
	"flag"
	"github.com/skycoin/messenger/misc"
)

// Config represents commandline arguments.
type Config struct {

	// [TEST MODE] enforces the following behaviours:
	// - `cxoMemoryMode = true` (disables modification to cxo database).
	// - `saveConfig = false` (disables modification to config files).

	Name                string // own name, can be changed
	testMode            bool   // Whether to enable test mode.
	testModeThreads     int    // Number of threads to use for test mode (will create them in test mode).
	testModeMinInterval int    // Minimum interval between simulated activity (in seconds).
	testModeMaxInterval int    // Maximum interval between simulated activity (in seconds).

	saveConfig      bool   // Whether to save and use messenger configuration files.
	configDir       string // Configuration directory.
	rpcServerPort   int    // RPC server port (master node only).
	rpcServerRemAdr string // RPC remote address (master node only).
	cxoPort         int    // Port of CXO Daemon.
	cxoMemoryMode   bool   // Whether to use in-memory database for CXO.
	cxoDir          string // Folder name to store db.
}

// NewConfig makes Config with default values.
func NewConfig() *Config {
	return &Config{
		Name:                misc.MakeRandomAlias(),
		testMode:            false,
		testModeThreads:     3,
		testModeMinInterval: 1,
		testModeMaxInterval: 10,

		saveConfig:      true,
		configDir:       ".",
		rpcServerPort:   6421,
		rpcServerRemAdr: "127.0.0.1:6421",
		cxoPort:         8998,
		cxoMemoryMode:   true,
		cxoDir:          "msg",
	}
}

// Parse fills the Config with commandline argument values.
func (c *Config) Parse() *Config {
	/*
		<<< TEST FLAGS >>>
	*/

	flag.BoolVar(&c.testMode,
		"test-mode", c.testMode,
		"whether to enable test mode")

	flag.IntVar(&c.testModeThreads,
		"test-mode-threads", c.testModeThreads,
		"number of threads to use for test mode")

	flag.IntVar(&c.testModeMinInterval,
		"test-mode-min", c.testModeMinInterval,
		"minimum interval in seconds between simulated activity")

	flag.IntVar(&c.testModeMaxInterval,
		"test-mode-max", c.testModeMaxInterval,
		"maximum interval in seconds between simulated activity")

	/*
		<<< Config FLAGS >>>
	*/

	flag.StringVar(&c.Name,
		"name", c.Name,
		"node name")

	flag.BoolVar(&c.saveConfig,
		"save-config", c.saveConfig,
		"whether to save and use configuration files")

	flag.StringVar(&c.configDir,
		"config-dir", c.configDir,
		"configuration directory")

	flag.IntVar(&c.rpcServerPort,
		"rpc-server-port", c.rpcServerPort,
		"port of rpc server for master node")

	flag.StringVar(&c.rpcServerRemAdr,
		"rpc-server-remote-address", c.rpcServerRemAdr,
		"remote address of rpc server for master node")

	flag.IntVar(&c.cxoPort,
		"cxo-port", c.cxoPort,
		"port of cxo daemon to connect to")

	flag.BoolVar(&c.cxoMemoryMode,
		"cxo-memory-mode", c.cxoMemoryMode,
		"whether to use in-memory database")

	flag.StringVar(&c.cxoDir,
		"cxo-dir", c.cxoDir,
		"folder to store cxo db files in")

	flag.Parse()
	return c
}

// PostProcess checks the validity and post processes the flags.
func (c *Config) PostProcess() (*Config, error) {
	// Action on test mode.
	if c.testMode {
		// Check test mode settings.
		if c.testModeThreads < 0 {
			return nil, errors.New("invalid number of test mode threads specified")
		}
		if c.testModeMinInterval < 1 {
			return nil, errors.New("invalid test mode minimum interval specified")
		}
		if c.testModeMaxInterval < 1 {
			return nil, errors.New("invalid test mode maximum interval specified")
		}
		if c.testModeMinInterval > c.testModeMaxInterval {
			return nil, errors.New("test mode minimum interval > maximum interval")
		}
		// Enforce behaviour.
		c.cxoMemoryMode = true
		c.saveConfig = false
	}
	return c, nil
}

/*
	These functions ensure that configuration values are not accidentally modified.
*/

func (c *Config) TestMode() bool           { return c.testMode }
func (c *Config) TestModeThreads() int     { return c.testModeThreads }
func (c *Config) TestModeMinInterval() int { return c.testModeMinInterval }
func (c *Config) TestModeMaxInterval() int { return c.testModeMaxInterval }

func (c *Config) SaveConfig() bool        { return c.saveConfig }
func (c *Config) ConfigDir() string       { return c.configDir }
func (c *Config) RPCServerPort() int      { return c.rpcServerPort }
func (c *Config) RPCServerRemAdr() string { return c.rpcServerRemAdr }
func (c *Config) CXOPort() int            { return c.cxoPort }
func (c *Config) CXOUseMemory() bool      { return c.cxoMemoryMode }
func (c *Config) CXODir() string          { return c.cxoDir }
