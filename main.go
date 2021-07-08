package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

var (
	serverAddr = flag.String("server", ":80", "Server addr (ip:port)")
)

func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		Handler(response, request)
	}))
	err := http.ListenAndServe(*serverAddr, nil)
	if err != nil {
		panic(err)
	}
}
func Handler(response http.ResponseWriter, request *http.Request) {
	domain := request.URL.Query().Get("domain")
	if domain == "" {
		write(response, result(404, "domain is null", nil))
		return
	}
	bt, err := find(domain)
	if err != nil {
		write(
			response,
			result(500, fmt.Sprintf("find subdomain error, error message:[%s]", err.Error()), nil),
		)
		return
	}
	write(
		response,
		result(200, "success", string(bt)),
	)
	return
}

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func result(code int, message string, data interface{}) (rt []byte) {
	r := Result{
		Code:    code,
		Message: message,
		Data:    data,
	}
	rt, _ = json.Marshal(&r)
	return
}

func write(w http.ResponseWriter, rt []byte) {
	_, _ = w.Write(rt)
	w.WriteHeader(200)
}
