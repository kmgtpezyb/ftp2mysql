package main

import (
	"gopkg.in/urfave/cli.v1"	
)

var (

	ftpCommand = cli.Command{
		Action:	ftp2Mysql,
		Name: "ftp",
		Usage: "only download files from ftp",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "date,d",
				Usage: "command ftp -d YYYYMMDD",
			},
		},
	}
	insertCommand = cli.Command{
		Action:	ftp2Mysql,
		Name: "insert",
		Usage: "only insert into mysql from files",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "date,d",
				Usage: "command ftp -d YYYYMMDD",
			},
		},
	}
)
