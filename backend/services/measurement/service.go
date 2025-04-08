package measurement

import (
	"FullStackOfYear/backend/internal/errors"
	"FullStackOfYear/backend/models"
	"FullStackOfYear/backend/types/measurement"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Service 定义测量服务接口
type Service interface {
	// 基础CRUD操作
	CreateMeasurement(ctx context.Context, userID string, request measurement.CreateMeasurementRequest) (*models.Measurement, error)
	GetMeasurement(ctx context.Context, userID, measurementID string) (*models.Measurement, error)
	UpdateMeasurement(ctx context.Context, userID, measurementID string, request measurement.UpdateMeasurementRequest) (*models.Measurement, error)
	DeleteMeasurement(ctx context.Context, userID, measurementID string) error
	ListMeasurements(ctx context.Context, userID string, request measurement.MeasurementListRequest) (*measurement.MeasurementListResponse, error)

	// 统计和分析
	GetUserStats(ctx context.Context, userID string) (*models.MeasurementUserStats, error)
	GetRecommendations(ctx context.Context, userID string) (*measurement.MeasurementRecommendationResponse, error)
}

// 服务实现
type service struct {
	db DatabaseInterface
}

// New 创建新的测量服务
func New(db *mongo.Database) Service {
	return &service{
		db: &DatabaseAdapter{DB: db},
	}
}

// CreateMeasurement 创建新的测量记录
func (s *service) CreateMeasurement(ctx context.Context, userID string, request measurement.CreateMeasurementRequest) (*models.Measurement, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的用户ID")
	}

	// 转换单位到毫米
	palm := request.Palm
	length := request.Length

	if request.Unit == "cm" {
		palm = palm * 10
		length = length * 10
	} else if request.Unit == "inch" {
		palm = palm * 25.4
		length = length * 25.4
	}

	// 创建新的测量记录
	now := time.Now()
	measurement := &models.Measurement{
		ID:     primitive.NewObjectID(),
		UserID: userObjID,
		Measurements: models.MeasurementData{
			Palm:   palm,
			Length: length,
			Unit:   request.Unit,
		},
		Metadata: models.MeasurementMetadata{
			Device: request.Device,
		},
		Quality: models.MeasurementQuality{
			Score: calculateQualityScore(request.Calibrated),
			Factors: models.MeasurementFactors{
				Calibration: request.Calibrated,
				Stability:   1.0, // 默认值
				Consistency: 1.0, // 默认值
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 插入数据库
	_, err = s.db.Collection(models.MeasurementsCollection).InsertOne(ctx, measurement)
	if err != nil {
		return nil, errors.NewInternalServerError("创建测量记录失败: " + err.Error())
	}

	// 更新用户统计
	go s.updateUserStats(context.Background(), userObjID)

	return measurement, nil
}

// GetMeasurement 获取单条测量记录
func (s *service) GetMeasurement(ctx context.Context, userID, measurementID string) (*models.Measurement, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的用户ID")
	}

	measurementObjID, err := primitive.ObjectIDFromHex(measurementID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的测量记录ID")
	}

	// 查询条件：匹配ID和用户ID，且未删除
	filter := bson.M{
		"_id":       measurementObjID,
		"userId":    userObjID,
		"deletedAt": bson.M{"$exists": false},
	}

	var result models.Measurement
	err = s.db.Collection(models.MeasurementsCollection).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("未找到测量记录")
		}
		return nil, errors.NewInternalServerError("获取测量记录失败: " + err.Error())
	}

	return &result, nil
}

// UpdateMeasurement 更新测量记录
func (s *service) UpdateMeasurement(ctx context.Context, userID, measurementID string, request measurement.UpdateMeasurementRequest) (*models.Measurement, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的用户ID")
	}

	measurementObjID, err := primitive.ObjectIDFromHex(measurementID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的测量记录ID")
	}

	// 查询现有记录
	filter := bson.M{
		"_id":       measurementObjID,
		"userId":    userObjID,
		"deletedAt": bson.M{"$exists": false},
	}

	var existingMeasurement models.Measurement
	err = s.db.Collection(models.MeasurementsCollection).FindOne(ctx, filter).Decode(&existingMeasurement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("未找到测量记录")
		}
		return nil, errors.NewInternalServerError("获取测量记录失败: " + err.Error())
	}

	// 准备更新数据
	updateData := bson.M{
		"updatedAt": time.Now(),
	}

	// 更新手掌宽度
	if request.Palm != nil {
		palm := *request.Palm
		// 单位转换
		if request.Unit != nil && *request.Unit != "mm" {
			if *request.Unit == "cm" {
				palm = palm * 10
			} else if *request.Unit == "inch" {
				palm = palm * 25.4
			}
		}
		updateData["measurements.palm"] = palm
	}

	// 更新手指长度
	if request.Length != nil {
		length := *request.Length
		// 单位转换
		if request.Unit != nil && *request.Unit != "mm" {
			if *request.Unit == "cm" {
				length = length * 10
			} else if *request.Unit == "inch" {
				length = length * 25.4
			}
		}
		updateData["measurements.length"] = length
	}

	// 更新单位
	if request.Unit != nil {
		updateData["measurements.unit"] = *request.Unit
	}

	// 更新校准状态
	if request.Calibrated != nil {
		updateData["quality.factors.calibration"] = *request.Calibrated
		updateData["quality.score"] = calculateQualityScore(*request.Calibrated)
	}

	// 执行更新
	update := bson.M{"$set": updateData}
	_, err = s.db.Collection(models.MeasurementsCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, errors.NewInternalServerError("更新测量记录失败: " + err.Error())
	}

	// 重新查询更新后的记录
	var updatedMeasurement models.Measurement
	err = s.db.Collection(models.MeasurementsCollection).FindOne(ctx, filter).Decode(&updatedMeasurement)
	if err != nil {
		return nil, errors.NewInternalServerError("获取更新后的测量记录失败: " + err.Error())
	}

	// 更新用户统计
	go s.updateUserStats(context.Background(), userObjID)

	return &updatedMeasurement, nil
}

// DeleteMeasurement 删除测量记录(软删除)
func (s *service) DeleteMeasurement(ctx context.Context, userID, measurementID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.NewBadRequestError("无效的用户ID")
	}

	measurementObjID, err := primitive.ObjectIDFromHex(measurementID)
	if err != nil {
		return errors.NewBadRequestError("无效的测量记录ID")
	}

	filter := bson.M{
		"_id":       measurementObjID,
		"userId":    userObjID,
		"deletedAt": bson.M{"$exists": false},
	}

	update := bson.M{
		"$set": bson.M{
			"deletedAt": time.Now(),
		},
	}

	result, err := s.db.Collection(models.MeasurementsCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewInternalServerError("删除测量记录失败: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return errors.NewNotFoundError("未找到测量记录或已被删除")
	}

	// 更新用户统计
	go s.updateUserStats(context.Background(), userObjID)

	return nil
}

// ListMeasurements 获取测量记录列表
func (s *service) ListMeasurements(ctx context.Context, userID string, request measurement.MeasurementListRequest) (*measurement.MeasurementListResponse, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的用户ID")
	}

	// 构建查询过滤条件
	filter := bson.M{
		"userId":    userObjID,
		"deletedAt": bson.M{"$exists": false},
	}

	// 添加日期过滤
	if request.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", request.StartDate)
		if err == nil {
			filter["createdAt"] = bson.M{"$gte": startDate}
		}
	}
	if request.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", request.EndDate)
		if err == nil {
			if _, exists := filter["createdAt"]; exists {
				filter["createdAt"].(bson.M)["$lte"] = endDate.Add(24 * time.Hour)
			} else {
				filter["createdAt"] = bson.M{"$lte": endDate.Add(24 * time.Hour)}
			}
		}
	}

	// 构建排序条件
	sortField := "createdAt"
	if request.SortBy != "" {
		switch request.SortBy {
		case "palm":
			sortField = "measurements.palm"
		case "length":
			sortField = "measurements.length"
		case "quality":
			sortField = "quality.score"
		}
	}

	sortOrder := -1 // 默认降序
	if request.SortOrder == "asc" {
		sortOrder = 1
	}

	// 设置分页
	page := 1
	if request.Page > 0 {
		page = request.Page
	}

	pageSize := 20
	if request.PageSize > 0 && request.PageSize <= 100 {
		pageSize = request.PageSize
	}

	skip := (page - 1) * pageSize

	// 查询选项
	findOptions := options.Find().
		SetSort(bson.D{{Key: sortField, Value: sortOrder}}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	// 执行查询
	cursor, err := s.db.Collection(models.MeasurementsCollection).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, errors.NewInternalServerError("查询测量记录失败: " + err.Error())
	}
	defer cursor.Close(ctx)

	// 统计总数
	total, err := s.db.Collection(models.MeasurementsCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.NewInternalServerError("计算总记录数失败: " + err.Error())
	}

	// 解析结果
	var measurements []models.Measurement
	if err = cursor.All(ctx, &measurements); err != nil {
		return nil, errors.NewInternalServerError("解析测量记录失败: " + err.Error())
	}

	// 构建响应
	var response measurement.MeasurementListResponse
	response.Total = int(total)
	response.Page = page
	response.PageSize = pageSize
	response.Measurements = make([]measurement.MeasurementResponse, len(measurements))

	for i, m := range measurements {
		response.Measurements[i] = measurement.MeasurementResponse{
			ID:        m.ID.Hex(),
			Palm:      m.Measurements.Palm,
			Length:    m.Measurements.Length,
			Unit:      m.Measurements.Unit,
			Quality:   &m.Quality,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}

	return &response, nil
}

// GetUserStats 获取用户测量统计信息
func (s *service) GetUserStats(ctx context.Context, userID string) (*models.MeasurementUserStats, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的用户ID")
	}

	filter := bson.M{"userId": userObjID}
	var stats models.MeasurementUserStats

	err = s.db.Collection(models.MeasurementUserStatsCollection).FindOne(ctx, filter).Decode(&stats)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 如果没有找到统计记录，则计算一个新的
			return s.calculateAndSaveUserStats(ctx, userObjID)
		}
		return nil, errors.NewInternalServerError("获取用户统计信息失败: " + err.Error())
	}

	return &stats, nil
}

// GetRecommendations 获取设备推荐
func (s *service) GetRecommendations(ctx context.Context, userID string) (*measurement.MeasurementRecommendationResponse, error) {
	// 获取用户统计信息
	stats, err := s.GetUserStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 如果没有足够的测量数据，返回错误
	if stats.MeasurementCount == 0 {
		return nil, errors.NewBadRequestError("没有足够的测量数据进行推荐")
	}

	// 基于手型推荐设备
	// 这里仅作为示例，实际应从设备数据库中查询匹配的设备
	var recommendations []measurement.DeviceRecommendation

	// 这里可以实现一个更复杂的推荐算法，但现在仅返回一个示例
	response := &measurement.MeasurementRecommendationResponse{
		HandSize: stats.HandSize,
		GripType: determineGripType(stats.Averages.Palm, stats.Averages.Length),
		Devices:  recommendations,
	}

	return response, nil
}

// 辅助方法

// calculateQualityScore 计算测量质量分数
func calculateQualityScore(calibrated bool) int {
	// 简单实现，实际应考虑更多因素
	if calibrated {
		return 85 // 校准过的给予较高分数
	}
	return 70 // 未校准给予较低分数
}

// updateUserStats 更新用户统计数据
func (s *service) updateUserStats(ctx context.Context, userID primitive.ObjectID) {
	s.calculateAndSaveUserStats(ctx, userID)
}

// calculateAndSaveUserStats 计算并保存用户统计数据
func (s *service) calculateAndSaveUserStats(ctx context.Context, userID primitive.ObjectID) (*models.MeasurementUserStats, error) {
	// 计算平均值
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"userId":    userID,
			"deletedAt": bson.M{"$exists": false},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":            nil,
			"avgPalm":        bson.M{"$avg": "$measurements.palm"},
			"avgLength":      bson.M{"$avg": "$measurements.length"},
			"count":          bson.M{"$sum": 1},
			"lastMeasuredAt": bson.M{"$max": "$createdAt"},
		}}},
	}

	cursor, err := s.db.Collection(models.MeasurementsCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.NewInternalServerError("计算用户统计数据失败: " + err.Error())
	}
	defer cursor.Close(ctx)

	type result struct {
		AvgPalm        float64   `bson:"avgPalm"`
		AvgLength      float64   `bson:"avgLength"`
		Count          int       `bson:"count"`
		LastMeasuredAt time.Time `bson:"lastMeasuredAt"`
	}

	var results []result
	if err = cursor.All(ctx, &results); err != nil {
		return nil, errors.NewInternalServerError("解析聚合结果失败: " + err.Error())
	}

	// 如果没有测量记录，返回默认值
	if len(results) == 0 {
		defaultStats := &models.MeasurementUserStats{
			UserID:           userID,
			Averages:         models.MeasurementData{Unit: "mm"},
			HandSize:         "unknown",
			MeasurementCount: 0,
			UpdatedAt:        time.Now(),
		}

		// 保存默认统计
		opts := options.Update().SetUpsert(true)
		filter := bson.M{"userId": userID}
		update := bson.M{"$set": defaultStats}
		_, _ = s.db.Collection(models.MeasurementUserStatsCollection).UpdateOne(ctx, filter, update, opts)

		return defaultStats, nil
	}

	// 计算手型分类
	handSize := classifyHandSize(results[0].AvgPalm, results[0].AvgLength)

	// 更新或创建统计记录
	stats := models.MeasurementUserStats{
		UserID: userID,
		Averages: models.MeasurementData{
			Palm:   results[0].AvgPalm,
			Length: results[0].AvgLength,
			Unit:   "mm",
		},
		HandSize:         handSize,
		LastMeasuredAt:   results[0].LastMeasuredAt,
		MeasurementCount: results[0].Count,
		UpdatedAt:        time.Now(),
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"userId": userID}
	update := bson.M{"$set": stats}

	// 更新用户统计数据
	_, err = s.db.Collection(models.MeasurementUserStatsCollection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, errors.NewInternalServerError("保存用户统计数据失败: " + err.Error())
	}

	return &stats, nil
}

// classifyHandSize 手型分类
func classifyHandSize(palm, length float64) string {
	// 简化的分类逻辑，实际应当更复杂
	if palm < 80 {
		return "small"
	} else if palm > 105 {
		return "large"
	}
	return "medium"
}

// determineGripType 确定握持类型
func determineGripType(palm, length float64) string {
	// 简化的握持类型判断
	ratio := length / palm

	if ratio > 1.2 {
		return "claw"
	} else if ratio < 0.9 {
		return "palm"
	}

	return "fingertip"
}
