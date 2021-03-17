package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Status struct {
	Address          string
	IsVesting        bool
	Balance          sdk.Coin
	Delegated        sdk.Coin
	BalanceLocked    sdk.Coin
	BalanceSpendable sdk.Coin
	BalanceVesting   sdk.Coin
	BalanceVested    sdk.Coin
	DelegatedFree    sdk.Coin
	DelegatedVesting sdk.Coin
}
