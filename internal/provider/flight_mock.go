package provider

import (
	"context"
	"fmt"
	"math/rand"
)

// MockFlightProvider Mock 航班 Provider
type MockFlightProvider struct{}

// NewMockFlightProvider 创建 Mock 航班 Provider
func NewMockFlightProvider() *MockFlightProvider {
	return &MockFlightProvider{}
}

// Search 查询航班/火车
func (p *MockFlightProvider) Search(_ context.Context, req FlightSearchRequest) (*FlightSearchResult, error) {
	// 内置常见城市对的 Mock 数据
	tickets := p.getMockTickets(req.DepartCity, req.ArriveCity, req.Type)
	if len(tickets) == 0 {
		// 不在覆盖范围内，返回通用估算数据
		tickets = p.getGenericTickets(req.DepartCity, req.ArriveCity, req.Type)
	}

	return &FlightSearchResult{
		Tickets: tickets,
		IsMock:  true,
	}, nil
}

// getMockTickets 获取内置 Mock 数据
func (p *MockFlightProvider) getMockTickets(depart, arrive, ticketType string) []TicketInfo {
	key := fmt.Sprintf("%s-%s", depart, arrive)

	// 内置常见城市对数据
	mockData := map[string][]TicketInfo{
		"北京-三亚": {
			{Type: "flight", Number: "CA1831", Depart: "07:30", Arrive: "11:00", Price: 680, Seat: "经济舱", IsMock: true},
			{Type: "flight", Number: "HU7181", Depart: "09:15", Arrive: "12:45", Price: 720, Seat: "经济舱", IsMock: true},
			{Type: "flight", Number: "CZ6712", Depart: "14:20", Arrive: "17:50", Price: 650, Seat: "经济舱", IsMock: true},
		},
		"北京-厦门": {
			{Type: "flight", Number: "MF8115", Depart: "08:00", Arrive: "11:20", Price: 580, Seat: "经济舱", IsMock: true},
			{Type: "train", Number: "G323", Depart: "09:30", Arrive: "21:45", Price: 650, Seat: "二等座", IsMock: true},
		},
		"北京-大理": {
			{Type: "flight", Number: "CA1471", Depart: "10:30", Arrive: "14:20", Price: 850, Seat: "经济舱", IsMock: true},
			{Type: "flight", Number: "MU5713", Depart: "13:45", Arrive: "17:35", Price: 820, Seat: "经济舱", IsMock: true},
		},
		"上海-三亚": {
			{Type: "flight", Number: "HU7651", Depart: "08:30", Arrive: "11:30", Price: 620, Seat: "经济舱", IsMock: true},
			{Type: "flight", Number: "MU5331", Depart: "12:00", Arrive: "15:00", Price: 650, Seat: "经济舱", IsMock: true},
		},
		"上海-杭州": {
			{Type: "train", Number: "G7375", Depart: "07:00", Arrive: "08:00", Price: 73, Seat: "二等座", IsMock: true},
			{Type: "train", Number: "G7377", Depart: "09:30", Arrive: "10:30", Price: 73, Seat: "二等座", IsMock: true},
		},
		"上海-厦门": {
			{Type: "flight", Number: "MF8501", Depart: "09:00", Arrive: "11:00", Price: 480, Seat: "经济舱", IsMock: true},
			{Type: "train", Number: "G1651", Depart: "10:15", Arrive: "16:30", Price: 380, Seat: "二等座", IsMock: true},
		},
		"广州-三亚": {
			{Type: "flight", Number: "CZ6712", Depart: "08:45", Arrive: "10:15", Price: 420, Seat: "经济舱", IsMock: true},
			{Type: "flight", Number: "HU7803", Depart: "14:30", Arrive: "16:00", Price: 450, Seat: "经济舱", IsMock: true},
		},
		"成都-大理": {
			{Type: "flight", Number: "3U8965", Depart: "09:20", Arrive: "10:50", Price: 380, Seat: "经济舱", IsMock: true},
			{Type: "flight", Number: "CA4513", Depart: "15:40", Arrive: "17:10", Price: 400, Seat: "经济舱", IsMock: true},
		},
	}

	tickets, ok := mockData[key]
	if !ok {
		return nil
	}

	// 根据 type 过滤
	if ticketType == "any" || ticketType == "" {
		return tickets
	}

	filtered := []TicketInfo{}
	for _, t := range tickets {
		if t.Type == ticketType {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

// getGenericTickets 生成通用估算数据
func (p *MockFlightProvider) getGenericTickets(depart, arrive, ticketType string) []TicketInfo {
	tickets := []TicketInfo{}

	if ticketType == "flight" || ticketType == "any" || ticketType == "" {
		tickets = append(tickets, TicketInfo{
			Type:   "flight",
			Number: "XX1234",
			Depart: "08:00",
			Arrive: "11:30",
			Price:  500 + rand.Intn(500),
			Seat:   "经济舱",
			IsMock: true,
		})
	}

	if ticketType == "train" || ticketType == "any" || ticketType == "" {
		tickets = append(tickets, TicketInfo{
			Type:   "train",
			Number: "G1234",
			Depart: "09:00",
			Arrive: "15:30",
			Price:  300 + rand.Intn(300),
			Seat:   "二等座",
			IsMock: true,
		})
	}

	return tickets
}
