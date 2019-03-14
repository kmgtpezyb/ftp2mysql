package main

import (
	"os"
//	"time"
	"errors"
	"path/filepath"
	"gopkg.in/urfave/cli.v1"
	"github.com/kataras/golog"
	"github.com/kmgtpezyb/ftp2mysql/cliftp"
	"github.com/kmgtpezyb/ftp2mysql/climysql"
)

func main() {
 
	app := cli.NewApp()
	app.Name = "ftp2mysql"
	app.Author = "WhLiu"
	app.Email = "97097725@qq.com"
	app.Version = "1.0"
	app.Usage = "Download From Ftp Server, Then Into Mysql"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "date,d",
			Usage: "date-type:yyyymmdd",
		},
	}
	app.Commands = []cli.Command{
		ftpCommand,
		insertCommand,
	}



	app.Action = ftp2Mysql
	app.Run(os.Args)
}

func getCommand (ctx *cli.Context) string {

	switch ctx.Command.Name {
	case "ftp":
		return "ftp"
	case "insert":
		return "insert"
	default:
		return "ftp2mysql"
	}
}

func ftp2Mysql(ctx *cli.Context) error {

	procC, err := NewProc(ctx)
	if err != nil {
		golog.Println("newProc error:", err)
		return err
	}
	defer procC.logfi.Close()
	
	command := getCommand(ctx)

	
	golog.Println("")
	golog.Println("-------------------------------------------------------")
	
	golog.Println("Beginning Command[",command,"]Proc Date[",procC.date,"]")
	
	if command == "ftp" || command == "ftp2mysql" {
		if err = procC.procFtp(); err != nil {
			return err
		}
		golog.Println("Success Download From Ftp Server :", procC.serverDir)
	}

	if command == "insert" || command == "ftp2mysql" {
		if err = procC.procMysql(); err != nil {
			return err
		}
		golog.Println("Success Into Mysql From:", procC.dataDir)
	}
	
	return nil
}

func (procC *procConfig) procFtp() error {

	opts := DefaultOptions()

	err := cliftp.NewFTPConn(opts)
	if err != nil {
		golog.Println("NewFTPConn Error:", err)
		return err
	}

	err = procC.checkFtpOk(opts)
	if err != nil {
		golog.Println("checkFtpOk error:", err)
		return err
	}

	if err = procC.getFtpFiles(opts); err != nil {
		return err
	}

	return nil
}

func (proc *procConfig) checkFtpOk(opts *cliftp.FTPOptions) error {

	pathFunc := func(opts *cliftp.FTPOptions, proc *procConfig) (bool, error) {
		deepdir := false
		fileok := false
	Loop:

		for {
			fileItems, err := cliftp.ListFiles(opts.Conn,proc.serverDir)
			if err != nil {
				golog.Println("ListFiles error:", err)
				return false, err
			}
			for _, fileitem := range fileItems {
				if fileitem.Name == proc.okfile && fileitem.Type == "file" {
					fileok = true
				}
				if fileitem.Name == proc.date && fileitem.Type == "directory" {
					deepdir = true
				}
			}
			if fileok {
				if deepdir {
					proc.serverDir = filepath.Join(proc.serverDir,proc.date)
				}
				break Loop
			}
			//time.Sleep(2*time.Minute)
			return false,errors.New("pathFunc error")
		}
		return deepdir, nil
	}

Loop1:
	for {
		dateItems, err := cliftp.ListFiles(opts.Conn,proc.serverDir)
		if err != nil {
			golog.Println("ListFiles error:", err)
			return err
		}
		for _, dateitem := range dateItems {
			if dateitem.Name == proc.date && dateitem.Type == "directory" {
				proc.serverDir = filepath.Join(proc.serverDir,proc.date)
				break Loop1
			}
		}
		// time.Sleep(10*time.Minute)
		golog.Println("directory :", proc.date, " is not exist!")
		return errors.New("date error")
	}

	deepdir, err := pathFunc(opts,proc)
	if err != nil {
		golog.Println("pathFunc error:", err)
		return err
	}
	if deepdir {
		_, err = pathFunc(opts,proc)
		if err != nil {
			golog.Println("pathFunc error:", err)
			return err
		}
	}

	return nil
}

func (proc *procConfig) getFtpFiles(opts *cliftp.FTPOptions) error {
	
	if err := os.MkdirAll(proc.dataDir,0755); err != nil {
		golog.Println("mkdir error:", err)
		return err
	}
	if err := os.Chdir(proc.dataDir); err != nil {
		golog.Println("os.Chdir error:", err)
		return err
	}

	return cliftp.DownloadFiles(opts.Conn,proc.serverDir)
}

func (proc *procConfig) procMysql() (error) {
	
	golog.Println("Beginning ToMysql... Date:", proc.date)

	server := DefaultServer()

	for _, table := range TableFiles {

		filename := filepath.Join(proc.dataDir,table.File)

		golog.Println("Begining Load file:", filename, "Into", table.Table)

		client, err := climysql.NewClient(server, table.Table, filename, proc.date)
		if err != nil {
			golog.Println("NewClient error:", err)
			return err
		}
		golog.Println("Begin Load...", table.Table)
		err = client.LoadFile()
		if err != nil {
			golog.Println("LoadFile error:", err)
			return err
		}
		golog.Println("Success Load", client.Count, "Line Into", table.Table)
	}	

	golog.Println("ToMysql End... Date:", proc.date)

	return nil
}
