package nsq

import (
	"log"
	"os"

	"github.com/tokopedia/gosample/src/config"

	"github.com/nsqio/go-nsq"
	"github.com/tokopedia/gosample/src/redis"
	logging "gopkg.in/tokopedia/logging.v1"
)

type NSQModule struct {
	cfg *config.Config
	q   *nsq.Consumer
}

func NewNSQModule() *NSQModule {

	var cfg *config.Config

	cfg = config.InitConfig()

	// this message only shows up if app is run with -debug option, so its great for debugging
	logging.Debug.Println("nsq init called", cfg.Server.Name)

	// contohnya: caranya ciptakan nsq consumer
	nsqCfg := nsq.NewConfig()
	q := createNewConsumer(nsqCfg, "TestTraining", "test", handler)
	q.SetLogger(log.New(os.Stderr, "nsq:", log.Ltime), nsq.LogLevelError)
	q.ConnectToNSQLookupd(cfg.NSQ.Lookupd)

	return &NSQModule{
		cfg: cfg,
		q:   q,
	}

}

func handler(msg *nsq.Message) error {
	redis.SetRedis(msg.Body)
	//log.Println("got message :", string(msg.Body))
	msg.Finish()
	return nil
}

func createNewConsumer(nsqCfg *nsq.Config, topic string, channel string, handler nsq.HandlerFunc) *nsq.Consumer {
	q, err := nsq.NewConsumer(topic, channel, nsqCfg)
	if err != nil {
		log.Fatal("failed to create consumer for ", topic, channel, err)
	}
	q.AddHandler(handler)
	return q
}
