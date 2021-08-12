package ecard

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

// ObtainBalance 获取饭卡余额
func (e *Ecard) ObtainBalance() (string, error) {
	resp, err := e.req.Get(e.URL+"/balance", e.Cookie)
	if err != nil {
		return "", err
	}
	reg := regexp.MustCompile("<p class=\"money\">(.*) <span>")
	money := reg.FindSubmatch(resp.Bytes())[1]
	return string(money), nil
}

// ObtainDormitoryElectricity 获取寝室电费余额
func (e *Ecard) ObtainDormitoryElectricity(areaNo string, buildingNo string, roomNo string) (string, error) {
	data := req.Param{
		"data": `{"itemNum":"1","areano":"` + areaNo + `","buildingno":"` + buildingNo + `","roomno":"` + roomNo + `"}`,
	}
	header := req.Header{
		"X-Requested-With": "XMLHttpRequest",
	}
	resp, err := e.req.Post(e.URL+"/payFee/getBalance", data, header, e.Cookie)
	if err != nil {
		return "", err
	}
	return gjson.GetBytes(resp.Bytes(), "feeDate.balance").String(), err
}

// ObtainIntervalBill 获取时间区间的账单信息
// typeFlag 1 消费 2 充值 3 补助 4 互转
// size 页面大小
// startTime 开始时间
// endTime 结束时间
// 时间格式应为 2006-01-02类型
func (e *Ecard) ObtainIntervalBill(typeFlag, size, startTime, endTime string) ([]Bill, error) {
	param := req.Param{
		"typeFlag":      typeFlag,
		"startdealTime": startTime,
		"enddealTime":   endTime,
		"start":         "1",
		"end":           size,
		"size":          size,
	}
	reader, err := e.req.Post(e.URL+"/bill", param, e.Cookie)
	if err != nil {
		return nil, err
	}
	billS := make([]Bill, 0)
	dom, err := goquery.NewDocumentFromReader(reader.Response().Body)
	if err != nil {
		return nil, err
	}
	dom.Find(".row tbody > tr").Each(func(i int, s *goquery.Selection) {
		time := s.Find(".text-muted").Text()
		content := s.Find(".time + td").Text()
		merchant := s.Find("td:nth-child(4)").Text()
		location := s.Find("td:nth-child(5)").Text()
		money, _ := strconv.ParseFloat(s.Find("td:nth-child(6)").Text(), 64)
		var balance float64
		if strings.Compare(typeFlag, "0") == 0 {
			balance, _ = strconv.ParseFloat(s.Find("td:last-child").Text(), 64)
		} else if typeFlag == "1" || typeFlag == "2" {
			balance, _ = strconv.ParseFloat(s.Find("td:nth-last-child(2)").Text(), 64)
		}
		bill := Bill{
			Time:     time,
			Content:  content,
			Merchant: merchant,
			Location: location,
			Money:    money,
			Balance:  balance,
		}
		billS = append(billS, bill)
	})
	return billS, nil

}

// ObtainTodayBill 获取今天的消费记录
// typeFlag 1 消费 2 充值 3 补助 4 互转
func (e *Ecard) ObtainTodayBill(typeFlag string, size string) ([]Bill, error) {
	return e.ObtainIntervalBill(typeFlag, size, time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"))
}
