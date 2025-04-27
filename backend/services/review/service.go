package review

import (
	"project/backend/internal/errors"
	"project/backend/models"
	"project/backend/types/review"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// ServiceImpl 评测服务实现
type ServiceImpl struct {
	db *mongo.Database
}

// DefaultService 默认评测服务实现
type DefaultService struct {
	db *mongo.Database
}

// New 创建新的评测服务
func New(db *mongo.Database) Service {
	return &ServiceImpl{
		db: db,
	}
}

// NewDefaultService 创建默认评测服务
func NewDefaultService(db *mongo.Database) *DefaultService {
	return &DefaultService{
		db: db,
	}
}

// ApproveReview 批准评测
func (s *DefaultService) ApproveReview(ctx context.Context, reviewerID, reviewID string, notes string) (*models.Review, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// CreateReview 创建评测
func (s *DefaultService) CreateReview(ctx context.Context, userID string, request review.CreateReviewRequest) (*models.Review, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetReview 获取评测
func (s *DefaultService) GetReview(ctx context.Context, reviewID string) (*models.Review, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// UpdateReview 更新评测
func (s *DefaultService) UpdateReview(ctx context.Context, userID, reviewID string, request review.UpdateReviewRequest) (*models.Review, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// DeleteReview 删除评测
func (s *DefaultService) DeleteReview(ctx context.Context, userID, reviewID string) error {
	// 空实现，仅为了满足接口
	return nil
}

// ListReviews 列出评测
func (s *DefaultService) ListReviews(ctx context.Context, request review.ReviewListRequest) (*review.ReviewListResponse, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// RejectReview 拒绝评测
func (s *DefaultService) RejectReview(ctx context.Context, reviewerID, reviewID string, notes string) (*models.Review, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// FeaturedReview 推荐评测
func (s *DefaultService) FeaturedReview(ctx context.Context, reviewerID, reviewID string, rank int) (*models.Review, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetPendingReviews 获取待审核评测
func (s *DefaultService) GetPendingReviews(ctx context.Context, reviewType string, page, pageSize int) (*review.ReviewListResponse, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// GetUserReviewStats 获取用户评测统计
func (s *DefaultService) GetUserReviewStats(ctx context.Context, userID string) (*review.UserReviewStats, error) {
	// 空实现，仅为了满足接口
	return nil, nil
}

// IncrementViewCount 增加查看次数
func (s *DefaultService) IncrementViewCount(ctx context.Context, reviewID string) error {
	// 空实现，仅为了满足接口
	return nil
}

// CreateReview 创建评测
func (s *ServiceImpl) CreateReview(ctx context.Context, userID string, request review.CreateReviewRequest) (*models.Review, error) {
	// 将用户ID转换为ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的用户ID")
	}
	
	// 将外部项目ID转换为ObjectID
	externalItemID, err := primitive.ObjectIDFromHex(request.ExternalItemID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的项目ID")
	}
	
	// 创建评测
	now := time.Now()
	reviewModel := models.Review{
		ID:             primitive.NewObjectID(),
		ExternalItemID: externalItemID,
		ItemType:       request.ItemType,
		UserID:         userObjectID,
		Content:        request.Content,
		Pros:           request.Pros,
		Cons:           request.Cons,
		Score:          request.Score,
		Usage:          request.Usage,
		Status:         models.ReviewStatusPending,
		Type:           models.ReviewType(request.Type),
		ContentType:    models.ReviewContentType(request.ContentType),
		ViewCount:      0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	
	// 保存到数据库
	_, err = s.db.Collection(models.ReviewsCollection).InsertOne(ctx, reviewModel)
	if err != nil {
		return nil, errors.NewInternalServerError("创建评测失败: " + err.Error())
	}
	
	return &reviewModel, nil
}

// GetReview 获取评测
func (s *ServiceImpl) GetReview(ctx context.Context, reviewID string) (*models.Review, error) {
	// 将评测ID转换为ObjectID
	id, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的评测ID")
	}
	
	// 查询评测
	var reviewModel models.Review
	err = s.db.Collection(models.ReviewsCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&reviewModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("评测不存在")
		}
		return nil, errors.NewInternalServerError("获取评测失败: " + err.Error())
	}
	
	return &reviewModel, nil
}

// ListReviews 列出评测
func (s *ServiceImpl) ListReviews(ctx context.Context, request review.ReviewListRequest) (*review.ReviewListResponse, error) {
	if request.ExternalItemID != "" {
		result, err := s.getReviewsByItemID(ctx, request.ExternalItemID, request.ItemType, request.Page, request.PageSize, request.SortBy)
		return result, err
	} else if request.UserID != "" {
		result, err := s.getReviewsByUserID(ctx, request.UserID, request.Page, request.PageSize)
		return result, err
	}
	
	// 默认返回全部评测
	return s.GetPendingReviews(ctx, request.ItemType, request.Page, request.PageSize)
}

// UpdateReview 更新评测
func (s *ServiceImpl) UpdateReview(ctx context.Context, userID, reviewID string, request review.UpdateReviewRequest) (*models.Review, error) {
	// 将ID转换为ObjectID
	reviewObjectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的评测ID")
	}
	
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的用户ID")
	}
	
	// 查找当前评测
	var currentReview models.Review
	err = s.db.Collection(models.ReviewsCollection).FindOne(ctx, bson.M{
		"_id":       reviewObjectID,
		"userId":    userObjectID,
		"deletedAt": nil,
	}).Decode(&currentReview)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("未找到评测，或者您无权修改此评测")
		}
		return nil, errors.NewInternalServerError("获取评测失败: " + err.Error())
	}
	
	// 检查评测状态，已审核通过的评测不能修改
	if currentReview.Status != models.ReviewStatusPending {
		return nil, errors.NewBadRequestError("已审核的评测不能修改")
	}
	
	// 更新字段
	if request.Content != nil && *request.Content != "" {
		currentReview.Content = *request.Content
	}
	
	if request.Score != nil && *request.Score > 0 {
		currentReview.Score = *request.Score
	}
	
	if request.Usage != nil && *request.Usage != "" {
		currentReview.Usage = *request.Usage
	}
	
	if request.Pros != nil && len(*request.Pros) > 0 {
		currentReview.Pros = *request.Pros
	}
	
	if request.Cons != nil && len(*request.Cons) > 0 {
		currentReview.Cons = *request.Cons
	}
	
	currentReview.UpdatedAt = time.Now()
	
	// 保存到数据库
	_, err = s.db.Collection(models.ReviewsCollection).ReplaceOne(ctx, bson.M{
		"_id":    reviewObjectID,
		"userId": userObjectID,
	}, currentReview)
	
	if err != nil {
		return nil, errors.NewInternalServerError("更新评测失败: " + err.Error())
	}
	
	return &currentReview, nil
}

// DeleteReview 删除评测
func (s *ServiceImpl) DeleteReview(ctx context.Context, userID, reviewID string) error {
	// 将ID转换为ObjectID
	reviewObjectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return errors.NewBadRequestError("无效的评测ID")
	}
	
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.NewBadRequestError("无效的用户ID")
	}
	
	// 查找当前评测
	var currentReview models.Review
	err = s.db.Collection(models.ReviewsCollection).FindOne(ctx, bson.M{
		"_id":       reviewObjectID,
		"userId":    userObjectID,
		"deletedAt": nil,
	}).Decode(&currentReview)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.NewNotFoundError("未找到评测，或者您无权删除此评测")
		}
		return errors.NewInternalServerError("获取评测失败: " + err.Error())
	}
	
	// 检查评测状态，已审核通过的评测不能删除
	if currentReview.Status != models.ReviewStatusPending {
		return errors.NewBadRequestError("已审核的评测不能删除")
	}
	
	// 软删除
	now := time.Now()
	_, err = s.db.Collection(models.ReviewsCollection).UpdateOne(
		ctx,
		bson.M{
			"_id":       reviewObjectID,
			"userId":    userObjectID,
			"deletedAt": nil,
		},
		bson.M{
			"$set": bson.M{
				"deletedAt": &now,
				"updatedAt": now,
			},
		},
	)
	
	if err != nil {
		return errors.NewInternalServerError("删除评测失败: " + err.Error())
	}
	
	return nil
}

// ApproveReview 批准评测
func (s *ServiceImpl) ApproveReview(ctx context.Context, reviewerID, reviewID string, notes string) (*models.Review, error) {
	// 将ID转换为ObjectID
	reviewObjectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的评测ID")
	}
	
	reviewerObjectID, err := primitive.ObjectIDFromHex(reviewerID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的审核员ID")
	}
	
	// 查找当前评测
	var currentReview models.Review
	err = s.db.Collection(models.ReviewsCollection).FindOne(ctx, bson.M{
		"_id":       reviewObjectID,
		"status":    models.ReviewStatusPending,
		"deletedAt": nil,
	}).Decode(&currentReview)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("未找到待审核的评测")
		}
		return nil, errors.NewInternalServerError("获取评测失败: " + err.Error())
	}
	
	// 更新评测状态
	now := time.Now()
	currentReview.Status = models.ReviewStatusApproved
	currentReview.ReviewerID = &reviewerObjectID
	currentReview.ReviewedAt = &now
	currentReview.PublishedAt = &now
	currentReview.UpdatedAt = now
	
	// 保存到数据库
	_, err = s.db.Collection(models.ReviewsCollection).ReplaceOne(ctx, bson.M{
		"_id":    reviewObjectID,
		"status": models.ReviewStatusPending,
	}, currentReview)
	
	if err != nil {
		return nil, errors.NewInternalServerError("更新评测状态失败: " + err.Error())
	}
	
	return &currentReview, nil
}

// RejectReview 拒绝评测
func (s *ServiceImpl) RejectReview(ctx context.Context, reviewerID, reviewID string, notes string) (*models.Review, error) {
	// 将ID转换为ObjectID
	reviewObjectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的评测ID")
	}
	
	reviewerObjectID, err := primitive.ObjectIDFromHex(reviewerID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的审核员ID")
	}
	
	// 查找当前评测
	var currentReview models.Review
	err = s.db.Collection(models.ReviewsCollection).FindOne(ctx, bson.M{
		"_id":       reviewObjectID,
		"status":    models.ReviewStatusPending,
		"deletedAt": nil,
	}).Decode(&currentReview)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("未找到待审核的评测")
		}
		return nil, errors.NewInternalServerError("获取评测失败: " + err.Error())
	}
	
	// 更新评测状态
	now := time.Now()
	currentReview.Status = models.ReviewStatusRejected
	currentReview.ReviewerID = &reviewerObjectID
	currentReview.ReviewerNotes = notes
	currentReview.ReviewedAt = &now
	currentReview.UpdatedAt = now
	
	// 保存到数据库
	_, err = s.db.Collection(models.ReviewsCollection).ReplaceOne(ctx, bson.M{
		"_id":    reviewObjectID,
		"status": models.ReviewStatusPending,
	}, currentReview)
	
	if err != nil {
		return nil, errors.NewInternalServerError("更新评测状态失败: " + err.Error())
	}
	
	return &currentReview, nil
}

// FeaturedReview 设置推荐评测
func (s *ServiceImpl) FeaturedReview(ctx context.Context, reviewerID, reviewID string, rank int) (*models.Review, error) {
	// 将ID转换为ObjectID
	reviewObjectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的评测ID")
	}
	
	reviewerObjectID, err := primitive.ObjectIDFromHex(reviewerID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的审核员ID")
	}
	
	// 查找当前评测
	var currentReview models.Review
	err = s.db.Collection(models.ReviewsCollection).FindOne(ctx, bson.M{
		"_id":       reviewObjectID,
		"status":    models.ReviewStatusApproved,
		"deletedAt": nil,
	}).Decode(&currentReview)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("未找到已批准的评测")
		}
		return nil, errors.NewInternalServerError("获取评测失败: " + err.Error())
	}
	
	// 更新评测状态
	now := time.Now()
	currentReview.Status = models.ReviewStatusFeatured
	currentReview.ReviewerID = &reviewerObjectID
	currentReview.FeaturedRank = &rank
	currentReview.UpdatedAt = now
	
	// 保存到数据库
	_, err = s.db.Collection(models.ReviewsCollection).ReplaceOne(ctx, bson.M{
		"_id":    reviewObjectID,
		"status": models.ReviewStatusApproved,
	}, currentReview)
	
	if err != nil {
		return nil, errors.NewInternalServerError("更新评测状态失败: " + err.Error())
	}
	
	return &currentReview, nil
}

// GetPendingReviews 获取待审核的评测
func (s *ServiceImpl) GetPendingReviews(ctx context.Context, reviewType string, page, pageSize int) (*review.ReviewListResponse, error) {
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
	filter := bson.M{
		"status":    models.ReviewStatusPending,
		"deletedAt": nil,
	}
	
	// 计算总数
	total, err := s.db.Collection(models.ReviewsCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.NewInternalServerError("计算待审核评测总数失败: " + err.Error())
	}
	
	// 查询评测列表
	cursor, err := s.db.Collection(models.ReviewsCollection).Find(ctx, filter, &options.FindOptions{
		Skip:  &skipInt64,
		Limit: &limitInt64,
		Sort:  bson.D{{Key: "createdAt", Value: 1}},
	})
	if err != nil {
		return nil, errors.NewInternalServerError("查询待审核评测列表失败: " + err.Error())
	}
	defer cursor.Close(ctx)
	
	// 解析评测列表
	var reviews []models.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, errors.NewInternalServerError("解析待审核评测列表失败: " + err.Error())
	}
	
	// 转换为响应格式
	response := &review.ReviewListResponse{
		Total:    int(total),
		Page:     page,
		PageSize: pageSize,
		Reviews:  make([]review.ReviewResponse, len(reviews)),
	}
	
	for i, r := range reviews {
		response.Reviews[i] = mapReviewToResponse(&r)
	}
	
	return response, nil
}

// GetUserReviewStats 获取用户评测统计
func (s *ServiceImpl) GetUserReviewStats(ctx context.Context, userID string) (*review.UserReviewStats, error) {
	// 实现统计逻辑 - 简单实现
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的用户ID")
	}
	
	// 构建基本过滤器
	filter := bson.M{
		"userId":    userObjectID,
		"deletedAt": nil,
	}
	
	// 计算总数
	totalCount, err := s.db.Collection(models.ReviewsCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.NewInternalServerError("计算评测总数失败: " + err.Error())
	}
	
	// 计算已批准的评测
	approvedFilter := bson.M{
		"userId":    userObjectID,
		"status":    models.ReviewStatusApproved,
		"deletedAt": nil,
	}
	approvedCount, err := s.db.Collection(models.ReviewsCollection).CountDocuments(ctx, approvedFilter)
	if err != nil {
		return nil, errors.NewInternalServerError("计算已批准评测总数失败: " + err.Error())
	}
	
	// 计算待审核的评测
	pendingFilter := bson.M{
		"userId":    userObjectID,
		"status":    models.ReviewStatusPending,
		"deletedAt": nil,
	}
	pendingCount, err := s.db.Collection(models.ReviewsCollection).CountDocuments(ctx, pendingFilter)
	if err != nil {
		return nil, errors.NewInternalServerError("计算待审核评测总数失败: " + err.Error())
	}
	
	// 计算已拒绝的评测
	rejectedFilter := bson.M{
		"userId":    userObjectID,
		"status":    models.ReviewStatusRejected,
		"deletedAt": nil,
	}
	rejectedCount, err := s.db.Collection(models.ReviewsCollection).CountDocuments(ctx, rejectedFilter)
	if err != nil {
		return nil, errors.NewInternalServerError("计算已拒绝评测总数失败: " + err.Error())
	}
	
	// 计算推荐的评测
	featuredFilter := bson.M{
		"userId":    userObjectID,
		"status":    models.ReviewStatusFeatured,
		"deletedAt": nil,
	}
	featuredCount, err := s.db.Collection(models.ReviewsCollection).CountDocuments(ctx, featuredFilter)
	if err != nil {
		return nil, errors.NewInternalServerError("计算推荐评测总数失败: " + err.Error())
	}
	
	// 返回统计结果
	return &review.UserReviewStats{
		TotalReviews:   int(totalCount),
		ApprovedCount:  int(approvedCount),
		PendingCount:   int(pendingCount),
		RejectedCount:  int(rejectedCount),
		FeaturedCount:  int(featuredCount),
		AverageScore:   0, // 暂不计算平均分
		TotalViewCount: 0, // 暂不计算总浏览量
	}, nil
}

// IncrementViewCount 增加评测浏览次数
func (s *ServiceImpl) IncrementViewCount(ctx context.Context, reviewID string) error {
	// 将评测ID转换为ObjectID
	id, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return errors.NewBadRequestError("无效的评测ID")
	}
	
	// 更新浏览次数
	_, err = s.db.Collection(models.ReviewsCollection).UpdateOne(
		ctx,
		bson.M{
			"_id": id,
			"deletedAt": nil,
			"status": bson.M{
				"$in": []string{string(models.ReviewStatusApproved), string(models.ReviewStatusFeatured)},
			},
		},
		bson.M{
			"$inc": bson.M{
				"viewCount": 1,
			},
		},
	)
	
	if err != nil {
		return errors.NewInternalServerError("增加查看次数失败: " + err.Error())
	}
	
	return nil
}

// 获取用户的评测列表 - 内部方法
func (s *ServiceImpl) getReviewsByUserID(ctx context.Context, userID string, page, pageSize int) (*review.ReviewListResponse, error) {
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
	filter := bson.M{
		"userId":    userObjectID,
		"deletedAt": nil,
	}
	
	// 计算总数
	total, err := s.db.Collection(models.ReviewsCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.NewInternalServerError("计算评测总数失败: " + err.Error())
	}
	
	// 查询评测列表
	cursor, err := s.db.Collection(models.ReviewsCollection).Find(ctx, filter, &options.FindOptions{
		Skip:  &skipInt64,
		Limit: &limitInt64,
		Sort:  bson.D{{Key: "createdAt", Value: -1}},
	})
	if err != nil {
		return nil, errors.NewInternalServerError("查询评测列表失败: " + err.Error())
	}
	defer cursor.Close(ctx)
	
	// 解析评测列表
	var reviews []models.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, errors.NewInternalServerError("解析评测列表失败: " + err.Error())
	}
	
	// 转换为响应格式
	response := &review.ReviewListResponse{
		Total:    int(total),
		Page:     page,
		PageSize: pageSize,
		Reviews:  make([]review.ReviewResponse, len(reviews)),
	}
	
	for i, r := range reviews {
		response.Reviews[i] = mapReviewToResponse(&r)
	}
	
	return response, nil
}

// 获取项目的评测列表 - 内部方法
func (s *ServiceImpl) getReviewsByItemID(ctx context.Context, itemID string, itemType string, page, pageSize int, sortBy string) (*review.ReviewListResponse, error) {
	// 将项目ID转换为ObjectID
	itemObjectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的项目ID")
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
	filter := bson.M{
		"externalItemId": itemObjectID,
		"itemType":       itemType,
		"status":         models.ReviewStatusApproved,
		"deletedAt":      nil,
	}
	
	// 计算总数
	total, err := s.db.Collection(models.ReviewsCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.NewInternalServerError("计算评测总数失败: " + err.Error())
	}
	
	// 确定排序方式
	sortOpt := bson.D{{Key: "createdAt", Value: -1}} // 默认按创建时间降序
	
	if sortBy == "score" {
		sortOpt = bson.D{{Key: "score", Value: -1}, {Key: "createdAt", Value: -1}}
	} else if sortBy == "views" {
		sortOpt = bson.D{{Key: "viewCount", Value: -1}, {Key: "createdAt", Value: -1}}
	}
	
	// 查询评测列表
	cursor, err := s.db.Collection(models.ReviewsCollection).Find(ctx, filter, &options.FindOptions{
		Skip:  &skipInt64,
		Limit: &limitInt64,
		Sort:  sortOpt,
	})
	if err != nil {
		return nil, errors.NewInternalServerError("查询评测列表失败: " + err.Error())
	}
	defer cursor.Close(ctx)
	
	// 解析评测列表
	var reviews []models.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, errors.NewInternalServerError("解析评测列表失败: " + err.Error())
	}
	
	// 转换为响应格式
	response := &review.ReviewListResponse{
		Total:    int(total),
		Page:     page,
		PageSize: pageSize,
		Reviews:  make([]review.ReviewResponse, len(reviews)),
	}
	
	for i, r := range reviews {
		response.Reviews[i] = mapReviewToResponse(&r)
	}
	
	return response, nil
}

// 获取推荐评测 - 内部方法
func (s *ServiceImpl) getFeaturedReviews(ctx context.Context, itemType string, limit int) (*review.ReviewListResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	// 构建查询
	filter := bson.M{
		"status":      models.ReviewStatusFeatured,
		"deletedAt":   nil,
		"featuredRank": bson.M{"$ne": nil},
	}
	
	if itemType != "" {
		filter["type"] = itemType
	}
	
	limitInt64 := int64(limit)
	
	// 查询评测列表
	cursor, err := s.db.Collection(models.ReviewsCollection).Find(ctx, filter, &options.FindOptions{
		Limit: &limitInt64,
		Sort:  bson.D{{Key: "featuredRank", Value: 1}, {Key: "createdAt", Value: -1}},
	})
	if err != nil {
		return nil, errors.NewInternalServerError("查询推荐评测失败: " + err.Error())
	}
	defer cursor.Close(ctx)
	
	// 解析评测列表
	var reviews []models.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, errors.NewInternalServerError("解析推荐评测失败: " + err.Error())
	}
	
	// 转换为响应格式
	response := &review.ReviewListResponse{
		Total:    len(reviews),
		Page:     1,
		PageSize: limit,
		Reviews:  make([]review.ReviewResponse, len(reviews)),
	}
	
	for i, r := range reviews {
		response.Reviews[i] = mapReviewToResponse(&r)
	}
	
	return response, nil
}

// mapReviewToResponse 将Review模型转换为ReviewResponse
func mapReviewToResponse(r *models.Review) review.ReviewResponse {
	response := review.ReviewResponse{
		ID:             r.ID.Hex(),
		ExternalItemID: r.ExternalItemID.Hex(),
		ItemType:       r.ItemType,
		UserID:         r.UserID.Hex(),
		Content:        r.Content,
		Pros:           r.Pros,
		Cons:           r.Cons,
		Score:          r.Score,
		Usage:          r.Usage,
		Status:         string(r.Status),
		Type:           string(r.Type),
		ContentType:    string(r.ContentType),
		ViewCount:      r.ViewCount,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
	
	if r.ReviewerID != nil {
		response.ReviewerID = r.ReviewerID.Hex()
	}
	
	response.ReviewerNotes = r.ReviewerNotes
	response.ReviewedAt = r.ReviewedAt
	response.PublishedAt = r.PublishedAt
	response.FeaturedRank = r.FeaturedRank
	
	return response
}