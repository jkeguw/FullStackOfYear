package device

import (
	"project/backend/internal/errors"
	"project/backend/models"
	"project/backend/services/similarity"
	"project/backend/types/device"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Service 外设服务接口
type Service interface {
	// 设备相关
	CreateMouseDevice(ctx context.Context, request device.CreateMouseRequest) (*models.MouseDevice, error)
	GetDeviceByID(ctx context.Context, deviceID string) (*models.HardwareDevice, error)
	GetDeviceList(ctx context.Context, deviceType string, page, pageSize int) (*device.DeviceListResponse, error)
	
	// 兼容层新增方法
	GetMouseDevice(ctx context.Context, deviceID string) (*models.MouseDevice, error)
	UpdateMouseDevice(ctx context.Context, deviceID string, request device.UpdateMouseRequest) (*models.MouseDevice, error)
	DeleteDevice(ctx context.Context, deviceID string) error
	ListDevices(ctx context.Context, filter device.DeviceListFilter) (*device.DeviceListResponse, error)
	
	// 相似度相关
	CompareMice(ctx context.Context, mouseIDs []string) (*similarity.ComparisonResult, error)
	FindSimilarMice(ctx context.Context, mouseID string, limit int) ([]models.MouseDevice, error)
	
	// 用户设备相关
	CreateUserDevice(ctx context.Context, userID string, request device.CreateUserDeviceRequest) (*models.UserDevice, error)
	UpdateUserDevice(ctx context.Context, userID, userDeviceID string, request device.UpdateUserDeviceRequest) (*models.UserDevice, error)
	GetUserDeviceByID(ctx context.Context, userID, userDeviceID string) (*models.UserDevice, error)
	GetUserDevices(ctx context.Context, userID string, page, pageSize int) (*device.UserDeviceListResponse, error)
	GetPublicUserDevices(ctx context.Context, page, pageSize int) (*device.UserDeviceListResponse, error)
	
	// 兼容层新增方法
	GetUserDevice(ctx context.Context, userID string, userDeviceID string) (*models.UserDevice, error)
	DeleteUserDevice(ctx context.Context, userID string, userDeviceID string) error
	ListUserDevices(ctx context.Context, request device.UserDeviceListRequest) (*device.UserDeviceListResponse, error)
	
	// 评测相关
	SubmitReview(ctx context.Context, userID string, request device.SubmitReviewRequest) (*models.DeviceReview, error)
	GetReviewsByDeviceID(ctx context.Context, deviceID string, page, pageSize int) (*device.DeviceReviewListResponse, error)
	GetReviewsByUserID(ctx context.Context, userID string, page, pageSize int) (*device.DeviceReviewListResponse, error)
	GetReviewByID(ctx context.Context, reviewID string) (*models.DeviceReview, error)
	UpdateReview(ctx context.Context, userID, reviewID string, request device.UpdateDeviceReviewRequest) (*models.DeviceReview, error)
	DeleteReview(ctx context.Context, userID, reviewID string) error
	ApproveReview(ctx context.Context, reviewerID, reviewID string) (*models.DeviceReview, error)
	RejectReview(ctx context.Context, reviewerID, reviewID string, notes string) (*models.DeviceReview, error)
	FeaturedReview(ctx context.Context, reviewerID, reviewID string, rank int) (*models.DeviceReview, error)
	GetPendingReviews(ctx context.Context, reviewType string, page, pageSize int) (*device.DeviceReviewListResponse, error)
	
	// 兼容层新增方法
	CreateDeviceReview(ctx context.Context, userID string, request device.CreateReviewRequest) (*models.DeviceReview, error)
	GetDeviceReview(ctx context.Context, reviewID string) (*models.DeviceReview, error)
	UpdateDeviceReview(ctx context.Context, userID string, reviewID string, request device.UpdateReviewRequest) (*models.DeviceReview, error)
	DeleteDeviceReview(ctx context.Context, userID string, reviewID string) error
	ListDeviceReviews(ctx context.Context, request device.DeviceReviewListRequest) (*device.DeviceReviewListResponse, error)
}

// 服务实现
type ServiceImpl struct {
	db *mongo.Database
}

// DefaultService 默认外设服务实现
type DefaultService struct {
	db *mongo.Database
}

// New 创建新的外设服务
func New(db *mongo.Database) Service {
	return &ServiceImpl{
		db: db,
	}
}

// NewDefaultService 创建默认外设服务
func NewDefaultService(db *mongo.Database) *DefaultService {
	return &DefaultService{
		db: db,
	}
}

// ApproveReview 批准评测
func (s *DefaultService) ApproveReview(ctx context.Context, reviewerID, reviewID string) (*models.DeviceReview, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// CreateMouseDevice 创建鼠标设备
func (s *DefaultService) CreateMouseDevice(ctx context.Context, request device.CreateMouseRequest) (*models.MouseDevice, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetDeviceByID 根据ID获取设备
func (s *DefaultService) GetDeviceByID(ctx context.Context, deviceID string) (*models.HardwareDevice, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetDeviceList 获取设备列表
func (s *DefaultService) GetDeviceList(ctx context.Context, deviceType string, page, pageSize int) (*device.DeviceListResponse, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetMouseDevice 获取鼠标设备
func (s *DefaultService) GetMouseDevice(ctx context.Context, deviceID string) (*models.MouseDevice, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// UpdateMouseDevice 更新鼠标设备
func (s *DefaultService) UpdateMouseDevice(ctx context.Context, deviceID string, request device.UpdateMouseRequest) (*models.MouseDevice, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// DeleteDevice 删除设备
func (s *DefaultService) DeleteDevice(ctx context.Context, deviceID string) error {
	// 空实现，仅为了满足接口
	return nil
}

// ListDevices 列出设备
func (s *DefaultService) ListDevices(ctx context.Context, filter device.DeviceListFilter) (*device.DeviceListResponse, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// CompareMice 比较鼠标
func (s *DefaultService) CompareMice(ctx context.Context, mouseIDs []string) (*similarity.ComparisonResult, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// FindSimilarMice 查找相似鼠标
func (s *DefaultService) FindSimilarMice(ctx context.Context, mouseID string, limit int) ([]models.MouseDevice, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// CreateUserDevice 创建用户设备
func (s *DefaultService) CreateUserDevice(ctx context.Context, userID string, request device.CreateUserDeviceRequest) (*models.UserDevice, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// UpdateUserDevice 更新用户设备
func (s *DefaultService) UpdateUserDevice(ctx context.Context, userID, userDeviceID string, request device.UpdateUserDeviceRequest) (*models.UserDevice, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetUserDeviceByID 根据ID获取用户设备
func (s *DefaultService) GetUserDeviceByID(ctx context.Context, userID, userDeviceID string) (*models.UserDevice, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetUserDevices 获取用户设备列表
func (s *DefaultService) GetUserDevices(ctx context.Context, userID string, page, pageSize int) (*device.UserDeviceListResponse, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetPublicUserDevices 获取公开用户设备
func (s *DefaultService) GetPublicUserDevices(ctx context.Context, page, pageSize int) (*device.UserDeviceListResponse, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetUserDevice 获取用户设备
func (s *DefaultService) GetUserDevice(ctx context.Context, userID string, userDeviceID string) (*models.UserDevice, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// DeleteUserDevice 删除用户设备
func (s *DefaultService) DeleteUserDevice(ctx context.Context, userID string, userDeviceID string) error {
	// 空实现，仅为了满足接口
	return nil
}

// ListUserDevices 列出用户设备
func (s *DefaultService) ListUserDevices(ctx context.Context, request device.UserDeviceListRequest) (*device.UserDeviceListResponse, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// SubmitReview 提交评测
func (s *DefaultService) SubmitReview(ctx context.Context, userID string, request device.SubmitReviewRequest) (*models.DeviceReview, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetReviewsByDeviceID 获取设备评测
func (s *DefaultService) GetReviewsByDeviceID(ctx context.Context, deviceID string, page, pageSize int) (*device.DeviceReviewListResponse, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetReviewsByUserID 获取用户评测
func (s *DefaultService) GetReviewsByUserID(ctx context.Context, userID string, page, pageSize int) (*device.DeviceReviewListResponse, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetReviewByID 获取评测
func (s *DefaultService) GetReviewByID(ctx context.Context, reviewID string) (*models.DeviceReview, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// UpdateReview 更新评测
func (s *DefaultService) UpdateReview(ctx context.Context, userID, reviewID string, request device.UpdateDeviceReviewRequest) (*models.DeviceReview, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// DeleteReview 删除评测
func (s *DefaultService) DeleteReview(ctx context.Context, userID, reviewID string) error {
	// 空实现，仅为了满足接口
	return nil
}

// RejectReview 拒绝评测
func (s *DefaultService) RejectReview(ctx context.Context, reviewerID, reviewID string, notes string) (*models.DeviceReview, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// FeaturedReview 推荐评测
func (s *DefaultService) FeaturedReview(ctx context.Context, reviewerID, reviewID string, rank int) (*models.DeviceReview, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetPendingReviews 获取待审核评测
func (s *DefaultService) GetPendingReviews(ctx context.Context, reviewType string, page, pageSize int) (*device.DeviceReviewListResponse, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// CreateDeviceReview 创建设备评测
func (s *DefaultService) CreateDeviceReview(ctx context.Context, userID string, request device.CreateReviewRequest) (*models.DeviceReview, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetDeviceReview 获取设备评测
func (s *DefaultService) GetDeviceReview(ctx context.Context, reviewID string) (*models.DeviceReview, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// UpdateDeviceReview 更新设备评测
func (s *DefaultService) UpdateDeviceReview(ctx context.Context, userID string, reviewID string, request device.UpdateReviewRequest) (*models.DeviceReview, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// DeleteDeviceReview 删除设备评测
func (s *DefaultService) DeleteDeviceReview(ctx context.Context, userID string, reviewID string) error {
	// 空实现，仅为了满足接口
	return nil
}

// ListDeviceReviews 列出设备评测
func (s *DefaultService) ListDeviceReviews(ctx context.Context, request device.DeviceReviewListRequest) (*device.DeviceReviewListResponse, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// CreateMouseDevice 创建鼠标设备
func (s *ServiceImpl) CreateMouseDevice(ctx context.Context, request device.CreateMouseRequest) (*models.MouseDevice, error) {
	now := time.Now()
	
	// 创建基础设备信息
	deviceInfo := models.HardwareDevice{
		ID:          primitive.NewObjectID(),
		Name:        request.Name,
		Brand:       request.Brand,
		Type:        models.DeviceTypeMouse,
		ImageURL:    request.ImageURL,
		Description: request.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	// 创建鼠标设备
	mouseDevice := &models.MouseDevice{
		HardwareDevice: deviceInfo,
		Dimensions:     request.Dimensions,
		Shape:          request.Shape,
		Technical:      request.Technical,
		Recommended:    request.Recommended,
	}
	
	// 保存到数据库
	_, err := s.db.Collection(models.DevicesCollection).InsertOne(ctx, mouseDevice)
	if err != nil {
		return nil, errors.NewInternalServerError("创建鼠标设备失败: " + err.Error())
	}
	
	return mouseDevice, nil
}

// GetDeviceByID 根据ID获取设备
func (s *ServiceImpl) GetDeviceByID(ctx context.Context, deviceID string) (*models.HardwareDevice, error) {
	id, err := primitive.ObjectIDFromHex(deviceID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的设备ID")
	}
	
	var deviceDoc models.HardwareDevice
	err = s.db.Collection(models.DevicesCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&deviceDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("设备不存在")
		}
		return nil, errors.NewInternalServerError("获取设备失败: " + err.Error())
	}
	
	return &deviceDoc, nil
}

// GetDeviceList 获取设备列表
func (s *ServiceImpl) GetDeviceList(ctx context.Context, deviceType string, page, pageSize int) (*device.DeviceListResponse, error) {
	// 构建查询条件
	filter := bson.M{}
	if deviceType != "" {
		filter["type"] = deviceType
	}
	
	// 设置分页
	if page <= 0 {
		page = 1
	}
	
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	
	skip := (page - 1) * pageSize
	skipInt64 := int64(skip)
	limitInt64 := int64(pageSize)
	
	// 查询设备总数
	total, err := s.db.Collection(models.DevicesCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.NewInternalServerError("计算设备总数失败: " + err.Error())
	}
	
	// 查询设备列表
	cursor, err := s.db.Collection(models.DevicesCollection).Find(ctx, filter, &options.FindOptions{
		Skip:  &skipInt64,
		Limit: &limitInt64,
		Sort:  bson.D{{Key: "createdAt", Value: -1}},
	})
	if err != nil {
		return nil, errors.NewInternalServerError("查询设备列表失败: " + err.Error())
	}
	defer cursor.Close(ctx)
	
	// 解析基本设备信息
	var devices []models.HardwareDevice
	if err = cursor.All(ctx, &devices); err != nil {
		return nil, errors.NewInternalServerError("解析设备列表失败: " + err.Error())
	}
	
	// 转换为响应格式
	var response device.DeviceListResponse
	response.Total = int(total)
	response.Page = page
	response.PageSize = pageSize
	response.Devices = make([]device.DevicePreview, len(devices))
	
	for i, d := range devices {
		response.Devices[i] = device.DevicePreview{
			ID:          d.ID.Hex(),
			Name:        d.Name,
			Brand:       d.Brand,
			Type:        string(d.Type),
			ImageURL:    d.ImageURL,
			Description: d.Description,
			CreatedAt:   d.CreatedAt,
		}
	}
	
	return &response, nil
}

// CreateUserDevice 创建用户设备配置
func (s *ServiceImpl) CreateUserDevice(ctx context.Context, userID string, request device.CreateUserDeviceRequest) (*models.UserDevice, error) {
	// 将用户ID转换为ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的用户ID")
	}
	
	// 创建用户设备
	now := time.Now()
	userDevice := &models.UserDevice{
		ID:          primitive.NewObjectID(),
		UserID:      userObjectID,
		Name:        request.Name,
		Description: request.Description,
		Devices:     make([]models.UserDeviceSettings, len(request.Devices)),
		IsPublic:    request.IsPublic,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	// 处理设备设置
	for i, d := range request.Devices {
		deviceID, err := primitive.ObjectIDFromHex(d.DeviceID)
		if err != nil {
			return nil, errors.NewBadRequestError("无效的设备ID: " + d.DeviceID)
		}
		
		userDevice.Devices[i] = models.UserDeviceSettings{
			DeviceID:   deviceID,
			DeviceType: models.DeviceTypeEnum(d.DeviceType),
			Settings:   d.Settings,
		}
	}
	
	// 保存到数据库
	_, err = s.db.Collection(models.UserDevicesCollection).InsertOne(ctx, userDevice)
	if err != nil {
		return nil, errors.NewInternalServerError("创建用户设备配置失败: " + err.Error())
	}
	
	return userDevice, nil
}

// UpdateUserDevice 更新用户设备配置
func (s *ServiceImpl) UpdateUserDevice(ctx context.Context, userID, userDeviceID string, request device.UpdateUserDeviceRequest) (*models.UserDevice, error) {
	// 将用户ID和设备ID转换为ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的用户ID")
	}
	
	deviceObjectID, err := primitive.ObjectIDFromHex(userDeviceID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的设备ID")
	}
	
	// 查找当前用户设备
	var userDevice models.UserDevice
	err = s.db.Collection(models.UserDevicesCollection).FindOne(ctx, bson.M{
		"_id":    deviceObjectID,
		"userId": userObjectID,
	}).Decode(&userDevice)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("未找到用户设备配置")
		}
		return nil, errors.NewInternalServerError("获取用户设备配置失败: " + err.Error())
	}
	
	// 更新字段
	if request.Name != "" {
		userDevice.Name = request.Name
	}
	
	if request.Description != "" {
		userDevice.Description = request.Description
	}
	
	userDevice.IsPublic = request.IsPublic
	
	// 处理设备设置
	if len(request.Devices) > 0 {
		userDevice.Devices = make([]models.UserDeviceSettings, len(request.Devices))
		
		for i, d := range request.Devices {
			deviceID, err := primitive.ObjectIDFromHex(d.DeviceID)
			if err != nil {
				return nil, errors.NewBadRequestError("无效的设备ID: " + d.DeviceID)
			}
			
			userDevice.Devices[i] = models.UserDeviceSettings{
				DeviceID:   deviceID,
				DeviceType: models.DeviceTypeEnum(d.DeviceType),
				Settings:   d.Settings,
			}
		}
	}
	
	userDevice.UpdatedAt = time.Now()
	
	// 保存到数据库
	_, err = s.db.Collection(models.UserDevicesCollection).ReplaceOne(ctx, bson.M{
		"_id":    deviceObjectID,
		"userId": userObjectID,
	}, userDevice)
	
	if err != nil {
		return nil, errors.NewInternalServerError("更新用户设备配置失败: " + err.Error())
	}
	
	return &userDevice, nil
}

// GetUserDeviceByID 获取用户设备配置
func (s *ServiceImpl) GetUserDeviceByID(ctx context.Context, userID, userDeviceID string) (*models.UserDevice, error) {
	// 将ID转换为ObjectID
	deviceObjectID, err := primitive.ObjectIDFromHex(userDeviceID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的设备ID")
	}
	
	// 构建查询
	filter := bson.M{"_id": deviceObjectID}
	
	// 如果提供了用户ID，确保只能查询该用户的设备
	if userID != "" {
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, errors.NewBadRequestError("无效的用户ID")
		}
		filter["userId"] = userObjectID
	} else {
		// 如果没有提供用户ID，确保只能查询公开的设备
		filter["isPublic"] = true
	}
	
	// 查询设备
	var userDevice models.UserDevice
	err = s.db.Collection(models.UserDevicesCollection).FindOne(ctx, filter).Decode(&userDevice)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("未找到用户设备配置")
		}
		return nil, errors.NewInternalServerError("获取用户设备配置失败: " + err.Error())
	}
	
	return &userDevice, nil
}

// GetUserDevices 获取用户设备列表
func (s *ServiceImpl) GetUserDevices(ctx context.Context, userID string, page, pageSize int) (*device.UserDeviceListResponse, error) {
	// 将用户ID转换为ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的用户ID")
	}
	
	// 设置分页
	if page <= 0 {
		page = 1
	}
	
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	
	skip := (page - 1) * pageSize
	skipInt64 := int64(skip)
	limitInt64 := int64(pageSize)
	
	// 构建查询
	filter := bson.M{"userId": userObjectID}
	
	// 计算总数
	total, err := s.db.Collection(models.UserDevicesCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.NewInternalServerError("计算用户设备配置总数失败: " + err.Error())
	}
	
	// 查询设备列表
	cursor, err := s.db.Collection(models.UserDevicesCollection).Find(ctx, filter, &options.FindOptions{
		Skip:  &skipInt64,
		Limit: &limitInt64,
		Sort:  bson.D{{Key: "updatedAt", Value: -1}},
	})
	
	if err != nil {
		return nil, errors.NewInternalServerError("查询用户设备配置列表失败: " + err.Error())
	}
	defer cursor.Close(ctx)
	
	// 解析设备列表
	var userDevices []models.UserDevice
	if err = cursor.All(ctx, &userDevices); err != nil {
		return nil, errors.NewInternalServerError("解析用户设备配置列表失败: " + err.Error())
	}
	
	// 转换为响应格式
	response := &device.UserDeviceListResponse{
		Total:       int(total),
		Page:        page,
		PageSize:    pageSize,
		UserDevices: make([]device.UserDeviceResponse, len(userDevices)),
	}
	
	for i, ud := range userDevices {
		response.UserDevices[i] = device.UserDeviceResponse{
			ID:          ud.ID.Hex(),
			UserID:      ud.UserID.Hex(),
			Name:        ud.Name,
			Description: ud.Description,
			IsPublic:    ud.IsPublic,
			CreatedAt:   ud.CreatedAt,
			UpdatedAt:   ud.UpdatedAt,
			Devices:     make([]device.UserDeviceSettingsResponse, len(ud.Devices)),
		}
		
		for j, d := range ud.Devices {
			response.UserDevices[i].Devices[j] = device.UserDeviceSettingsResponse{
				DeviceID:   d.DeviceID.Hex(),
				DeviceType: string(d.DeviceType),
				Settings:   d.Settings,
			}
		}
	}
	
	return response, nil
}

// GetPublicUserDevices 获取公开的用户设备列表
func (s *ServiceImpl) GetPublicUserDevices(ctx context.Context, page, pageSize int) (*device.UserDeviceListResponse, error) {
	// 设置分页
	if page <= 0 {
		page = 1
	}
	
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	
	skip := (page - 1) * pageSize
	skipInt64 := int64(skip)
	limitInt64 := int64(pageSize)
	
	// 构建查询
	filter := bson.M{"isPublic": true}
	
	// 计算总数
	total, err := s.db.Collection(models.UserDevicesCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.NewInternalServerError("计算公开设备配置总数失败: " + err.Error())
	}
	
	// 查询设备列表
	cursor, err := s.db.Collection(models.UserDevicesCollection).Find(ctx, filter, &options.FindOptions{
		Skip:  &skipInt64,
		Limit: &limitInt64,
		Sort:  bson.D{{Key: "updatedAt", Value: -1}},
	})
	
	if err != nil {
		return nil, errors.NewInternalServerError("查询公开设备配置列表失败: " + err.Error())
	}
	defer cursor.Close(ctx)
	
	// 解析设备列表
	var userDevices []models.UserDevice
	if err = cursor.All(ctx, &userDevices); err != nil {
		return nil, errors.NewInternalServerError("解析公开设备配置列表失败: " + err.Error())
	}
	
	// 转换为响应格式
	response := &device.UserDeviceListResponse{
		Total:       int(total),
		Page:        page,
		PageSize:    pageSize,
		UserDevices: make([]device.UserDeviceResponse, len(userDevices)),
	}
	
	for i, ud := range userDevices {
		response.UserDevices[i] = device.UserDeviceResponse{
			ID:          ud.ID.Hex(),
			UserID:      ud.UserID.Hex(),
			Name:        ud.Name,
			Description: ud.Description,
			IsPublic:    ud.IsPublic,
			CreatedAt:   ud.CreatedAt,
			UpdatedAt:   ud.UpdatedAt,
			Devices:     make([]device.UserDeviceSettingsResponse, len(ud.Devices)),
		}
		
		for j, d := range ud.Devices {
			response.UserDevices[i].Devices[j] = device.UserDeviceSettingsResponse{
				DeviceID:   d.DeviceID.Hex(),
				DeviceType: string(d.DeviceType),
				Settings:   d.Settings,
			}
		}
	}
	
	return response, nil
}

// 为兼容性提供的方法
// GetMouseDevice 获取鼠标设备详情
func (s *ServiceImpl) GetMouseDevice(ctx context.Context, deviceID string) (*models.MouseDevice, error) {
	device, err := s.GetDeviceByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	
	// 将基础设备转换为鼠标设备
	mouseDevice := &models.MouseDevice{
		HardwareDevice: *device,
	}
	
	return mouseDevice, nil
}

// UpdateMouseDevice 更新鼠标设备
func (s *ServiceImpl) UpdateMouseDevice(ctx context.Context, deviceID string, request device.UpdateMouseRequest) (*models.MouseDevice, error) {
	// 实现待完成
	return nil, errors.NewInternalServerError("功能待实现")
}

// DeleteDevice 删除设备
func (s *ServiceImpl) DeleteDevice(ctx context.Context, deviceID string) error {
	// 实现待完成
	return errors.NewInternalServerError("功能待实现")
}

// ListDevices 列出设备
func (s *ServiceImpl) ListDevices(ctx context.Context, filter device.DeviceListFilter) (*device.DeviceListResponse, error) {
	return s.GetDeviceList(ctx, filter.Type, filter.Page, filter.PageSize)
}

// CreateDeviceReview 创建设备评测
func (s *ServiceImpl) CreateDeviceReview(ctx context.Context, userID string, request device.CreateReviewRequest) (*models.DeviceReview, error) {
	return s.SubmitReview(ctx, userID, device.SubmitReviewRequest(request))
}

// GetDeviceReview 获取设备评测
func (s *ServiceImpl) GetDeviceReview(ctx context.Context, reviewID string) (*models.DeviceReview, error) {
	return s.GetReviewByID(ctx, reviewID)
}

// UpdateDeviceReview 更新设备评测
func (s *ServiceImpl) UpdateDeviceReview(ctx context.Context, userID string, reviewID string, request device.UpdateReviewRequest) (*models.DeviceReview, error) {
	return s.UpdateReview(ctx, userID, reviewID, device.UpdateDeviceReviewRequest(request))
}

// DeleteDeviceReview 删除设备评测
func (s *ServiceImpl) DeleteDeviceReview(ctx context.Context, userID string, reviewID string) error {
	return s.DeleteReview(ctx, userID, reviewID)
}

// ListDeviceReviews 列出设备评测
func (s *ServiceImpl) ListDeviceReviews(ctx context.Context, request device.DeviceReviewListRequest) (*device.DeviceReviewListResponse, error) {
	return s.GetReviewsByDeviceID(ctx, request.DeviceID, request.Page, request.PageSize)
}

// GetUserDevice 获取用户设备
func (s *ServiceImpl) GetUserDevice(ctx context.Context, userID string, userDeviceID string) (*models.UserDevice, error) {
	return s.GetUserDeviceByID(ctx, userID, userDeviceID)
}

// DeleteUserDevice 删除用户设备
func (s *ServiceImpl) DeleteUserDevice(ctx context.Context, userID string, userDeviceID string) error {
	// 将ID转换为ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.NewBadRequestError("无效的用户ID")
	}
	
	deviceObjectID, err := primitive.ObjectIDFromHex(userDeviceID)
	if err != nil {
		return errors.NewBadRequestError("无效的设备ID")
	}
	
	// 删除用户设备配置
	result, err := s.db.Collection(models.UserDevicesCollection).DeleteOne(ctx, bson.M{
		"_id":    deviceObjectID,
		"userId": userObjectID,
	})
	
	if err != nil {
		return errors.NewInternalServerError("删除用户设备配置失败: " + err.Error())
	}
	
	if result.DeletedCount == 0 {
		return errors.NewNotFoundError("未找到用户设备配置")
	}
	
	return nil
}

// ListUserDevices 获取用户设备列表
func (s *ServiceImpl) ListUserDevices(ctx context.Context, request device.UserDeviceListRequest) (*device.UserDeviceListResponse, error) {
	if request.UserID == "" {
		return nil, errors.NewBadRequestError("用户ID不能为空")
	}
	return s.GetUserDevices(ctx, request.UserID, request.Page, request.PageSize)
}

// 以下是评测相关方法（在此省略，有需要时再实现）
func (s *ServiceImpl) SubmitReview(ctx context.Context, userID string, request device.SubmitReviewRequest) (*models.DeviceReview, error) {
	// 实现待完成
	return nil, errors.NewInternalServerError("功能待实现")
}

func (s *ServiceImpl) GetReviewsByDeviceID(ctx context.Context, deviceID string, page, pageSize int) (*device.DeviceReviewListResponse, error) {
	// 实现待完成
	return nil, errors.NewInternalServerError("功能待实现")
}

func (s *ServiceImpl) GetReviewsByUserID(ctx context.Context, userID string, page, pageSize int) (*device.DeviceReviewListResponse, error) {
	// 实现待完成
	return nil, errors.NewInternalServerError("功能待实现")
}

func (s *ServiceImpl) GetReviewByID(ctx context.Context, reviewID string) (*models.DeviceReview, error) {
	// 实现待完成
	return nil, errors.NewInternalServerError("功能待实现")
}

func (s *ServiceImpl) UpdateReview(ctx context.Context, userID, reviewID string, request device.UpdateDeviceReviewRequest) (*models.DeviceReview, error) {
	// 实现待完成
	return nil, errors.NewInternalServerError("功能待实现")
}

func (s *ServiceImpl) DeleteReview(ctx context.Context, userID, reviewID string) error {
	// 实现待完成
	return errors.NewInternalServerError("功能待实现")
}

func (s *ServiceImpl) ApproveReview(ctx context.Context, reviewerID, reviewID string) (*models.DeviceReview, error) {
	// 实现待完成
	return nil, errors.NewInternalServerError("功能待实现")
}

func (s *ServiceImpl) RejectReview(ctx context.Context, reviewerID, reviewID string, notes string) (*models.DeviceReview, error) {
	// 实现待完成
	return nil, errors.NewInternalServerError("功能待实现")
}

func (s *ServiceImpl) FeaturedReview(ctx context.Context, reviewerID, reviewID string, rank int) (*models.DeviceReview, error) {
	// 实现待完成
	return nil, errors.NewInternalServerError("功能待实现")
}

func (s *ServiceImpl) GetPendingReviews(ctx context.Context, reviewType string, page, pageSize int) (*device.DeviceReviewListResponse, error) {
	// 实现待完成
	return nil, errors.NewInternalServerError("功能待实现")
}

// 这个方法的实际实现在 similarity.go 文件中

// 这个方法的实际实现在 similarity.go 文件中