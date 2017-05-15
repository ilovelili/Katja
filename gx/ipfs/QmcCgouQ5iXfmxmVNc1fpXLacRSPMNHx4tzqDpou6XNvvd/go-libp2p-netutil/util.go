package testutil

import (
	"context"
	"testing"

	pstore "gx/ipfs/QmNUVzEjq3XWJ89hegahPvyfJbTXgTaom48pLb7YBD9gHQ/go-libp2p-peerstore"
	inet "gx/ipfs/QmVHSBsn8LEeay8m5ERebgUVuhzw838PsyTttCmP6GMJkg/go-libp2p-net"
	ma "gx/ipfs/QmcyqRMCAXVtYPS4DiBrA7sezL9rRGfW8Ctx7cywL4TXJj/go-multiaddr"
	tu "gx/ipfs/QmdVnYKahrvndXhyreWhY8YT3a5chJoWv8b4wVyH9JG2KB/go-testutil"
	swarm "gx/ipfs/Qmeo7oJxR65PLPx68KPFi8rjzcEmmWN2dL66fPuq9nVMv8/go-libp2p-swarm"
)

func GenSwarmNetwork(t *testing.T, ctx context.Context) *swarm.Network {
	p := tu.RandPeerNetParamsOrFatal(t)
	ps := pstore.NewPeerstore()
	ps.AddPubKey(p.ID, p.PubKey)
	ps.AddPrivKey(p.ID, p.PrivKey)
	n, err := swarm.NewNetwork(ctx, []ma.Multiaddr{p.Addr}, p.ID, ps, nil)
	if err != nil {
		t.Fatal(err)
	}
	ps.AddAddrs(p.ID, n.ListenAddresses(), pstore.PermanentAddrTTL)
	return n
}

func DivulgeAddresses(a, b inet.Network) {
	id := a.LocalPeer()
	addrs := a.Peerstore().Addrs(id)
	b.Peerstore().AddAddrs(id, addrs, pstore.PermanentAddrTTL)
}

var RandPeerID = tu.RandPeerID
