package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"token-hub/internal/model"
	"token-hub/internal/repository"
	"token-hub/pkg/jwt"

	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(req *model.UserRegisterRequest) (*model.User, error) {
	var existingUser model.User
	if err := repository.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	user := &model.User{
		Username: req.Username,
		Password: s.hashPassword(req.Password),
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Role:     "user",
		Status:   1,
	}

	if err := repository.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(req *model.UserLoginRequest) (string, *model.User, error) {
	var user model.User
	if err := repository.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("用户名或密码错误")
		}
		return "", nil, err
	}

	if user.Password != s.hashPassword(req.Password) {
		return "", nil, errors.New("用户名或密码错误")
	}

	if user.Status != 1 {
		return "", nil, errors.New("用户已被禁用")
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

func (s *UserService) GetByID(userID uint) (*model.User, error) {
	var user model.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Update(userID uint, req *model.UserUpdateRequest) (*model.User, error) {
	var user model.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) > 0 {
		if err := repository.DB.Model(&user).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserService) GetUserPoints(userID uint) (*model.UserPointsResponse, error) {
	var user model.User
	if err := repository.DB.Select("points", "total_points", "used_points").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &model.UserPointsResponse{
		Points:      user.Points,
		TotalPoints: user.TotalPoints,
		UsedPoints:  user.UsedPoints,
	}, nil
}

func (s *UserService) GetList(page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := repository.DB.Model(&model.User{})

	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *UserService) UpdateStatus(userID uint, status int) error {
	return repository.DB.Model(&model.User{}).Where("id = ?", userID).Update("status", status).Error
}

func (s *UserService) hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
