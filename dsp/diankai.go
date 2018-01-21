package dsp

import (
	"strconv"
	"crypto/md5"
	"net/url"
	"fmt"
	"encoding/json"
	"github.com/json-iterator/go/extra"
	"dsp_demo/util"
	"dsp_demo/model"
	"dsp_demo/converter"
	"time"
	"github.com/astaxie/beego/logs"
)

const (
	CPC = iota
	CPD
	CPM
	ALL
)

func init() {
	extra.RegisterFuzzyDecoders()
}

type DianKaiDsp struct {
	Dsp
}

var diankai_address = "http://sa.a.alldk.com/"
var default_ua = "Mozilla/5.0 (Linux; Android 5.1.1; Nexus 5 Build/LMY48B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.93 Mobile Safari/537.36"

type DianKaiApp struct {
	pid string
	ad_space string
	ad_type	int
	ad_category string
	package_name string
}

type DianKaiDevice struct {
	device_name string
	device_band string
	os string
	os_version string
	screen_width int
	screen_hight int
	sn string
	ua string
	sim string
	android_id string
	mac string
	is_tablet bool
}

type DianKaiNet struct {
	cell_ids string
	carrier string
	lat string
	lng string
	is_wifi bool
	ip string
	client int
	access_points string
}

type DianKaiVersion struct {
	Api_ver string  //2.2.0
	Token string  //md5(ad_space + pid + package_name)
}


type DianKaiRequest struct {
	App DianKaiApp `json:"app"`
	Device DianKaiDevice `json:"device"`
	Network DianKaiNet `json:"network"`
	Version DianKaiVersion `json:"version"`
}

type DianKaiAd struct {
	Mark string `json:"mark"` //unique id
	Ad_logo string `json:"ad_logo"`
	Ad_text string `json:"ad_text"`
	Icon []string `json:"icon"`
	Description []string `json:"description"`
	Pic string `json:"pic"`
	App_name string `json:"app_name"`
	Package_name string `json:"package_name"`
	Click_url string `json:"click_url"`
	Brand_name string `json:"brand_name"`
	Title string `json:"title"`
	Memo string `json:"memo"`
	Ad_category string `json:"ad_category"`
	Union string `json:"union"`
	S int `json:"s"`
	View_report []string `json:"view_report"`
	Click_report []string `json:"click_report"`
	Download_report []string `json:"download_report"`
}

type DianKaiResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Version string `json:"version"`
	Ads []DianKaiAd `json:"ads"`
}

//construct request string
func (r *DianKaiRequest) toString() string {
	req_url := diankai_address + "?"
	req_url += "pid=" + r.App.pid + "&ad_space=" + r.App.ad_space + "&ad_type=" + strconv.Itoa(r.App.ad_type) + "&package_name=" + r.App.package_name
	if r.Device.device_name != "" {
		req_url += "&device_name=" + r.Device.device_name
	}
	if r.Device.device_band != "" {
		req_url += "&device_brand=" + r.Device.device_band
	}
	if r.Device.os != "" {
		req_url += "&os=" + r.Device.os
	}
	if r.Device.os_version != "" {
		req_url += "&os_version=" + r.Device.os_version
	}
	req_url += "&screen_width=" + strconv.Itoa(r.Device.screen_width) + "&screen_height=" + strconv.Itoa(r.Device.screen_hight)
	ua := r.Device.ua
	if ua == "" {
		ua = default_ua
	}
	if r.Network.ip != "" {
		req_url += "&ip=" + r.Network.ip
	}
	if r.Device.android_id != "" {
		req_url += "&android_id=" + r.Device.android_id
	}
	req_url += "&ua=" + url.QueryEscape(ua)
	if r.Device.sn != "" {
		req_url += "&sn=" + r.Device.sn
	}
	if r.Device.mac != "" {
		req_url += "&mac=" + r.Device.mac
	}
	req_url += "&api_ver=" + r.Version.Api_ver + "&token=" + r.Version.Token

	return req_url
}

//convert req to diankai request
func ConvertRequestToDianKaiRequest(req *model.HMRequest, conf *model.DspConf) *DianKaiRequest {
	pid := conf.Token
	ad_space_id :=conf.SlotId
	tp := conf.Type

	tokenStr := ad_space_id + pid +  req.GetPackageName()
	logs.Info("TOKEN STR: ", tokenStr)
	token := md5.Sum([]byte(tokenStr))
	tokenMd5Str := fmt.Sprintf("%x", token)
	logs.Info("Token: ", tokenMd5Str)
	return &DianKaiRequest{
		App: DianKaiApp{
			pid: pid,
			ad_space: ad_space_id,
			ad_type: tp,
			package_name: req.GetPackageName(),
		},
		Device: DianKaiDevice{
			device_name: req.Device.Model,
			device_band: req.Device.Make,
			os: req.Device.OS,
			os_version: req.Device.Osv,
			screen_hight: int(req.Device.H),
			screen_width: int(req.Device.W),
			android_id: req.Device.AndroidId,
			mac: req.Device.ImeiMD5,
			sn: req.Device.ImeiMD5,
			is_tablet: false,
		},
		Network: DianKaiNet{
			ip: req.Device.Ip,
			is_wifi: req.Device.Connectiontype == 2,
			carrier: req.Device.Carrier,
		},
		Version: DianKaiVersion{
			Api_ver: "2.2.0",
			Token: tokenMd5Str,
		},

	}
}

//convert diankai response to response
func ConvertDiankaiResponseToNormalResponse(resp *DianKaiResponse, req *model.HMRequest) *model.HMResponse {
	var hmResp = &model.HMResponse{
		ID: req.ID,
		SeatBids: []model.HaoMaiSeatBid{
			model.HaoMaiSeatBid{
				Seat: "HaoMai",
			},
		},
	}

	var bids []model.HaoMaiBid
	if len(resp.Ads) == 0 {

	}
	for _, ad:=range resp.Ads {
		bid:=convertDiankaiAdsToHaoMaiBid(&ad, req)
		if bid == nil {
			continue
		}
		bids = append(bids, *bid)
	}

	hmResp.SeatBids[0].Bids = bids
	return hmResp
}

func convertDiankaiAdsToHaoMaiBid(ads *DianKaiAd, req *model.HMRequest) *model.HaoMaiBid{
	//not support
	if ads.S==1 {
		return nil
	}
	return &model.HaoMaiBid{
		ID: req.ID,
		Impid: req.Imps[0].ID,
		Price: req.Imps[0].Bidfloor,
		Adid: ads.Mark,
		Adm: model.Adm{
			Imgurl: ads.Pic,
			Landingurl: ads.Click_url,
			Title: ads.Title,
			Source: ads.Union,
			Description: ads.Memo,
		},
		Tagid: req.Imps[0].TagId,
		Billingtype: 1, //CPM
		Impurl: ads.View_report,
		Curl: ads.Click_report,
	}
}

func (this *DianKaiDsp)GetAds(req *model.HMRequest, conf *model.DspConf) *model.HMResponse {
	//1. convert normal request to DianKai ads
	diankai_req := ConvertRequestToDianKaiRequest(req, conf)
	//2. send get ads request
	req_url := diankai_req.toString()
	logs.Info(req_url)
	//3. parse DianKai response to normal response
	var data DianKaiResponse
	start := time.Now()
	util.HttpGetResponse(req_url, &data)
	cost:= time.Now().Sub(start)
	logs.Info("Real cost: ", cost)
	//for _,v := range rr {
	logs.Info(data)
	return ConvertDiankaiResponseToNormalResponse(&data, req)
	//4. return response or error
}

func DianKaiTT() {
	req := `{"id":"86cbf2b0ccf612d02ab61c3316cca533","imp":[{"id":"1","tagid":"1.3.a.1","nativead":{"request":"","w":1000,"h":500},"admtype":2,"bidfloor":100,"templates":[{"id":"2.2","width":1000,"height":500}]}],"app":{"name":"yidiannews","bundle":"com.yidian.org"},"device":{"ip":"106.39.75.132","make":"XIAOMI","model":"Mi3","os":"android","osv":"4.4.4","h":1920,"w":1080,"js":0,"language":"zh","connectiontype":2,"didmd5":"df54274061357627c9a5bd3d70a63610","dpidmd5":"46d02327677a49d76f28bb80750ccc7d","dpid":"d210c1e671045b5c"},"test":0,"at":2}`
	var dt model.MiRequest
	json.Unmarshal([]byte(req), &dt)

	hmRequest := converter.ConvertMiRequestToHMRequest(&dt)
	_, dspConfs := model.GetDspConfBySrcAndSlotId("xiaomi", hmRequest.Imps[0].TagId)


	hmb, _ :=json.Marshal(*hmRequest)
	logs.Info("HMRequest: ", string(hmb))

	var dsp = DianKaiDsp{
		Dsp: Dsp{
			Name: "diankai",
			Host: "",
		},
	}

	var tt IDsp = &dsp

	res := tt.GetAds(hmRequest, dspConfs[0])
	b, _ := json.Marshal(*res)
	logs.Info(string(b))
	//formatResponse()
}

func formatResponse() {
	resp := DianKaiResponse{
		Code: 0,
		Message: "ok",
		Version: "2.2.0",
		Ads: []DianKaiAd {
			 {
				Mark: "d13301b22c23b6dcc40ed4defc2ed6bf",
				Ad_logo: "https://cpro.baidustatic.com/cpro/ui/noexpire/img/2.0.1/bd-logo4.png",
				Ad_text: "https://cpro.baidustatic.com/cpro/ui/noexpire/img/mob_adicon.png",
				Icon: []string{},
				Pic: "http://ubmcmm.baidustatic.com/media/v1/0f000ZOV9pMJdVgTL0-Y6s.jpg",
				App_name: "贷款免费在线申请",
				Package_name: "",
				Click_url: "http://m.baidu.com/mobads.php?K60000K6ZiyXBnabrk8X4WsoyRdXO1xretBIu22DTnWHE-bZxHCTFWNGWr31mgsx7mYivAYi9p2UajujWOeDo9ILT6nVGS5gIlAJYm-2yOTc5okPXnzF03iooeGSxIEHyt2g2vBP6StJ2KTxYH50Q_MxaA7pz5DPnC8DgrcOLYfWHFJuM6.7Y_if_lTUQqRg-nePklUDfwGYsUJ3PE5E6CpXyPvap7Q7erQKdsRP5QGHTOKGm8WkqT7jHzs_lTUQqRH7XHT_HjrOkvucILWzs34-9h9mlX1F_L2.mLKGujY0uh-zTLwxIZF9uARqnWD0TvNWUv4bgLwzmyw-5Hcknjb0TZFEIh-8mvRqnfKWpgw45Hn3PfKWpAdb5HD0TvPCUyfqnfK1TjYk0ZI_5HD0mv99Uh4-UjYs0ZN1ugFxIZ-suHYs0AN3mv99UhI-gLw4TARqn0KVmgwWp7qWUvw-5H00mywWUA71T1YknsKYIHYdPjf3PH640ZwV5HcLnWcLP1ndr0KYujYs0Zw9TWYzPHm0TLPs5Hc0uh3qmHbLrjFWPHIxmLKz0ZwYTjY10A4YIZ0qnsKBTWYv0Aw-Ih-WuNqGujYz0AwYTjYz0ZNzUjY0mv6q0Au8TWYs0ZP85HRYuyD3mWw9mhmkPvFbmyD4nWf3Pvn4mHuBm1nduym3QfK9IAqVgv-b5H00UyNbpy7xIZ-suHYs0A4ETjYs0Aw-I-qYTjdmiR75HRbVHyb10AwWgvPdpyfq0ZPbp1YdQW08n0K1uAVxIWYY0AwWgLF-TNqspvTqmvqVQM-GuA-9UB43py7EUyb0myPWIjd9TAb0my4bThqGu7qGujdbnWDsm17-PWTknjfdmWNW0A4-I7qYXgK-gv-b5HD0mgKxpyfqn0K9TZKGujd9rHT3nhndPsKETdqGujYzQHf8P0KBTh78ujd3py7EUyb0Iv7sgvPEUvVGuHY0TLPsnWYs0APYp1Y0IvbqP6K9u7q15HI-nHf3nWP9njcvm10Ym1Kxn0K9TZKWmgw-gv-b5fKdThkxmgKsmv7YuNqGujY0UAF1gv-b5fKBIgPGUhN1TdqGujY0uA-1IZFGmLwxpyfq0A7WIA-vuNq9TZ0q0A7WIA-vuNqWmgw-uvqzXHYzP6KhpgF1I7qGIjY0TvNWUv4bgv-Y5fKspAq8uHYs0AP_pv7YIZ0qmv99ThI-0ZPCULTqn0KWUA-Wp1Yk0A4Y5Hc0ULKYgLw4TARqP0K9mLfqnfKkIAYqnHRkPHc4njm3r0Kzm1Y0Th78p1Ys0Z6qn0KsThqb5yu-uyf0mhIYTjYs0A-vIZcqn0K-TMKWUvw-5N0zP1n40A7WIA-vmgw-5H00uZPspyfqP6K1ThPxIvbqnHm0TAPvTWYznHTznjT0uvN8uANz5H00myI-5H00mLFW5HR3nWDk0A",
				Brand_name: "融360",
				Description: []string{"1-10万贷款在线申请,月入3000快速搞定贷款,利息0.18%"},
				Memo: "1-10万贷款在线申请,月入3000快速搞定贷款,利息0.18%",
				Title: "贷款免费在线申请",
				S: 0,
				Union: "baidu",
				View_report: []string{"http://wn.pos.baidu.com/adx.php?c=d25pZD03ZTE0ODIzYTAyNmMwNGMwXzAAcz03ZTE0ODIzYTAyNmMwNGMwAHQ9MTUxNTI5MDY4OABzZT0yAGJ1PTgAcHJpY2U9V2xHQVFBQUdDM0o3akVwZ1c1SUE4cnpjazVPVElzUGJ3RGtvVUEAY2hhcmdlX3ByaWNlPTQwMQBzaGFyaW5nX3ByaWNlPTQwMTAwMAB3aW5fZHNwPTgAY2htZD0xAGJkaWQ9AGNwcm9pZD0Ad2Q9MTY3OTQwMTEAcG9zPTAAYmNobWQ9MAB2PTEAaT01NmE1YzNiYQ&ext=cHJpY2U9NDImcGxhbmlkPTY1MDQxOA"},
				Click_report: []string{},
				Download_report: []string{},
				Ad_category: "c",
			},
		},
	}

	b,err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(b))
	fmt.Println("Content: ", string(b))
}
//pid:        64fbb58824c045f8dde149d170913db0

//首页推荐：   3e92ae868e9eb8468ebb608696e59889: 7

//文章页底部：  3ec87b897114a4b0ab200b4de59d091c

//slotid:

