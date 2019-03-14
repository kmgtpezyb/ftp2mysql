package climysql

import (
	"fmt"
	"os"
	"io"
	"bufio"
	"errors"
        "database/sql"
	"github.com/kataras/golog"
        _ "github.com/go-sql-driver/mysql"
)

var (
	ErrColumns      = errors.New("columns count error")
	ErrHelp         = errors.New("flags error")
)

type ServerConfig struct {
	User		string
	Pass		string
	Host		string
	DbName		string
}

type TableColumn struct {
	ColName string
	DataType string
}

type ClientConfig struct {
	ServerConfig	*ServerConfig
	Table		string
	TableColumns	[]TableColumn
	File		string
	Date		string
	Cols		int
	Count		int
}

func (config *ServerConfig) GetDsn(dbname string) string {

	var dbName string		

	if dbname == "" {
		dbName = config.DbName
	} else {
		dbName = dbname
	}
	
	dbDsn := fmt.Sprint(config.User,":",config.Pass,"@tcp(",config.Host,")/",dbName)
	
	return dbDsn
}

func NewClient(sconfig *ServerConfig, Table, File, Date string) (*ClientConfig, error) {

	config := &ClientConfig {
			ServerConfig: sconfig,
			Table: Table,
			File: File,
			Date: Date,
			Count: 0,
		}

	if err := config.getTableCols(); err != nil {
		return nil, err	
	}
	return config, nil
}

func (config *ClientConfig) getTableCols( ) error {
	
	var column TableColumn

	db, err := sql.Open("mysql",config.ServerConfig.GetDsn("information_schema"))
	if err != nil {
		golog.Println("Open information_schema error")
		return err
	}
	defer db.Close()

	sqlcmd := fmt.Sprint("select column_name, data_type from columns where table_schema='",config.ServerConfig.DbName,"'"," and table_name='",config.Table,"'")

	rows, err := db.Query(sqlcmd)
	if err != nil {
		golog.Println("db query error :",err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&column.ColName, &column.DataType)
		config.Cols ++
		config.TableColumns = append(config.TableColumns,column)
	}

	err = rows.Err() // get any error encountered during iteration

	return err
}

func (config *ClientConfig) LoadFile( ) error {

	file, err := os.Open(config.File)
	if err != nil {
		golog.Println("Open file error",config.File)
		return err		
	}
	defer file.Close()

	db, err := sql.Open("mysql",config.ServerConfig.GetDsn(""))
	if err != nil {
		golog.Println("Sql Open error")
		return err
	}
	defer db.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		sslice := slineToSslice(line)
		if len(sslice) != config.Cols + 1 {
			golog.Println("line:",config.Count+1,"error! columns:",len(sslice)-2,"-",config.Cols)
			return ErrColumns
		}

	        sql := fmt.Sprint("REPLACE INTO ", config.Table," VALUES(")
		for i := 0; i<config.Cols-1; i++ {
			sql = fmt.Sprint(sql,TypeCol(config.TableColumns[i].DataType, sslice[i]), ",")
		}
		sql = fmt.Sprint(sql,TypeCol(config.TableColumns[config.Cols-1].DataType, sslice[config.Cols-1]), ")")
		_, err = db.Exec(sql)
		if err != nil {
			golog.Println("Sql Exec error", err)
			return err
		}

		config.Count ++
	}
	
	golog.Println("Successed!")

	return nil
}
