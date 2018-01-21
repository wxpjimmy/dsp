package util

import (
	"testing"
	"dsp_demo/model"
	"fmt"
	"strings"
)

func TestJsonEncodeWithoutEscape(t *testing.T) {
	adm := model.Adm{
		Imgurl: "http://img1.360buyimg.com/pop/jfs/t10558/15/201167676/178266/f138e2e1/59c8a58bN3e922d95.jpg?adModifyId=138403019",
		Title: "信阳毛尖铁罐装3罐冰点价仅149",
		Source: "精品信阳毛尖",
		Landingurl: "http://ccc-x.jd.com/dsp/nc?ext=aHR0cDovL21laS53b2NoYXcuY29tL3Byb2R1Y3RzL3h5bWo5OTlfMTQvMzAwMWhsc2owMDU&log=B4ZHNweusEbAOvPYea_kJBU3CaEJozHl0atEgL2tR1Ej5tG7ZZV0XR6NeoI4gpUHZzauuUwnOYW9eOhgkUoaZR2jOhHsORFqaNkE0smAQFgg_ewhbwO6tCtkr4bA-FotftPn-YoNniG-Ln-2Ra56C6BSGtydhaEkCRx9lyjQkh-WnP8iWFXlHaL8LhajJhYRvrKCWAly3ye2vnYlqzJN-2Okyh-rgxl5mS1Sc4ue0MztuE4NXDUZinTCUGT_rTRdRywwBz14mZ5ZWf9qSC77RLgGbQK4lcIrs6JJh6biTUgo8P44LqEQQQHHqccmxgZ-R2OQqC8Eqj4UhBzwRo-s23Co4-OOwOsBg1aD99cRTdRdcn0zC4ouPBMnG1_vyFqdw5EjcvM1fWXF3m-24pp95NwWXg5iqq8yf4NLdjtgCN8yWN5gdjJRrJKFb3yfX5qE&v=404",
	}

	b := JsonEncodeWithoutEscape(adm)
	fmt.Println(strings.TrimSpace(string(b)))
	//b = []byte(strings.TrimSpace(string(b)))
	var decode model.Adm
	JsonDecode(b, &decode)
	fmt.Println(decode)
	if adm.Imgurl != decode.Imgurl || adm.Landingurl != decode.Landingurl || adm.Source != decode.Source || adm.Title != decode.Title {
		t.Error("Decode content changed!")
	}
}
