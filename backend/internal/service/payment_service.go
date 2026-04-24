package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"token-hub/internal/model"
	"token-hub/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentService struct {
	pointsService *PointsService
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		pointsService: NewPointsService(),
	}
}

func (s *PaymentService) GetPaymentMethods(includeDisabled bool) ([]model.PaymentMethod, error) {
	var methods []model.PaymentMethod
	query := repository.DB.Order("sort ASC")

	if !includeDisabled {
		query = query.Where("status = ?", 1)
	}

	if err := query.Find(&methods).Error; err != nil {
		return nil, err
	}

	return methods, nil
}

func (s *PaymentService) CreatePaymentMethod(method *model.PaymentMethod) error {
	return repository.DB.Create(method).Error
}

func (s *PaymentService) UpdatePaymentMethod(id uint, method *model.PaymentMethod) error {
	return repository.DB.Model(&model.PaymentMethod{}).Where("id = ?", id).Updates(method).Error
}

func (s *PaymentService) DeletePaymentMethod(id uint) error {
	return repository.DB.Delete(&model.PaymentMethod{}, id).Error
}

func (s *PaymentService) CreateRechargeOrder(userID uint, req *model.RechargeRequest) (*model.RechargeRecord, error) {
	var method model.PaymentMethod
	if err := repository.DB.Where("id = ? AND status = ?", req.PaymentMethodID, 1).First(&method).Error; err != nil {
		return nil, errors.New("支付方式不存在或已禁用")
	}

	pointsRate, err := s.pointsService.GetPointsRate()
	if err != nil {
		return nil, err
	}

	points := req.Amount * pointsRate

	orderNo := generateOrderNo()
	expiredAt := time.Now().Add(30 * time.Minute)

	record := &model.RechargeRecord{
		OrderNo:         orderNo,
		UserID:          userID,
		PaymentMethodID: req.PaymentMethodID,
		Amount:          req.Amount,
		Points:          points,
		PointsRate:      pointsRate,
		Status:          model.PaymentStatusPending,
		ExpiredAt:       &expiredAt,
	}

	if err := repository.DB.Create(record).Error; err != nil {
		return nil, err
	}

	return record, nil
}

func (s *PaymentService) ProcessPaymentNotify(req *model.PaymentNotifyRequest) error {
	return repository.DB.Transaction(func(tx *gorm.DB) error {
		var record model.RechargeRecord
		if err := tx.Where("order_no = ?", req.OrderNo).First(&record).Error; err != nil {
			return errors.New("订单不存在")
		}

		if record.Status != model.PaymentStatusPending {
			return nil
		}

		if record.Amount != req.Amount {
			return errors.New("金额不匹配")
		}

		now := time.Now()
		updates := map[string]interface{}{
			"status":   model.PaymentStatusPaid,
			"trade_no": req.TradeNo,
			"paid_at":  &now,
		}

		if err := tx.Model(&record).Updates(updates).Error; err != nil {
			return err
		}

		if err := s.pointsService.AddPoints(
			record.UserID,
			record.Points,
			fmt.Sprintf("充值%.2f元，获得%.2f积分", record.Amount, record.Points),
			record.OrderNo,
			"recharge",
			tx,
		); err != nil {
			return err
		}

		return nil
	})
}

func (s *PaymentService) GetRechargeRecordList(userID uint, page, pageSize int, status string) ([]model.RechargeRecord, int64, error) {
	var records []model.RechargeRecord
	var total int64

	query := repository.DB.Model(&model.RechargeRecord{}).
		Preload("PaymentMethod").
		Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (s *PaymentService) GetAllRechargeRecords(page, pageSize int, status string) ([]model.RechargeRecord, int64, error) {
	var records []model.RechargeRecord
	var total int64

	query := repository.DB.Model(&model.RechargeRecord{}).
		Preload("User").
		Preload("PaymentMethod")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (s *PaymentService) SimulatePayment(orderNo string) error {
	var record model.RechargeRecord
	if err := repository.DB.Where("order_no = ?", orderNo).First(&record).Error; err != nil {
		return errors.New("订单不存在")
	}

	if record.Status != model.PaymentStatusPending {
		return errors.New("订单状态不正确")
	}

	notifyReq := &model.PaymentNotifyRequest{
		OrderNo: orderNo,
		TradeNo: "SIM_" + generateOrderNo(),
		Amount:  record.Amount,
		Status:  model.PaymentStatusPaid,
	}

	return s.ProcessPaymentNotify(notifyReq)
}

func (s *PaymentService) GetPaymentConfig(methodCode string) (interface{}, error) {
	var method model.PaymentMethod
	if err := repository.DB.Where("code = ?", methodCode).First(&method).Error; err != nil {
		return nil, err
	}

	if method.Config == "" {
		return nil, nil
	}

	switch methodCode {
	case "alipay":
		var config model.AlipayConfig
		if err := json.Unmarshal([]byte(method.Config), &config); err != nil {
			return nil, err
		}
		return &config, nil
	case "wechat":
		var config model.WechatConfig
		if err := json.Unmarshal([]byte(method.Config), &config); err != nil {
			return nil, err
		}
		return &config, nil
	default:
		return nil, errors.New("不支持的支付方式")
	}
}

func generateOrderNo() string {
	uuidStr := uuid.New().String()
	return fmt.Sprintf("TH%s%d", uuidStr[:8], time.Now().Unix())
}
