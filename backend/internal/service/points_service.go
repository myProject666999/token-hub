package service

import (
	"errors"
	"fmt"
	"token-hub/internal/model"
	"token-hub/internal/repository"

	"gorm.io/gorm"
)

type PointsService struct{}

func NewPointsService() *PointsService {
	return &PointsService{}
}

func (s *PointsService) GetPointsRate() (float64, error) {
	var config model.PointsConfig
	if err := repository.DB.Where("key = ?", model.PointsKeyRate).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 100.0, nil
		}
		return 0, err
	}

	var rate float64
	if _, err := fmt.Sscanf(config.Value, "%f", &rate); err != nil {
		return 100.0, nil
	}

	return rate, nil
}

func (s *PointsService) SetPointsRate(rate float64) error {
	var config model.PointsConfig
	err := repository.DB.Where("key = ?", model.PointsKeyRate).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			config = model.PointsConfig{
				Key:         model.PointsKeyRate,
				Value:       string(rune(rate)),
				Description: "积分汇率：1元兑换多少积分",
			}
			return repository.DB.Create(&config).Error
		}
		return err
	}

	return repository.DB.Model(&config).Update("value", rate).Error
}

func (s *PointsService) AddPoints(userID uint, amount float64, description, relatedID, relatedType string, tx ...*gorm.DB) error {
	exec := func(db *gorm.DB) error {
		var user model.User
		if err := db.First(&user, userID).Error; err != nil {
			return err
		}

		newPoints := user.Points + amount
		newTotalPoints := user.TotalPoints + amount

		if err := db.Model(&user).Updates(map[string]interface{}{
			"points":       newPoints,
			"total_points": newTotalPoints,
		}).Error; err != nil {
			return err
		}

		record := &model.PointsRecord{
			UserID:      userID,
			Type:        model.PointsTypeRecharge,
			Amount:      amount,
			Balance:     newPoints,
			Description: description,
			RelatedID:   relatedID,
			RelatedType: relatedType,
		}

		return db.Create(record).Error
	}

	if len(tx) > 0 && tx[0] != nil {
		return exec(tx[0])
	}

	return repository.DB.Transaction(func(tx *gorm.DB) error {
		return exec(tx)
	})
}

func (s *PointsService) ConsumePoints(userID uint, amount float64, description, relatedID, relatedType string, tx ...*gorm.DB) error {
	exec := func(db *gorm.DB) error {
		var user model.User
		if err := db.First(&user, userID).Error; err != nil {
			return err
		}

		if user.Points < amount {
			return errors.New("积分不足")
		}

		newPoints := user.Points - amount
		newUsedPoints := user.UsedPoints + amount

		if err := db.Model(&user).Updates(map[string]interface{}{
			"points":      newPoints,
			"used_points": newUsedPoints,
		}).Error; err != nil {
			return err
		}

		record := &model.PointsRecord{
			UserID:      userID,
			Type:        model.PointsTypeConsume,
			Amount:      -amount,
			Balance:     newPoints,
			Description: description,
			RelatedID:   relatedID,
			RelatedType: relatedType,
		}

		return db.Create(record).Error
	}

	if len(tx) > 0 && tx[0] != nil {
		return exec(tx[0])
	}

	return repository.DB.Transaction(func(tx *gorm.DB) error {
		return exec(tx)
	})
}

func (s *PointsService) GetPointsRecordList(userID uint, page, pageSize int, recordType string) ([]model.PointsRecord, int64, error) {
	var records []model.PointsRecord
	var total int64

	query := repository.DB.Model(&model.PointsRecord{}).Where("user_id = ?", userID)

	if recordType != "" {
		query = query.Where("type = ?", recordType)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (s *PointsService) GetUserStatistics(userID uint) (*model.PointsStatistics, error) {
	var user model.User
	if err := repository.DB.Select("points", "total_points", "used_points").First(&user, userID).Error; err != nil {
		return nil, err
	}

	var rechargeCount int64
	repository.DB.Model(&model.PointsRecord{}).Where("user_id = ? AND type = ?", userID, model.PointsTypeRecharge).Count(&rechargeCount)

	var callCount int64
	repository.DB.Model(&model.CallLog{}).Where("user_id = ?", userID).Count(&callCount)

	return &model.PointsStatistics{
		UserID:        userID,
		TotalPoints:   user.TotalPoints,
		UsedPoints:    user.UsedPoints,
		Balance:       user.Points,
		RechargeCount: rechargeCount,
		CallCount:     callCount,
	}, nil
}

func (s *PointsService) GetAllStatistics(page, pageSize int) ([]model.PointsStatistics, int64, error) {
	var users []model.User
	var total int64

	query := repository.DB.Model(&model.User{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Select("id", "points", "total_points", "used_points").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var stats []model.PointsStatistics
	for _, user := range users {
		var rechargeCount int64
		repository.DB.Model(&model.PointsRecord{}).Where("user_id = ? AND type = ?", user.ID, model.PointsTypeRecharge).Count(&rechargeCount)

		var callCount int64
		repository.DB.Model(&model.CallLog{}).Where("user_id = ?", user.ID).Count(&callCount)

		stats = append(stats, model.PointsStatistics{
			UserID:        user.ID,
			TotalPoints:   user.TotalPoints,
			UsedPoints:    user.UsedPoints,
			Balance:       user.Points,
			RechargeCount: rechargeCount,
			CallCount:     callCount,
		})
	}

	return stats, total, nil
}
