package types

import (
	"encoding/json"
	"testing"

	"github.com/ledgerwatch/erigon-lib/common"
	"github.com/stretchr/testify/require"
)

func TestKnownAccounts(t *testing.T) {
	t.Parallel()

	requestRaw := []byte(`{"0xadd1add1add1add1add1add1add1add1add1add1": "0x000000000000000000000000313aadca1750caadc7bcb26ff08175c95dcf8e38", "0xadd2add2add2add2add2add2add2add2add2add2": {"0x0000000000000000000000000000000000000000000000000000000000000aaa": "0x0000000000000000000000000000000000000000000000000000000000000bbb", "0x0000000000000000000000000000000000000000000000000000000000000ccc": "0x0000000000000000000000000000000000000000000000000000000000000ddd"}}`)

	accs := &KnownAccountStorageConditions{}

	err := json.Unmarshal(requestRaw, accs)
	require.NoError(t, err)

	expected := &KnownAccountStorageConditions{
		common.HexToAddress("0xadd1add1add1add1add1add1add1add1add1add1"): NewKnownAccountStorageConditionWithRootHash("0x000000000000000000000000313aadca1750caadc7bcb26ff08175c95dcf8e38"),
		common.HexToAddress("0xadd2add2add2add2add2add2add2add2add2add2"): NewKnownAccountStorageConditionWithSlotHashes(map[string]string{
			"0x0000000000000000000000000000000000000000000000000000000000000aaa": "0x0000000000000000000000000000000000000000000000000000000000000bbb",
			"0x0000000000000000000000000000000000000000000000000000000000000ccc": "0x0000000000000000000000000000000000000000000000000000000000000ddd",
		}),
	}

	require.Equal(t, expected, accs)
}
