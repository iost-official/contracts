package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/iost-official/go-iost/account"
	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/crypto"
	"github.com/iost-official/go-iost/ilog"
	"github.com/iost-official/go-iost/verifier"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

var (
	admin    *account.Account
	adminKey *account.KeyPair

	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

func init() {
    var err error
    adminKey, err = account.NewKeyPair(common.Base58Decode("EhNiaU4DzUmjCrvynV3gaUeuj2VjB1v2DCmbGD5U2nSE"), crypto.Secp256k1)
	if err != nil {
		panic(err)
	}
	admin = account.NewInitAccount("admin", adminKey.ID, adminKey.ID)
}

func main() {
    ilog.SetLevel(ilog.LevelInfo)
	ilog.Info("start!!")
    s := verifier.NewSimulator()
	defer s.Clear()
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

	web2 := s.Visitor.Get(cname + "@" + admin.ID + "-w0")
	ilog.Info(web2)

	router := fasthttprouter.New()
	router.POST("/getContractStorage/", func(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
		res := make(map[string]string)
		res["jsonStr"] = web2[1:]

		ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

		ctx.Response.SetStatusCode(200)
		json.NewEncoder(ctx.Response.BodyWriter()).Encode(res)
	})
	if err := fasthttp.ListenAndServe("0.0.0.0:20001", router.Handler); err != nil {
		fmt.Println("start fasthttp fail:", err.Error())
	}
}
