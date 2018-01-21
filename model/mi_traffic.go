package model

//request interface with xiaomi
type Template struct {
	ID string `json:"id"`
	Width int32 `json:"width"`
	Height int32 `json:"height"`
}

type Imp struct {
	ID string `json:"id"`
	TagId string `json:"tagid"`
	Admtype int32 `json:"admtype"`
	Bidfloor float32 `json:"bidfloor"`
	Templates []Template `json:"templates"`
}

type App struct {
	Name string `json:"name"`
	Bundle string `json:"bundle"`
}

type Device struct {
	Ip string `json:"ip"`
	Make string `json:"make"`
	Model string `json:"model"`
	OS string `json:"os"`
	Osv string `json:"osv"`
	H int32 `json:"h"`
	W int32 `json:"w"`
	JS int32 `json:"js"`
	Language string `json:"language"`
	Connectiontype int32 `json:"connectiontype"`
	ImeiMD5 string `json:"didmd5"`
	AndroidId string `json:"dpid"`
	AndroidIdMD5 string `json:"dpidmd5"`
	Carrier string `json:"carrier"`
}

type FrequencyCapping struct {
	Global int32 `json:"global"`
	Weekly int32 `json:"weekly"`
	Daily int32 `json:"daily"`
	Hourly int32 `json:"hourly"`
}

//response interface with xiaomi
type Bid struct {
	ID string `json:"id"`
	Impid string `json:"impid"`
	Price float32 `json:"price"`
	Adid string `json:"adid"`
	Adm string `json:"adm"`
	Tagid string `json:"tagid"`
	Templateid string `json:"templateid"`
	Billingtype int32 `json:"billingtype"`
	Adomain []string `json:"adomain"`
	Cat []int32 `json:"cat"`
	H int32 `json:"h"`
	W int32 `json:"w"`
	Impurl []string `json:"impurl"`
	Curl []string `json:"curl"`
	Frequencycapping FrequencyCapping `json:"frequencycapping"`
}

type SeatBid struct {
	Bids []Bid `json:"bid"`
	Seat string `json:"seat"`
	Cm int32 `json:"cm"`
	Group int32 `json:"group"`
}

type MiRequest struct {
	ID string `json:"id"`
	Imps []Imp `json:"imp"`
	App App `json:"app"`
	Device Device `json:"device"`
	Test int32 `json:"test"`
	At int32 `json:"at"`
}
type MiResponse struct {
	ID string `json:"id"`
	SeatBids []SeatBid `json:"seatbid"`
	Nbr int32 `json:"nbr"`
}

func (req *MiRequest) GetPackageName() string {
	if req.App.Bundle == "" {
		return DEFAULT_PACKAGE_NAME
	}
	return req.App.Bundle
}