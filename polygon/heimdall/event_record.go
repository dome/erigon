package heimdall

import (
	"bytes"
	"fmt"
	"math/big"
	"time"

	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/hexutility"
	"github.com/ledgerwatch/erigon/accounts/abi"
	"github.com/ledgerwatch/erigon/rlp"
)

// EventRecord represents state record
type EventRecord struct {
	ID       uint64            `json:"id" yaml:"id"`
	Contract libcommon.Address `json:"contract" yaml:"contract"`
	Data     hexutility.Bytes  `json:"data" yaml:"data"`
	TxHash   libcommon.Hash    `json:"tx_hash" yaml:"tx_hash"`
	LogIndex uint64            `json:"log_index" yaml:"log_index"`
	ChainID  string            `json:"bor_chain_id" yaml:"bor_chain_id"`
}

type EventRecordWithTime struct {
	EventRecord
	Time time.Time `json:"record_time" yaml:"record_time"`
}

// String returns the string representatin of a state record
func (e *EventRecordWithTime) String() string {
	return fmt.Sprintf(
		"id %v, contract %v, data: %v, txHash: %v, logIndex: %v, chainId: %v, time %s",
		e.ID,
		e.Contract.String(),
		e.Data.String(),
		e.TxHash.Hex(),
		e.LogIndex,
		e.ChainID,
		e.Time.Format(time.RFC3339),
	)
}

func (e *EventRecordWithTime) BuildEventRecord() *EventRecord {
	return &EventRecord{
		ID:       e.ID,
		Contract: e.Contract,
		Data:     e.Data,
		TxHash:   e.TxHash,
		LogIndex: e.LogIndex,
		ChainID:  e.ChainID,
	}
}

func UnpackEventRecordWithTime(stateContract abi.ABI, encodedEvent rlp.RawValue) *EventRecordWithTime {
	commitStateInputs := stateContract.Methods["commitState"].Inputs
	methodId := stateContract.Methods["commitState"].ID

	if bytes.Equal(methodId, encodedEvent[0:4]) {
		t := time.Unix((&big.Int{}).SetBytes(encodedEvent[4:36]).Int64(), 0)
		args, _ := commitStateInputs.Unpack(encodedEvent[4:])

		if len(args) == 2 {
			var eventRecord EventRecord
			if err := rlp.DecodeBytes(args[1].([]byte), &eventRecord); err == nil {
				return &EventRecordWithTime{EventRecord: eventRecord, Time: t}
			}
		}
	}

	return nil
}

type StateSyncEventsResponse struct {
	Height string                 `json:"height"`
	Result []*EventRecordWithTime `json:"result"`
}
