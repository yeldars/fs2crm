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
		rowCount   int `json:"rowCount"`
		Error      string `json:"error"`
		Items [] orm.Params`json:"items"`
	}

	//alter table cdr add status varchar(50) default '0';

	o := orm.NewOrm()
	o.Using("default")

	_,err := o.Raw("update cdr set status='1'  "+
	" where id in (select id from cdr where status='0' limit 1000)	").Exec()
	if err!=nil {
		panic(err)
	}
	var arr [] orm.Params
	_,err = o.Raw("select c.id,c.local_ip_v4,caller_id_name,destination_number,context"+
	",to_char(start_stamp,'YYYY-DD-MM HH24:MI:SS TZ')as start_stamp"+
	",to_char(answer_stamp,'YYYY-DD-MM HH24:MI:SS TZ')as answer_stamp"+
	",to_char(end_stamp,'YYYY-DD-MM HH24:MI:SS TZ')as end_stamp "+
	",duration,billsec,hangup_cause,uuid,bleg_uuid,accountcode,read_codec,write_codec,sip_hangup_disposition,ani"+
	"  from cdr c where status='1'").Values(&arr)
	if err!=nil {
		panic(err)
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
	respO.rowCount = len(arr)
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
