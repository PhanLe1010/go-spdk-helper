package basic

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/longhorn/go-spdk-helper/pkg/spdk/client"
	spdktypes "github.com/longhorn/go-spdk-helper/pkg/spdk/types"
)

func BdevLvolCmd() cli.Command {
	return cli.Command{
		Name:      "bdev-lvol",
		ShortName: "lvol",
		Subcommands: []cli.Command{
			BdevLvolCreateCmd(),
			BdevLvolDeleteCmd(),
			BdevLvolGetCmd(),
			BdevLvolSnapshotCmd(),
			BdevLvolCloneCmd(),
			BdevLvolDecoupleParentCmd(),
			BdevLvolResizeCmd(),
		},
	}
}

func BdevLvolCreateCmd() cli.Command {
	return cli.Command{
		Name: "create",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "lvs-name",
				Usage: "Required.",
			},
			cli.StringFlag{
				Name:  "lvol-name",
				Usage: "Required.",
			},
			cli.Uint64Flag{
				Name:  "size",
				Usage: "Required. Specify bdev lvol size in MiB.",
			},
			cli.StringFlag{
				Name:  "uuid",
				Usage: "Optional.",
			},
		},
		Usage: "create a bdev lvol on a lvstore: \"create --lvs-name <LVSTORE NAME> --lvol-name <LVOL NAME> --size <LVOL SIZE in MIB>\"",
		Action: func(c *cli.Context) {
			if err := bdevLvolCreate(c); err != nil {
				logrus.WithError(err).Fatalf("Error running create bdev lvol command")
			}
		},
	}
}

func bdevLvolCreate(c *cli.Context) error {
	spdkCli, err := client.NewClient()
	if err != nil {
		return err
	}

	lvsName := c.String("lvs-name")
	lvolName := c.String("lvol-name")

	uuid, err := spdkCli.BdevLvolCreate(lvsName, lvolName, c.String("uuid"), c.Uint64("size"),
		spdktypes.BdevLvolClearMethodUnmap, true)
	if err != nil {
		return err
	}

	bdevLvolCreateRespJSON, err := json.MarshalIndent(map[string]string{"uuid": uuid, "alias": fmt.Sprintf("%s/%s", lvsName, lvolName)}, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(bdevLvolCreateRespJSON))

	return nil
}

func BdevLvolDeleteCmd() cli.Command {
	return cli.Command{
		Name: "delete",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "alias",
				Usage: "Optional. The alias of a lvol is <LVSTORE NAME>/<LVOL NAME>. Specify this or uuid.",
			},
			cli.StringFlag{
				Name:  "uuid",
				Usage: "Optional. Specify this or alias",
			},
		},
		Usage: "delete a bdev lvol using a block device: \"delete --alias <LVSTORE NAME>/<LVOL NAME>\" or \"delete --uuid <UUID>\"",
		Action: func(c *cli.Context) {
			if err := bdevLvolDelete(c); err != nil {
				logrus.WithError(err).Fatalf("Error running delete bdev lvol command")
			}
		},
	}
}

func bdevLvolDelete(c *cli.Context) error {
	spdkCli, err := client.NewClient()
	if err != nil {
		return err
	}

	name := c.String("alias")
	if name == "" {
		name = c.String("uuid")
	}

	deleted, err := spdkCli.BdevLvolDelete(name)
	if err != nil {
		return err
	}

	bdevLvolDeleteRespJSON, err := json.Marshal(deleted)
	if err != nil {
		return err
	}
	fmt.Println(string(bdevLvolDeleteRespJSON))

	return nil
}

func BdevLvolGetCmd() cli.Command {
	return cli.Command{
		Name: "get",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "alias",
				Usage: "Optional. The alias of a lvol is <LVSTORE NAME>/<LVOL NAME>. If you want to get one specific Lvol info, please input this or uuid.",
			},
			cli.StringFlag{
				Name:  "uuid",
				Usage: "Optional. If you want to get one specific Lvol info, please input this or alias",
			},
			cli.Uint64Flag{
				Name:  "timeout, t",
				Usage: "Optional. Determine the timeout of the execution",
				Value: 0,
			},
		},
		Usage: "get all bdev lvol if the info is not specified: \"get\", or \"get --alias <LVSTORE NAME>/<LVOL NAME>\", or \"get --uuid <UUID>\"",
		Action: func(c *cli.Context) {
			if err := bdevLvolGet(c); err != nil {
				logrus.WithError(err).Fatalf("Error running get bdev lvol command")
			}
		},
	}
}

func bdevLvolGet(c *cli.Context) error {
	spdkCli, err := client.NewClient()
	if err != nil {
		return err
	}

	name := c.String("alias")
	if name == "" {
		name = c.String("uuid")
	}

	bdevLvolGetResp, err := spdkCli.BdevLvolGet(name, c.Uint64("timeout"))
	if err != nil {
		return err
	}

	bdevLvolGetRespJSON, err := json.MarshalIndent(bdevLvolGetResp, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(bdevLvolGetRespJSON))

	return nil
}

func BdevLvolSnapshotCmd() cli.Command {
	return cli.Command{
		Name: "snapshot",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "alias",
				Usage: "Optional. The alias of a lvol is <LVSTORE NAME>/<LVOL NAME>. Specify this or uuid.",
			},
			cli.StringFlag{
				Name:  "uuid",
				Usage: "Optional. Specify this or alias",
			},
			cli.StringFlag{
				Name:  "snapshot-name",
				Usage: "Required. The snapshot lvol name.",
			},
		},
		Usage: "create a snapshot as a new bdev lvol based on an existing one: \"snapshot --alias <LVSTORE NAME>/<LVOL NAME> --snapshot-name <SNAPSHOT NAME>\", or \"snapshot --uuid <UUID> --snapshot-name <SNAPSHOT NAME>\"",
		Action: func(c *cli.Context) {
			if err := bdevLvolSnapshot(c); err != nil {
				logrus.WithError(err).Fatalf("Error running snapshot bdev lvol command")
			}
		},
	}
}

func bdevLvolSnapshot(c *cli.Context) error {
	spdkCli, err := client.NewClient()
	if err != nil {
		return err
	}

	name := c.String("alias")
	if name == "" {
		name = c.String("uuid")
	}

	uuid, err := spdkCli.BdevLvolSnapshot(name, c.String("snapshot-name"))
	if err != nil {
		return err
	}

	fmt.Println(uuid)

	return nil
}

func BdevLvolCloneCmd() cli.Command {
	return cli.Command{
		Name: "clone",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "alias",
				Usage: "Optional. The alias of a snapshot lvol is <LVSTORE NAME>/<LVOL NAME>. Specify this or uuid.",
			},
			cli.StringFlag{
				Name:  "uuid",
				Usage: "Optional. Specify this or alias",
			},
			cli.StringFlag{
				Name:  "clone-name",
				Usage: "Required. The cloned lvol name.",
			},
		},
		Usage: "create a clone lvol based on an existing snapshot lvol: \"clone --alias <LVSTORE NAME>/<SNAPSHOT LVOL NAME> --clone-name <CLONE NAME>\", or \"clone --uuid <SNAPSHOT LVOL UUID> --clone-name <CLONE NAME>\"",
		Action: func(c *cli.Context) {
			if err := bdevLvolClone(c); err != nil {
				logrus.WithError(err).Fatalf("Error running clone bdev lvol command")
			}
		},
	}
}

func bdevLvolClone(c *cli.Context) error {
	spdkCli, err := client.NewClient()
	if err != nil {
		return err
	}

	name := c.String("alias")
	if name == "" {
		name = c.String("uuid")
	}

	uuid, err := spdkCli.BdevLvolClone(name, c.String("clone-name"))
	if err != nil {
		return err
	}

	fmt.Println(uuid)

	return nil
}

func BdevLvolDecoupleParentCmd() cli.Command {
	return cli.Command{
		Name: "decouple",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "alias",
				Usage: "Optional. The alias of a lvol is <LVSTORE NAME>/<LVOL NAME>. Specify this or uuid.",
			},
			cli.StringFlag{
				Name:  "uuid",
				Usage: "Optional. Specify this or alias",
			},
		},
		Usage: "decouple a lvol from its parent lvol: \"decouple --alias <LVSTORE NAME>/<LVOL NAME>\", or \"decouple --uuid <LVOL UUID>\"",
		Action: func(c *cli.Context) {
			if err := bdevLvolDecoupleParent(c); err != nil {
				logrus.WithError(err).Fatalf("Error running decouple parent bdev lvol command")
			}
		},
	}
}

func bdevLvolDecoupleParent(c *cli.Context) error {
	spdkCli, err := client.NewClient()
	if err != nil {
		return err
	}

	name := c.String("alias")
	if name == "" {
		name = c.String("uuid")
	}

	decoupled, err := spdkCli.BdevLvolDecoupleParent(name)
	if err != nil {
		return err
	}

	bdevLvolDecoupleParentRespJSON, err := json.Marshal(decoupled)
	if err != nil {
		return err
	}
	fmt.Println(string(bdevLvolDecoupleParentRespJSON))

	return nil
}

func BdevLvolResizeCmd() cli.Command {
	return cli.Command{
		Name: "resize",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "alias",
				Usage: "Optional. The alias of a snapshot lvol is <LVSTORE NAME>/<LVOL NAME>. Specify this or uuid.",
			},
			cli.StringFlag{
				Name:  "uuid",
				Usage: "Optional. Specify this or alias",
			},
			cli.Uint64Flag{
				Name:  "size",
				Usage: "Required.",
			},
		},
		Usage: "resize a lvol to a new size: \"resize --alias <LVSTORE NAME>/<LVOL NAME> --size <SIZE>\", or \"resize --uuid <LVOL UUID> --size <SIZE>\"",
		Action: func(c *cli.Context) {
			if err := bdevLvolResize(c); err != nil {
				logrus.WithError(err).Fatalf("Error running resize bdev lvol command")
			}
		},
	}
}

func bdevLvolResize(c *cli.Context) error {
	spdkCli, err := client.NewClient()
	if err != nil {
		return err
	}

	name := c.String("alias")
	if name == "" {
		name = c.String("uuid")
	}

	resized, err := spdkCli.BdevLvolResize(name, c.Uint64("size"))
	if err != nil {
		return err
	}

	bdevLvolResizeRespJSON, err := json.Marshal(resized)
	if err != nil {
		return err
	}
	fmt.Println(string(bdevLvolResizeRespJSON))

	return nil
}