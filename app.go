package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/tokopedia/gosample/src/deposit"
	"github.com/tokopedia/gosample/src/order"

	"github.com/google/gops/agent"

	"github.com/tokopedia/gosample/src/config"
	"github.com/tokopedia/gosample/src/hello"
	grace "gopkg.in/tokopedia/grace.v1"
	"gopkg.in/tokopedia/logging.v1"
)

func printWord(input string) (int, int, bool, string) {
	fmt.Println(input)
	return 1, 0, false, "sqwe"
}
func init() {
	fmt.Println("jalan pertama")
	//var1 := 1

	type datastruct struct {
		angka   int
		kalimat string
	}

	var temp datastruct
	temp.angka = 1
	temp.kalimat = "asd"

}
func main() {
	fmt.Println("jalan kedua")
	flag.Parse()
	logging.LogInit()

	debug := logging.Debug.Println

	debug("app started") // message will not appear unless run with -debug switch

	var configs *config.Config

	configs = config.InitConfig()

	if err := agent.Listen(agent.Options{
		ShutdownCleanup: true, // automatically closes on os.Interrupt
	}); err != nil {
		log.Fatal(err)
	}

	hwm := hello.NewHelloWorldModule(configs)
	dm := deposit.NewDepositModule(configs)
	od := order.NewOrderModule(configs)
	http.HandleFunc("/hello", hwm.SayHelloWorld)
	http.HandleFunc("/tambah", hwm.SumNumber)
	http.HandleFunc("/add", dm.AddDeposit)
	http.HandleFunc("/tugas", od.GetOrder)
	//go logging.StatsLog()
	log.Fatal(grace.Serve(":8080", nil))
}
