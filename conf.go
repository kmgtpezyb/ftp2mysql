package main

import (
	"fmt"
	"time"
	"os"
	"errors"
	"os/user"
	"path/filepath"
	"github.com/kataras/golog"
	"github.com/kmgtpezyb/ftp2mysql/cliftp"
	"github.com/kmgtpezyb/ftp2mysql/climysql"
	"gopkg.in/urfave/cli.v1"
)

type TableFile struct {
	Table string
	File string
}

type procConfig struct {
	date string
	serverDir string
	dataDir string
	logfi *os.File
	okfile string
}

var (
	TableFiles = []TableFile {
		{"9901_0400_01客存款每日余额明细表_通讯","oa_kckmryemxb_0400_01.unl"},
		{"9901_0400_02客户信用等级评价","oa_khxxb.unl"},
	}
)

func DefaultOptions() (*cliftp.FTPOptions) {

	return &cliftp.FTPOptions{
		User:"axxxxx",
		Word:"xxxxxx",
		Server:"172.xxx.xx.xx",
		Port:"21",
	}
}

func DefaultServer() *climysql.ServerConfig {

	return &climysql.ServerConfig {
		User:"xxxxx",
		Pass:"xxxxx",
		Host:"172.xx.xx:3000",
		DbName:"xx",
	}
}

func NewProc(ctx *cli.Context) (*procConfig, error) {
	
	var procdate string	

	if ctx.IsSet("date") {
		procdate = ctx.String("date")
	} else {
		procdate = YesDay(time.Now(),0,0,-1)
	}

	if procdate == "" {
		return nil, errors.New("date is null")
	}

	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	datadir := filepath.Join(user.HomeDir,"ftp2mysql/txt",procdate)
	logdir := filepath.Join(user.HomeDir,"ftp2mysql/log")

	logfi,err:= os.OpenFile(fmt.Sprint(logdir,"/",procdate,".log"),os.O_CREATE|os.O_APPEND|os.O_WRONLY,0775)
	if err != nil {
		return nil,err
	}


	pConfig := &procConfig{
			date:procdate,
			serverDir:"/ods2oa/ods2oa",
			dataDir:datadir,
			logfi:logfi,
			okfile:fmt.Sprint("ok_", procdate, ".txt"),
		}

	golog.SetOutput(pConfig.logfi)

	return pConfig,nil
}
