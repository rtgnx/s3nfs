package main

import (
	"log"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/rtgnx/s3nfs/internal/fs/nfsd"
	"gopkg.in/yaml.v3"
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
		addr    = cmd.StringOpt("addr", ":6969", "addr to listen to")
		cfgFile = cmd.StringOpt("cfg", "/etc/s3nfsd.yml", "s3nfsd config")
	)

	cmd.Action = func() {
		cfg := new(nfsd.Config)
		fd, err := os.Open(*cfgFile)
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()

		if err := yaml.NewDecoder(fd).Decode(cfg); err != nil {
			log.Fatal(err)
		}
		fs := nfsd.NewMapperFS(*cfg)
		if err := nfsd.Serve(*addr, fs); err != nil {
			log.Fatal(err)
		}
	}

}
