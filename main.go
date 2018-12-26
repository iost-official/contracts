package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net"

	"github.com/iost-official/go-iost/account"
	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/core/tx"
	"github.com/iost-official/go-iost/crypto"
	"github.com/iost-official/go-iost/ilog"
	"github.com/iost-official/go-iost/verifier"
	"github.com/iost-official/go-iost/vm/native"
	"github.com/iost-official/playground/rpc"
	"github.com/iost-official/playground/rpc/pb"
	"golang.org/x/net/netutil"
	"google.golang.org/grpc"
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

func createToken(s *verifier.Simulator) error {
	// create token
	r, err := s.Call("token.iost", "create", fmt.Sprintf(`["%v", "%v", %v, {}]`, "iost", admin.ID, 1000000), admin.ID, adminKey)
	if err != nil || r.Status.Code != tx.Success {
		return fmt.Errorf("err %v, receipt: %v", err, r)
	}
	// issue token
	r, err = s.Call("token.iost", "issue", fmt.Sprintf(`["%v", "%v", "%v"]`, "iost", admin.ID, "1000"), admin.ID, adminKey)
	if err != nil || r.Status.Code != tx.Success {
		return fmt.Errorf("err %v, receipt: %v", err, r)
	}
	if 1e11 != s.Visitor.TokenBalance("iost", admin.ID) {
		return fmt.Errorf("err %v, receipt: %v", err, r)
	}
	s.Visitor.Commit()
	return nil
}

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

func main() {

	s, cname := testDice()
	fmt.Println("contract id is:", cname)

	grpcServer := grpc.NewServer()
	apiService := rpc.NewServer(s)
	rpcpb.RegisterApiServiceServer(grpcServer, apiService)

	lis, err := net.Listen("tcp", "0.0.0.0:30002")
	if err != nil {
		panic(err)
	}
	lis = netutil.LimitListener(lis, 10000)
	if err := grpcServer.Serve(lis); err != nil {
		ilog.Fatalf("start grpc failed. err=%v", err)
	}

}
