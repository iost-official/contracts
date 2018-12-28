package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/iost-official/go-iost/ilog"
	"github.com/iost-official/go-iost/verifier"
)

func testWeb() (*verifier.Simulator, string) {
	ilog.SetLevel(ilog.LevelInfo)
	ilog.Info("start!!")
	s := verifier.NewSimulator()
	s.SetAccount(admin)
	s.SetGas(admin.ID, 1000000)
	s.SetRAM(admin.ID, 30000)

	c, err := s.Compile("web.demo", "./demos/web_on_blockchain", "./demos/web_on_blockchain")
	if err != nil {
		panic(err)
	}

	cname, _, err := s.DeployContract(c, admin.ID, adminKey)
	if err != nil {
		panic(err)
	}
	ilog.Info(cname)
	ilog.Info(admin.ID)

	web, err := ioutil.ReadFile("./demos/web_test.html")
	if err != nil {
		panic(err)
	}

	args, err := json.Marshal([]string{string(web)})
	if err != nil {
		panic(err)
	}

	r, err := s.Call(cname, "deploy", string(args), admin.ID, adminKey)
	ilog.Info(r)

	ilog.Info("cname is ", cname)
	ilog.Info("owner is ", admin.ID)

	return s, cname
}
