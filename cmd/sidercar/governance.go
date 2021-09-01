package main

import (
	"fmt"

	rpcx "github.com/link33/sidercar/hub/client"
	"github.com/link33/sidercar/model/constant"
	"github.com/urfave/cli"
)

var governanceCMD = cli.Command{
	Name:  "proposals",
	Usage: "proposals manage command",
	Subcommands: cli.Commands{
		cli.Command{
			Name:  "withdraw",
			Usage: "withdraw a proposal",
			Flags: []cli.Flag{
				adminKeyPathFlag,
				cli.StringFlag{
					Name:     "id",
					Usage:    "proposal id",
					Required: true,
				},
			},
			Action: withdraw,
		},
	},
}

// TODO
func withdraw(ctx *cli.Context) error {
	chainAdminKeyPath := ctx.String("admin-key")
	id := ctx.String("id")

	client, _, err := initClientWithKeyPath(ctx, chainAdminKeyPath)
	if err != nil {
		return fmt.Errorf("load client: %w", err)
	}

	// TODO modify hub client
	receipt, err := client.InvokeBVMContract(
		constant.GovernanceContractAddr.Address(),
		"WithdrawProposal", nil, rpcx.String(id),
	)
	if err != nil {
		return fmt.Errorf("invoke bvm contract: %w", err)
	}

	if !receipt.IsSuccess() {
		return fmt.Errorf("invoke withdraw proposal: %s", receipt.Ret)
	}

	return nil
}
