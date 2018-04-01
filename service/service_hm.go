package service

import (
	"time"
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"encoding/json"
	"dsp_demo/model"
	"dsp_demo/converter"
	"dsp_demo/dsp"
	"dsp_demo/util"
	"github.com/astaxie/beego/logs"
)

var DispatchNumControl = make(chan bool, 1000)

type DSPService struct {
	*Service
}

func Limit() bool {
	select {
	case <-time.After(time.Millisecond * 20):
		logs.Warn("Queue full, the server is busy now!")
		return false
	case DispatchNumControl <- true:
	// 任务放入任务队列channal
		return true
	}
}

func (s *DSPService) AdsHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	c := make(chan interface{})
	//get token
	if Limit() {
		//release the token
		defer func() {
			<-DispatchNumControl
		}()

		go s.HandleAdsRequest(r, c)

		for {
			select {
			case <-time.After(1000 * time.Millisecond):
				logs.Info("HandleAdsRequest Timeout!")
				close(c)
				goto FOREND
			case re, ok := <-c:
				if ok {
					close(c)
					if re == nil {
						w.WriteHeader(http.StatusNoContent)
						return
					}
					b, err := json.Marshal(re)
					if err != nil {
						util.JsonStrResult(500, "Json encode failed: " + err.Error(), w)
					} else {
						util.JsonResult(http.StatusOK, b, w)
					}
					return
				} else {
					log.Println("Process failed!", ok)
					goto FOREND
				}
			}
		}

		FOREND:
			util.JsonStrResult(500, "Internal Timeout!", w)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *DSPService) ExposeHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Expose received!")
}

func (s *DSPService) ClickHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Click received!")
}

func (s *DSPService) ConfHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		result, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logs.Error("Read config data failed!", err)
			panic(err)
		}
		temp := new(model.HMConf)
		err = json.Unmarshal(result, &temp)
		if err != nil {
			logs.Error("decode config data failed!", err)
			panic(err)
		}
		if temp == nil || len(temp.Templates) == 0 || len (temp.Traffic) == 0 {
			logs.Warn("Config May be wrong!", string(result))
			util.JsonStrResult(500, "Config update failed due to wrong config, please check!", w)
		}else {
			err := ioutil.WriteFile(model.GetConfigPath(), result, 0644)
			if err != nil {
				logs.Error("update local disk config failed!", err)
				util.JsonStrResult(500, "update local disk config failed! " + err.Error(), w)
			} else {
				model.SetConfig(temp)
				util.JsonStrResult(200, "Config update succeed!", w)
			}
		}
		return
	}
	//return the conf
	b, err := ioutil.ReadFile("conf/app.conf")
	if err != nil {
		logs.Error("load config error: ", err)
		//return false
		util.JsonStrResult(404, err.Error(), w)
		return
	}

	util.JsonResult(200, b, w)
}

func (s *DSPService) HandleAdsRequest(r *http.Request, ch chan interface{}) {
	s.WaitGroup.Add(1)
	defer func() {
		if err:=recover();err!=nil{
			logs.Error("Panic recover: ", err) // 这里的err其实就是panic传入的内容，55
		}
	}()

	defer s.WaitGroup.Done()

	r.ParseForm()
	var src = "mi"
	srcs, ok := r.Form["src"]
	if ok {
		src = srcs[0]
	}

	logs.Info("Traffic Src: ", src)

	if src=="mi" {
		var req model.MiRequest
		result, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Read data failed!")
			panic(err)
		}
		if err := json.Unmarshal(result, &req); err == nil {
			log.Println("Parse request succeed!  %v", req)
			hmRequest := converter.ConvertMiRequestToHMRequest(&req)
			_, dspConfs := model.GetDspConfBySrcAndSlotId("xiaomi", hmRequest.Imps[0].TagId)
			start := time.Now()
			dsp := dsp.GetDsp(dspConfs[0].Name)

			re := dsp.GetAds(hmRequest, dspConfs[0])
			var dt *model.MiResponse
			if re!=nil {
				dt = converter.ConvertHaoMaiResponseToMiResponse(re)
			}
			cost := time.Now().Sub(start)
			logs.Info("Cost: ", cost)

			if err != nil {
				logs.Error("Unmarshal failed", err)
				close(ch)
			} else {
				ch <- &dt
			}
		} else {
			logs.Error("Parse request failed!", err)
		}
	} else {
		logs.Warn("Found unknown src: ", src)
		close(ch)
	}
}
