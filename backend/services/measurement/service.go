package measurement

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"project/backend/internal/errors"
	"project/backend/models"
	"project/backend/types/measurement"
	"time"
)

// Collection names
const (
	MeasurementsCollection        = "measurements"
	MeasurementUserStatsCollection = "measurement_user_stats"
)

// Service 是测量服务
type Service struct {
	db DatabaseInterface
}

// NewService 创建测量服务实例
func NewService(client *mongo.Client) *Service {
	var db DatabaseInterface
	if client != nil {
		dbAdapter := &DatabaseAdapter{DB: client.Database("cpc")}
		db = dbAdapter
	}
	
	return &Service{
		db: db,
	}
}

// CreateMeasurement 创建测量记录
func (s *Service) CreateMeasurement(ctx context.Context, userID string, req measurement.CreateMeasurementRequest) (*models.Measurement, error) {
	// 验证用户ID
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperrors.NewBadRequestError("无效的用户ID")
	}

	// 单位转换 - 统一为毫米 (mm)
	palmMm, lengthMm := req.Palm, req.Length
	switch req.Unit {
	case "cm":
		palmMm *= 10
		lengthMm *= 10
	case "inch":
		palmMm *= 25.4
		lengthMm *= 25.4
	}

	// 计算测量质量分数
	qualityScore := 70 // 基础分数
	if req.Calibrated {
		qualityScore += 15 // 校准加分
	}

	qualityLevel := models.QualityMedium
	if qualityScore >= 85 {
		qualityLevel = models.QualityHigh
	} else if qualityScore < 60 {
		qualityLevel = models.QualityLow
	}

	// 创建测量记录
	now := time.Now()
	measurement := &models.Measurement{
		ID:        primitive.NewObjectID(),
		UserID:    uid,
		Palm:      palmMm, 
		Length:    lengthMm,
		Unit:      "mm", // 存储始终用毫米
		Device:    req.Device,
		Calibrated: req.Calibrated,
		Quality:   &qualityLevel,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 插入数据库
	_, err = s.db.Collection(MeasurementsCollection).InsertOne(ctx, measurement)
	if err != nil {
		return nil, apperrors.NewInternalServerError("创建测量记录失败: " + err.Error())
	}

	// 异步更新用户统计
	go s.updateUserStats(context.Background(), uid)

	return measurement, nil
}

// GetMeasurement 获取单条测量记录
func (s *Service) GetMeasurement(ctx context.Context, userID string, measurementID string) (*models.Measurement, error) {
	// 验证用户ID
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperrors.NewBadRequestError("无效的用户ID")
	}

	// 验证测量记录ID
	mid, err := primitive.ObjectIDFromHex(measurementID)
	if err != nil {
		return nil, apperrors.NewBadRequestError("无效的测量记录ID")
	}

	// 查询条件
	filter := bson.M{
		"_id":    mid,
		"userId": uid,
	}

	// 查询数据库
	var measurement models.Measurement
	err = s.db.Collection(MeasurementsCollection).FindOne(ctx, filter).Decode(&measurement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperrors.NewNotFoundError("未找到测量记录")
		}
		return nil, apperrors.NewInternalServerError("获取测量记录失败: " + err.Error())
	}

	return &measurement, nil
}

// ListMeasurements 获取测量记录列表
func (s *Service) ListMeasurements(ctx context.Context, userID string, req measurement.MeasurementListRequest) (*measurement.MeasurementListResponse, error) {
	// 验证用户ID
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperrors.NewBadRequestError("无效的用户ID")
	}

	// 查询条件
	filter := bson.M{"userId": uid}

	// 分页
	skip := int64((req.Page - 1) * req.PageSize)
	limit := int64(req.PageSize)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"createdAt": -1}) // 最新的记录优先

	// 查询数据库
	cursor, err := s.db.Collection(MeasurementsCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, apperrors.NewInternalServerError("获取测量记录列表失败: " + err.Error())
	}
	defer cursor.Close(ctx)

	// 解析结果
	var measurements []models.Measurement
	if err = cursor.All(ctx, &measurements); err != nil {
		return nil, apperrors.NewInternalServerError("解析测量记录失败: " + err.Error())
	}

	// 获取总记录数
	total, err := s.db.Collection(MeasurementsCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, apperrors.NewInternalServerError("计算总记录数失败: " + err.Error())
	}

	// 转换为响应格式
	response := &measurement.MeasurementListResponse{
		Total:        int(total),
		Page:         req.Page,
		PageSize:     req.PageSize,
		Measurements: make([]measurement.MeasurementResponse, 0, len(measurements)),
	}

	for _, m := range measurements {
		response.Measurements = append(response.Measurements, measurement.MeasurementResponse{
			ID:        m.ID.Hex(),
			Palm:      m.Palm,
			Length:    m.Length,
			Unit:      m.Unit,
			Device:    m.Device,
			Quality:   m.Quality,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		})
	}

	return response, nil
}

// UpdateMeasurement 更新测量记录
func (s *Service) UpdateMeasurement(ctx context.Context, userID string, measurementID string, req measurement.UpdateMeasurementRequest) (*models.Measurement, error) {
	// 验证用户ID
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperrors.NewBadRequestError("无效的用户ID")
	}

	// 验证测量记录ID
	mid, err := primitive.ObjectIDFromHex(measurementID)
	if err != nil {
		return nil, apperrors.NewBadRequestError("无效的测量记录ID")
	}

	// 查询现有记录
	filter := bson.M{
		"_id":    mid,
		"userId": uid,
	}

	var existingMeasurement models.Measurement
	err = s.db.Collection(MeasurementsCollection).FindOne(ctx, filter).Decode(&existingMeasurement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperrors.NewNotFoundError("未找到测量记录")
		}
		return nil, apperrors.NewInternalServerError("获取测量记录失败: " + err.Error())
	}

	// 准备更新内容
	updateFields := bson.M{
		"updatedAt": time.Now(),
	}

	// 更新palm字段（如果提供）
	if req.Palm != nil {
		// 单位转换
		palmMm := *req.Palm
		switch req.Unit {
		case "cm":
			palmMm *= 10
		case "inch":
			palmMm *= 25.4
		}
		updateFields["palm"] = palmMm
	}

	// 更新length字段（如果提供）
	if req.Length != nil {
		// 单位转换
		lengthMm := *req.Length
		switch req.Unit {
		case "cm":
			lengthMm *= 10
		case "inch":
			lengthMm *= 25.4
		}
		updateFields["length"] = lengthMm
	}

	// 更新设备字段（如果提供）
	if req.Device != nil {
		updateFields["device"] = *req.Device
	}

	// 更新校准字段（如果提供）
	if req.Calibrated != nil {
		updateFields["calibrated"] = *req.Calibrated

		// 重新计算质量分数
		qualityScore := 70 // 基础分数
		if *req.Calibrated {
			qualityScore += 15 // 校准加分
		}

		qualityLevel := models.QualityMedium
		if qualityScore >= 85 {
			qualityLevel = models.QualityHigh
		} else if qualityScore < 60 {
			qualityLevel = models.QualityLow
		}

		updateFields["quality"] = qualityLevel
	}

	// 更新数据库
	update := bson.M{"$set": updateFields}
	_, err = s.db.Collection(MeasurementsCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, apperrors.NewInternalServerError("更新测量记录失败: " + err.Error())
	}

	// 获取更新后的记录
	var updatedMeasurement models.Measurement
	err = s.db.Collection(MeasurementsCollection).FindOne(ctx, filter).Decode(&updatedMeasurement)
	if err != nil {
		return nil, apperrors.NewInternalServerError("获取更新后的测量记录失败: " + err.Error())
	}

	// 异步更新用户统计
	go s.updateUserStats(context.Background(), uid)

	return &updatedMeasurement, nil
}

// DeleteMeasurement 删除测量记录
func (s *Service) DeleteMeasurement(ctx context.Context, userID string, measurementID string) error {
	// 验证用户ID
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperrors.NewBadRequestError("无效的用户ID")
	}

	// 验证测量记录ID
	mid, err := primitive.ObjectIDFromHex(measurementID)
	if err != nil {
		return apperrors.NewBadRequestError("无效的测量记录ID")
	}

	// 使用软删除 - 标记为已删除而不是真正删除
	filter := bson.M{
		"_id":    mid,
		"userId": uid,
	}

	update := bson.M{
		"$set": bson.M{
			"deleted":   true,
			"updatedAt": time.Now(),
		},
	}

	result, err := s.db.Collection(MeasurementsCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return apperrors.NewInternalServerError("删除测量记录失败: " + err.Error())
	}

	// 检查是否有记录被更新
	if result.ModifiedCount == 0 {
		return apperrors.NewNotFoundError("未找到测量记录")
	}

	// 异步更新用户统计
	go s.updateUserStats(context.Background(), uid)

	return nil
}

// GetUserStats 获取用户测量统计
func (s *Service) GetUserStats(ctx context.Context, userID string) (*models.MeasurementStats, error) {
	// 验证用户ID
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperrors.NewBadRequestError("无效的用户ID")
	}

	// 查询用户统计信息
	filter := bson.M{"userId": uid}
	var stats models.MeasurementStats
	err = s.db.Collection(MeasurementUserStatsCollection).FindOne(ctx, filter).Decode(&stats)
	if err != nil {
		// 如果没找到，计算并保存新的统计信息
		if err == mongo.ErrNoDocuments {
			stats, err = s.calculateAndSaveUserStats(ctx, uid)
			if err != nil {
				return nil, apperrors.NewInternalServerError("计算用户统计信息失败: " + err.Error())
			}
			return &stats, nil
		}
		return nil, apperrors.NewInternalServerError("获取用户统计信息失败: " + err.Error())
	}

	return &stats, nil
}

// GetRecommendations 获取设备推荐
func (s *Service) GetRecommendations(ctx context.Context, userID string) (*measurement.RecommendationResponse, error) {
	// 获取用户统计信息
	stats, err := s.GetUserStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 检查是否有足够的测量数据
	if stats.MeasurementCount == 0 || stats.AveragePalm == 0 || stats.AverageLength == 0 {
		return nil, apperrors.NewBadRequestError("没有足够的测量数据进行推荐")
	}

	// 确定握持类型
	ratio := stats.AverageLength / stats.AveragePalm
	gripType := models.GripStylePalm
	if ratio > 0.9 {
		gripType = models.GripStyleClaw
		if ratio > 0.95 {
			gripType = models.GripStyleFingertip
		}
	}

	// 构建推荐响应
	response := &measurement.RecommendationResponse{
		Palm:      stats.AveragePalm,
		Length:    stats.AverageLength,
		HandSize:  string(stats.HandSize),
		GripType:  string(gripType),
		Devices:   []measurement.DeviceRecommendation{},
	}

	// 设备推荐逻辑 - TODO: 实现具体的设备推荐算法
	// 这里仅返回手型和握持类型的信息，具体设备推荐需要另外实现

	return response, nil
}

// updateUserStats 更新用户统计信息
func (s *Service) updateUserStats(ctx context.Context, userID primitive.ObjectID) {
	// 防止在计算统计信息过程中出现错误导致程序崩溃
	defer func() {
		if r := recover(); r != nil {
			// 记录恢复的错误，但不中断主逻辑
			// 这里应该有日志记录，但是为了保持简单，暂不实现
		}
	}()

	// 计算并保存统计信息
	_, _ = s.calculateAndSaveUserStats(ctx, userID)
}

// calculateAndSaveUserStats 计算并保存用户统计信息
func (s *Service) calculateAndSaveUserStats(ctx context.Context, userID primitive.ObjectID) (models.MeasurementStats, error) {
	// 查询条件 - 只查询未删除的记录
	filter := bson.M{
		"userId":  userID,
		"deleted": bson.M{"$ne": true},
	}

	// 聚合管道
	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id":          "$userId",
			"averagePalm":  bson.M{"$avg": "$palm"},
			"averageLength": bson.M{"$avg": "$length"},
			"count":        bson.M{"$sum": 1},
			"lastMeasured": bson.M{"$max": "$createdAt"},
		}},
	}

	// 执行聚合
	cursor, err := s.db.Collection(MeasurementsCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return models.MeasurementStats{}, err
	}
	defer cursor.Close(ctx)

	// 解析结果
	type aggregateResult struct {
		ID            primitive.ObjectID `bson:"_id"`
		AveragePalm   float64           `bson:"averagePalm"`
		AverageLength float64           `bson:"averageLength"`
		Count         int               `bson:"count"`
		LastMeasured  time.Time         `bson:"lastMeasured"`
	}

	var results []aggregateResult
	if err = cursor.All(ctx, &results); err != nil {
		return models.MeasurementStats{}, err
	}

	// 初始化默认值
	stats := models.MeasurementStats{
		UserID:           userID,
		AveragePalm:      0,
		AverageLength:    0,
		HandSize:         "unknown",
		MeasurementCount: 0,
		LastMeasuredAt:   time.Time{},
		UpdatedAt:        time.Now(),
	}

	// 如果有测量数据，设置统计信息
	if len(results) > 0 {
		agg := results[0]
		stats.AveragePalm = agg.AveragePalm
		stats.AverageLength = agg.AverageLength
		stats.MeasurementCount = agg.Count
		stats.LastMeasuredAt = agg.LastMeasured

		// 根据手掌宽度确定手型大小（单位: mm）
		if agg.AveragePalm < 80 {
			stats.HandSize = models.HandSizeSmall
		} else if agg.AveragePalm > 95 {
			stats.HandSize = models.HandSizeLarge
		} else {
			stats.HandSize = models.HandSizeMedium
		}
	}

	// 保存或更新统计信息
	filter = bson.M{"userId": userID}
	update := bson.M{"$set": stats}
	opts := options.Update().SetUpsert(true)

	_, err = s.db.Collection(MeasurementUserStatsCollection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return models.MeasurementStats{}, err
	}

	return stats, nil
}