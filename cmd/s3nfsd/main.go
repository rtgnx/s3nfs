package main

import (
	"log"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/rtgnx/s3nfs/internal/fs/nfsd"
	"github.com/rtgnx/s3nfs/internal/fs/s3"
)

func main() {
	app := cli.App("s3nfsd", "")
	app.Command("serve", "Start NFS bridge", cmdServe)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func cmdServe(cmd *cli.Cmd) {

	var (
		addr = cmd.StringOpt("addr", ":6969", "addr to listen to")
	)

	cmd.Action = func() {
		s3, err := s3.FromEnv()
		if err != nil {
			log.Fatal(err)
		}
		if err := nfsd.Serve(*addr, s3); err != nil {
			log.Fatal(err)
		}
	}

}
