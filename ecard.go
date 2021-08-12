package ecard

import (
	"bytes"
	"fmt"
	"net/http"
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
	// Store the cookie
	Cookie *http.Cookie
	req    *req.Req
}

// Bill 账单
// Time 交易时间
// Content 交易内容
// Money 交易金额
// Balance 账户余额
type Bill struct {
	Time     string
	Content  string
	Merchant string
	Location string
	Money    float64
	Balance  float64
}

// New method will return the Ecard object
func New(Username, Password, URL string) *Ecard {
	e := &Ecard{
		Username: Username,
		Password: Password,
		URL:      URL,
	}
	req := req.New()
	// req.EnableCookie(false)
	e.req = req
	return e
}

// 获取rsa的公钥
func (e *Ecard) getKeyMap(cookie *http.Cookie) (exponent string, modulus string, err error) {
	resp, err := e.req.Post(e.URL+"/publiccombo/keyPair", cookie)
	if err != nil {
		return
	}
	keyMap := gjson.GetBytes(resp.Bytes(), "publicKeyMap")
	exponent = keyMap.Get("exponent").String()
	modulus = keyMap.Get("modulus").String()
	return
}

// 获取验证码
func (e *Ecard) getCodeImg(cookie *http.Cookie) (code string, err error) {
	resp, err := req.Get(e.URL+"/jcaptcha.jpg?"+strconv.FormatInt(time.Now().Unix(), 10), cookie)
	if err != nil {
		return
	}
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(resp.Bytes())
	code, _ = client.Text()
	return
}

// readCookies cookie from array
func readCookies(name string, cookies []*http.Cookie) *http.Cookie {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// Login 登陆系统
func (e *Ecard) Login() (err error) {
	if e.Cookie != nil {
		if ok, _ := e.IsCookieOverDue(); !ok {
			return
		}
	}
	resp, err := e.req.Get(e.URL + "/login")
	if err != nil {
		return err
	}
	// cookie like this
	// sid=8a3aeab3-16e8-48ac-a3d1-5c776cb638b3
	cookie := readCookies("sid", resp.Response().Cookies())

	// get the website rsa exponent and modulus
	exponent, modulus, err := e.getKeyMap(cookie)
	if err != nil {
		return
	}
	header := req.Header{
		"X-Requested-With": "XMLHttpRequest",
	}
	var state int64
OUT:
	for state != 3 {
		code, err := e.getCodeImg(cookie)
		if err != nil {
			return err
		}
		data := req.Param{
			"username":     utils.RsaEncrypt(exponent, modulus, e.Username),
			"password":     utils.RsaEncrypt(exponent, modulus, e.Password),
			"jcaptchacode": code,
		}
		resp, _ := e.req.Post(e.URL+"/login", data, header, cookie)
		state = gjson.GetBytes(resp.Bytes(), "ajaxState").Int()
		msg := gjson.GetBytes(resp.Bytes(), "msg").String()
		if state != 3 {
			fmt.Printf("%s %s\n", msg, e.Username)
		} else {
			// store the cookie into ecard
			e.Cookie = cookie
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
	if e.Cookie == nil {
		return true, nil
	}
	resp, err := e.req.Get(e.URL, e.Cookie)
	if err != nil {
		return true, err
	}
	reg := regexp.MustCompile("<title>(.*)</title>")
	ans := reg.FindSubmatch(resp.Bytes())
	return bytes.Equal(ans[1], []byte("智慧一卡通－登录")), nil
}

// GetCookie will return the cookie
func (e *Ecard) GetCookie() *http.Cookie {
	return e.Cookie
}
