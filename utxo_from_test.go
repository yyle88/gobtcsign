package gobtcsign

import (
	"testing"
)

func TestUtxoFromClient_GetUtxoFrom(t *testing.T) {
	var _ GetUtxoFromInterface = &UtxoFromClient{}
}

func TestUtxoFromOutMap_GetUtxoFrom(t *testing.T) {
	var _ GetUtxoFromInterface = &OutPointUtxoSenderAmountMap{}
}
