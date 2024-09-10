package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	tmjson "github.com/cometbft/cometbft/libs/json"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func writeGenesis(genesisFile string, airdrop airdrop) error {
	bz, err := os.ReadFile(genesisFile)
	if err != nil {
		return fmt.Errorf("readfile %s: %w", genesisFile, err)
	}
	var genesisState map[string]json.RawMessage
	if err := tmjson.Unmarshal(bz, &genesisState); err != nil {
		return fmt.Errorf("unmarshal genesis: %w", err)
	}
	var appState map[string]json.RawMessage
	if err := tmjson.Unmarshal(genesisState["app_state"], &appState); err != nil {
		return fmt.Errorf("unmarshal appstate: %w", err)
	}
	var authGen authtypes.GenesisState
	if err := tmjson.Unmarshal(appState["auth"], &authGen); err != nil {
		return fmt.Errorf("umarshal auth genesis: %w", err)
	}
	var bankGen banktypes.GenesisState
	if err := tmjson.Unmarshal(appState["bank"], &bankGen); err != nil {
		return fmt.Errorf("umarshal bank genesis: %w", err)
	}

	// Reset supply, balances and accounts
	bankGen.Supply = sdk.NewCoins()
	bankGen.Balances = nil
	authGen.Accounts = nil
	// Add airdrop.addresses to balances and accounts
	const ticker = "atone"
	for _, addr := range slices.Sorted(maps.Keys(airdrop.addresses)) {
		// update bank genesis
		amt := airdrop.addresses[addr]
		coins := sdk.NewCoins(sdk.NewCoin("u"+ticker, amt))
		bankGen.Balances = append(bankGen.Balances, banktypes.Balance{
			Address: addr,
			Coins:   coins,
		})
		bankGen.Supply = bankGen.Supply.Add(coins...)

		// update auth genesis
		acc := &authtypes.BaseAccount{Address: addr}
		any, err := codectypes.NewAnyWithValue(acc)
		if err != nil {
			return fmt.Errorf("newAny from base account: %w", err)
		}
		authGen.Accounts = append(authGen.Accounts, any)
	}
	// Add reserved address
	// hex:    0x0000000000000000000000000000000000000da0
	// bech32: atone1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqrdqzf7whr
	reservedAddrBz := []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0d\xa0")
	reservedAddrCoins := sdk.NewCoins(sdk.NewCoin("u"+ticker, airdrop.reservedAddr.RoundInt()))
	reservedAddr := sdk.MustBech32ifyAddressBytes("atone", reservedAddrBz)
	bankGen.Balances = append(bankGen.Balances, banktypes.Balance{
		Address: reservedAddr,
		Coins:   reservedAddrCoins,
	})
	bankGen.Supply = bankGen.Supply.Add(reservedAddrCoins...)
	// FIXME: add reserved address as a base account or not ?
	// maybe it's useless due to the nature of this reserved address: being just a way to hold the tokens with no account.
	// reservedAcc := &authtypes.BaseAccount{Address: reservedAddr}
	// any, err := codectypes.NewAnyWithValue(reservedAcc)
	// if err != nil {
	// return fmt.Errorf("newAny from base account: %w", err)
	// }
	// authGen.Accounts = append(authGen.Accounts, any)

	bankGen.Params = banktypes.Params{
		DefaultSendEnabled: true,
		SendEnabled:        []*banktypes.SendEnabled{},
	}
	bankGen.DenomMetadata = []banktypes.Metadata{
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
	}

	//-----------------------------------------
	// Update the  genesis
	appState["bank"], err = tmjson.Marshal(bankGen)
	if err != nil {
		return err
	}
	// Must use `marshaler` here because tmjson throws an error because of the
	// Any types.
	var b bytes.Buffer
	err = marshaler.Marshal(&b, &authGen)
	if err != nil {
		return err
	}
	appState["auth"] = b.Bytes()
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
