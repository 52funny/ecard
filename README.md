# Ecard

智慧一卡通 Golang

## Usage

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

**URL 是智慧一卡通的地址，要加 http:// 末尾不用加/**
