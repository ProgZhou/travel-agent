package provider

import "context"

// FlightProvider 航班/火车查询 Provider
type FlightProvider interface {
	Search(ctx context.Context, req FlightSearchRequest) (*FlightSearchResult, error)
}

// FlightSearchRequest 航班查询请求
type FlightSearchRequest struct {
	DepartCity string `json:"depart_city"`
	ArriveCity string `json:"arrive_city"`
	Date       string `json:"date"`
	Type       string `json:"type"` // flight | train | any
}

// FlightSearchResult 航班查询结果
type FlightSearchResult struct {
	Tickets []TicketInfo `json:"tickets"`
	IsMock  bool         `json:"is_mock"`
}

// TicketInfo 票务信息
type TicketInfo struct {
	Type   string `json:"type"`
	Number string `json:"number"`
	Depart string `json:"depart"`
	Arrive string `json:"arrive"`
	Price  int    `json:"price"`
	Seat   string `json:"seat"`
	IsMock bool   `json:"is_mock"`
}

// HotelProvider 酒店查询 Provider
type HotelProvider interface {
	Search(ctx context.Context, req HotelSearchRequest) (*HotelSearchResult, error)
}

// HotelSearchRequest 酒店查询请求
type HotelSearchRequest struct {
	City     string `json:"city"`
	CheckIn  string `json:"check_in"`
	CheckOut string `json:"check_out"`
	Level    string `json:"level"` // economy | comfort | luxury
}

// HotelSearchResult 酒店查询结果
type HotelSearchResult struct {
	Hotels []HotelInfo `json:"hotels"`
	IsMock bool        `json:"is_mock"`
}

// HotelInfo 酒店信息
type HotelInfo struct {
	Name          string   `json:"name"`
	Location      string   `json:"location"`
	Level         string   `json:"level"`
	PricePerNight int      `json:"price_per_night"`
	Rating        float64  `json:"rating"`
	RoomTypes     []string `json:"room_types"`
}

// WeatherProvider 天气查询 Provider
type WeatherProvider interface {
	Query(ctx context.Context, req WeatherQueryRequest) (*WeatherQueryResult, error)
}

// WeatherQueryRequest 天气查询请求
type WeatherQueryRequest struct {
	City      string `json:"city"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// WeatherQueryResult 天气查询结果
type WeatherQueryResult struct {
	City    string        `json:"city"`
	Daily   []DailyWeather `json:"daily"`
	Summary WeatherSummary `json:"summary"`
	IsMock  bool          `json:"is_mock"`
}

// DailyWeather 每日天气
type DailyWeather struct {
	Date      string  `json:"date"`
	TempHigh  int     `json:"temp_high"`
	TempLow   int     `json:"temp_low"`
	Condition string  `json:"condition"`
	RainProb  float64 `json:"rain_prob"`
}

// WeatherSummary 天气汇总
type WeatherSummary struct {
	TempRange string  `json:"temp_range"`
	Condition string  `json:"condition"`
	RainProb  float64 `json:"rain_prob"`
	IsMock    bool    `json:"is_mock"`
}
