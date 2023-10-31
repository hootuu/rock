package stone

import (
	"context"
	"github.com/hootuu/tome/bk"
	"github.com/hootuu/utils/errors"
	"github.com/hootuu/utils/logger"
	"github.com/hootuu/utils/sys"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/ipfs/kubo/core"
	"github.com/multiformats/go-multihash"
	"go.uber.org/zap"
)

var gIpfsNode *core.IpfsNode

func Initialize(ipfsNode *core.IpfsNode) {
	gIpfsNode = ipfsNode
}

func mustGetIpfsNode() *core.IpfsNode {
	if gIpfsNode == nil {
		sys.Error("must init gIpfsNode first.")
	}
	return gIpfsNode
}

func RegisterDataType(tpl interface{}) {
	cbor.RegisterCborType(tpl)
}

func Inscribe(data bk.Invariable, pin bool) (bk.BID, uint64, *errors.Error) {
	node, nErr := cbor.WrapObject(data, multihash.SHA2_256, -1)
	if nErr != nil {
		return bk.NilBID, 0, errors.Sys("wrap cbor node failed: " + nErr.Error())
	}
	nErr = mustGetIpfsNode().DAG.Add(context.Background(), node)
	if nErr != nil {
		return bk.NilBID, 0, errors.Sys("inscribe data failed: " + nErr.Error())
	}
	if pin {
		nErr = mustGetIpfsNode().Pinning.Pin(context.Background(), node, true)
		if nErr != nil {
			return bk.NilBID, 0, errors.Sys("pin data failed: " + nErr.Error())
		}
	}
	size, nErr := node.Size()
	if nErr != nil {
		return bk.NilBID, 0, errors.Sys("calc node.size failed: " + nErr.Error())
	}
	return bk.BID(node.Cid().String()), size, nil
}

func Get(bid bk.BID, v interface{}) (uint64, *errors.Error) {
	cid, nErr := cid.Decode(bid.S())
	if nErr != nil {
		return 0, errors.Verify("invalid BID")
	}
	node, nErr := mustGetIpfsNode().DAG.Get(context.Background(), cid)
	if nErr != nil {
		logger.Logger.Error("ipfs.DAG.Get failed", zap.String("bid", bid.S()),
			zap.Error(nErr))
		return 0, errors.Sys("load data failed: " + nErr.Error())
	}
	nErr = cbor.DecodeInto(node.RawData(), v)
	if nErr != nil {
		logger.Logger.Error("cbor.DecodeInto(node.RawData(), v) failed", zap.String("bid", bid.S()),
			zap.Error(nErr))
		return 0, errors.Sys("decode data failed: " + nErr.Error())
	}
	size, nErr := node.Size()
	if nErr != nil {
		return 0, errors.Sys("calc node.size failed: " + nErr.Error())
	}
	return size, nil
}
