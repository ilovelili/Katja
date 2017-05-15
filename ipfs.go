package katja

import (
	"context"
	"fmt"
	cid "gx/ipfs/QmYhQaCYEcaPPjxJX7YcPcVKkQfRy6sJ7B3XmGFk82XYdQ/go-cid"

	"runtime"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/merkledag"
	"github.com/ipfs/go-ipfs/path"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

func defaultPath() string {
	if runtime.GOOS == "windows" {
		return "C:\\Users\\username\\.ipfs"
	}
	return "~/.ipfs"
}

// StartNode Start IPFS Node
func StartNode() (*core.IpfsNode, error) {
	// Assume the user has run 'ipfs init'
	repo := defaultPath()
	r, err := fsrepo.Open(repo)
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

// GetStrings get strings by cid
func GetStrings(node *core.IpfsNode, cid *cid.Cid) (stringArr []string, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	nodeGetter := node.DAG
	defer cancel()
	// merkledag proto Node
	nd, err := nodeGetter.Get(ctx, cid)
	fmt.Println("the node is", nd)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("bout to crash")
	fmt.Printf("%s", nd.String())
	fmt.Println("not crashed ")

	for {
		var err error
		if len(nd.Links()) == 0 {
			break
		}

		nd, err = nd.Links()[0].GetNode(ctx, nodeGetter)
		if err != nil {
			fmt.Println(err)
			break
		}

		data := nd.String()
		fmt.Println(data)
		stringArr = append(stringArr, data)
	}

	return stringArr, nil
}

// AddString add input string to ipfs node
func AddString(node *core.IpfsNode, inputString string) (*cid.Cid, error) {
	pointsTo, err := node.Namesys.Resolve(node.Context(), node.Identity.Pretty())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//If there is an error, user is new and hasn't yet created a DAG.
	if err != nil {
		newProtoNode := makeStringNode(inputString)
		cid, err := node.DAG.Add(newProtoNode)
		if err != nil {
			return nil, err
		}

		err = node.Namesys.Publish(ctx, node.PrivateKey, path.FromCid(cid))
		if err != nil {
			return nil, err
		}

		return cid, nil
	}

	// Else user has already creatd a DAG
	newProtoNode := makeStringNode(inputString)
	cid, err := cid.Decode(pointsTo.String())
	if err != nil {
		return nil, err
	}

	oldProtoNode, err := node.DAG.Get(ctx, cid)
	if err != nil {
		return nil, err
	}

	err = newProtoNode.AddNodeLink("next", oldProtoNode)
	if err != nil {
		return nil, err
	}

	node.DAG.Add(newProtoNode)
	err = node.Namesys.Publish(ctx, node.PrivateKey, pointsTo)
	if err != nil {
		return nil, err
	}

	return cid, nil
}

func makeStringNode(s string) *merkledag.ProtoNode {
	nd := new(merkledag.ProtoNode)
	nd.SetData([]byte(s))
	return nd
}
