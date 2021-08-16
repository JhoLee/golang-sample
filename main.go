package main

import (
	"flag"
	"net/http"

	"github.com/jholee/golang-sample/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

type User struct {
	Email    string
	Username string
}

func main() {
	var bind string
	flag.StringVar(&bind, "bind", "0.0.0.0:9100", "bind")
	flag.Parse()

	err := prometheus.Register(version.NewCollector("query_exporter"))
	if err != nil {
		log.Errorf("Failed to register collector.", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<a href='/metrics'>metrics</a>"))
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.HandlerFor(prometheus.Gatherers{
			prometheus.DefaultGatherer,
		}, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	})

	log.Info("Users: ")
	token := utils.GetToken("aaa", "aaa")
	users, err := utils.GetUsers("aaa", token)
	if err != nil {
		log.Errorf("Failed to get users. %v", err)
	}
	log.Info(users)

	log.Info("hI")
	email, err := utils.GetUserEmail("jhlee", "https://docker.siadev.kr", token)
	if err != nil {
		log.Panicf("Failed to get User's email. %v", err)
	}
	log.Info(email)



	//log.Infof("Starting http server - %s", bind)
	//if err := http.ListenAndServe(bind, nil); err != nil {
	//	log.Errorf("Failed to start http server: %s", err)
	//}

}
