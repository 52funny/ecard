package ecard

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

// ObtainBalance 获取饭卡余额
func (e *Ecard) ObtainBalance() string {
	resp, err := req.Get(e.URL + "/balance")
	if err != nil {
		fmt.Println(err)
	}
	reg := regexp.MustCompile("<p class=\"money\">(.*) <span>")
	money := reg.FindSubmatch(resp.Bytes())[1]
	return string(money)
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

// ObtainTodayBill 获取今天的消费记录
func (e *Ecard) ObtainTodayBill(size string) []Bill {
	param := req.Param{
		"startdealTime": time.Now().Format("2006-01-02"),
		"enddealTime":   time.Now().Format("2006-01-02"),
		"start":         "1",
		"end":           size,
		"size":          size,
	}
	reader, err := req.Post(e.URL+"/bill", param)
	if err != nil {
		fmt.Println(err)
	}
	billS := make([]Bill, 0)
	dom, err := goquery.NewDocumentFromReader(reader.Response().Body)
	dom.Find(".row tbody > tr").Each(func(i int, s *goquery.Selection) {
		time := s.Find(".text-muted").Text()
		content := s.Find(".time + td").Text()
		money, _ := strconv.ParseFloat(s.Find(".price").Text(), 10)
		balance, _ := strconv.ParseFloat(s.Find("td:last-child").Text(), 10)
		bill := Bill{
			Time:    time,
			Content: content,
			Money:   money,
			Balance: balance,
		}
		billS = append(billS, bill)
	})
	return billS
}
