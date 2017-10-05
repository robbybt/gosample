package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Product struct {
	ProductId           int    `json:"product_id"`
	ProductName         string `json:"product_name"`
	Product_description string `json:"product_description"`
	Product_status      int    `json:"product_status"`
}

type ProductRawData struct {
	Data InfoData `json:"data"`
}

type InfoData struct {
	Info Product `json:"info"`
}

const (
	Product_URL = "http://devel-go.tkpd:3002/v1/web-service/product/get_detail"
)

func GetProduct(id string) Product {
	url := fmt.Sprintf("%s?product_id=%s", Product_URL, id)
	//log.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
	}
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	//log.Println(string(body))
	var product ProductRawData

	err = json.Unmarshal(body, &product)
	if err != nil {
		log.Println(err.Error())
	}
	return product.Data.Info
}
