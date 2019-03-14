package climysql

import ( 
	"strings"
	"strconv"
	"github.com/axgle/mahonia"	
	"github.com/kataras/golog"
)

func GbkToUtf8(str string) (string) {

	dec:=mahonia.NewDecoder("gbk")
	ret := dec.ConvertString(str)
	return ret 
}

func Utf8ToGbk(str string) (string) {

	enc:=mahonia.NewEncoder("gbk")
	ret := enc.ConvertString(str)
	return ret 
}

func slineToSslice(line string) []string {

	return strings.Split(GbkToUtf8(line),",")
}

func TypeCol(DateType string, strcol string) interface{} {

        if strings.Contains("int,tinyint,bigint,smallint", DateType) {

		if strings.HasPrefix(strcol,"\"") && strings.HasSuffix(strcol,"\"") {

			strcol = strings.Replace(strcol,"\"\"","0",1)

			outi, err := strconv.Atoi( strings.Replace(strcol,"\"","",-1) )
			if err != nil {
				golog.Println("strconv.Atoi error:", err)
				panic("strconv.Atoi error")
			}
                	return outi
		}
               	return strcol
        }
        return strcol
}
