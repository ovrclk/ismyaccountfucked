package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/cosmos/cosmos-sdk/client"
	tmrpc "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	denom   = "uakt"
	rpcNode = "http://135.181.60.250:26657"
)

type runctx struct {
	cctx  client.Context
	denom string
}

type queryCmd struct {
	Address string `arg help:"Address"`
	Verbose bool   `help:"Verbose output" default:"false"`
}

func (c *queryCmd) Run(ctx *runctx) error {
	status, err := getStatus(ctx.cctx, ctx.denom, c.Address)
	if err != nil {
		return err
	}

	if c.Verbose {
		return c.runVerbose(status)
	}

	fmt.Print(status.Address)
	fmt.Print("\t")

	switch {
	case !status.IsVesting:
		fmt.Print("unaffected")
	case status.BalanceSpendable.IsZero():
		fmt.Print("locked")
	case status.IsCorrupted():
		fmt.Print("corrupted")
	default:
		fmt.Print("uncorrupted")
	}

	fmt.Print("\n")

	return nil
}

func (c *queryCmd) runVerbose(status Status) error {
	if !status.IsVesting {
		fmt.Printf("%s is NOT FUCKED, it is not a vesting account\n", status.Address)
		return nil
	}

	if status.BalanceSpendable.IsZero() {
		fmt.Printf("\n%s is *KINDA* FUCKED. There are no spendable coins until we update.\n", status.Address)
	} else {
		fmt.Printf("\n%s is NOT *YET* FUCKED, do not modify delegation until we update so it stays that way\n", status.Address)
	}

	fmt.Print("\n\n")

	fmt.Printf("Balance                  %18s\t\n", formatAmount(status.Balance))
	fmt.Printf("Delegated                %18s\t\n", formatAmount(status.Delegated))
	fmt.Printf("BalanceLocked            %18s\t\n", formatAmount(status.BalanceLocked))
	fmt.Printf("BalanceSpendable         %18s\t\n", formatAmount(status.BalanceSpendable))
	fmt.Printf("BalanceVesting           %18s\t\n", formatAmount(status.BalanceVesting))
	fmt.Printf("BalanceVested            %18s\t\n", formatAmount(status.BalanceVested))
	fmt.Printf("DelegatedFree            %18s\t\n", formatAmount(status.DelegatedFree))
	fmt.Printf("DelegatedVesting         %18s\t\n", formatAmount(status.DelegatedVesting))
	fmt.Printf("ExpectedDelegatedFree    %18s\t\n", formatAmount(status.ExpectedDelegatedFree))
	fmt.Printf("ExpectedDelegatedVesting %18s\t\n", formatAmount(status.ExpectedDelegatedVesting))

	return nil
}

func main() {
	var cli struct {
		Node  string `help:"RPC URI" default:"http://135.181.60.250:26657"`
		Denom string `help:"Denomination" default:"uakt"`

		Query  queryCmd  `cmd help:"Query account" default:"1"`
		Server serverCmd `cmd help:"Run server"`
	}

	ctx := kong.Parse(&cli)

	tmclient, err := tmrpc.New(cli.Node, "")
	ctx.FatalIfErrorf(err)

	cctx := createContext()
	cctx = cctx.WithClient(tmclient)

	rctx := runctx{
		cctx:  cctx,
		denom: cli.Denom,
	}

	ctx.FatalIfErrorf(ctx.Run(&rctx))
}
