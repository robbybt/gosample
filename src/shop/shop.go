package shop

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ShopRawData struct {
	Data []ShopData `json:"data"`
}
type ShopData struct {
	Shop_id     int    `db:"shop_id" json:"shop_id"`
	Shop_name   string `db:"shop_name" json:"shop_name"`
	Shop_domain string `db:"shop_domain" json:"domain"`
	Shop_status int    `db:"shop_status" json:"status"`
}

const (
	Shop_URL = "http://devel-go.tkpd:3002/v1/shop/get_summary"
)

type Shop struct {
	Shop_id     int    `db:"shop_id" json:"shop_id"`
	Shop_name   string `db:"shop_name" json:"shop_name"`
	Shop_domain string `db:"shop_domain" json:"shop_domain"`
	Shop_status int    `db:"shop_status" json:"shop_status"`
}

func GetShop(id string) Shop {
	url := fmt.Sprintf("%s?shop_id=%s", Shop_URL, id)
	//log.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
	}
	//log.Println(string(body))
	var shopData ShopRawData
	err = json.Unmarshal(body, &shopData)
	if err != nil {
		log.Println(err.Error())
	}

	var result Shop
	result.Shop_domain = shopData.Data[0].Shop_domain
	result.Shop_id = shopData.Data[0].Shop_id
	result.Shop_name = shopData.Data[0].Shop_name
	result.Shop_status = shopData.Data[0].Shop_status
	return result
}
