package stone

import (
	"github.com/hootuu/tome/bk"
	"github.com/hootuu/utils/errors"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/multiformats/go-multihash"
)

func BuildBID(data bk.Invariable) (bk.BID, *errors.Error) {
	node, nErr := cbor.WrapObject(data, multihash.SHA2_256, -1)
	if nErr != nil {
		return bk.NilBID, errors.Sys("wrap cbor node failed: " + nErr.Error())
	}
	return bk.BID(node.Cid().String()), nil
}
