package session

import (
	"time"
)

// Phase 会话阶段枚举
type Phase string

const (
	PhaseCollecting   Phase = "COLLECTING"   // 收集用户需求
	PhaseRecommending Phase = "RECOMMENDING" // 目的地推荐
	PhasePlanning     Phase = "PLANNING"     // 行程规划
	PhaseBooking      Phase = "BOOKING"      // 预订引导
	PhaseCompleted    Phase = "COMPLETED"    // 流程完成
)

// Session 旅行规划会话，是整个业务的根聚合对象
type Session struct {
	ID              string           `json:"id"`
	Phase           Phase            `json:"phase"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	UserRequest     *UserRequest     `json:"user_request"`
	Recommendations []Recommendation `json:"recommendations"`
	Itinerary       *Itinerary       `json:"itinerary"`
	BookingItems    []BookingItem    `json:"booking_items"`
	ChatHistory     []ChatMessage    `json:"chat_history"`
}

// UserRequest 用户原始需求，在 COLLECTING 阶段由 Agent 逐步解析填充
type UserRequest struct {
	RawInput      string   `json:"raw_input"`
	DepartureCity string   `json:"departure_city"`
	TravelStart   string   `json:"travel_start"`
	TravelEnd     string   `json:"travel_end"`
	Budget        int      `json:"budget"`
	Preferences   []string `json:"preferences"`
	Travelers     int      `json:"travelers"`
}

// Recommendation 单个目的地推荐方案
type Recommendation struct {
	ID            string        `json:"id"`
	Destination   string        `json:"destination"`
	Reason        string        `json:"reason"`
	Weather       WeatherInfo   `json:"weather"`
	Transport     TransportInfo `json:"transport"`
	HotelRange    HotelRange    `json:"hotel_range"`
	TotalEstimate CostRange     `json:"total_estimate"`
	SampleTickets []TicketInfo  `json:"sample_tickets"`
	IsMock        bool          `json:"is_mock"`
}

// WeatherInfo 天气信息
type WeatherInfo struct {
	TempRange string  `json:"temp_range"`
	Condition string  `json:"condition"`
	RainProb  float64 `json:"rain_prob"`
	IsMock    bool    `json:"is_mock"`
}

// TransportInfo 交通信息
type TransportInfo struct {
	Mode     string `json:"mode"`
	Duration string `json:"duration"`
	Price    int    `json:"price"`
	IsMock   bool   `json:"is_mock"`
}

// HotelRange 酒店价格范围
type HotelRange struct {
	Economy int  `json:"economy"`
	Comfort int  `json:"comfort"`
	IsMock  bool `json:"is_mock"`
}

// CostRange 费用区间
type CostRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// TicketInfo 票务示例
type TicketInfo struct {
	Type   string `json:"type"`
	Number string `json:"number"`
	Depart string `json:"depart"`
	Arrive string `json:"arrive"`
	Price  int    `json:"price"`
	Seat   string `json:"seat"`
	IsMock bool   `json:"is_mock"`
}

// Itinerary 完整行程表，PLANNING 阶段生成
type Itinerary struct {
	Destination string       `json:"destination"`
	Outbound    FlightDetail `json:"outbound"`
	Hotel       HotelDetail  `json:"hotel"`
	DailyPlans  []DayPlan    `json:"daily_plans"`
	Return      FlightDetail `json:"return"`
	CostSummary CostSummary  `json:"cost_summary"`
	TotalCost   int          `json:"total_cost"`
}

// FlightDetail 具体交通详情
type FlightDetail struct {
	Type       string `json:"type"`
	Number     string `json:"number"`
	DepartCity string `json:"depart_city"`
	ArriveCity string `json:"arrive_city"`
	DepartTime string `json:"depart_time"`
	ArriveTime string `json:"arrive_time"`
	Duration   string `json:"duration"`
	Price      int    `json:"price"`
	Seat       string `json:"seat"`
	IsMock     bool   `json:"is_mock"`
}

// HotelDetail 酒店详情
type HotelDetail struct {
	Name          string `json:"name"`
	Location      string `json:"location"`
	RoomType      string `json:"room_type"`
	CheckIn       string `json:"check_in"`
	CheckOut      string `json:"check_out"`
	Nights        int    `json:"nights"`
	PricePerNight int    `json:"price_per_night"`
	TotalPrice    int    `json:"total_price"`
	IsMock        bool   `json:"is_mock"`
}

// DayPlan 每日行程
type DayPlan struct {
	Date       string     `json:"date"`
	DayNumber  int        `json:"day_number"`
	Title      string     `json:"title"`
	Activities []Activity `json:"activities"`
}

// Activity 活动/景点
type Activity struct {
	Time      string `json:"time"`
	Name      string `json:"name"`
	Transport string `json:"transport"`
	Cost      int    `json:"cost"`
	Note      string `json:"note"`
}

// CostSummary 费用汇总
type CostSummary struct {
	OutboundCost  int `json:"outbound_cost"`
	ReturnCost    int `json:"return_cost"`
	HotelCost     int `json:"hotel_cost"`
	ActivityCost  int `json:"activity_cost"`
	MealCost      int `json:"meal_cost"`
	TransportCost int `json:"transport_cost"`
}

// BookingItem 预订项
type BookingItem struct {
	ID         string        `json:"id"`
	ItemType   BookingType   `json:"item_type"`
	ItemName   string        `json:"item_name"`
	Detail     string        `json:"detail"`
	Platform   string        `json:"platform"`
	Link       string        `json:"link"`
	AltLinks   []PlatformLink `json:"alt_links"`
	Price      int           `json:"price"`
	Status     BookingStatus `json:"status"`
	ActualCost int           `json:"actual_cost"`
	Order      int           `json:"order"`
}

// PlatformLink 平台链接
type PlatformLink struct {
	Platform string `json:"platform"`
	Link     string `json:"link"`
}

// BookingType 预订类型
type BookingType string

const (
	BookingOutbound BookingType = "OUTBOUND"
	BookingHotel    BookingType = "HOTEL"
	BookingReturn   BookingType = "RETURN"
	BookingTicket   BookingType = "TICKET"
)

// BookingStatus 预订状态
type BookingStatus string

const (
	BookingPending BookingStatus = "PENDING"
	BookingBooked  BookingStatus = "BOOKED"
	BookingSkipped BookingStatus = "SKIPPED"
)

// ChatMessage 单条对话消息
type ChatMessage struct {
	ID        string         `json:"id"`
	Role      string         `json:"role"`
	Content   string         `json:"content"`
	Blocks    []ContentBlock `json:"blocks"`
	Timestamp int64          `json:"timestamp"`
}

// ContentBlock 消息中的内容块
type ContentBlock struct {
	Type     string      `json:"type"`
	Content  string      `json:"content"`
	CardType string      `json:"card_type"`
	Data     interface{} `json:"data"`
}

// BookingSummary 预订汇总
type BookingSummary struct {
	Destination  string        `json:"destination"`
	Items        []BookingItem `json:"items"`
	TotalBooked  int           `json:"total_booked"`
	TotalSkipped int           `json:"total_skipped"`
	BookedCount  int           `json:"booked_count"`
}
