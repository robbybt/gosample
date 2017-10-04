package hello

import (
	"expvar"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tokopedia/gosample/src/config"

	logging "gopkg.in/tokopedia/logging.v1"
)

type HelloWorldModule struct {
	cfg       *config.Config
	something string
	stats     *expvar.Int
}

func NewHelloWorldModule(cfgs *config.Config) *HelloWorldModule {
	// this message only shows up if app is run with -debug option, so its great for debugging
	logging.Debug.Println("hello init called", cfgs.Server.Name)

	return &HelloWorldModule{
		cfg:       cfgs,
		something: "John Doe",
		stats:     expvar.NewInt("rpsStats"),
	}

}

func (hlm *HelloWorldModule) SayHelloWorld(w http.ResponseWriter, r *http.Request) {
	hlm.stats.Add(1)
	w.Write([]byte("Hello " + hlm.something))
}

func (hlm *HelloWorldModule) SumNumber(w http.ResponseWriter, r *http.Request) {
	//buat baca parameter
	r.ParseForm()
	val1, _ := strconv.ParseInt(r.FormValue("angka1"), 10, 64)
	val2, _ := strconv.ParseInt(r.FormValue("angka2"), 10, 64)
	op := r.FormValue("operasi")
	var sum int64
	if op == "kali" {
		sum = val1 * val2
	} else if op == "tambah" {
		sum = val1 + val2
	} else if op == "bagi" {
		sum = val1 / val2
	}

	strSum := strconv.FormatInt(sum, 10)
	fmt.Println(strSum)
	w.Write([]byte("hasil: " + strSum))
}
