package main

import (
	// "fmt"
	// "errors"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	// "github.com/skycoin/skycoin/src/cipher/encoder"
	"log"
	"strconv"
	"time"
)

type Container struct {
	c      *node.Container
	client *node.Client
	server *node.Server
	config *Config
	// msgs   chan *node.Msg
}

func NewContainer(config *Config) (c *Container, err error) {
	c = &Container{
		config: config,
		// msgs:   make(chan *Msg),
	}
	pk, _ := cipher.GenerateDeterministicKeyPair([]byte(config.Name))

	r := skyobject.NewRegistry()
	r.Register("Node", Node{})
	r.Register("Message", Message{})
	r.Register("NodeContainer", NodeContainer{})
	r.Done()

	cc := node.NewClientConfig()
	cc.InMemoryDB = config.CXOUseMemory()
	cc.DataDir = config.CXODir()
	cc.Log.Debug = true

	feeds := []cipher.PubKey{pk}

	// To handle the syncing use ClientConfig.OnAddFeed callback
	// Run cxo server and client.
	if c.client, err = node.NewClient(cc, r); err != nil {
		return
	}

	sc := node.NewServerConfig()
	sc.Log.Debug = true
	sc.EnableRPC = true
	sc.Listen = "[::]:" + strconv.Itoa(c.config.CXOPort())
	sc.InMemoryDB = config.CXOUseMemory()
	sc.DataDir = config.CXODir()

	c.server, err = node.NewServer(sc)
	if err != nil {
		return
	}

	if err = c.server.Start(); err != nil {
		c.server.Close()
		return
	}

	for _, f := range feeds {
		c.server.AddFeed(f)
	}
	if err = c.client.Start("[::]:" + strconv.Itoa(c.config.CXOPort())); err != nil {
		c.server.Close()
		c.client.Close()
		return
	}

	if c.client.Subscribe(pk) == false {
		log.Fatal("Unable to subscribe")
	}

	// Set Container.
	c.c = c.client.Container()

	// Wait.
	time.Sleep(5 * time.Second)
	return
}

func (c *Container) Close() error                      { return c.client.Close() }
func (c *Container) Connected() bool                   { return c.client.IsConnected() }
func (c *Container) Feeds() []cipher.PubKey            { return c.client.Feeds() }
func (c *Container) Subscribe(pk cipher.PubKey) bool   { return c.client.Subscribe(pk) }
func (c *Container) Unsubscribe(pk cipher.PubKey) bool { return c.client.Unsubscribe(pk) }

func makeNodeContainerFinder(r *node.Root) func(_ int, dRef skyobject.Dynamic) bool {
	return func(_ int, dRef skyobject.Dynamic) bool {
		schema, e := r.SchemaByReference(dRef.Schema)
		if e != nil {
			return false
		}
		return schema.Name() == "NodeContainer"
	}
}
