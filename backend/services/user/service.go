package user

import (
	"context"
	"project/backend/internal/errors"
	"project/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Service 用户服务接口
type Service interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	UpdateLastLogin(ctx context.Context, userID string) error
	GetUserProfile(ctx context.Context, userID string) (*models.User, error)
	UpdateUserProfile(ctx context.Context, userID string, user *models.User) error
}

// ServiceImpl 用户服务实现
type ServiceImpl struct {
	db *mongo.Database
}

// DefaultService 默认用户服务实现
type DefaultService struct {
	db *mongo.Database
}

// NewService 创建新的用户服务
func NewService(db *mongo.Database) Service {
	return &ServiceImpl{
		db: db,
	}
}

// NewDefaultService 创建默认用户服务
func NewDefaultService(db *mongo.Database) *DefaultService {
	return &DefaultService{
		db: db,
	}
}

// GetUserByID 通过ID获取用户
func (s *DefaultService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// UpdateLastLogin 更新最后登录时间
func (s *DefaultService) UpdateLastLogin(ctx context.Context, userID string) error {
	// 空实现，仅为了满足接口
	return nil
}

// GetUserProfile 获取用户资料
func (s *DefaultService) GetUserProfile(ctx context.Context, userID string) (*models.User, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// UpdateUserProfile 更新用户资料
func (s *DefaultService) UpdateUserProfile(ctx context.Context, userID string, user *models.User) error {
	// 空实现，仅为了满足接口
	return nil
}

// GetUserByID 通过ID获取用户
func (s *ServiceImpl) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的用户ID")
	}

	user := &models.User{}
	err = s.db.Collection("users").FindOne(ctx, bson.M{"_id": objID}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("用户不存在")
		}
		return nil, errors.NewInternalServerError("获取用户失败: " + err.Error())
	}

	return user, nil
}

// UpdateLastLogin 更新最后登录时间
func (s *ServiceImpl) UpdateLastLogin(ctx context.Context, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.NewBadRequestError("无效的用户ID")
	}

	now := time.Now()
	_, err = s.db.Collection("users").UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"stats.lastLoginAt": now}},
	)
	if err != nil {
		return errors.NewInternalServerError("更新登录时间失败: " + err.Error())
	}

	return nil
}

// GetUserProfile 获取用户资料
func (s *ServiceImpl) GetUserProfile(ctx context.Context, userID string) (*models.User, error) {
	return s.GetUserByID(ctx, userID)
}

// UpdateUserProfile 更新用户资料
func (s *ServiceImpl) UpdateUserProfile(ctx context.Context, userID string, user *models.User) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.NewBadRequestError("无效的用户ID")
	}

	_, err = s.db.Collection("users").ReplaceOne(ctx, bson.M{"_id": objID}, user)
	if err != nil {
		return errors.NewInternalServerError("更新用户资料失败: " + err.Error())
	}

	return nil
}