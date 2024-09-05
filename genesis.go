package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	tmjson "github.com/tendermint/tendermint/libs/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func writeGenesis(genesisFile string, airdrop airdrop) error {
	bz, err := os.ReadFile(genesisFile)
	if err != nil {
		return fmt.Errorf("readfile %s: %w", genesisFile, err)
	}
	var genesisState map[string]json.RawMessage
	if err := tmjson.Unmarshal(bz, &genesisState); err != nil {
		return fmt.Errorf("tmjson.Unmarshal: %w", err)
	}
	// bz, err = tmjson.MarshalIndent(genesisState, "", "  ")
	// if err != nil {
	// return err
	// }
	// fmt.Println(string(bz))
	// return nil
	var appState map[string]json.RawMessage
	if err := tmjson.Unmarshal(genesisState["app_state"], &appState); err != nil {
		return fmt.Errorf("tmjson.Unmarshal appstate: %w", err)
	}

	//-----------------------------------------
	// Update bank genesis
	const ticker = "atone"
	var (
		balances    []banktypes.Balance
		totalSupply sdk.Coins
	)
	for addr, amt := range airdrop.addresses {
		coins := sdk.NewCoins(sdk.NewCoin("u"+ticker, amt))
		balances = append(balances, banktypes.Balance{
			Address: addr,
			Coins:   coins,
		})
		totalSupply = totalSupply.Add(coins...)
	}
	bankGen := banktypes.GenesisState{
		Supply: totalSupply,
		Params: banktypes.Params{
			DefaultSendEnabled: true,
			SendEnabled:        []*banktypes.SendEnabled{},
		},
		DenomMetadata: []banktypes.Metadata{
			{
				Display:     ticker,
				Symbol:      strings.ToUpper(ticker),
				Base:        "u" + ticker,
				Name:        "AtomOne Atone",
				Description: "The native staking token of AtomOne Hub",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Aliases:  []string{"micro" + ticker},
						Denom:    "u" + ticker,
						Exponent: 0,
					},
					{
						Aliases:  []string{"milli" + ticker},
						Denom:    "m" + ticker,
						Exponent: 3,
					},
					{
						Aliases:  []string{ticker},
						Denom:    ticker,
						Exponent: 6,
					},
				},
			},
		},
		Balances: balances,
	}

	//-----------------------------------------
	// Update the whole genesis
	appState["bank"], err = tmjson.Marshal(bankGen)
	if err != nil {
		return err
	}
	genesisState["app_state"], err = tmjson.Marshal(appState)
	if err != nil {
		return err
	}
	bz, err = tmjson.MarshalIndent(genesisState, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(bz))
	return nil
}
