package ecard

import (
	"regexp"

	"github.com/imroc/req"
)

// PayElectricity 支付电费
func (e *Ecard) PayElectricity(areaNo string, buildingNo string, roomNo string, payMoney string) (string, error) {
	resp, err := req.Post(e.URL+"/payFee/showItemListPayPage", req.Param{
		"itemNum": "1",
	})
	if err != nil {
		return "", err
	}
	// fmt.Println(resp.String())
	reg := regexp.MustCompile("<input type=\"hidden\" id=\"itemNum\" value=\"(.*)\">")
	itemNum := string(reg.FindSubmatch(resp.Bytes())[1])

	reg = regexp.MustCompile("<input type=\"hidden\" id=\"typeNum\" value=\"(.*)\">")
	typeNum := string(reg.FindSubmatch(resp.Bytes())[1])

	reg = regexp.MustCompile("<input type=\"hidden\" id=\"itemName\" value=\"(.*)\">")
	itemName := string(reg.FindSubmatch(resp.Bytes())[1])

	reg = regexp.MustCompile("<li  onclick=\"compareItemMoney\\(1\\);\" balance=(.*) cardAccNum =\"(.*)\"  moneyMin =\"0\"  id=\"(.*)\">(.*)</li>")
	wallet := reg.FindSubmatch(resp.Bytes())
	eWalletMoney := string(wallet[1])
	cardAccNum := string(wallet[2])
	eWalletID := string(wallet[3])
	eWalletName := string(wallet[4])

	reg = regexp.MustCompile("<input type=\"hidden\" name=\"token\"  id =\"token\" value=\"(.*)\">")
	token := string(reg.FindSubmatch(resp.Bytes())[1])

	param := req.Param{
		"itemNum":      itemNum,
		"eWalletId":    eWalletID,
		"typeNum":      typeNum,
		"itemName":     itemName,
		"payMenoy":     payMoney,
		"ewalletMenoy": eWalletMoney,
		"eWalletName":  eWalletName,
		"areano":       areaNo,
		"buildingno":   buildingNo,
		"roomno":       roomNo,
		// "buildingname": buildingNo + "号楼",
		"cardaccNum": cardAccNum,
		"token":      token,
	}
	resp, err = req.Post(e.URL+"/payFee/payItemlist", param)
	if err != nil {
		return "", nil
	}
	reg = regexp.MustCompile("<p.*><strong>(.*)</strong></p.*>")
	result := string(reg.FindSubmatch(resp.Bytes())[1])

	return result, nil
}
