package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func d(str string) sdk.Dec { return sdk.MustNewDecFromStr(str) }

func TestGetSharesPercentage(t *testing.T) {
	tests := []struct {
		args []sdk.Int
		want sdk.Dec
	}{
		{[]sdk.Int{sdk.NewInt(50), sdk.NewInt(500)}, d("0.1")},
		{[]sdk.Int{sdk.NewInt(500), sdk.NewInt(500)}, d("1.0")},
		{[]sdk.Int{sdk.NewInt(100), sdk.NewInt(300)}, d("0.33333333333333")},
		{[]sdk.Int{sdk.NewInt(200), sdk.NewInt(300)}, d("0.66666666666667")},
		{[]sdk.Int{sdk.NewInt(100000000), sdk.NewInt(500000000)}, d("0.2")},
		{[]sdk.Int{sdk.NewInt(10000000000), sdk.NewInt(30000000000)}, d("0.33333333333333")},
		{[]sdk.Int{sdk.NewInt(2000000000000), sdk.NewInt(3000000000000)}, d("0.66666666666667")},
	}
	for i, tc := range tests {
		res := CalculateSharesPercentage(tc.args[0], tc.args[1])
		require.Equal(t, tc.want, res, "unexpected result for test case %d, input: %v, got: %v", i, tc.args, res)
	}
}
