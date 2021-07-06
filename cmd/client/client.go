package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	nacos2 "go-micro-service/pkg/nacos"
	"go-micro-service/plugin/pkg/consul"
	"go-micro-service/plugin/pkg/logging"
	"net/http"
)

var logger log.Logger

func main() {

	logger = logging.SetLogging(logger, "", "info")

	svcHost := "localhost"
	svcPort := ":8500"
	_, err := consul.NewConsulClient("http://" + svcHost + svcPort)
	if err != nil {
		println(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/nacos/", nacos2.MakeHandler("GET", "/nacos", logger))

	handlers := make(map[string]string, 3)
	http.Handle("/", accessControl(mux, logger, handlers))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		println(err)
	}

	//tgt, err := url.Parse("http://localhost:8880")
	//client := httpTransport.NewClient("GET", tgt, getInfo)
	//getInfo := client.Endpoint()
	//_, _ = getInfo(context.Background(), nacos.NacosRequest{MetaData: svcHost})
}

//func getInfo(_ context.Context, httpRequest *http.Request, itf interface{}) error {
//
//	request := itf.(nacos.NacosRequest)
//
//}

func accessControl(h http.Handler, logger log.Logger, headers map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range headers {
			w.Header().Set(key, val)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Connection", "keep-alive")

		if r.Method == "OPTIONS" {
			return
		}

		//requestId := r.Header.Get("X-Request-Id")
		_ = level.Info(logger).Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)
		h.ServeHTTP(w, r)
	})

}
