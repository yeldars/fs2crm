package main

/*
Integration between FreeSwitch and BAPPS CRM
*/
import (
	_ "github.com/lib/pq"

	"github.com/astaxie/beego/orm"
	"github.com/yeldars/fs2crm/cdrutils"
	"os"
)

func init() {
	err := orm.RegisterDriver("postgres", orm.DRPostgres)
	if err!=nil{
		panic(err)
	}
	//url :="host=192.168.1.105 dbname=fs user=fs password=fs connect_timeout=10 sslmode=disable"
	uri := os.Getenv("BAPPSKZ_FS2_CRM_DBCONNECTION")

	orm.RegisterDataBase("default", "postgres", uri)


}

func main()  {

	cdrutils.UploadCdr()

}
