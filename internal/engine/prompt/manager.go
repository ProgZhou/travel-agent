package prompt

import (
	"fmt"
	"travel-agent/internal/session"
)

// Manager Prompt 管理器
type Manager struct{}

// NewManager 创建 Prompt 管理器
func NewManager() *Manager {
	return &Manager{}
}

// BuildSystemPrompt 构建完整的 System Prompt
func (m *Manager) BuildSystemPrompt(sess *session.Session) string {
	base := SystemPrompt
	phasePrompt := m.getPhasePrompt(sess.Phase)
	context := m.buildContext(sess)

	return base + "\n\n" + phasePrompt + "\n\n" + context
}

// getPhasePrompt 获取阶段 Prompt
func (m *Manager) getPhasePrompt(phase session.Phase) string {
	switch phase {
	case session.PhaseCollecting:
		return CollectingPrompt
	case session.PhaseRecommending:
		return RecommendingPrompt
	case session.PhasePlanning:
		return PlanningPrompt
	case session.PhaseBooking:
		return BookingPrompt
	case session.PhaseCompleted:
		return "## 当前阶段：完成 (COMPLETED)\n\n旅行规划已完成，祝用户旅途愉快！"
	default:
		return ""
	}
}

// buildContext 构建上下文信息
func (m *Manager) buildContext(sess *session.Session) string {
	ctx := "## 当前会话上下文\n\n"

	// 用户需求
	if sess.UserRequest != nil {
		ctx += "### 用户需求\n"
		req := sess.UserRequest
		if req.DepartureCity != "" {
			ctx += fmt.Sprintf("- 出发城市：%s\n", req.DepartureCity)
		}
		if req.TravelStart != "" && req.TravelEnd != "" {
			ctx += fmt.Sprintf("- 出行时间：%s 至 %s\n", req.TravelStart, req.TravelEnd)
		}
		if req.Budget > 0 {
			ctx += fmt.Sprintf("- 预算：%d 元\n", req.Budget)
		}
		if len(req.Preferences) > 0 {
			ctx += fmt.Sprintf("- 偏好：%v\n", req.Preferences)
		}
		if req.Travelers > 0 {
			ctx += fmt.Sprintf("- 出行人数：%d 人\n", req.Travelers)
		}
		ctx += "\n"
	}

	// 推荐列表
	if len(sess.Recommendations) > 0 {
		ctx += "### 已推荐目的地\n"
		for _, rec := range sess.Recommendations {
			ctx += fmt.Sprintf("- %s：%s\n", rec.Destination, rec.Reason)
		}
		ctx += "\n"
	}

	// 行程表
	if sess.Itinerary != nil {
		ctx += "### 已规划行程\n"
		ctx += fmt.Sprintf("- 目的地：%s\n", sess.Itinerary.Destination)
		ctx += fmt.Sprintf("- 总费用：%d 元\n", sess.Itinerary.TotalCost)
		ctx += "\n"
	}

	// 预订状态
	if len(sess.BookingItems) > 0 {
		ctx += "### 预订进度\n"
		for _, item := range sess.BookingItems {
			ctx += fmt.Sprintf("- %s：%s\n", item.ItemName, item.Status)
		}
		ctx += "\n"
	}

	return ctx
}
