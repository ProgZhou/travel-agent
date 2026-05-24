package provider

import (
	"context"
	"math/rand"
)

// MockHotelProvider Mock 酒店 Provider
type MockHotelProvider struct{}

// NewMockHotelProvider 创建 Mock 酒店 Provider
func NewMockHotelProvider() *MockHotelProvider {
	return &MockHotelProvider{}
}

// Search 查询酒店
func (p *MockHotelProvider) Search(_ context.Context, req HotelSearchRequest) (*HotelSearchResult, error) {
	hotels := p.getMockHotels(req.City, req.Level)
	return &HotelSearchResult{
		Hotels: hotels,
		IsMock: true,
	}, nil
}

// getMockHotels 获取 Mock 酒店数据
func (p *MockHotelProvider) getMockHotels(city, level string) []HotelInfo {
	// 根据城市和档次生成合理的酒店数据
	cityHotels := map[string]map[string][]HotelInfo{
		"三亚": {
			"economy": {
				{Name: "如家快捷酒店(三亚湾店)", Location: "三亚湾", Level: "economy", PricePerNight: 180, Rating: 4.2, RoomTypes: []string{"标准间", "大床房"}},
				{Name: "汉庭酒店(三亚市区店)", Location: "市区", Level: "economy", PricePerNight: 200, Rating: 4.1, RoomTypes: []string{"标准间"}},
			},
			"comfort": {
				{Name: "三亚湾海景度假酒店", Location: "三亚湾", Level: "comfort", PricePerNight: 380, Rating: 4.5, RoomTypes: []string{"海景房", "豪华房"}},
				{Name: "亚龙湾舒适酒店", Location: "亚龙湾", Level: "comfort", PricePerNight: 420, Rating: 4.6, RoomTypes: []string{"园景房", "海景房"}},
			},
			"luxury": {
				{Name: "三亚亚龙湾丽思卡尔顿酒店", Location: "亚龙湾", Level: "luxury", PricePerNight: 1200, Rating: 4.9, RoomTypes: []string{"豪华海景房", "套房"}},
				{Name: "三亚海棠湾君悦酒店", Location: "海棠湾", Level: "luxury", PricePerNight: 1500, Rating: 4.8, RoomTypes: []string{"海景套房", "别墅"}},
			},
		},
		"厦门": {
			"economy": {
				{Name: "7天连锁酒店(厦门中山路店)", Location: "中山路", Level: "economy", PricePerNight: 150, Rating: 4.0, RoomTypes: []string{"标准间"}},
				{Name: "如家酒店(厦门鼓浪屿店)", Location: "鼓浪屿", Level: "economy", PricePerNight: 220, Rating: 4.3, RoomTypes: []string{"标准间", "大床房"}},
			},
			"comfort": {
				{Name: "厦门海景花园酒店", Location: "环岛路", Level: "comfort", PricePerNight: 350, Rating: 4.5, RoomTypes: []string{"海景房", "豪华房"}},
				{Name: "厦门鼓浪屿精品酒店", Location: "鼓浪屿", Level: "comfort", PricePerNight: 400, Rating: 4.6, RoomTypes: []string{"园景房", "海景房"}},
			},
			"luxury": {
				{Name: "厦门威斯汀酒店", Location: "环岛路", Level: "luxury", PricePerNight: 900, Rating: 4.8, RoomTypes: []string{"豪华海景房", "套房"}},
			},
		},
		"大理": {
			"economy": {
				{Name: "大理古城青年旅舍", Location: "古城", Level: "economy", PricePerNight: 120, Rating: 4.2, RoomTypes: []string{"标准间", "多人间"}},
			},
			"comfort": {
				{Name: "大理洱海景观酒店", Location: "洱海边", Level: "comfort", PricePerNight: 320, Rating: 4.6, RoomTypes: []string{"湖景房", "豪华房"}},
				{Name: "大理古城精品客栈", Location: "古城", Level: "comfort", PricePerNight: 280, Rating: 4.5, RoomTypes: []string{"庭院房", "标准房"}},
			},
			"luxury": {
				{Name: "大理洱海天域英迪格酒店", Location: "洱海边", Level: "luxury", PricePerNight: 1100, Rating: 4.9, RoomTypes: []string{"湖景套房", "别墅"}},
			},
		},
	}

	// 查找对应城市和档次的酒店
	if cityData, ok := cityHotels[city]; ok {
		if hotels, ok := cityData[level]; ok {
			return hotels
		}
		// 如果没有指定档次，返回 comfort 档次
		if hotels, ok := cityData["comfort"]; ok {
			return hotels
		}
	}

	// 通用估算数据
	return p.getGenericHotels(level)
}

// getGenericHotels 生成通用估算酒店数据
func (p *MockHotelProvider) getGenericHotels(level string) []HotelInfo {
	priceMap := map[string]int{
		"economy": 150 + rand.Intn(100),
		"comfort": 300 + rand.Intn(200),
		"luxury":  800 + rand.Intn(500),
	}

	price, ok := priceMap[level]
	if !ok {
		price = 300
	}

	return []HotelInfo{
		{
			Name:          "当地精选酒店",
			Location:      "市区",
			Level:         level,
			PricePerNight: price,
			Rating:        4.0 + rand.Float64()*0.5,
			RoomTypes:     []string{"标准间", "大床房"},
		},
	}
}
