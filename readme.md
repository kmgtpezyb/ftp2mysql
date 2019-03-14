一、功能描述：将Ftp服务端的文件下载，并导入Mysql数据库。 

    核心系统每天将行方需要数据按日期目录上传至Ftp服务器，每个文件与Mysql表结构一一对应，文件中的字段用“,”分割

二、用到的第三方包

    CLI命令行工具	gopkg.in/urfave/cli.v1
    日志包		github.com/kataras/golog
    Ftp包		github.com/jlaffaye/ftp
    Mysql驱动		github.com/go-sql-driver/mysql

三、编译

    执行make，会将可执行文件ftp2mysql安装到$HOME/bin/gobin下

四、运行

    参数：-d date(yyyymmdd) 处理date这天的数据，如果不传该参数，默认取当前系统日期的前一天

    命令：ftp2mysql ftp [-d date]只从Ftp服务器下载文件到本地$HOME/ftp2mysql/txt/date下
          ftp2mysql insert [-d date]将本地$HOME/ftp2mysql/txt/date下的文件导入Mysql
          ftp2mysql [-d date]下载文件到本地$HOME/ftp2mysql/txt/date下然后导入Mysql

五、参数配置

    目前参数直接写到程序conf.go中

    1、文件与表名对照关系：导入Mysql环节，会range TableFile，然后将一一导入

    var (
        TableFiles = []TableFile {
                {"9901_0400_01客存款每日余额明细表_通讯","oa_kckmryemxb_0400_01.unl"},
                {"9901_0400_02客户信用等级评价","oa_khxxb.unl"},
		...
        }
    )

    2、Ftp服务端配置

    func DefaultOptions() (*cliftp.FTPOptions) {

        return &cliftp.FTPOptions{
                User:"FTP_tongxun",	//Ftp用户名
                Word:"aaaaaaaa",	//Ftp密码
                Server:"172.168.98.xx",	//Ftp地址
                Port:"21",		//Ftp端口
        }
    }

    3、MYsql服务端配置

    DefaultServer() *climysql.ServerConfig {

        return &climysql.ServerConfig {
                User:"useruser",		//Mysql用户名
                Pass:"aaaaxxxxx",	//Mysql密码
                Host:"172.xx.xx.xx:3000",//Mysql地址及端口
                DbName:"办公系统",	//Mysql数据库
        }
    }
