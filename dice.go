package main

import (
	"github.com/iost-official/go-iost/ilog"
	"github.com/iost-official/go-iost/verifier"
	"github.com/iost-official/go-iost/vm/native"
)

func testDice() (*verifier.Simulator, string) {
	ilog.SetLevel(ilog.LevelInfo)
	s := verifier.NewSimulator()
	s.SetAccount(admin)
	s.SetAccount(playerb)
	s.SetGas(admin.ID, 100000000)
	s.SetRAM(admin.ID, 30000)
	s.SetGas(playerb.ID, 100000000)
	s.SetRAM(playerb.ID, 30000)

	s.SetContract(native.TokenABI())
	s.Visitor.Commit()
	createToken(s)

	c, err := s.Compile("gobang.demo", "./demos/dice", "./demos/dice")
	if err != nil {
		panic(err)
	}

	cname, _, err := s.DeployContract(c, admin.ID, adminKey)
	if err != nil {
		panic(err)
	}
	return s, cname
}
