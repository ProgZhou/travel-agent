package provider

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// MockWeatherProvider Mock 天气 Provider
type MockWeatherProvider struct{}

// NewMockWeatherProvider 创建 Mock 天气 Provider
func NewMockWeatherProvider() *MockWeatherProvider {
	return &MockWeatherProvider{}
}

// Query 查询天气
func (p *MockWeatherProvider) Query(_ context.Context, req WeatherQueryRequest) (*WeatherQueryResult, error) {
	daily := p.generateDailyWeather(req.City, req.StartDate, req.EndDate)
	summary := p.generateSummary(daily)

	return &WeatherQueryResult{
		City:    req.City,
		Daily:   daily,
		Summary: summary,
		IsMock:  true,
	}, nil
}

// generateDailyWeather 生成每日天气数据
func (p *MockWeatherProvider) generateDailyWeather(city, startDate, endDate string) []DailyWeather {
	// 根据城市和月份生成合理的天气数据
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	daily := []DailyWeather{}
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		weather := p.getCityWeather(city, d.Month())
		daily = append(daily, DailyWeather{
			Date:      d.Format("2006-01-02"),
			TempHigh:  weather.tempHigh + rand.Intn(3) - 1,
			TempLow:   weather.tempLow + rand.Intn(3) - 1,
			Condition: weather.condition,
			RainProb:  weather.rainProb + (rand.Float64()-0.5)*0.1,
		})
	}

	return daily
}

// generateSummary 生成天气汇总
func (p *MockWeatherProvider) generateSummary(daily []DailyWeather) WeatherSummary {
	if len(daily) == 0 {
		return WeatherSummary{IsMock: true}
	}

	minTemp := daily[0].TempLow
	maxTemp := daily[0].TempHigh
	totalRainProb := 0.0
	conditionCount := make(map[string]int)

	for _, d := range daily {
		if d.TempLow < minTemp {
			minTemp = d.TempLow
		}
		if d.TempHigh > maxTemp {
			maxTemp = d.TempHigh
		}
		totalRainProb += d.RainProb
		conditionCount[d.Condition]++
	}

	// 找出最常见的天气状况
	mainCondition := "晴"
	maxCount := 0
	for cond, count := range conditionCount {
		if count > maxCount {
			maxCount = count
			mainCondition = cond
		}
	}

	return WeatherSummary{
		TempRange: fmt.Sprintf("%d-%d°C", minTemp, maxTemp),
		Condition: mainCondition,
		RainProb:  totalRainProb / float64(len(daily)),
		IsMock:    true,
	}
}

// cityWeatherTemplate 城市天气模板
type cityWeatherTemplate struct {
	tempHigh  int
	tempLow   int
	condition string
	rainProb  float64
}

// getCityWeather 根据城市和月份获取天气模板
func (p *MockWeatherProvider) getCityWeather(city string, month time.Month) cityWeatherTemplate {
	// 内置常见城市的季节性天气数据
	templates := map[string]map[time.Month]cityWeatherTemplate{
		"三亚": {
			time.January:   {tempHigh: 26, tempLow: 20, condition: "晴", rainProb: 0.1},
			time.February:  {tempHigh: 27, tempLow: 21, condition: "晴", rainProb: 0.15},
			time.March:     {tempHigh: 29, tempLow: 23, condition: "晴", rainProb: 0.2},
			time.April:     {tempHigh: 31, tempLow: 25, condition: "晴", rainProb: 0.25},
			time.May:       {tempHigh: 32, tempLow: 26, condition: "晴", rainProb: 0.3},
			time.June:      {tempHigh: 33, tempLow: 27, condition: "多云", rainProb: 0.4},
			time.July:      {tempHigh: 33, tempLow: 27, condition: "多云", rainProb: 0.45},
			time.August:    {tempHigh: 33, tempLow: 27, condition: "多云", rainProb: 0.5},
			time.September: {tempHigh: 32, tempLow: 26, condition: "多云", rainProb: 0.4},
			time.October:   {tempHigh: 30, tempLow: 25, condition: "晴", rainProb: 0.3},
			time.November:  {tempHigh: 28, tempLow: 23, condition: "晴", rainProb: 0.2},
			time.December:  {tempHigh: 26, tempLow: 21, condition: "晴", rainProb: 0.15},
		},
		"厦门": {
			time.January:   {tempHigh: 18, tempLow: 12, condition: "多云", rainProb: 0.2},
			time.February:  {tempHigh: 19, tempLow: 13, condition: "多云", rainProb: 0.25},
			time.March:     {tempHigh: 22, tempLow: 16, condition: "多云", rainProb: 0.3},
			time.April:     {tempHigh: 26, tempLow: 20, condition: "晴", rainProb: 0.25},
			time.May:       {tempHigh: 29, tempLow: 23, condition: "晴", rainProb: 0.3},
			time.June:      {tempHigh: 31, tempLow: 26, condition: "多云", rainProb: 0.4},
			time.July:      {tempHigh: 33, tempLow: 27, condition: "晴", rainProb: 0.35},
			time.August:    {tempHigh: 33, tempLow: 27, condition: "晴", rainProb: 0.4},
			time.September: {tempHigh: 31, tempLow: 25, condition: "晴", rainProb: 0.35},
			time.October:   {tempHigh: 27, tempLow: 22, condition: "晴", rainProb: 0.25},
			time.November:  {tempHigh: 23, tempLow: 18, condition: "多云", rainProb: 0.2},
			time.December:  {tempHigh: 19, tempLow: 14, condition: "多云", rainProb: 0.2},
		},
		"大理": {
			time.January:   {tempHigh: 15, tempLow: 3, condition: "晴", rainProb: 0.1},
			time.February:  {tempHigh: 17, tempLow: 5, condition: "晴", rainProb: 0.15},
			time.March:     {tempHigh: 20, tempLow: 8, condition: "晴", rainProb: 0.2},
			time.April:     {tempHigh: 23, tempLow: 11, condition: "晴", rainProb: 0.25},
			time.May:       {tempHigh: 25, tempLow: 14, condition: "多云", rainProb: 0.35},
			time.June:      {tempHigh: 26, tempLow: 16, condition: "多云", rainProb: 0.5},
			time.July:      {tempHigh: 25, tempLow: 16, condition: "阴", rainProb: 0.6},
			time.August:    {tempHigh: 25, tempLow: 16, condition: "阴", rainProb: 0.6},
			time.September: {tempHigh: 23, tempLow: 14, condition: "多云", rainProb: 0.5},
			time.October:   {tempHigh: 21, tempLow: 11, condition: "晴", rainProb: 0.3},
			time.November:  {tempHigh: 18, tempLow: 7, condition: "晴", rainProb: 0.15},
			time.December:  {tempHigh: 15, tempLow: 4, condition: "晴", rainProb: 0.1},
		},
	}

	if cityData, ok := templates[city]; ok {
		if weather, ok := cityData[month]; ok {
			return weather
		}
	}

	// 通用估算数据
	return cityWeatherTemplate{
		tempHigh:  20 + rand.Intn(10),
		tempLow:   15 + rand.Intn(5),
		condition: "晴",
		rainProb:  0.2 + rand.Float64()*0.2,
	}
}
