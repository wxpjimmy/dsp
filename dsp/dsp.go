package dsp

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"errors"
	"dsp_demo/model"
	"github.com/astaxie/beego/logs"
)

type Dsp struct {
	Name string
	Host string
}

type IDsp interface {
	GetAds(req *model.HMRequest, conf *model.DspConf) *model.HMResponse
}

var dspMap map[string]IDsp

func init() {
	dspMap = make(map[string]IDsp)
	dspMap["diankai"] = &DianKaiDsp{
		Dsp{
			Name: "diankai",
			Host: "",
		},
	}
}

func GetDsp(name string) IDsp {
	re, ok := dspMap[name]
	if ok {
		return re
	}
	logs.Error("Found unknown dsp: ", name)
	return nil
}

func GetAdsInternal(dsp Dsp, req []byte, ch chan model.MiResponse) error {
	req_new := bytes.NewBuffer(req)
	response, _ := http.Post(dsp.Host, "application/json", req_new)
	var resp model.MiResponse

	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		logs.Info(string(body))
		bodyByte := []byte(string(body))

		err := json.Unmarshal(bodyByte, &resp)
		if err != nil {
			logs.Error("parse response failed!", err)
			return err
		} else {
			logs.Info("Parse succeed! %v", resp)
			ch<-resp
			return nil
		}
	}

	logs.Error(response.StatusCode)
	err := errors.New(string(response.StatusCode))
	return err
}



func GetResponse(dsp Dsp) (model.MiResponse, error) {
	client := &http.Client{}
	req := `{"id":"86cbf2b0ccf612d02ab61c3316cca533","imp":[{"id":"1","tagid":"1.3.a.1","nativead":{"request":"","w":1000,"h":500},"admtype":2,"bidfloor":100,"templates":[{"id":"2.2","width":1000,"height":500}]}],"app":{"name":"yidiannews","bundle":"com.yidian.org"},"device":{"ip":"106.39.75.132","make":"XIAOMI","model":"Mi3","os":"android","osv":"4.4.4","h":1920,"w":1080,"js":0,"language":"zh","connectiontype":2,"didmd5":"df54274061357627c9a5bd3d70a63610","dpidmd5":"46d02327677a49d76f28bb80750ccc7d","dpid":"d210c1e671045b5c"},"test":0,"at":2}`
	req_new := bytes.NewBuffer([]byte(req))
	request, _ := http.NewRequest("POST", dsp.Host, req_new)
	request.Header.Set("Content-type", "application/json")
	response, _ := client.Do(request)
	var resp model.MiResponse

	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		logs.Info(string(body))
		bodyByte := []byte(string(body))

		err := json.Unmarshal(bodyByte, &resp)
		if err != nil {
			logs.Error("parse response failed!", err)
			return resp, err
		} else {
			logs.Info("Parse succeed! %v", resp)
			return resp, nil
		}
	}

	logs.Error(response.StatusCode)
	err := errors.New(string(response.StatusCode))
	return resp, err
}
