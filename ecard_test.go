package ecard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func getEcard() *Ecard {
	f, err := os.Open("user.json")
	if err != nil {
		panic(err)
	}
	dataBytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	u := new(User)
	json.Unmarshal(dataBytes, u)

	e := &Ecard{
		URL:      "http://60.171.203.79:8090/easytong_portal",
		Username: u.Username,
		Password: u.Password,
	}
	return e
}
func TestLogin(t *testing.T) {
	e := getEcard()
	e.Login()
}

func TestDormitoryElectricity(t *testing.T) {
	e := getEcard()
	e.Login()
	fmt.Println(e.ObtainDormitoryElectricity("0", "5", "237"))
	fmt.Println(e.ObtainDormitoryElectricity("0", "5", "k237"))
}

func TestCookie(t *testing.T) {
	e := getEcard()
	fmt.Println(e.IsCookieOverDue())
}

func TestBalance(t *testing.T) {
	e := getEcard()
	e.Login()
	fmt.Println(e.ObtainBalance())
}

func TestPay(t *testing.T) {
	e := getEcard()
	e.Login()
	for i := 0; i < 5; i++ {
		result, err := e.PayElectricity("0", "5", "k237", "0.01")
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}
}

func TestBill(t *testing.T) {
	e := getEcard()
	e.Login()
	billS, err := e.ObtainTodayBill("2", "40")
	if err != nil {
		panic(err)
	}
	fmt.Println(billS)
}
