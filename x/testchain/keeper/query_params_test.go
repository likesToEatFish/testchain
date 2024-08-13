package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "testchain/testutil/keeper"
	"testchain/x/testchain/keeper"
	"testchain/x/testchain/types"
)

func TestParamsQuery(t *testing.T) {
	k, ctx, _ := keepertest.TestchainKeeper(t)

	qs := keeper.NewQueryServerImpl(k)
	params := types.DefaultParams()
	require.NoError(t, k.Params.Set(ctx, params))

	response, err := qs.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
