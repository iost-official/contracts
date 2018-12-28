package main

import (
	"encoding/json"
	"fmt"

	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/core/tx"
	"github.com/iost-official/go-iost/verifier"
	"github.com/iost-official/go-sdk/pb"
)

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
