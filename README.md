# Ecard

智慧一卡通 Golang

## Usage

**URL 是智慧一卡通的地址，要加 http:// 末尾不用加/**

### 获取电费

```go
e := Ecard{
		URL:      "http://60.171.203.79:8090/easytong_portal",
		Username: "xxx",
		Password: "xxx",
	}
e.Login()
fmt.Println(e.ObtainDormitoryElectricity("0", "5", "237"))
```

### 获取余额

```go
e := Ecard{
		URL:      "http://60.171.203.79:8090/easytong_portal",
		Username: "xxx",
		Password: "xxx",
	}
e.Login()
fmt.Println(e.ObtainBalance())
```

### 获取今天消费记录

```go
e := Ecard{
		URL:      "http://60.171.203.79:8090/easytong_portal",
		Username: "xxx",
		Password: "xxx",
	}
e.Login()
billS, err := e.ObtainTodayBill("40")
if err != nil {
	panic(err)
}
fmt.Println(billS)
```

### 支付寝室电费

```go
e := Ecard{
		URL:      "http://60.171.203.79:8090/easytong_portal",
		Username: "xxx",
		Password: "xxx",
	}
e.Login()
for i := 0; i < 5; i++ {
	result, err := e.PayElectricity("0", "5", "k237", "0.01")
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
```
