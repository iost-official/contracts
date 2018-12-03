package main

import (
	"testing"

	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/ilog"
	"github.com/iost-official/go-iost/verifier"
)

func TestGobang(t *testing.T) {
	ilog.SetLevel(ilog.LevelInfo)
	s := verifier.NewSimulator()
	defer s.Clear()
	s.SetAccount(admin)
	s.SetGas(admin.ID, 100000000)
	s.SetRAM(admin.ID, 30000)

	c, err := s.Compile("gobang.demo", "./demos/gobang", "./demos/gobang")
	if err != nil {
		panic(err)
	}

	cname, _, err := s.DeployContract(c, admin.ID, adminKey)
	if err != nil {
		panic(err)
	}

	r, err := s.Call(cname, "newGameWith", makeArgs([]interface{}{"playerb"}), admin.ID, adminKey)
	ilog.Info(r)

	txhash := common.Base58Encode(r.TxHash)
	ilog.Info(txhash)
	r, err = s.Call(cname, "move", makeArgs([]interface{}{0, 1, 1, txhash}), admin.ID, adminKey)
	ilog.Info(r)

	ilog.Info("cname is ", cname)
	ilog.Info("owner is ", admin.ID)
}
