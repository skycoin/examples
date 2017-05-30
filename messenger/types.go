package main

import (
	"errors"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/messenger/misc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"strings"
	"time"
)

// NodeContainer represents the first branch object from root, in cxo.
type NodeContainer struct {
	Subscribed skyobject.References `skyobject:"schema=Node"`
	Connected  skyobject.References `skyobject:"schema=Node"`
	Queue      skyobject.References `skyobject:"schema=Message"`
}

// Node represents a node stored in cxo.
type Node struct {
	Name   string `json:"name"`
	PubKey string `json:"public_key"`
}

// Check checks the validity of the node, outputs it's public key and updates it's fields.
func (n *Node) Check() (cipher.PubKey, error) {
	npk, e := cipher.PubKeyFromHex(n.PubKey)
	if e != nil {
		return npk, e
	}
	n.Name = strings.TrimSpace(n.Name)
	if len(n.Name) < 3 {
		n.Name = misc.MakeRandomAlias()
	}
	return npk, nil
}

// TouchWithSeed updates the node's public key and returns the key pair
func (n *Node) TouchWithSeed(seed []byte) (cipher.PubKey, cipher.SecKey) {
	npk, nsk := cipher.GenerateDeterministicKeyPair(seed)
	n.PubKey = npk.Hex()
	return npk, nsk
}

// Message represents a message as stored in cxo.
type Message struct {
	Body      string     `json:"body"`
	Origin    string     `json:"origin"`
	Created   int64      `json:"created"`
	Signature cipher.Sig `json:"-"`
	Ref       string     `json:"hash" enc:"-"`
}

func (m *Message) checkContent() error {
	body := strings.TrimSpace(m.Body)
	if len(body) < 2 {
		return errors.New("message content too short")
	}
	return nil
}

func (m *Message) checkOrigin() (cipher.PubKey, error) {
	if m.Origin == (cipher.PubKey{}.Hex()) {
		return cipher.PubKey{}, errors.New("empty origin public key")
	}
	return misc.GetPubKey(m.Origin)
}

// Sign checks and signs the message.
func (m *Message) Sign(pk cipher.PubKey, sk cipher.SecKey) error {
	if e := m.checkContent(); e != nil {
		return e
	}
	m.Origin = pk.Hex()
	m.Created = 0
	m.Signature = cipher.Sig{}
	m.Signature = cipher.SignHash(cipher.SumSHA256(encoder.Serialize(*m)), sk)
	return nil
}

// Verify checks the legitimacy of the message.
func (m Message) Verify() error {
	// Check title and body.
	if e := m.checkContent(); e != nil {
		return e
	}
	// Check origin.
	originPK, e := m.checkOrigin()
	if e != nil {
		return e
	}
	// Check signature.
	sig := m.Signature
	m.Signature = cipher.Sig{}
	m.Created = 0

	return cipher.VerifySignature(
		originPK, sig,
		cipher.SumSHA256(encoder.Serialize(m)))
}

// Touch updates the timestamp of Message.
func (m *Message) Touch() {
	m.Created = time.Now().UnixNano()
}

func (m *Message) Deserialize(data []byte) error {
	return encoder.DeserializeRaw(data, m)
}
