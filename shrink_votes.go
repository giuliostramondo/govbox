package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	govtypes "github.com/atomone-hub/atomone/x/gov/types/v1"

	tmjson "github.com/cometbft/cometbft/libs/json"
	tmtypes "github.com/cometbft/cometbft/types"
)

func shrinkVotes(_ context.Context, genesisFile string, high int) error {
	// Read input genesis and update it
	bz, err := os.ReadFile(genesisFile)
	if err != nil {
		return fmt.Errorf("readfile %s: %w", genesisFile, err)
	}
	var genesisState tmtypes.GenesisDoc
	if err := tmjson.Unmarshal(bz, &genesisState); err != nil {
		return fmt.Errorf("unmarshal genesis doc: %w", err)
	}
	var appState map[string]json.RawMessage
	if err := json.Unmarshal(genesisState.AppState, &appState); err != nil {
		return fmt.Errorf("unmarshal appstate: %w", err)
	}
	var govGen govtypes.GenesisState
	if err := cdc.UnmarshalJSON(appState["gov"], &govGen); err != nil {
		return fmt.Errorf("umarshal gov genesis: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Shrinking votes from %d to %d\n", len(govGen.Votes), high)
	govGen.Votes = govGen.Votes[:high]
	appState["gov"], err = cdc.MarshalJSON(&govGen)
	if err != nil {
		return fmt.Errorf("marshal gov genesis: %w", err)
	}
	genesisState.AppState, err = json.MarshalIndent(appState, "", "  ")
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
