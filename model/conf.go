package model

import (
	"sync"
	"os"
	"os/signal"
	"syscall"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"dsp_demo/util"
	"github.com/astaxie/beego/logs"
)

type DspConf struct {
	Name string `json:"name"`
	Token string `json:"token"`
	SlotId string `json:"slot_id"`
	Type int  `json:"type"`
}

type TemplateStyle struct {
	Type string `json:"type"`
	H int32 `json:"h"`
	W int32 `json:"w"`
} 

type TemplateConf struct {
	Id int `json:"id"`
	Style TemplateStyle `json:"style"`
	//Dsp []DspConf `json:"dsp"`
	Dsps map[string]DspConf `json:"dsps"`
}

type TrafficTemplate struct {
	Id int `json:"id"`
	Dsp []string `json:"dsp"`
}

type TrafficConf struct {
	SlotId string `json:"slot_id"`
	SspTemplateId string `json:"ssp_template_id"`
	Price int64 `json:"price"`
	Templates []*TrafficTemplate `json:"templates"`
}

type TrafficSSP struct {
	Src string `json:"src"`
	Conf []TrafficConf `json:"conf"`
}

type HMConf struct {
	Templates []TemplateConf `json:"templates"`
	Traffic []TrafficSSP `json:"traffic"`
}

var (
	config *HMConf
	idToTemplates map[int]TemplateConf
	idToTrafficSsp map[string]map[string]TrafficConf
	configLock = new(sync.RWMutex)
)

func GetConfigPath() string {
	var AppPath string
	var err error
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	appConfigPath := filepath.Join(workPath, "conf", "app.conf")
	if !util.FileExists(appConfigPath) {
		logs.Info("File not exist! ", config )
		appConfigPath = filepath.Join(AppPath, "conf", "app.conf")
	}

	logs.Info("Config Path: ", appConfigPath)
	return appConfigPath
}

func loadConfig() bool {
	//path := GetConfigPath()
	f, err := ioutil.ReadFile("conf/app.conf")
	if err != nil {
		logs.Error("load config error: ", err)
		return false
	}
	logs.Info("Loaded Config: ", string(f))
	//不同的配置规则，解析复杂度不同
	temp := new(HMConf)
	err = json.Unmarshal(f, &temp)
	logs.Info("Loaded config object: ", temp)
	if err != nil {
		logs.Error("Para config failed: ", err)
		return false
	}

	SetConfig(temp)
	return true
}

func GetConfig() *HMConf {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func SetConfig(temp *HMConf) {
	logs.Info("Config updated...")
	var tempIdTemplateMap = map[int]TemplateConf{}
	var tempIdTrafficMap  = map[string]map[string]TrafficConf{}

	if temp.Templates != nil {
		for _, v := range temp.Templates {
			tempIdTemplateMap[v.Id] = v
		}
	}
	if temp.Traffic != nil {
		for _, v := range temp.Traffic {
			var tt = map[string]TrafficConf{}
			if len(v.Conf) > 0 {
				for _, conf := range v.Conf {
					tt[conf.SlotId] = conf
				}
			}
			tempIdTrafficMap[v.Src] = tt
		}
	}

	configLock.Lock()
	config = temp
	idToTemplates = tempIdTemplateMap
	idToTrafficSsp = tempIdTrafficMap
	configLock.Unlock()
}


func GetDspConfBySrcAndSlotId(src, slot_id string) (*TemplateStyle, []*DspConf){
	configLock.RLock()
	defer configLock.RUnlock()
	tf := idToTrafficSsp[src][slot_id]
	template := tf.Templates[0]
	cnt := len(template.Dsp)
	templateConf, ok := idToTemplates[template.Id]
	if !ok {
		logs.Error("Traffic Conf not found! Dsp: %s, SlotId: %s", src, slot_id)
		return nil, nil
	}
	var chosenDsps = make([]*DspConf, cnt, cnt)
	for i,v := range template.Dsp{
		t := templateConf.Dsps[v]
		chosenDsps[i] = &t
	}
	style := templateConf.Style
	return &style, chosenDsps
}

func GetTrafficConfBySrcAndSlotId(src, slot_id string) *TrafficConf {
	configLock.RLock()
	defer configLock.RUnlock()
	tf, ok := idToTrafficSsp[src][slot_id]
	if ok {
		return &tf
	}
	logs.Warn("Traffic Conf not found! Dsp: %s, SlotId: %s", src, slot_id)
	return nil
}

func init() {
	if !loadConfig() {
		os.Exit(1)
	}

	//热更新配置可能有多种触发方式，这里使用系统信号量sigusr1实现
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR1)
	go func() {
		for {
			<-s
			logs.Info("Reloaded config:", loadConfig())
		}
	}()
}


