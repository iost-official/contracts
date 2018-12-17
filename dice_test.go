package main

import (
	"testing"

	"github.com/iost-official/go-iost/common"
)

func TestDice(t *testing.T) {
	s, cname := testDice()
	hostSeed := "12345678"
	hostHash := common.Base58Encode(common.Sha3([]byte(hostSeed)))
	r, err := s.Call(cname, "bet", makeArgs([]interface{}{1, 50, "abc", hostHash, ""}), admin.ID, adminKey)
	t.Log(r, err)
	r, err = s.Call(cname, "reveal", makeArgs([]interface{}{0, hostSeed, false}), admin.ID, adminKey)
	t.Log(r, err)
}
