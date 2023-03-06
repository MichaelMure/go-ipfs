package ipns

import (
	"context"

	pin "github.com/ipfs/go-ipfs-pinner"
	nsys "github.com/ipfs/go-namesys"
	path "github.com/ipfs/go-path"
	ft "github.com/ipfs/go-unixfs"
	ci "github.com/libp2p/go-libp2p/core/crypto"

	"github.com/ipfs/kubo/core"
)

// InitializeKeyspace sets the ipns record for the given key to
// point to an empty directory.
func InitializeKeyspace(n *core.IpfsNode, key ci.PrivKey) error {
	ctx, cancel := context.WithCancel(n.Context())
	defer cancel()

	emptyDir := ft.EmptyDirNode()

	err := func() error {
		defer n.Blockstore.PinLock(ctx).Unlock(ctx)

		err := n.DAG.Add(ctx, emptyDir)
		if err != nil {
			return err
		}

		err = n.Pinning.Pin(ctx, emptyDir.Cid(), pin.Direct)
		if err != nil {
			return err
		}

		return n.Pinning.Flush(ctx)
	}()
	if err != nil {
		return err
	}

	pub := nsys.NewIpnsPublisher(n.Routing, n.Repo.Datastore())

	return pub.Publish(ctx, key, path.FromCid(emptyDir.Cid()))
}
