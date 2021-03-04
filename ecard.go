package ecard

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/52funny/ecard/utils"
	"github.com/imroc/req"
	"github.com/otiai10/gosseract"
	"github.com/tidwall/gjson"
)

// Ecard 是一个配置项目
type Ecard struct {
	URL      string
	Username string
	Password string
}

// Bill 账单
// Time 交易时间
// Content 交易内容
// Money 交易金额
// Balance 账户余额
type Bill struct {
	Time    string
	Content string
	Money   float64
	Balance float64
}

// 获取rsa的公钥
func (e *Ecard) getKeyMap() (exponent string, modulus string, err error) {
	resp, err := req.Post(e.URL + "/publiccombo/keyPair")
	if err != nil {
		return
	}
	keyMap := gjson.GetBytes(resp.Bytes(), "publicKeyMap")
	exponent = keyMap.Get("exponent").String()
	modulus = keyMap.Get("modulus").String()
	return
}

// 获取验证码
func (e *Ecard) getCodeImg() (code string, err error) {
	resp, err := req.Get(e.URL + "/jcaptcha.jpg?" + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return
	}
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(resp.Bytes())
	code, _ = client.Text()
	return
}

// Login 登陆系统
func (e *Ecard) Login() (err error) {
	exponent, modulus, err := e.getKeyMap()
	if err != nil {
		return
	}
	header := req.Header{
		"X-Requested-With": "XMLHttpRequest",
	}
	var state int64
OUT:
	for state != 3 {
		code, err := e.getCodeImg()
		if err != nil {
			return err
		}
		data := req.Param{
			"username":     utils.RsaEncrypt(exponent, modulus, e.Username),
			"password":     utils.RsaEncrypt(exponent, modulus, e.Password),
			"jcaptchacode": code,
		}
		resp, _ := req.Post(e.URL+"/login", data, header)
		state = gjson.GetBytes(resp.Bytes(), "ajaxState").Int()
		msg := gjson.GetBytes(resp.Bytes(), "msg").String()
		if state != 3 {
			fmt.Println(msg)
		}
		switch msg {
		case "账号不存在":
			break OUT
		case "用户名或密码错误":
			break OUT
		}
	}
	if state == 3 {
		fmt.Println("登陆成功:" + e.Username)
	}
	return
}

// MustLogin 会panic error
func (e *Ecard) MustLogin() {
	err := e.Login()
	if err != nil {
		panic(err)
	}
}

//IsCookieOverDue 判断cookie是否过期
func (e *Ecard) IsCookieOverDue() (b bool, err error) {
	resp, err := req.Get(e.URL)
	if err != nil {
		return true, err
	}
	reg := regexp.MustCompile("<title>(.*)</title>")
	ans := reg.FindSubmatch(resp.Bytes())
	return bytes.Equal(ans[1], []byte("智慧一卡通－登录")), nil
}
