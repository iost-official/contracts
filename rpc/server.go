package rpc

import (
	"context"

	"errors"

	"encoding/json"
	"fmt"
	"reflect"

	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/verifier"
	"github.com/iost-official/go-iost/vm/host"
	"github.com/iost-official/playground/rpc/pb"
)

var errShouldNotUse = errors.New("you should not call this api in PLAYGROUND environment")

type Server struct {
	txs map[string]*rpcpb.Transaction
	s   *verifier.Simulator
}

func NewServer(s *verifier.Simulator) *Server {
	return &Server{
		s:   s,
		txs: make(map[string]*rpcpb.Transaction),
	}
}

func (s *Server) GetNodeInfo(context.Context, *rpcpb.EmptyRequest) (*rpcpb.NodeInfoResponse, error) {
	return &rpcpb.NodeInfoResponse{
		BuildTime: "as your wish",
		GitHash:   "anyway",
		Mode:      "playground",
		Network: &rpcpb.NetworkInfo{
			Id:        "playground",
			PeerCount: 1,
			PeerInfo: []*rpcpb.PeerInfo{
				{
					Id:   "playground",
					Addr: "you know",
				},
			},
		},
	}, nil
}
func (s *Server) GetChainInfo(context.Context, *rpcpb.EmptyRequest) (*rpcpb.ChainInfoResponse, error) {
	return &rpcpb.ChainInfoResponse{}, nil
}
func (s *Server) GetRAMInfo(context.Context, *rpcpb.EmptyRequest) (*rpcpb.RAMInfoResponse, error) {
	return &rpcpb.RAMInfoResponse{}, nil
}
func (s *Server) GetTxByHash(c context.Context, req *rpcpb.TxHashRequest) (*rpcpb.TransactionResponse, error) {
	return &rpcpb.TransactionResponse{
		Status:      1,
		Transaction: s.txs[req.Hash],
	}, nil
}
func (s *Server) GetTxReceiptByTxHash(c context.Context, req *rpcpb.TxHashRequest) (*rpcpb.TxReceipt, error) {
	return s.txs[req.Hash].TxReceipt, nil
}
func (s *Server) GetBlockByHash(context.Context, *rpcpb.GetBlockByHashRequest) (*rpcpb.BlockResponse, error) {
	return nil, errShouldNotUse
}
func (s *Server) GetBlockByNumber(context.Context, *rpcpb.GetBlockByNumberRequest) (*rpcpb.BlockResponse, error) {
	return nil, errShouldNotUse
}
func (s *Server) GetAccount(context.Context, *rpcpb.GetAccountRequest) (*rpcpb.Account, error) {
	return nil, errShouldNotUse
}
func (s *Server) GetTokenBalance(context.Context, *rpcpb.GetTokenBalanceRequest) (*rpcpb.GetTokenBalanceResponse, error) {
	return nil, nil
}
func (s *Server) GetGasRatio(context.Context, *rpcpb.EmptyRequest) (*rpcpb.GasRatioResponse, error) {
	return nil, nil
}
func (s *Server) GetContract(context.Context, *rpcpb.GetContractRequest) (*rpcpb.Contract, error) {
	return nil, nil
}
func (s *Server) GetContractStorage(c context.Context, req *rpcpb.GetContractStorageRequest) (*rpcpb.GetContractStorageResponse, error) {
	h := host.NewHost(host.NewContext(nil), s.s.Visitor, nil, nil)
	var value interface{}
	if req.GetField() == "" {
		value, _ = h.GlobalGet(req.GetId(), req.GetKey())
	} else {
		value, _ = h.GlobalMapGet(req.GetId(), req.GetKey(), req.GetField())
	}
	var data string
	if value != nil && reflect.TypeOf(value).Kind() == reflect.String {
		data = value.(string)
	} else {
		bytes, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal %v", value)
		}
		data = string(bytes)
	}
	return &rpcpb.GetContractStorageResponse{
		Data: data,
	}, nil
}
func (s *Server) SendTransaction(c context.Context, req *rpcpb.TransactionRequest) (*rpcpb.SendTransactionResponse, error) {
	t := toCoreTx(req)

	// todo check tx

	hash := common.Base58Encode(t.Hash())

	r, err := s.s.RunTx(t)
	if err != nil {
		return nil, err
	}
	s.txs[hash] = toPbTx(t, r)

	return &rpcpb.SendTransactionResponse{
		Hash: hash,
	}, nil
}
func (s *Server) ExecTransaction(context.Context, *rpcpb.TransactionRequest) (*rpcpb.TxReceipt, error) {
	return nil, errShouldNotUse
}
func (s *Server) Subscribe(*rpcpb.SubscribeRequest, rpcpb.ApiService_SubscribeServer) error {
	return errShouldNotUse
}
