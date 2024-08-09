package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "testchain/testutil/keeper"
	"testchain/x/testchain/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.TestchainKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
