package keeper

import (
	"testchain/x/testchain/types"
)

var _ types.QueryServer = Keeper{}
