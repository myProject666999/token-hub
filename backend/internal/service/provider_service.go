package service

import (
	"errors"
	"token-hub/internal/model"
	"token-hub/internal/repository"

	"gorm.io/gorm"
)

type ProviderService struct{}

func NewProviderService() *ProviderService {
	return &ProviderService{}
}

func (s *ProviderService) CreateProvider(req *model.ProviderCreateRequest) (*model.Provider, error) {
	var existing model.Provider
	if err := repository.DB.Where("code = ?", req.Code).First(&existing).Error; err == nil {
		return nil, errors.New("服务商编码已存在")
	}

	provider := &model.Provider{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Logo:        req.Logo,
		Website:     req.Website,
		APIEndpoint: req.APIEndpoint,
		APIKey:      req.APIKey,
		APISecret:   req.APISecret,
		Status:      req.Status,
		Sort:        req.Sort,
	}

	if err := repository.DB.Create(provider).Error; err != nil {
		return nil, err
	}

	return provider, nil
}

func (s *ProviderService) UpdateProvider(id uint, req *model.ProviderCreateRequest) (*model.Provider, error) {
	var provider model.Provider
	if err := repository.DB.First(&provider, id).Error; err != nil {
		return nil, errors.New("服务商不存在")
	}

	updates := map[string]interface{}{
		"name":         req.Name,
		"description":  req.Description,
		"logo":         req.Logo,
		"website":      req.Website,
		"api_endpoint": req.APIEndpoint,
		"status":       req.Status,
		"sort":         req.Sort,
	}

	if req.APIKey != "" {
		updates["api_key"] = req.APIKey
	}
	if req.APISecret != "" {
		updates["api_secret"] = req.APISecret
	}

	if err := repository.DB.Model(&provider).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &provider, nil
}

func (s *ProviderService) DeleteProvider(id uint) error {
	return repository.DB.Delete(&model.Provider{}, id).Error
}

func (s *ProviderService) GetProviderByID(id uint) (*model.Provider, error) {
	var provider model.Provider
	if err := repository.DB.Preload("Models").First(&provider, id).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}

func (s *ProviderService) GetAllProviders(includeDisabled bool) ([]model.Provider, error) {
	var providers []model.Provider
	query := repository.DB.Preload("Models", func(db *gorm.DB) *gorm.DB {
		if !includeDisabled {
			return db.Where("status = ?", 1)
		}
		return db
	}).Order("sort ASC")

	if !includeDisabled {
		query = query.Where("status = ?", 1)
	}

	if err := query.Find(&providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}

func (s *ProviderService) GetProviderList(page, pageSize int, keyword string) ([]model.Provider, int64, error) {
	var providers []model.Provider
	var total int64

	query := repository.DB.Model(&model.Provider{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&providers).Error; err != nil {
		return nil, 0, err
	}

	return providers, total, nil
}

func (s *ProviderService) CreateModel(req *model.ModelCreateRequest) (*model.Model, error) {
	var provider model.Provider
	if err := repository.DB.First(&provider, req.ProviderID).Error; err != nil {
		return nil, errors.New("服务商不存在")
	}

	model := &model.Model{
		Name:              req.Name,
		Code:              req.Code,
		ProviderID:        req.ProviderID,
		Description:       req.Description,
		ModelType:         req.ModelType,
		PointsPer1KInput:  req.PointsPer1KInput,
		PointsPer1KOutput: req.PointsPer1KOutput,
		MaxTokens:         req.MaxTokens,
		ContextLimit:      req.ContextLimit,
		SupportsVision:    req.SupportsVision,
		SupportsFunction:  req.SupportsFunction,
		Status:            req.Status,
		Sort:              req.Sort,
	}

	if err := repository.DB.Create(model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (s *ProviderService) UpdateModel(id uint, req *model.ModelCreateRequest) (*model.Model, error) {
	var m model.Model
	if err := repository.DB.First(&m, id).Error; err != nil {
		return nil, errors.New("模型不存在")
	}

	updates := map[string]interface{}{
		"name":                req.Name,
		"description":         req.Description,
		"model_type":          req.ModelType,
		"points_per_1k_input": req.PointsPer1KInput,
		"points_per_1k_output": req.PointsPer1KOutput,
		"max_tokens":          req.MaxTokens,
		"context_limit":       req.ContextLimit,
		"supports_vision":     req.SupportsVision,
		"supports_function":   req.SupportsFunction,
		"status":              req.Status,
		"sort":                req.Sort,
	}

	if err := repository.DB.Model(&m).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *ProviderService) DeleteModel(id uint) error {
	return repository.DB.Delete(&model.Model{}, id).Error
}

func (s *ProviderService) GetModelByID(id uint) (*model.Model, error) {
	var m model.Model
	if err := repository.DB.Preload("Provider").First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *ProviderService) GetModelByCode(code string) (*model.Model, error) {
	var m model.Model
	if err := repository.DB.Preload("Provider").Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *ProviderService) GetModelList(page, pageSize int, providerID uint, keyword string) ([]model.Model, int64, error) {
	var models []model.Model
	var total int64

	query := repository.DB.Model(&model.Model{}).Preload("Provider")

	if providerID > 0 {
		query = query.Where("provider_id = ?", providerID)
	}

	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return models, total, nil
}
