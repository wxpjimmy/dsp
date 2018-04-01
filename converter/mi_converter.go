package converter

import (
	"encoding/json"
	"dsp_demo/model"
	"bytes"
	"strings"
	"github.com/astaxie/beego/logs"
)

func ConvertMiRequestToHMRequest(req *model.MiRequest) *model.HMRequest {
	size := len(req.Imps)
	var hmImps = make([]model.HMImp, size, size)
	for i,v:=range req.Imps {
		d := ConvertMiImpToHMImp(&v)
		hmImps[i] = *d
	}
	return &model.HMRequest{
		ID: req.ID,
		Imps: hmImps,
		App: req.App,
		Device: req.Device,
		Test: req.Test,
		At: req.At,
	}
}

func ConvertMiImpToHMImp(imp *model.Imp) *model.HMImp {
	style, _ := model.GetDspConfBySrcAndSlotId("xiaomi", imp.TagId)
	if style != nil {
		return &model.HMImp{
			ID: imp.ID,
			TagId: imp.TagId,
			Bidfloor: imp.Bidfloor,
			Template: *style,
		}
	}
	return nil
}

func ConvertHaoMaiResponseToMiResponse(ir *model.HMResponse) *model.MiResponse {
	if ir ==nil {
		return nil
	}
	var seats = make([]model.SeatBid, len(ir.SeatBids))
	for idx,seat := range ir.SeatBids {
		d := convertHaoMaiSeatToMiSeat(&seat)
		seats[idx] = *d
	}
	return model.MiResponse{
		ID: ir.ID,
		SeatBids: seats,
		Nbr: ir.Nbr,
	}
}

func convertHaoMaiSeatToMiSeat(is *model.HaoMaiSeatBid) *model.SeatBid {
	if is == nil {
		return nil
	}
	var bids = make([]model.Bid, len(is.Bids))
	for idx,bid := range is.Bids {
		d := convertHaoMaiBidToMiBid(&bid)
		bids[idx] = *d
	}
	return &model.SeatBid{
		Bids: bids,
		Seat: is.Seat,
		Cm: is.Cm,
		Group: is.Group,
	}
}

func convertHaoMaiBidToMiBid(internal *model.HaoMaiBid) *model.Bid {
	var adm string
	buffer := bytes.NewBuffer(make([]byte, 0))
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(internal.Adm)
	//content := fmt.Sprintf("%s", buffer)
	//fmt.Println("Data: ", content)
	if err != nil {
		logs.Error("encode Adm error!", err)
	} else {
		logs.Info(buffer.Cap())
		adm = strings.TrimSpace(string(buffer.Bytes()))
	}
	conf := model.GetTrafficConfBySrcAndSlotId("xiaomi", internal.Tagid)
	tid := conf.SspTemplateId
	price := internal.Price
	if conf.Price != 0 {
		price = float32(conf.Price)
	}
	return &model.Bid{
		ID: internal.ID,
		Impid: internal.Impid,
		Price: price,
		Adid: internal.Adid,
		Adm: adm,
		Tagid: internal.Tagid,
		Templateid: tid,
		Billingtype: internal.Billingtype,
		Adomain: internal.Adomain,
		Cat: internal.Cat,
		H: internal.Template.H,
		W: internal.Template.W,
		Impurl: internal.Impurl,
		Curl: internal.Curl,
		Frequencycapping: internal.Frequencycapping,
	}
}

func convertMiBidToHaoMaiBid(bid *model.Bid) *model.HaoMaiBid {
	var adm model.Adm;
	b := []byte(bid.Adm)
	err := json.Unmarshal(b, &adm)
	if err != nil {
		logs.Error("parse Adm error!", err)
	}
	style, _ := model.GetDspConfBySrcAndSlotId("xiaomi", bid.Tagid)
	return &model.HaoMaiBid{
		ID: bid.ID,
		Impid: bid.Impid,
		Price: bid.Price,
		Adid: bid.Adid,
		Adm: adm,
		Tagid: bid.Tagid,
		Template: *style,
		Billingtype: bid.Billingtype,
		Adomain: bid.Adomain,
		Cat: bid.Cat,
		Impurl: bid.Impurl,
		Curl: bid.Curl,
		Frequencycapping: bid.Frequencycapping,
	}
}

