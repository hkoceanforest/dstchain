package types

func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

func (gs GenesisState) Validate() error {

	return nil
}
