package query

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// @desc get measurement tags from influxdb deps on ns name
// @router /tags [get]
func postDataHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" && req.Method != "POST" {
		errResp(resp, http.StatusMethodNotAllowed, "Get or POST please!")
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("read body fail")
	} else {
		fmt.Println("read body:", string(body))
	}
	// just return the origin influxdb rs
	resp.Header().Add("Content-Type", "application/json")
	succResp(resp, "OK", "{\"data\":\"data\"}")
	resp.WriteHeader(200)
}
