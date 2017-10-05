package order

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	nsq "github.com/nsqio/go-nsq"
	"github.com/tokopedia/gosample/src/config"
	"github.com/tokopedia/gosample/src/product"
	"github.com/tokopedia/gosample/src/redis"
	"github.com/tokopedia/gosample/src/shop"
	logging "gopkg.in/tokopedia/logging.v1"
)

type OrderModule struct {
	cfg      *config.Config
	dbClient *sqlx.DB
}

type ResponseData struct {
	Data Order
}
type Order struct {
	List []OrderDesc
}
type OrderDesc struct {
	Order_id        int64       `db:"order_id" json:"order_id"`
	Invoice_ref_num string      `db:"invoice_ref_num"`
	Order_status    int         `db:"order_status"`
	Shop            shop.Shop   `db:"shop_id" json:"shop"`
	Order_detail    OrderDetail `db:"order" json:"order_detail"`
}

type OrderDetail struct {
	Quantity_delivered int             `db:"quantity_deliver"`
	Product            product.Product `json:"product"`
}

type OrderData struct {
	Order_id           int64  `db:"order_id""`
	Invoice_ref_num    string `db:"invoice_ref_num"`
	Order_status       int    `db:"order_status"`
	Shop_id            int    `db:"shop_id"`
	Quantity_delivered int    `db:"quantity_deliver"`
	Product_id         int    `db:"product_id"`
}

func NewOrderModule(cfgs *config.Config) *OrderModule {
	// this message only shows up if app is run with -debug option, so its great for debugging
	logging.Debug.Println("hello init called", cfgs.Server.Name)

	connString := cfgs.Database["order"].Master
	orderDBClient, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatalln("failed to open connection db deposit : ", err)
	}
	return &OrderModule{
		cfg:      cfgs,
		dbClient: orderDBClient,
	}
}

func (od *OrderModule) GetOrder(w http.ResponseWriter, r *http.Request) {
	user := redis.GetRedis("tkpd:5Ms_gkx1AcUoSMok4me7qG0CY3aeO90_qgVAeSOzqoWQnzLk1a7YlgAttb14y5qT_APLGQtNmalvcClUzlMvykND0lgCHiLy0YPTdR8jq-vOxv9qj1eeDmXGQQMf2vqT:session_json")
	log.Println("user : ", user)
	if user == nil {
		w.WriteHeader(403)
		return
	}
	pq := `select w.order_id,w.invoice_ref_num,w.order_status,w.shop_id,wd.quantity_deliver,wd.product_id
			from ws_order w
			join ws_order_dtl wd
			on w.order_id=wd.order_id
			limit 10`

	//query := od.dbClient.Rebind(pq)
	var orderData []OrderData
	err := od.dbClient.Select(&orderData, pq)
	if err != nil {
		fmt.Println(err)
	}
	var order []OrderDesc
	for _, v := range orderData {
		var temp OrderDesc
		temp.Order_id = v.Order_id
		temp.Invoice_ref_num = v.Invoice_ref_num
		temp.Order_status = v.Order_status
		temp.Shop = shop.GetShop(strconv.Itoa(v.Shop_id))
		temp.Order_detail.Quantity_delivered = v.Quantity_delivered
		temp.Order_detail.Product = product.GetProduct(strconv.Itoa(v.Product_id))
		order = append(order, temp)
	}
	//fmt.Println(order)
	var response ResponseData
	var resOrder Order
	resOrder.List = order
	response.Data = resOrder
	res, err := json.Marshal(response)
	publisNSQ(res)

	if err != nil {
		fmt.Println(err)
	}
	w.Write(res)
}

func publisNSQ(response []byte) {
	cfg := config.InitConfig()
	producer, err := nsq.NewProducer(cfg.NSQ.NSQD, nsq.NewConfig())
	if err != nil {
		panic(err)
	}
	err = producer.Publish("TestTraining", response)
	if err != nil {
		log.Println("error :", err)
	}
}
