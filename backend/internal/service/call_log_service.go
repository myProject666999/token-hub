package service

import (
	"token-hub/internal/model"
	"token-hub/internal/repository"
)

type CallLogService struct{}

func NewCallLogService() *CallLogService {
	return &CallLogService{}
}

func (s *CallLogService) CreateLog(log *model.CallLog) error {
	return repository.DB.Create(log).Error
}

func (s *CallLogService) GetUserCallLogs(userID uint, page, pageSize int, status string) ([]model.CallLog, int64, error) {
	var logs []model.CallLog
	var total int64

	query := repository.DB.Model(&model.CallLog{}).
		Preload("Provider").
		Preload("Model").
		Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (s *CallLogService) GetAllCallLogs(page, pageSize int, userID uint, providerID uint, modelID uint, status string) ([]model.CallLog, int64, error) {
	var logs []model.CallLog
	var total int64

	query := repository.DB.Model(&model.CallLog{}).
		Preload("User").
		Preload("Provider").
		Preload("Model")

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if providerID > 0 {
		query = query.Where("provider_id = ?", providerID)
	}
	if modelID > 0 {
		query = query.Where("model_id = ?", modelID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (s *CallLogService) GetUserStatistics(userID uint) (*model.CallStatistics, error) {
	var stats model.CallStatistics
	stats.UserID = userID

	repository.DB.Model(&model.CallLog{}).Where("user_id = ?", userID).Count(&stats.TotalCalls)
	repository.DB.Model(&model.CallLog{}).Where("user_id = ? AND status = ?", userID, model.CallStatusSuccess).Count(&stats.SuccessCalls)
	repository.DB.Model(&model.CallLog{}).Where("user_id = ? AND status = ?", userID, model.CallStatusFailed).Count(&stats.FailedCalls)

	var tokens struct {
		Input  int64
		Output int64
		Points float64
	}
	repository.DB.Model(&model.CallLog{}).
		Select("COALESCE(SUM(input_tokens), 0) as input, COALESCE(SUM(output_tokens), 0) as output, COALESCE(SUM(points_consumed), 0) as points").
		Where("user_id = ?", userID).
		Scan(&tokens)

	stats.TotalInputTokens = tokens.Input
	stats.TotalOutputTokens = tokens.Output
	stats.TotalTokens = tokens.Input + tokens.Output
	stats.TotalPoints = tokens.Points

	return &stats, nil
}

func (s *CallLogService) GetDailyStatistics(userID uint, providerID uint, modelID uint, startDate, endDate string) ([]model.DailyStatistics, error) {
	var stats []model.DailyStatistics

	query := repository.DB.Model(&model.CallLog{}).
		Select(`
			DATE(created_at) as date,
			COUNT(*) as total_calls,
			SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as success_calls,
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed_calls,
			COALESCE(SUM(input_tokens), 0) as total_input_tokens,
			COALESCE(SUM(output_tokens), 0) as total_output_tokens,
			COALESCE(SUM(input_tokens + output_tokens), 0) as total_tokens,
			COALESCE(SUM(points_consumed), 0) as total_points
		`)

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if providerID > 0 {
		query = query.Where("provider_id = ?", providerID)
	}
	if modelID > 0 {
		query = query.Where("model_id = ?", modelID)
	}
	if startDate != "" {
		query = query.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(created_at) <= ?", endDate)
	}

	query = query.Group("DATE(created_at)").Order("date DESC")

	err := query.Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *CallLogService) GetOverviewStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalUsers int64
	repository.DB.Model(&model.User{}).Count(&totalUsers)
	stats["total_users"] = totalUsers

	var totalCalls int64
	repository.DB.Model(&model.CallLog{}).Count(&totalCalls)
	stats["total_calls"] = totalCalls

	var totalRecharge float64
	repository.DB.Model(&model.RechargeRecord{}).Where("status = ?", model.PaymentStatusPaid).Select("COALESCE(SUM(amount), 0)").Scan(&totalRecharge)
	stats["total_recharge"] = totalRecharge

	var totalPointsUsed float64
	repository.DB.Model(&model.User{}).Select("COALESCE(SUM(used_points), 0)").Scan(&totalPointsUsed)
	stats["total_points_used"] = totalPointsUsed

	var todayCalls int64
	repository.DB.Model(&model.CallLog{}).Where("DATE(created_at) = DATE('now')").Count(&todayCalls)
	stats["today_calls"] = todayCalls

	var todayRecharge float64
	repository.DB.Model(&model.RechargeRecord{}).Where("status = ? AND DATE(created_at) = DATE('now')", model.PaymentStatusPaid).Select("COALESCE(SUM(amount), 0)").Scan(&todayRecharge)
	stats["today_recharge"] = todayRecharge

	return stats, nil
}
