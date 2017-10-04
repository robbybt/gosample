package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/tokopedia/gosample/src/deposit"

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
	configs = InitConfig()

	if err := agent.Listen(agent.Options{
		ShutdownCleanup: true, // automatically closes on os.Interrupt
	}); err != nil {
		log.Fatal(err)
	}

	hwm := hello.NewHelloWorldModule(configs)
	dm := deposit.NewDepositModule(configs)
	http.HandleFunc("/hello", hwm.SayHelloWorld)
	http.HandleFunc("/tambah", hwm.SumNumber)
	http.HandleFunc("/add", dm.AddDeposit)
	go logging.StatsLog()

	log.Fatal(grace.Serve(":9000", nil))
	fmt.Println(printWord("test123"))
}

func InitConfig() *config.Config {
	var cfg config.Config

	ok := logging.ReadModuleConfig(&cfg, "config", "hello") || logging.ReadModuleConfig(&cfg, "files/etc/gosample", "hello")
	if !ok {
		// when the app is run with -e switch, this message will automatically be redirected to the log file specified
		log.Fatalln("failed to read config")
	}

	return &cfg
}
