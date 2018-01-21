package converter

import (
	"fmt"
	"testing"
	"dsp_demo/model"
	"encoding/json"
)

func TestRequest(test *testing.T) {
	var m model.Request

	fmt.Println("start parse~!")
	b := []byte(`{"id":"86cbf2b0ccf612d02ab61c3316cca533","imp":[{"id":"1","tagid":"1.3.a.1","nativead":{"request":"","w":1000,"h":500},"admtype":2,"bidfloor":100,"templates":[{"id":"2.2","width":1000,"height":500}]}],"app":{"name":"yidiannews","bundle":"com.yidian.org"},"device":{"ip":"106.39.75.132","make":"XIAOMI","model":"Mi3","os":"android","osv":"4.4.4","h":1920,"w":1080,"js":0,"language":"zh","connectiontype":2,"didmd5":"df54274061357627c9a5bd3d70a63610","dpidmd5":"46d02327677a49d76f28bb80750ccc7d","dpid":"d210c1e671045b5c"},"test":0,"at":2}`)
	err := json.Unmarshal(b, &m)
	if err != nil {
		test.Error("Unmarshal failed!")
	}
	if m.Device.AndroidId != "d210c1e671045b5c" {
		test.Error("Unmarshal wrong message")
	}
}

func TestResponse(test *testing.T) {
	var m model.MiResponse
	b := []byte(`{"id":"86cbf2b0ccf612d02ab61c3316cca533","bidid":"22_056756c9077649f1b80bdb5c81901cf6_1","seatbid":[{"cm":0,"group":0,"seat":"jingdong","bid":[{"id":"22_056756c9077649f1b80bdb5c81901cf6_11","impid":"1","tagid":"1.3.a.1","domain":"https://www.jd.com/","crid":"2720","cat":[20],"adid":"2720_0_138403020","price":186.7345869541168,"w":480,"h":320,"landingurl":"http://ccc-x.jd.com/dsp/nc?ext=aHR0cDovL21laS53b2NoYXcuY29tL3Byb2R1Y3RzL3h5bWo5OTlfMTQvMzAwMWhsc2owMDU&log=B4ZHNweusEbAOvPYea_kJBU3CaEJozHl0atEgL2tR1Ej5tG7ZZV0XR6NeoI4gpUHZzauuUwnOYW9eOhgkUoaZR2jOhHsORFqaNkE0smAQFgg_ewhbwO6tCtkr4bA-FotftPn-YoNniG-Ln-2Ra56C6BSGtydhaEkCRx9lyjQkh-WnP8iWFXlHaL8LhajJhYRvrKCWAly3ye2vnYlqzJN-2Okyh-rgxl5mS1Sc4ue0MztuE4NXDUZinTCUGT_rTRdRywwBz14mZ5ZWf9qSC77RLgGbQK4lcIrs6JJh6biTUgo8P44LqEQQQHHqccmxgZ-R2OQqC8Eqj4UhBzwRo-s23Co4-OOwOsBg1aD99cRTdRdcn0zC4ouPBMnG1_vyFqdw5EjcvM1fWXF3m-24pp95NwWXg5iqq8yf4NLdjtgCN8yWN5gdjJRrJKFb3yfX5qE&v=404","impurl":["http://im-x.jd.com/dsp/np?log=B4ZHNweusEbAOvPYea_kJBU3CaEJozHl0atEgL2tR1Ej5tG7ZZV0XR6NeoI4gpUHQJ6yevsE-fDHRfuv5khsr9yhtdkwN8EEE2KqdZkNvBkL7oMMP0-stnIT_hrDtL4_Pcs1z6_L34Ck3jZUAfU7_VSlPueOJmRaAbpTENqfud1nwyoMKQiL3BJumzIHa-ZqqRlIjXWwJpZFAz1PdPHEltqXROymL59LSN0_1NwLabzjp1r3GvoLG2F_pus81kdp4kXlFcf2TZSj4yd96FYag5__OIAAwMg2uAMfC7ndp3bRe_8HCoGWaTe86KUB500ly-3rtfq4H1PmAX4a04BhhJQHES6nccQMF6KYEvSx4V1KoFVNLRd5RqhNTsH0tdjupFUaWMmfRFIxymd-ifmEKTgAQ5a-uwGblrAVLlUOKlJbjULzRuZxIXaL2eICJJ2tb_Yku6qv-wXbzxn9MY3iUg~~&v=404&seq=1","http://wn.x.jd.com/adx/nurl/xiaomi?price={WIN_PRICE}&v=100&ad=2720&info=KlcKJTIyXzA1Njc1NmM5MDc3NjQ5ZjFiODBiZGI1YzgxOTAxY2Y2XzEgFii6ATIgODZjYmYyYjBjY2Y2MTJkMDJhYjYxYzMzMTZjY2E1MzM4AEBASKAVUAE~"],"adm":"{\"imgurl\":\"http://img1.360buyimg.com/pop/jfs/t10558/15/201167676/178266/f138e2e1/59c8a58bN3e922d95.jpg?adModifyId=138403019\",\"landingurl\":\"http://ccc-x.jd.com/dsp/nc?ext=aHR0cDovL21laS53b2NoYXcuY29tL3Byb2R1Y3RzL3h5bWo5OTlfMTQvMzAwMWhsc2owMDU&log=B4ZHNweusEbAOvPYea_kJBU3CaEJozHl0atEgL2tR1Ej5tG7ZZV0XR6NeoI4gpUHZzauuUwnOYW9eOhgkUoaZR2jOhHsORFqaNkE0smAQFgg_ewhbwO6tCtkr4bA-FotftPn-YoNniG-Ln-2Ra56C6BSGtydhaEkCRx9lyjQkh-WnP8iWFXlHaL8LhajJhYRvrKCWAly3ye2vnYlqzJN-2Okyh-rgxl5mS1Sc4ue0MztuE4NXDUZinTCUGT_rTRdRywwBz14mZ5ZWf9qSC77RLgGbQK4lcIrs6JJh6biTUgo8P44LqEQQQHHqccmxgZ-R2OQqC8Eqj4UhBzwRo-s23Co4-OOwOsBg1aD99cRTdRdcn0zC4ouPBMnG1_vyFqdw5EjcvM1fWXF3m-24pp95NwWXg5iqq8yf4NLdjtgCN8yWN5gdjJRrJKFb3yfX5qE&v=404\",\"source\":\"精品信阳毛尖\",\"title\":\"信阳毛尖铁罐装3罐冰点价仅149\"}","templateid":"2.1"}]}]}`)
	fmt.Println("start parse response!")
	err := json.Unmarshal(b, &m)

	if err != nil {
		test.Error("parse response error!", err)
	}

	bids := m.SeatBids[0].Bids
	for _,bid := range bids {
		t := convertMiBidToHaoMaiBid(&bid)
		//fmt.Println(t.Adm.Landingurl)
		bd := convertHaoMaiBidToMiBid(t)
		fmt.Println(len(bid.Adm), bid.Adm, "Hello")
		fmt.Println(len(bd.Adm), bd.Adm, "Hello")
		if len(bd.Adm) != len(bid.Adm) {
			test.Error("Convert not equals!")
		}
	}
}
