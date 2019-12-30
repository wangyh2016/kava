package types

// GenesisState is the state that must be provided at genesis.
type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

// DefaultGenesisState returns a default genesis state
// TODO pick better values
func DefaultGenesisState() GenesisState {
	return GenesisState{
		DefaultParams(),
	}
}

// ValidateGenesis performs basic validation of genesis gs returning an error for any failed validation criteria.
func ValidateGenesis(gs GenesisState) error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	return nil
}
