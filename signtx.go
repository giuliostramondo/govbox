package main

import (
	"encoding/json"
	"fmt"
	"os"

	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
)

func signTx(unsignedTxFile string) error {
	f, err := os.Open(unsignedTxFile)
	if err != nil {
		return err
	}
	defer f.Close()
	var tx txtypes.Tx
	// NOTE: Unsurprisingly doesn't work with legacy encoding/json
	if err := json.NewDecoder(f).Decode(&tx); err != nil {
		return err
	}
	fmt.Println(tx)
	return nil
}
