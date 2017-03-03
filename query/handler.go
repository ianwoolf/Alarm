package query

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type AlertData struct {
	ID       string        `json:"id"`
	Message  string        `json:"message"`
	Details  string        `json:"details"`
	Time     time.Time     `json:"time"`
	Duration time.Duration `json:"duration"`
	Level    string        `json:"level"`
	Data     Result        `json:"data"`
}

type Result struct {
	// StatementID is just the statement's position in the query. It's used
	// to combine statement results if they're being buffered in memory.
	StatementID int `json:"-"`
	Series      Rows
	Messages    []*Message
	Err         error
}

type Message struct {
	Level string `json:"level"`
	Text  string `json:"text"`
}

type Rows []*Row

type Row struct {
	Name    string            `json:"name,omitempty"`
	Tags    map[string]string `json:"tags,omitempty"`
	Columns []string          `json:"columns,omitempty"`
	Values  [][]interface{}   `json:"values,omitempty"`
}

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

	var alertData AlertData
	if err = json.Unmarshal(body, &alertData); err != nil {
		fmt.Println("json unmarshal error:", err.Error())
	}
	fmt.Printf("%+v", alertData)
	// just return the origin influxdb rs
	resp.Header().Add("Content-Type", "application/json")
	succResp(resp, "OK", alertData)
	resp.WriteHeader(200)
}
