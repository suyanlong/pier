package main

import (
	"fmt"

	rpcx "github.com/link33/sidercar/hub/client"
	"github.com/link33/sidercar/internal/repo"
	"github.com/link33/sidercar/model/constant"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/urfave/cli"
)

var interchainCMD = cli.Command{
	Name:  "interchain",
	Usage: "Query interchain info",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "key",
			Usage:    "Specific key.json path",
			Required: true,
		},
	},
	Subcommands: []cli.Command{
		{
			Name:  "ibtp",
			Usage: "Query ibtp by id",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "id",
					Usage:    "Specific ibtp id",
					Required: true,
				},
			},
			Action: getIBTP,
		},
	},
}

func getIBTP(ctx *cli.Context) error {
	id := ctx.String("id")

	repoRoot, err := repo.PathRootWithDefault(ctx.GlobalString("repo"))
	if err != nil {
		return err
	}

	config, err := repo.UnmarshalConfig(repoRoot)
	if err != nil {
		return fmt.Errorf("init config error: %s", err)
	}

	client, err := loadClient(repo.KeyPath(repoRoot), config.Peer.Peers, ctx)
	if err != nil {
		return fmt.Errorf("load client: %w", err)
	}

	receipt, err := client.InvokeBVMContract(
		constant.InterchainContractAddr.Address(),
		"GetIBTPByID", nil,
		rpcx.String(id),
	)
	if err != nil {
		return err
	}

	hash := types.NewHash(receipt.Ret)

	fmt.Printf("Tx hash: %s\n", hash.String())

	response, err := client.GetTransaction(hash.String())
	if err != nil {
		return err
	}

	fmt.Println(response)

	return nil
}
