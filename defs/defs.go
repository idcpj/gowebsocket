package defs

import "encoding/json"

type SendFormat struct {
	Type     int    `json:"type"`
	Name     string `json:"name"` //用于登录
	ClientId string `json:"client_id"`
	Pwd      string `json:"pwd"`   //用于登录
	SendId   string `json:"send_id"`
	Param    string `json:"param"`  //用于传递参数
	Content  string `json:"content"`
	ResCode	string `json:"res_code"`
	Response string `json:"response"`
}

func (s *SendFormat) Unmarshal(message string) (e error){
	 e = json.Unmarshal([]byte(message), &s)
	return e
}

func (s *SendFormat) Marshal()([]byte){
	msg, _ := json.Marshal(s)
	return msg
}
func NewSendError(msg string)SendFormat{
	return SendFormat{
		ResCode:"0",
		Response:msg,
	}
}
func NewSendSuccess(msg string)SendFormat{
	return SendFormat{
		ResCode:"1",
		Response:msg,
	}
}
const (
	//error
	ERROR_LOGIN          = "登录失败"
	ERROR_NO_CLIENT_ID   = "没有对应的id"
	ERROR_BAD_MSG_FORMAT ="发送数据格式不正确"
	//ok
	OKEY_LOGIN="登录成功"
	OKEY_MSG ="消息发送成功"

	// send type
	TYPE_LOGIN = 1
	TYPE_SEND_ID = 2
)


func Checkuser(user,pwd string)  string{
	type userDb struct {
		Id string
		Name string
		Pwd string
	}

	a :=[]userDb{}

	b :=userDb{
		Id:"1",
		Name:"c",
		Pwd:"123",
	}
	c :=userDb{
		Id:"2",
		Name:"p",
		Pwd:"123",
	}
	d :=userDb{
		Id:"3",
		Name:"j",
		Pwd:"123",
	}

	a =append(a, b)
	a =append(a, c)
	a =append(a, d)
	for _,v :=range a{
		if v.Name==user && v.Pwd==pwd {
			return v.Id
		}
	}
	return ""

}