package main

import (
	"fmt"
	"os"
	"text/tabwriter"

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
}

func (c *queryCmd) Run(ctx *runctx) error {
	status, err := getStatus(ctx.cctx, ctx.denom, c.Address)
	if err != nil {
		return err
	}

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

	w := tabwriter.NewWriter(os.Stdout, 20, 8, 0, '\t', 0)
	// fmt.Fprintf(w, "Address\t%s\t\n", status.Address)
	// fmt.Fprintf(w, "IsVesting\t%v\t\n", status.IsVesting)

	fmt.Fprintf(w, "Balance\t%13s\t\n", formatAmount(status.Balance))
	fmt.Fprintf(w, "Delegated\t%13s\t\n", formatAmount(status.Delegated))
	fmt.Fprintf(w, "BalanceLocked\t%13s\t\n", formatAmount(status.BalanceLocked))
	fmt.Fprintf(w, "BalanceSpendable\t%13s\t\n", formatAmount(status.BalanceSpendable))
	fmt.Fprintf(w, "BalanceVesting\t%13s\t\n", formatAmount(status.BalanceVesting))
	fmt.Fprintf(w, "BalanceVested\t%13s\t\n", formatAmount(status.BalanceVested))
	fmt.Fprintf(w, "DelegatedFree\t%13s\t\n", formatAmount(status.DelegatedFree))
	fmt.Fprintf(w, "DelegatedVesting\t%13s\t\n", formatAmount(status.DelegatedVesting))

	w.Flush()

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
