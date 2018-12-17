package main

import (
	"testing"

	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/ilog"
)

func TestGobang(t *testing.T) {
	s, cname := testGobang()

	r, err := s.Call(cname, "newGameWith", makeArgs([]interface{}{"playerb"}), admin.ID, adminKey)
	ilog.Info(r, err)
	txhash := common.Base58Encode(r.TxHash)
	ilog.Info(txhash)
	r, err = s.Call(cname, "move", makeArgs([]interface{}{0, 7, 7, txhash}), admin.ID, adminKey)
	ilog.Info(r)
	txhash = common.Base58Encode(r.TxHash)
	r, err = s.Call(cname, "move", makeArgs([]interface{}{0, 7, 8, txhash}), playerb.ID, adminKey)
	ilog.Info(r, err)
}
