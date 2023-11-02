package rock

import (
	"github.com/hootuu/rock/stone"
	"github.com/hootuu/tome/bk"
	"github.com/hootuu/utils/errors"
	"github.com/ipfs/kubo/core"
)

func Initialize(ipfsNode *core.IpfsNode, dataTypes []interface{}) *errors.Error {
	stone.Initialize(ipfsNode)
	if len(dataTypes) > 0 {
		for _, dt := range dataTypes {
			stone.RegisterDataType(dt)
		}
	}
	return nil
}

func Inscribe(data bk.Invariable, pin bool) (bk.BID, uint64, *errors.Error) {
	return stone.Inscribe(data, pin)
}

func Get(bid bk.BID, v interface{}) (uint64, *errors.Error) {
	return stone.Get(bid, v)
}

func BuildBID(data bk.Invariable) (bk.BID, *errors.Error) {
	return stone.BuildBID(data)
}
