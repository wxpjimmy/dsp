package model


//interface with downstream parts
type HMImp struct {
	ID string `json:"id"`
	TagId string `json:"tagid"`
	//Admtype int32 `json:"admtype"`
	Bidfloor float32 `json:"bidfloor"`
	Template TemplateStyle `json:"template"`
}

type HMRequest struct {
	ID string `json:"id"`
	Imps []HMImp `json:"imps"`
	App App `json:"app"`
	Device Device `json:"device"`
	Test int32 `json:"test"`
	At int32 `json:"at"`
}

type Adm struct {
	Title string `json:"title"`
	Description string `json:"description,omitempty"`
	Source string `json:"source"`
	Landingurl string `json:"landingurl"`
	Imgurl string `json:"imgurl"`
	Actionurl string `json:"actionurl,omitempty"`
	Videourl string `json:"videourl,omitempty"`
}

type HaoMaiBid struct {
	ID string `json:"id"`
	Impid string `json:"impid"`
	Price float32 `json:"price"`
	Adid string `json:"adid"`
	Adm Adm `json:"adm"`
	Tagid string `json:"tagid"`
	Template TemplateStyle `json:"template"`
	//Templateid string `json:"templateid"`
	Billingtype int32 `json:"billingtype"`
	Adomain []string `json:"adomain"`
	Cat []int32 `json:"cat"`
	//H int32 `json:"h"`
	//W int32 `json:"w"`
	Impurl []string `json:"impurl"`
	Curl []string `json:"curl"`
	Frequencycapping FrequencyCapping `json:"frequencycapping"`
}

type HaoMaiSeatBid struct {
	Bids []HaoMaiBid `json:"bid"`
	Seat string `json:"seat"`
	Cm int32 `json:"cm"`
	Group int32 `json:"group"`
}

type HMResponse struct {
	ID string `json:"id"`
	SeatBids []HaoMaiSeatBid `json:"seatbid"`
	Nbr int32 `json:"nbr"`
}

type Request interface {
	GetPackageName() string
}

func (req *HMRequest) GetPackageName() string {
	if req.App.Bundle == "" {
		return DEFAULT_PACKAGE_NAME
	}
	return req.App.Bundle
}