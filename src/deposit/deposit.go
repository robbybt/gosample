package deposit

import (
	"bytes"
	"errors"
	"log"
	"net/http"

	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tokopedia/gosample/src/config"
	logging "gopkg.in/tokopedia/logging.v1"
)

type DepositModule struct {
	cfg      *config.Config
	dbClient *sqlx.DB
}

func NewDepositModule(cfgs *config.Config) *DepositModule {
	// this message only shows up if app is run with -debug option, so its great for debugging
	logging.Debug.Println("hello init called", cfgs.Server.Name)

	connString := cfgs.Database["deposit"].Master
	depositDBClient, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatalln("failed to open connection db deposit : ", err)
	}
	return &DepositModule{
		cfg:      cfgs,
		dbClient: depositDBClient,
	}
}

func (dm *DepositModule) AddDeposit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, _ := strconv.Atoi(r.FormValue("id"))
	amount, _ := strconv.Atoi(r.FormValue("amount"))
	err := dm.doAddDeposit(userId, amount)
	if err != nil {
		log.Println(err.Error())
	}
	w.Write([]byte("tambah " + strconv.Itoa(amount)))
}

func (dm *DepositModule) doAddDeposit(userId int, amount int) error {
	if userId == 0 {
		return errors.New("please insert user id")
	}
	if amount == 0 {
		return errors.New("please enter positive amount")
	}

	buff := bytes.NewBufferString(`INSERT INTO ws_deposit (user_id,deposit_type,amount,notes,create_by) VALUES ($1,14,$2,'for devel purpose',$1)`)
	query := dm.dbClient.Rebind(buff.String())
	_, err := dm.dbClient.Exec(query, userId, amount)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
