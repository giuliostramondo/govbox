package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"slices"
	"strings"

	tmjson "github.com/cometbft/cometbft/libs/json"
	tmtypes "github.com/cometbft/cometbft/types"

	govtypes "github.com/atomone-hub/atomone/x/gov/types/v1"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

const constitutionLink = "https://raw.githubusercontent.com/atomone-hub/genesis/af652e0bc2bf1579350648770bf1f7b2d51d4884/CONSTITUTION.md"

func writeGenesis(genesisFile string, airdrop airdrop) error {
	bz, err := os.ReadFile(genesisFile)
	if err != nil {
		return fmt.Errorf("readfile %s: %w", genesisFile, err)
	}
	var genesisState tmtypes.GenesisDoc
	if err := tmjson.Unmarshal(bz, &genesisState); err != nil {
		return fmt.Errorf("unmarshal genesis: %w", err)
	}
	var appState map[string]json.RawMessage
	if err := tmjson.Unmarshal(genesisState.AppState, &appState); err != nil {
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
	var distrGen distrtypes.GenesisState
	// FIXME check how different modules unmarshal their genesis, might be better
	// using cdc than tmson
	if err := tmjson.Unmarshal(appState["distribution"], &distrGen); err != nil {
		return fmt.Errorf("umarshal distribution genesis: %w", err)
	}
	var govGen govtypes.GenesisState
	if err := cdc.UnmarshalJSON(appState["gov"], &govGen); err != nil {
		return fmt.Errorf("umarshal gov genesis: %w", err)
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
	// add auth reserved address
	reservedAcc := &authtypes.BaseAccount{Address: reservedAddr}
	any, err := codectypes.NewAnyWithValue(reservedAcc)
	if err != nil {
		return fmt.Errorf("newAny from base account: %w", err)
	}
	authGen.Accounts = append(authGen.Accounts, any)

	// setup community pool
	communityPoolCoins := sdk.NewCoins(sdk.NewCoin("u"+ticker, airdrop.communityPool.RoundInt()))
	distrGen.FeePool = distrtypes.FeePool{
		CommunityPool: sdk.NewDecCoinsFromCoins(communityPoolCoins...),
	}
	// same amount must be distributed to the distribution module account
	distrModuleAddr := sdk.MustBech32ifyAddressBytes("atone", authtypes.NewModuleAddress(distrtypes.ModuleName))
	bankGen.Balances = append(bankGen.Balances, banktypes.Balance{
		Address: distrModuleAddr,
		Coins:   communityPoolCoins,
	})

	// setup bank params and denoms
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

	// Update constitution
	resp, err := http.Get(constitutionLink)
	if err != nil {
		return err
	}
	bz, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	govGen.Constitution = string(bz)

	//-----------------------------------------
	// Update the  genesis
	appState["bank"], err = tmjson.Marshal(bankGen)
	if err != nil {
		return err
	}
	appState["distribution"], err = tmjson.Marshal(distrGen)
	if err != nil {
		return err
	}
	appState["gov"], err = tmjson.Marshal(govGen)
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
	genesisState.AppState, err = json.MarshalIndent(appState, "", " ")
	if err != nil {
		return err
	}
	bz, err = tmjson.MarshalIndent(genesisState, "", " ")
	if err != nil {
		return err
	}
	fmt.Println(string(bz))
	return nil
}
