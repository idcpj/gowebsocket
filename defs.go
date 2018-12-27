package main

type userDb struct {
	Id string
	Name string
	Pwd string
}

const (
	ERROR_VALIDATE = "验证失败"
)


func Checkuser(user,pwd string)  string{
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