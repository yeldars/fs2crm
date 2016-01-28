package cdrutils
import (
	_ "github.com/lib/pq"
	"fmt"
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"net/http"
	"strings"
	"os"
)

func UploadCdr()  {


	//cdrCRMUrl := "http://dev.XXX.XXX.kz/cdr"
	cdrCRMUrl := os.Getenv("BAPPSKZ_FS2_CRM_CDRCRMURL")


	type queryGetResponse struct {
		RowCount   int `json:"rowCount"`
		Error      string `json:"error"`
		Items [] orm.Params`json:"items"`
	}

	//alter table cdr add status varchar(50) default '0';

	o := orm.NewOrm()
	o.Using("default")

	o.Raw("update cdr set status='0'  ").Exec() //TEMP
	_,err := o.Raw("update cdr set status='1'  "+
	" where id in (select id from cdr where status='0' limit 1000)	").Exec()
	if err!=nil {
		panic(err)
	}
	var arr [] orm.Params
	sql := "select c.id,c.local_ip_v4,caller_id_name,caller_id_number,destination_number,context"+
	",extract(epoch from start_stamp) start_stamp "+
	",extract(epoch from answer_stamp) answer_stamp " +
	",extract(epoch from end_stamp) end_stamp " +
	",duration,billsec,hangup_cause,uuid,bleg_uuid,accountcode,read_codec,write_codec,sip_hangup_disposition,ani"+
	"  from cdr c where status='1'"
	fmt.Println(sql)
	_,err = o.Raw(sql).Values(&arr)
	if err!=nil {
		//panic(err)
	}
	//	id |  local_ip_v4  | caller_id_name | caller_id_number | destination_number | context |      start_stamp       |      answer_stamp      |       e
	//	nd_stamp        | duration | billsec |   hangup_cause    |                 uuid                 |              bleg_uuid               | accountco
	//	de | read_codec | write_codec | sip_hangup_disposition | ani

	if len(arr)==0{
		fmt.Println("Nothing to upload")
		return
	}
	respO := queryGetResponse{}
	respO.Items = arr
	respO.Error = "0"
	respO.RowCount = len(arr)
	jsonData, err := json.Marshal(respO)

	if err==nil {
		fmt.Print(string(jsonData))
	}


	_, err = http.Post(cdrCRMUrl, "text/json", strings.NewReader(string(jsonData)))

	if err!=nil{
		panic(err)
	}
	o.Raw("update cdr set status='2' where status='1'").Exec()


}
