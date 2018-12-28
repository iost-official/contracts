package main

import (
	"fmt"

	"net"

	"github.com/iost-official/go-iost/account"
	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/crypto"
	"github.com/iost-official/go-iost/ilog"
	"github.com/iost-official/playground/rpc"
	"github.com/iost-official/playground/rpc/pb"
	"golang.org/x/net/netutil"
	"google.golang.org/grpc"
)

var (
	admin    *account.Account
	adminKey *account.KeyPair
	playerb  *account.Account
)

func init() {
	var err error
	adminKey, err = account.NewKeyPair(common.Base58Decode("EhNiaU4DzUmjCrvynV3gaUeuj2VjB1v2DCmbGD5U2nSE"), crypto.Secp256k1)
	if err != nil {
		panic(err)
	}
	admin = account.NewInitAccount("admin", adminKey.ID, adminKey.ID)
	playerb = account.NewInitAccount("playerb", adminKey.ID, adminKey.ID)
}

func main() {

	s, cname := testGobang()
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
