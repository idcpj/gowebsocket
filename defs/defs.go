package defs

import "encoding/json"

type SendFormat struct {
	Type int `json:"type"`
	Name string `json:"name"` //用于登录
	Pwd string `json:"pwd"`   //用于登录
	SendId string `json:"send_id"`
	Param string `json:"param"`  //用于传递参数
	Content string `json:"content"`
}

func (s *SendFormat) Form(message string) (e error){
	 e = json.Unmarshal([]byte(message), &s)
	return e
}

const (
	ERROR_VALIDATE = "验证失败"
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