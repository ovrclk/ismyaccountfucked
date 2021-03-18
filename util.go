package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func getStatus(cctx client.Context, denom string, address string) (Status, error) {

	acc, err := getAccount(cctx, address)
	if err != nil {
		return Status{}, err
	}

	vacc, ok := acc.(vestexported.VestingAccount)
	if !ok {
		return Status{
			Address:   address,
			IsVesting: false,
		}, nil
	}

	balances, err := getBalances(cctx, address)
	if err != nil {
		return Status{}, err
	}

	delegations, err := getDelegations(cctx, address)
	if err != nil {
		return Status{}, err
	}
	now := time.Now()

	locked := vacc.LockedCoins(now)

	status := Status{
		Address:          address,
		IsVesting:        true,
		Balance:          getDenom(balances, denom),
		BalanceLocked:    getDenom(locked, denom),
		BalanceSpendable: getDenom(spendableCoins(locked, balances), denom),
		Delegated:        getDenom(delegations, denom),
		BalanceVesting:   getDenom(vacc.GetVestingCoins(now), denom),
		BalanceVested:    getDenom(vacc.GetVestedCoins(now), denom),
		DelegatedFree:    getDenom(vacc.GetDelegatedFree(), denom),
		DelegatedVesting: getDenom(vacc.GetDelegatedVesting(), denom),
	}

	status.ExpectedDelegatedFree = getExpectedDelegatedFree(
		status.Delegated, status.BalanceVesting, status.BalanceVested)
	status.ExpectedDelegatedVesting = getExpectedDelegatedVesting(
		status.Delegated, status.BalanceVesting)

	return status, nil
}

func getExpectedDelegatedVesting(delegated, vesting sdk.Coin) sdk.Coin {
	if delegated.IsLT(vesting) {
		return delegated
	}
	return vesting
}

func getExpectedDelegatedFree(delegated, vesting, vested sdk.Coin) sdk.Coin {
	notvesting := delegated.Amount.Sub(getExpectedDelegatedVesting(delegated, vesting).Amount)
	return sdk.NewCoin(delegated.Denom, notvesting)
}

func getDenom(coins sdk.Coins, denom string) sdk.Coin {
	return sdk.NewCoin(denom, coins.AmountOf(denom))
}

func formatAmount(coin sdk.Coin) string {
	denom := strings.ToUpper(coin.Denom[1:])

	whole := coin.Amount.QuoRaw(1000000).Uint64()
	frac := coin.Amount.ModRaw(1000000).QuoRaw(10000).Uint64()

	return fmt.Sprintf("%d.%02d %s", whole, frac, denom)
}

func getAccount(cctx client.Context, address string) (authtypes.AccountI, error) {
	authclient := authtypes.NewQueryClient(cctx)
	res, err := authclient.Account(context.Background(), &authtypes.QueryAccountRequest{Address: address})
	if err != nil {
		return nil, err
	}

	var account authtypes.AccountI
	if err := cctx.InterfaceRegistry.UnpackAny(res.Account, &account); err != nil {
		return nil, err
	}
	return account, nil
}

func getBalances(cctx client.Context, address string) (sdk.Coins, error) {
	bankclient := banktypes.NewQueryClient(cctx)
	res, err := bankclient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{
		Address: address,
	})
	if err != nil {
		return sdk.NewCoins(), err
	}
	return res.Balances, nil
}

func getDelegations(cctx client.Context, address string) (sdk.Coins, error) {

	stakingclient := stakingtypes.NewQueryClient(cctx)
	res, err := stakingclient.DelegatorDelegations(context.Background(), &stakingtypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: address,
	})
	if err != nil {
		return sdk.NewCoins(), nil
	}

	delegations := sdk.NewCoins()

	for _, delegation := range res.DelegationResponses {
		delegations = delegations.Add(delegation.Balance)
	}

	return delegations, nil
}

func spendableCoins(locked sdk.Coins, balances sdk.Coins) sdk.Coins {
	spendable, hasNeg := balances.SafeSub(locked)
	if hasNeg {
		return sdk.NewCoins()
	}
	return spendable
}

func createContext() client.Context {
	iregistry := codectypes.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(iregistry)
	banktypes.RegisterInterfaces(iregistry)
	stakingtypes.RegisterInterfaces(iregistry)
	vestingtypes.RegisterInterfaces(iregistry)
	cryptocodec.RegisterInterfaces(iregistry)

	cctx := client.Context{}
	cctx = cctx.WithOffline(false)
	cctx = cctx.WithInterfaceRegistry(iregistry)

	return cctx
}
