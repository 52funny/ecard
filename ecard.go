package ecard

import (
	"fmt"
	"log"
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

// 获取rsa的公钥
func (e *Ecard) getKeyMap() (exponent string, modulus string) {
	resp, err := req.Post(e.URL + "/publiccombo/keyPair")
	if err != nil {
		log.Println(err)
	}
	keyMap := gjson.GetBytes(resp.Bytes(), "publicKeyMap")
	exponent = keyMap.Get("exponent").String()
	modulus = keyMap.Get("modulus").String()
	return
}

// 获取验证码
func (e *Ecard) getCodeImg() (code string) {
	resp, err := req.Get(e.URL + "/jcaptcha.jpg?" + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		log.Println(err)
	}
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(resp.Bytes())
	code, _ = client.Text()
	return
}

// Login 登陆系统
func (e *Ecard) Login() {
	exponent, modulus := e.getKeyMap()
	header := req.Header{
		"X-Requested-With": "XMLHttpRequest",
	}
	var state int64
OUT:
	for state != 3 {
		code := e.getCodeImg()
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
}

// ObtainDormitoryElectricity 获取寝室电费余额
func (e *Ecard) ObtainDormitoryElectricity(areaNo string, buildingNo string, roomNo string) string {
	data := req.Param{
		"data": `{"itemNum":"1","areano":"` + areaNo + `","buildingno":"` + buildingNo + `","roomno":"` + roomNo + `"}`,
	}
	header := req.Header{
		"X-Requested-With": "XMLHttpRequest",
	}
	resp, _ := req.Post(e.URL+"/payFee/getBalance", data, header)
	return gjson.GetBytes(resp.Bytes(), "feeDate.balance").String()
}

//IsCookieOverDue 判断cookie是否过期
func (e *Ecard) IsCookieOverDue() bool {
	resp, err := req.Get(e.URL)
	if err != nil {
		fmt.Println(err)
	}
	reg := regexp.MustCompile("<title>(.*)</title>")
	ans := reg.FindSubmatch(resp.Bytes())
	return string(ans[1]) == "智慧一卡通－登录"
}
