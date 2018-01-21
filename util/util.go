package util

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"bytes"
	"github.com/astaxie/beego/logs"
	"os"
)

func HttpGetResponse(url string, obj interface{}) {
	resp, err := http.Get(url)
	if err != nil {
		panic("Error occurred during Http Get with url: " + url + ", Error: " + err.Error())
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	logs.Debug(string(b))
	err = json.Unmarshal(b, obj)

	if err != nil {
		logs.Error("GetResponse failed!", err)
		panic(err)
	}
}

func HttpPostJsonResponse(url string, data []byte, obj interface{}) {
	client := &http.Client{}
	req_data := bytes.NewBuffer(data)
	request, _ := http.NewRequest("POST", url, req_data)
	request.Header.Set("Content-type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(b, obj)

	if err != nil {
		logs.Error("PostJsonResponse failed!", err)
		panic(err)
	}
}

func JsonEncodeWithoutEscape(v interface{}) []byte{
	buffer := bytes.NewBuffer(make([]byte, 0))
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		logs.Warn("encode message failed!", err)
	}
	return buffer.Bytes()
}

func JsonDecode(b []byte, v interface{}) {
	json.Unmarshal(b, v)
}

func JsonResult(code int, data []byte, w http.ResponseWriter) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func JsonStrResult(code int, data string, w http.ResponseWriter) {
	var dt = map[string]interface{} {
		"status": code,
		"msg": data,
	}

	d, err := json.Marshal(dt)
	if err != nil {
		panic(err)
	}
	JsonResult(code, d, w)
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}