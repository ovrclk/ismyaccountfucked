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
	w := tabwriter.NewWriter(os.Stdout, 20, 8, 0, '\t', 0)

	fmt.Fprintf(w, "Address\t%s\t\n", status.Address)
	fmt.Fprintf(w, "IsVesting\t%v\t\n", status.IsVesting)

	if !status.IsVesting {
		w.Flush()
		fmt.Print("\nNOT FUCKED\n")
		return nil
	}

	fmt.Fprintf(w, "Balance\t%v\t\n", formatAmount(status.Balance))
	fmt.Fprintf(w, "Delegated\t%v\t\n", formatAmount(status.Delegated))
	fmt.Fprintf(w, "BalanceLocked\t%v\t\n", formatAmount(status.BalanceLocked))
	fmt.Fprintf(w, "BalanceSpendable\t%v\t\n", formatAmount(status.BalanceSpendable))
	fmt.Fprintf(w, "BalanceVesting\t%v\t\n", formatAmount(status.BalanceVesting))
	fmt.Fprintf(w, "BalanceVested\t%v\t\n", formatAmount(status.BalanceVested))
	fmt.Fprintf(w, "DelegatedFree\t%v\t\n", formatAmount(status.DelegatedFree))
	fmt.Fprintf(w, "DelegatedVesting\t%v\t\n", formatAmount(status.DelegatedVesting))

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
