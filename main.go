package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/jsonpb"
	"github.com/iost-official/go-iost/account"
	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/core/tx"
	"github.com/iost-official/go-iost/core/tx/pb"
	"github.com/iost-official/go-iost/crypto"
	"github.com/iost-official/go-iost/ilog"
	"github.com/iost-official/go-iost/rpc/pb"
	"github.com/iost-official/go-iost/verifier"
	"github.com/iost-official/go-iost/vm/database"
	"github.com/iost-official/playground/handlers"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

var (
	admin    *account.Account
	adminKey *account.KeyPair
	playerb  *account.Account

	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")

	receiptDB map[string]*rpcpb.TxReceipt
)

func init() {
	var err error
	adminKey, err = account.NewKeyPair(common.Base58Decode("EhNiaU4DzUmjCrvynV3gaUeuj2VjB1v2DCmbGD5U2nSE"), crypto.Secp256k1)
	if err != nil {
		panic(err)
	}
	admin = account.NewInitAccount("admin", adminKey.ID, adminKey.ID)
	playerb = account.NewInitAccount("playerb", adminKey.ID, adminKey.ID)

	receiptDB = make(map[string]*rpcpb.TxReceipt)
}

func makeArgs(s []interface{}) string {
	args, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(args)
}

func toPbTxReceipt(tr *tx.TxReceipt) *rpcpb.TxReceipt {
	if tr == nil {
		return nil
	}
	ret := &rpcpb.TxReceipt{
		TxHash:     common.Base58Encode(tr.TxHash),
		GasUsage:   float64(tr.GasUsage) / 100,
		RamUsage:   tr.RAMUsage,
		StatusCode: rpcpb.TxReceipt_StatusCode(tr.Status.Code),
		Message:    tr.Status.Message,
		Returns:    tr.Returns,
	}
	for _, r := range tr.Receipts {
		ret.Receipts = append(ret.Receipts, &rpcpb.TxReceipt_Receipt{
			FuncName: r.FuncName,
			Content:  r.Content,
		})
	}
	return ret
}

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

func testGobang() (*verifier.Simulator, string) {
	ilog.SetLevel(ilog.LevelInfo)
	s := verifier.NewSimulator()
	s.SetAccount(admin)
	s.SetAccount(playerb)
	s.SetGas(admin.ID, 100000000)
	s.SetRAM(admin.ID, 30000)
	s.SetGas(playerb.ID, 100000000)
	s.SetRAM(playerb.ID, 30000)

	c, err := s.Compile("gobang.demo", "./demos/gobang", "./demos/gobang")
	if err != nil {
		panic(err)
	}

	cname, _, err := s.DeployContract(c, admin.ID, adminKey)
	if err != nil {
		panic(err)
	}
	return s, cname
}

func main() {

	s, cname := testGobang()

	go func() {
		r, err := s.Call(cname, "newGameWith", makeArgs([]interface{}{"playerb"}), admin.ID, adminKey)
		ilog.Info(r, err)
		txhash := common.Base58Encode(r.TxHash)
		ilog.Info(txhash)
		r, err = s.Call(cname, "move", makeArgs([]interface{}{0, 7, 7, txhash}), admin.ID, adminKey)
		ilog.Info(r)
		txhash = common.Base58Encode(r.TxHash)
		r, err = s.Call(cname, "move", makeArgs([]interface{}{0, 7, 8, txhash}), playerb.ID, adminKey)
		ilog.Info(r, err)
	}()

	router := fasthttprouter.New()
	router.OPTIONS("/getContractStorage/", handlers.OptionsHandler)
	router.POST("/getContractStorage/", func(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
		body := ctx.PostBody()
		var req map[string]interface{}
		err := json.Unmarshal(body, &req)
		if err != nil {
			ilog.Error(err)
		}

		cid := cname

		res := make(map[string]string)

		ret := s.Visitor.Get(cid + "-" + req["key"].(string))
		r2, ok := database.MustUnmarshal(ret).(string)
		if ok {
			res["jsonStr"] = r2
		} else {
			res["jsonStr"] = ""
		}

		ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

		ctx.Response.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(res)
	})

	router.OPTIONS("/sendTx/", handlers.OptionsHandler)
	router.POST("/sendTx/", func(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
		body := ctx.PostBody()
		var req map[string]interface{}
		err := json.Unmarshal(body, &req)
		if err != nil {
			ilog.Error(err)
		}

		ilog.Info(req)

		var tpb txpb.Tx

		err = json.Unmarshal(body, &tpb)
		if err != nil {
			ilog.Error(err)
		}

		ilog.Info(tpb)

		cid := cname
		//cid := tpb.

		res := make(map[string]string)

		r, err := s.Call(cid, tpb.Actions[0].ActionName, tpb.Actions[0].Data, tpb.Publisher, adminKey)
		ilog.Info(r, err)

		if err != nil {
			res["hash"] = ""
		} else {
			res["hash"] = common.Base58Encode(r.TxHash)
			receiptDB[res["hash"]] = toPbTxReceipt(r)
		}

		ret := s.Visitor.Get(cid + "-games0")
		r2 := database.MustUnmarshal(ret)
		ilog.Error(r2)

		ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

		ctx.Response.SetStatusCode(200)
		json.NewEncoder(ctx).Encode(res)
	})

	router.OPTIONS("/getTxReceiptByTxHash/:hash", handlers.OptionsHandler)

	router.GET("/getTxReceiptByTxHash/:hash", func(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {

		hash := params.ByName("hash")

		ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

		ctx.Response.SetStatusCode(200)
		var ma = jsonpb.Marshaler{OrigName: true, EmitDefaults: true}
		ma.Marshal(ctx, receiptDB[hash])
	})

	if err := fasthttp.ListenAndServe("0.0.0.0:20001", router.Handler); err != nil {
		fmt.Println("start fasthttp fail:", err.Error())
	}
}
