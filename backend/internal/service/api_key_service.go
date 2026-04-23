package service

import (
	"errors"
	"time"
	"token-hub/internal/model"
	"token-hub/internal/repository"
	"token-hub/pkg/jwt"
)

type APIKeyService struct{}

func NewAPIKeyService() *APIKeyService {
	return &APIKeyService{}
}

func (s *APIKeyService) CreateAPIKey(userID uint, name string) (*model.APIKey, error) {
	key := &model.APIKey{
		UserID: userID,
		Key:    jwt.GenerateAPIKey(),
		Name:   name,
		Status: 1,
	}

	if err := repository.DB.Create(key).Error; err != nil {
		return nil, err
	}

	return key, nil
}

func (s *APIKeyService) GetAPIKeys(userID uint) ([]model.APIKey, error) {
	var keys []model.APIKey
	if err := repository.DB.Where("user_id = ?", userID).Order("id DESC").Find(&keys).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

func (s *APIKeyService) GetAPIKeyByID(userID, keyID uint) (*model.APIKey, error) {
	var key model.APIKey
	if err := repository.DB.Where("id = ? AND user_id = ?", keyID, userID).First(&key).Error; err != nil {
		return nil, errors.New("API密钥不存在")
	}
	return &key, nil
}

func (s *APIKeyService) DeleteAPIKey(userID, keyID uint) error {
	result := repository.DB.Where("id = ? AND user_id = ?", keyID, userID).Delete(&model.APIKey{})
	if result.RowsAffected == 0 {
		return errors.New("API密钥不存在")
	}
	return result.Error
}

func (s *APIKeyService) UpdateAPIKeyStatus(userID, keyID uint, status int) error {
	result := repository.DB.Model(&model.APIKey{}).
		Where("id = ? AND user_id = ?", keyID, userID).
		Update("status", status)
	if result.RowsAffected == 0 {
		return errors.New("API密钥不存在")
	}
	return result.Error
}

func (s *APIKeyService) UpdateLastUsed(keyID uint) error {
	return repository.DB.Model(&model.APIKey{}).
		Where("id = ?", keyID).
		Update("last_used", time.Now()).Error
}
