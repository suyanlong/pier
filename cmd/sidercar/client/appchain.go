package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/link33/sidercar/internal/repo"
	appchainmgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/urfave/cli"
)

type Approve struct {
	Id         string `json:"id"`
	IsApproved int32  `json:"is_approved"`
	Desc       string `json:"desc"`
}

var clientCMD = cli.Command{
	Name:  "client",
	Usage: "Command about appchain in sidercar",
	Subcommands: []cli.Command{
		{
			Name:  "register",
			Usage: "Register appchain in sidercar",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "sidercar-id",
					Usage:    "Specify target sidercar id",
					Required: true,
				},
				cli.StringFlag{
					Name:     "name",
					Usage:    "Specify appchain name",
					Required: true,
				},
				cli.StringFlag{
					Name:     "type",
					Usage:    "Specify appchain type",
					Required: true,
				},
				cli.StringFlag{
					Name:     "desc",
					Usage:    "Specify appchain description",
					Required: true,
				},
				cli.StringFlag{
					Name:     "version",
					Usage:    "Specify appchain version",
					Required: true,
				},
				cli.StringFlag{
					Name:     "validators",
					Usage:    "Specify appchain validators path",
					Required: true,
				},
				cli.StringFlag{
					Name:     "consensus-type",
					Usage:    "Specify appchain consensus type",
					Required: true,
				},
			},
			Action: registerSidercarAppchain,
		},
		{
			Name:  "update",
			Usage: "Update appchain in sidercar",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "sidercar-id",
					Usage:    "Specify target sidercar id",
					Required: false,
				},
				cli.StringFlag{
					Name:     "name",
					Usage:    "Specify appchain name",
					Required: false,
				},
				cli.StringFlag{
					Name:     "type",
					Usage:    "Specify appchain type",
					Required: false,
				},
				cli.StringFlag{
					Name:     "desc",
					Usage:    "Specify appchain description",
					Required: false,
				},
				cli.StringFlag{
					Name:     "version",
					Usage:    "Specify appchain version",
					Required: false,
				},
				cli.StringFlag{
					Name:     "validators",
					Usage:    "Specify appchain validators path",
					Required: false,
				},
				cli.StringFlag{
					Name:     "consensus-type",
					Usage:    "Specify appchain consensus type",
					Required: false,
				},
			},
			Action: updateSidercarAppchain,
		},
		{
			Name:  "audit",
			Usage: "Audit appchain in sidercar",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "id",
					Usage:    "Specific appchain id",
					Required: true,
				},
				cli.StringFlag{
					Name:     "is-approved",
					Usage:    "Specific approved signal",
					Required: true,
				},
				cli.StringFlag{
					Name:     "desc",
					Usage:    "Specific audit description",
					Required: true,
				},
			},
			Action: auditSidercarAppchain,
		},
		{
			Name:  "get",
			Usage: "Get appchain info",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "sidercar-id",
					Usage:    "Specific target sidercar id",
					Required: true,
				},
			},
			Action: getSidercarAppchain,
		},
	},
}

func LoadClientCMD() cli.Command {
	return clientCMD
}

func registerSidercarAppchain(ctx *cli.Context) error {
	return saveSidercarAppchain(ctx, RegisterAppchainUrl)
}

func updateSidercarAppchain(ctx *cli.Context) error {
	return saveSidercarAppchain(ctx, UpdateAppchainUrl)
}

func auditSidercarAppchain(ctx *cli.Context) error {
	id := ctx.String("id")
	isApproved := ctx.String("is-approved")
	desc := ctx.String("desc")

	ia, err := strconv.ParseInt(isApproved, 0, 64)
	if err != nil {
		return fmt.Errorf("isApproved must be 0 or 1: %w", err)
	}
	approve := &Approve{
		Id:         id,
		IsApproved: int32(ia),
		Desc:       desc,
	}

	data, err := json.Marshal(approve)
	if err != nil {
		return err
	}
	url, err := getURL(ctx, AuditAppchainUrl)
	if err != nil {
		return err
	}

	_, err = httpPost(url, data)
	if err != nil {
		return err
	}

	fmt.Printf("audit appchain %s successfully\n", id)

	return nil
}

func saveSidercarAppchain(ctx *cli.Context, path string) error {
	sidercar := ctx.String("sidercar-id")
	name := ctx.String("name")
	typ := ctx.String("type")
	desc := ctx.String("desc")
	version := ctx.String("version")
	validatorsPath := ctx.String("validators")
	consensusType := ctx.String("consensus-type")

	url, err := getURL(ctx, fmt.Sprintf("%s?sidercar_id=%s", path, sidercar))
	if err != nil {
		return err
	}
	res, err := httpGet(url)
	if err != nil {
		return err
	}

	appchainInfo := appchainmgr.Appchain{}
	if err = json.Unmarshal(res, &appchainInfo); err != nil {
		return err
	}
	if name == "" {
		name = appchainInfo.Name
	}
	if typ == "" {
		typ = appchainInfo.ChainType
	}
	if desc == "" {
		desc = appchainInfo.Desc
	}
	if version == "" {
		version = appchainInfo.Version
	}
	validators := ""
	if validatorsPath == "" {
		validators = appchainInfo.Validators
	} else {
		data, err := ioutil.ReadFile(validatorsPath)
		if err != nil {
			return fmt.Errorf("read validators file: %w", err)
		}
		validators = string(data)
	}
	if consensusType == "" {
		consensusType = appchainInfo.ConsensusType
	}

	repoRoot, err := repo.PathRootWithDefault(ctx.GlobalString("repo"))
	if err != nil {
		return err
	}

	pubKey, err := getPubKey(repo.KeyPath(repoRoot))
	if err != nil {
		return fmt.Errorf("get public key: %w", err)
	}
	addr, _ := pubKey.Address()
	pubKeyBytes, _ := pubKey.Bytes()
	appchain := &appchainmgr.Appchain{
		ID:            addr.String(),
		Name:          name,
		Validators:    validators,
		ConsensusType: consensusType,
		ChainType:     typ,
		Desc:          desc,
		Version:       version,
		PublicKey:     string(pubKeyBytes),
	}

	data, err := json.Marshal(appchain)
	if err != nil {
		return fmt.Errorf("marshal appchain error: %w", err)
	}

	url, err = getURL(ctx, fmt.Sprintf("%s?sidercar_id=%s", path, sidercar))
	if err != nil {
		return err
	}
	resp, err := httpPost(url, data)
	if err != nil {
		return err
	}

	fmt.Println(parseResponse(resp))

	return nil
}

func getSidercarAppchain(ctx *cli.Context) error {
	targetSidercarID := ctx.String("sidercar-id")

	url, err := getURL(ctx, fmt.Sprintf("%s?sidercar_id=%s", GetAppchainUrl, targetSidercarID))
	if err != nil {
		return err
	}
	res, err := httpGet(url)
	if err != nil {
		return err
	}
	fmt.Println(parseResponse(res))

	return nil
}

func getPubKey(keyPath string) (crypto.PublicKey, error) {
	privKey, err := asym.RestorePrivateKey(keyPath, "bitxhub")
	if err != nil {
		return nil, err
	}

	return privKey.PublicKey(), nil
}
