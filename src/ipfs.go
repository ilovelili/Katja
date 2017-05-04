package katja

import (
	"context"
	"fmt"
	cid "gx/ipfs/QmYhQaCYEcaPPjxJX7YcPcVKkQfRy6sJ7B3XmGFk82XYdQ/go-cid"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

// StartNode Start Node
func StartNode(repo string) (*core.IpfsNode, error) {
	// Assume the user has run 'ipfs init'
	r, err := fsrepo.Open(repo) // "~/.ipfs" in linux and "C:\Users\username\.ipfs" in windows
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &core.BuildCfg{
		Repo:   r,
		Online: true,
	}

	return core.NewNode(ctx, cfg)
}

// GetStrings get strings by userID
func GetStrings(node *core.IpfsNode, userID string) ([]string, error) {
	// base58: https://kikobeats.com/base58-is-base64-for-humans/
	cid, err := cid.Decode(userID)
	if err != nil {
		return nil, err
	}

	return resolveAllInOrder(node, cid), nil
}

func resolveAllInOrder(nd *core.IpfsNode, cid *cid.Cid) (stringArr []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	node, err := nd.DAG.Get(ctx, cid)
	fmt.Println("the node is", node)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("bout to crash")
	fmt.Printf("%s", string(node.RawData()[:]))
	fmt.Println("not crashed ")

	for {
		var err error
		if len(node.Links()) == 0 {
			break
		}

		node, err = node.Links()[0].GetNode(ctx, nd.DAG)
		if err != nil {
			fmt.Println(err)
		}

		data := string(node.RawData()[:])
		fmt.Println(data)
		stringArr = append(stringArr, data)
	}

	return stringArr
}
