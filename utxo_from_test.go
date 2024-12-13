package gobtcsign

import (
	"testing"
)

func TestUtxoFromClient_GetUtxoFrom(t *testing.T) {
	var _ GetUtxoFromInterface = &SenderAmountUtxoClient{}
}

func TestSenderAmountUtxoCache_GetUtxoFrom(t *testing.T) {
	var _ GetUtxoFromInterface = &SenderAmountUtxoCache{}
}
