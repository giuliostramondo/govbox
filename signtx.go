package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/davecgh/go-spew/spew"

	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
)

func signTx(unsignedTxFile string) error {
	f, err := os.Open(unsignedTxFile)
	if err != nil {
		return err
	}
	defer f.Close()
	var rawTx struct {
		Body struct {
			Messages []map[string]any
		}
	}
	if err := json.NewDecoder(f).Decode(&rawTx); err != nil {
		return fmt.Errorf("JSON decode %s: %v", unsignedTxFile, err)
	}
	var types []string
	for _, msg := range rawTx.Body.Messages {
		findUnregisteredTypes(msg, &types)
	}
	fmt.Println("TODO: NEED TO REGISTER THESE TYPES", types)
	var tx txtypes.Tx
	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("seek %s: %v", unsignedTxFile, err)
	}
	if err := unmarshaler.Unmarshal(f, &tx); err != nil {
		return fmt.Errorf("unmarshal tx: %v", err)
	}
	spew.Dump(tx)
	return nil
}

func findUnregisteredTypes(m map[string]any, types *[]string) {
	for k, v := range m {
		if k == "@type" {
			typeURL := v.(string)
			if _, err := registry.Resolve(typeURL); err != nil {
				if !slices.Contains(*types, typeURL) {
					// type not registered, add to the list
					*types = append(*types, typeURL)
				}
			}
			continue
		}
		switch x := v.(type) {
		case []map[string]any:
			for _, m := range x {
				findUnregisteredTypes(m, types)
			}
		case map[string]any:
			findUnregisteredTypes(x, types)
		}
	}
}
