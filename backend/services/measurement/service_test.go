package measurement

import (
	"FullStackOfYear/backend/models"
	"FullStackOfYear/backend/types/measurement"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

// 模拟数据库接口实现
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Collection(name string) CollectionInterface {
	args := m.Called(name)
	return args.Get(0).(CollectionInterface)
}

// 模拟集合接口实现
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResultInterface {
	args := m.Called(ctx, filter)
	return args.Get(0).(SingleResultInterface)
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (CursorInterface, error) {
	args := m.Called(ctx, filter, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(CursorInterface), args.Error(1)
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCollection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (CursorInterface, error) {
	args := m.Called(ctx, pipeline)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(CursorInterface), args.Error(1)
}

// 模拟单条结果接口实现
type MockSingleResult struct {
	mock.Mock
}

func (m *MockSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	// 检查特殊的decoder参数
	if decoder, ok := args.Get(0).(func(interface{}) error); ok && decoder != nil {
		return decoder(v)
	}
	return args.Error(1) // 这里可能是正确的，如果Called返回两个参数
}

// 模拟游标接口实现
type MockCursor struct {
	mock.Mock
}

func (m *MockCursor) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockCursor) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0) // 使用索引0
}

func (m *MockCursor) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0) // 使用索引0
}

func (m *MockCursor) All(ctx context.Context, results interface{}) error {
	args := m.Called(ctx, results)
	return args.Error(0) // 使用索引0
}

// 设置异步更新的模拟行为的辅助函数
func setupAsyncUpdateMocks(mockDb *MockDatabase, mockCollection *MockCollection) *MockCursor {
	mockCursor := new(MockCursor)

	// 为Aggregate操作设置模拟行为
	mockDb.On("Collection", models.MeasurementsCollection).Return(mockCollection).Maybe()
	mockCollection.On("Aggregate", mock.Anything, mock.Anything).Return(mockCursor, nil).Maybe()

	// 为游标操作设置模拟行为
	mockCursor.On("All", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	// 添加 Close 方法的模拟行为 - 这是关键修改
	mockCursor.On("Close", mock.Anything).Return(nil).Maybe()

	// 为UpdateOne操作设置模拟行为
	mockDb.On("Collection", models.MeasurementUserStatsCollection).Return(mockCollection).Maybe()
	mockCollection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
		&mongo.UpdateResult{UpsertedCount: 1},
		nil,
	).Maybe()

	return mockCursor
}

// 测试用例
func TestCreateMeasurement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// 创建模拟对象
		mockDb := new(MockDatabase)
		mockCollection := new(MockCollection)

		// 设置基本模拟行为
		mockDb.On("Collection", models.MeasurementsCollection).Return(mockCollection)
		mockCollection.On("InsertOne", mock.Anything, mock.Anything).Return(
			&mongo.InsertOneResult{InsertedID: primitive.NewObjectID()},
			nil,
		)

		// 设置异步更新的模拟行为
		mockCursor := setupAsyncUpdateMocks(mockDb, mockCollection)

		// 创建服务实例
		service := &service{db: mockDb}

		// 创建请求
		request := measurement.CreateMeasurementRequest{
			Palm:       85.5,
			Length:     70.2,
			Unit:       "mm",
			Device:     "test-device",
			Calibrated: true,
		}

		// 设置有效的userID
		userID := primitive.NewObjectID().Hex()

		// 调用服务
		result, err := service.CreateMeasurement(context.Background(), userID, request)

		// 验证结果
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 85.5, result.Measurements.Palm)
		assert.Equal(t, 70.2, result.Measurements.Length)
		assert.Equal(t, "mm", result.Measurements.Unit)
		assert.True(t, result.Quality.Factors.Calibration)
		assert.Equal(t, 85, result.Quality.Score) // 校准后的分数

		// 等待异步操作完成
		time.Sleep(100 * time.Millisecond)

		// 验证模拟对象被正确调用
		mockDb.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
		mockCursor.AssertExpectations(t)
	})

	t.Run("unit_conversion_cm", func(t *testing.T) {
		// 创建模拟对象
		mockDb := new(MockDatabase)
		mockCollection := new(MockCollection)

		// 设置基本模拟行为
		mockDb.On("Collection", models.MeasurementsCollection).Return(mockCollection)
		mockCollection.On("InsertOne", mock.Anything, mock.Anything).Return(
			&mongo.InsertOneResult{InsertedID: primitive.NewObjectID()},
			nil,
		)

		// 设置异步更新的模拟行为
		mockCursor := setupAsyncUpdateMocks(mockDb, mockCollection)

		// 创建服务实例
		service := &service{db: mockDb}

		// 创建请求 - 厘米单位
		request := measurement.CreateMeasurementRequest{
			Palm:       8.55,
			Length:     7.02,
			Unit:       "cm",
			Device:     "test-device",
			Calibrated: true,
		}

		userID := primitive.NewObjectID().Hex()

		// 调用服务
		result, err := service.CreateMeasurement(context.Background(), userID, request)

		// 验证结果
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.InDelta(t, 85.5, result.Measurements.Palm, 0.1)   // 8.55 cm = 85.5 mm
		assert.InDelta(t, 70.2, result.Measurements.Length, 0.1) // 7.02 cm = 70.2 mm

		// 等待异步操作完成
		time.Sleep(100 * time.Millisecond)

		// 验证模拟对象被正确调用
		mockDb.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
		mockCursor.AssertExpectations(t)
	})

	t.Run("unit_conversion_inch", func(t *testing.T) {
		// 创建模拟对象
		mockDb := new(MockDatabase)
		mockCollection := new(MockCollection)

		// 设置基本模拟行为
		mockDb.On("Collection", models.MeasurementsCollection).Return(mockCollection)
		mockCollection.On("InsertOne", mock.Anything, mock.Anything).Return(
			&mongo.InsertOneResult{InsertedID: primitive.NewObjectID()},
			nil,
		)

		// 设置异步更新的模拟行为
		mockCursor := setupAsyncUpdateMocks(mockDb, mockCollection)

		// 创建服务实例
		service := &service{db: mockDb}

		// 创建请求 - 英寸单位
		request := measurement.CreateMeasurementRequest{
			Palm:       3.0,
			Length:     2.5,
			Unit:       "inch",
			Device:     "test-device",
			Calibrated: false,
		}

		userID := primitive.NewObjectID().Hex()

		// 调用服务
		result, err := service.CreateMeasurement(context.Background(), userID, request)

		// 验证结果
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.InDelta(t, 76.2, result.Measurements.Palm, 0.1)   // 3.0 inch = 76.2 mm
		assert.InDelta(t, 63.5, result.Measurements.Length, 0.1) // 2.5 inch = 63.5 mm
		assert.Equal(t, 70, result.Quality.Score)                // 未校准的分数

		// 等待异步操作完成
		time.Sleep(100 * time.Millisecond)

		// 验证模拟对象被正确调用
		mockDb.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
		mockCursor.AssertExpectations(t)
	})

	t.Run("invalid_user_id", func(t *testing.T) {
		// 创建模拟对象
		mockDb := new(MockDatabase)

		// 创建服务实例
		service := &service{db: mockDb}

		// 创建请求
		request := measurement.CreateMeasurementRequest{
			Palm:       85.5,
			Length:     70.2,
			Unit:       "mm",
			Device:     "test-device",
			Calibrated: true,
		}

		// 设置无效的userID
		userID := "invalid-id"

		// 调用服务
		result, err := service.CreateMeasurement(context.Background(), userID, request)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "无效的用户ID")
	})

	t.Run("database_error", func(t *testing.T) {
		// 创建模拟对象
		mockDb := new(MockDatabase)
		mockCollection := new(MockCollection)

		// 设置模拟行为 - 数据库插入错误
		mockDb.On("Collection", models.MeasurementsCollection).Return(mockCollection)
		mockCollection.On("InsertOne", mock.Anything, mock.Anything).Return(
			nil,
			errors.New("database error"),
		)

		// 创建服务实例
		service := &service{db: mockDb}

		// 创建请求
		request := measurement.CreateMeasurementRequest{
			Palm:       85.5,
			Length:     70.2,
			Unit:       "mm",
			Device:     "test-device",
			Calibrated: true,
		}

		userID := primitive.NewObjectID().Hex()

		// 调用服务
		result, err := service.CreateMeasurement(context.Background(), userID, request)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "创建测量记录失败")

		// 验证模拟对象被正确调用
		mockDb.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
	})
}

func TestGetMeasurement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// 创建模拟对象
		mockDb := new(MockDatabase)
		mockCollection := new(MockCollection)
		mockSingleResult := new(MockSingleResult)

		// 准备测试数据
		userID := primitive.NewObjectID()
		measurementID := primitive.NewObjectID()
		now := time.Now()

		expectedMeasurement := models.Measurement{
			ID:     measurementID,
			UserID: userID,
			Measurements: models.MeasurementData{
				Palm:   85.5,
				Length: 70.2,
				Unit:   "mm",
			},
			Quality: models.MeasurementQuality{
				Score: 85,
				Factors: models.MeasurementFactors{
					Calibration: true,
					Stability:   1.0,
					Consistency: 1.0,
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		}

		// 设置模拟行为
		mockDb.On("Collection", models.MeasurementsCollection).Return(mockCollection)
		mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockSingleResult)

		// 设置Decode行为来填充数据
		mockSingleResult.On("Decode", mock.Anything).Return(func(v interface{}) error {
			// 将预期数据复制到传入的变量中
			measurement, ok := v.(*models.Measurement)
			if !ok {
				return errors.New("invalid type")
			}
			*measurement = expectedMeasurement
			return nil
		}, nil)

		// 创建服务实例
		service := &service{db: mockDb}

		// 调用服务
		result, err := service.GetMeasurement(context.Background(), userID.Hex(), measurementID.Hex())

		// 验证结果
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, measurementID, result.ID)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, 85.5, result.Measurements.Palm)
		assert.Equal(t, 70.2, result.Measurements.Length)

		// 验证模拟对象被正确调用
		mockDb.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
		mockSingleResult.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		// 创建模拟对象
		mockDb := new(MockDatabase)
		mockCollection := new(MockCollection)
		mockSingleResult := new(MockSingleResult)

		// 设置模拟行为 - 未找到记录
		mockDb.On("Collection", models.MeasurementsCollection).Return(mockCollection)
		mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockSingleResult)
		mockSingleResult.On("Decode", mock.Anything).Return(nil, mongo.ErrNoDocuments)

		// 创建服务实例
		service := &service{db: mockDb}

		userID := primitive.NewObjectID().Hex()
		measurementID := primitive.NewObjectID().Hex()

		// 调用服务
		result, err := service.GetMeasurement(context.Background(), userID, measurementID)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "未找到测量记录")

		// 验证模拟对象被正确调用
		mockDb.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
		mockSingleResult.AssertExpectations(t)
	})

	t.Run("invalid_id", func(t *testing.T) {
		// 创建模拟对象
		mockDb := new(MockDatabase)

		// 创建服务实例
		service := &service{db: mockDb}

		// 使用无效ID
		userID := "invalid-id"
		measurementID := primitive.NewObjectID().Hex()

		// 调用服务
		result, err := service.GetMeasurement(context.Background(), userID, measurementID)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "无效的用户ID")

		// 无需验证模拟对象，因为应该在ID验证失败时提前返回
	})
}

func TestUpdateMeasurement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// 创建模拟对象
		mockDb := new(MockDatabase)
		mockCollection := new(MockCollection)
		mockFindOneResult := new(MockSingleResult)
		mockUpdateOneResult := new(MockSingleResult)

		// 设置异步更新的模拟行为
		mockCursor := setupAsyncUpdateMocks(mockDb, mockCollection)

		// 准备测试数据
		userID := primitive.NewObjectID()
		measurementID := primitive.NewObjectID()
		now := time.Now()

		existingMeasurement := models.Measurement{
			ID:     measurementID,
			UserID: userID,
			Measurements: models.MeasurementData{
				Palm:   85.5,
				Length: 70.2,
				Unit:   "mm",
			},
			Quality: models.MeasurementQuality{
				Score: 85,
				Factors: models.MeasurementFactors{
					Calibration: true,
					Stability:   1.0,
					Consistency: 1.0,
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		}

		updatedMeasurement := existingMeasurement
		updatedMeasurement.Measurements.Palm = 90.0 // 更新的值
		updatedMeasurement.UpdatedAt = now.Add(time.Hour)

		// 设置模拟行为
		mockDb.On("Collection", models.MeasurementsCollection).Return(mockCollection)

		// 第一次FindOne - 获取现有记录
		mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockFindOneResult).Once()
		mockFindOneResult.On("Decode", mock.Anything).Return(func(v interface{}) error {
			measurement, ok := v.(*models.Measurement)
			if !ok {
				return errors.New("invalid type")
			}
			*measurement = existingMeasurement
			return nil
		}, nil)

		// UpdateOne操作
		mockCollection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(
			&mongo.UpdateResult{ModifiedCount: 1},
			nil,
		)

		// 第二次FindOne - 获取更新后的记录
		mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockUpdateOneResult).Once()
		mockUpdateOneResult.On("Decode", mock.Anything).Return(func(v interface{}) error {
			measurement, ok := v.(*models.Measurement)
			if !ok {
				return errors.New("invalid type")
			}
			*measurement = updatedMeasurement
			return nil
		}, nil)

		// 创建服务实例
		service := &service{db: mockDb}

		// 创建更新请求
		palmValue := 90.0
		request := measurement.UpdateMeasurementRequest{
			Palm: &palmValue,
		}

		// 调用服务
		result, err := service.UpdateMeasurement(context.Background(), userID.Hex(), measurementID.Hex(), request)

		// 验证结果
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 90.0, result.Measurements.Palm) // 检查更新后的值

		// 等待异步操作完成
		time.Sleep(100 * time.Millisecond)

		// 验证模拟对象被正确调用
		mockDb.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
		mockFindOneResult.AssertExpectations(t)
		mockUpdateOneResult.AssertExpectations(t)
		mockCursor.AssertExpectations(t)
	})
}

func TestListMeasurements(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// 创建模拟对象
		mockDb := new(MockDatabase)
		mockCollection := new(MockCollection)
		mockCursor := new(MockCursor)

		// 准备测试数据
		userID := primitive.NewObjectID()
		measurement1 := primitive.NewObjectID()
		measurement2 := primitive.NewObjectID()
		now := time.Now()

		measurements := []models.Measurement{
			{
				ID:     measurement1,
				UserID: userID,
				Measurements: models.MeasurementData{
					Palm:   85.5,
					Length: 70.2,
					Unit:   "mm",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:     measurement2,
				UserID: userID,
				Measurements: models.MeasurementData{
					Palm:   90.0,
					Length: 75.0,
					Unit:   "mm",
				},
				CreatedAt: now.Add(time.Hour),
				UpdatedAt: now.Add(time.Hour),
			},
		}

		// 设置模拟行为
		mockDb.On("Collection", models.MeasurementsCollection).Return(mockCollection)

		// Find操作
		mockCollection.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(mockCursor, nil)

		// CountDocuments操作
		mockCollection.On("CountDocuments", mock.Anything, mock.Anything).Return(int64(2), nil)

		// 直接操作参数的方式
		mockCursor.On("All", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			// 第二个参数是结果指针
			if ptr, ok := args.Get(1).(*[]models.Measurement); ok {
				*ptr = measurements
			}
		}).Return(nil)

		// 添加Close方法的模拟行为
		mockCursor.On("Close", mock.Anything).Return(nil)

		// 创建服务实例
		service := &service{db: mockDb}

		// 创建列表请求
		request := measurement.MeasurementListRequest{
			Page:     1,
			PageSize: 10,
		}

		// 调用服务
		result, err := service.ListMeasurements(context.Background(), userID.Hex(), request)

		// 验证结果
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.PageSize)

		// 先检查长度，避免索引越界
		assert.Len(t, result.Measurements, 2, "应该有2个测量记录")

		// 只有长度正确时才继续检查内容
		if len(result.Measurements) == 2 {
			assert.Equal(t, measurement1.Hex(), result.Measurements[0].ID)
			assert.Equal(t, measurement2.Hex(), result.Measurements[1].ID)
		}

		// 验证模拟对象被正确调用
		mockDb.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
		mockCursor.AssertExpectations(t)
	})
}

func TestDeleteMeasurement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// 创建模拟对象
		mockDb := new(MockDatabase)
		mockCollection := new(MockCollection)
		mockCursor := new(MockCursor)

		// 设置删除操作的模拟行为
		mockDb.On("Collection", models.MeasurementsCollection).Return(mockCollection)
		mockCollection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(
			&mongo.UpdateResult{ModifiedCount: 1},
			nil,
		)

		// 手动设置异步更新的模拟行为，而不是使用setupAsyncUpdateMocks
		mockDb.On("Collection", models.MeasurementsCollection).Return(mockCollection)
		mockCollection.On("Aggregate", mock.Anything, mock.Anything).Return(mockCursor, nil)
		mockCursor.On("All", mock.Anything, mock.Anything).Return(nil)
		mockCursor.On("Close", mock.Anything).Return(nil)
		mockDb.On("Collection", models.MeasurementUserStatsCollection).Return(mockCollection)
		mockCollection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
			&mongo.UpdateResult{UpsertedCount: 1},
			nil,
		)

		// 创建服务实例
		service := &service{db: mockDb}

		userID := primitive.NewObjectID().Hex()
		measurementID := primitive.NewObjectID().Hex()

		// 调用服务
		err := service.DeleteMeasurement(context.Background(), userID, measurementID)

		// 验证结果
		assert.NoError(t, err)

		// 等待异步操作完成
		time.Sleep(100 * time.Millisecond)

		// 验证模拟对象被正确调用
		mockDb.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
		mockCursor.AssertExpectations(t)
	})
}

func TestGetUserStats(t *testing.T) {
	t.Run("not_found_calculate_new", func(t *testing.T) {
		// 创建模拟对象
		mockDb := new(MockDatabase)
		mockCollection := new(MockCollection)
		mockSingleResult := new(MockSingleResult)
		mockCursor := new(MockCursor)

		// 设置模拟行为 - 没有找到用户统计，需要计算新的
		mockDb.On("Collection", models.MeasurementUserStatsCollection).Return(mockCollection).Maybe()
		mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockSingleResult).Maybe()
		mockSingleResult.On("Decode", mock.Anything).Return(nil, mongo.ErrNoDocuments).Maybe()

		// 为calculateAndSaveUserStats设置模拟行为
		mockDb.On("Collection", models.MeasurementsCollection).Return(mockCollection).Maybe()
		mockCollection.On("Aggregate", mock.Anything, mock.Anything).Return(mockCursor, nil).Maybe()

		// 为空结果设置模拟行为
		// 注意：这里不要使用返回函数，直接使用Return()
		mockCursor.On("All", mock.Anything, mock.Anything).Return(nil).Maybe()
		mockCursor.On("Close", mock.Anything).Return(nil).Maybe()

		// 为保存默认统计信息设置模拟行为
		mockDb.On("Collection", models.MeasurementUserStatsCollection).Return(mockCollection).Maybe()
		mockCollection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
			&mongo.UpdateResult{UpsertedCount: 1},
			nil,
		).Maybe()

		// 创建服务实例
		service := &service{db: mockDb}

		userID := primitive.NewObjectID().Hex()

		// 调用服务
		result, err := service.GetUserStats(context.Background(), userID)

		// 验证结果
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "unknown", result.HandSize)
		assert.Equal(t, 0, result.MeasurementCount)

		// 验证模拟对象被正确调用
		mockDb.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
		mockSingleResult.AssertExpectations(t)
		mockCursor.AssertExpectations(t)
	})
}

func TestGetRecommendations(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// 创建模拟对象用于GetUserStats
		mockDb := new(MockDatabase)
		mockCollection := new(MockCollection)
		mockSingleResult := new(MockSingleResult)

		// 准备测试数据
		userID := primitive.NewObjectID()
		now := time.Now()

		userStats := models.MeasurementUserStats{
			UserID: userID,
			Averages: models.MeasurementData{
				Palm:   87.5,
				Length: 72.0,
				Unit:   "mm",
			},
			HandSize:         "medium",
			LastMeasuredAt:   now,
			MeasurementCount: 5,
			UpdatedAt:        now,
		}

		// 设置模拟行为
		mockDb.On("Collection", models.MeasurementUserStatsCollection).Return(mockCollection)
		mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockSingleResult)
		mockSingleResult.On("Decode", mock.Anything).Return(func(v interface{}) error {
			stats, ok := v.(*models.MeasurementUserStats)
			if !ok {
				return errors.New("invalid type")
			}
			*stats = userStats
			return nil
		}, nil)

		// 创建服务实例
		service := &service{db: mockDb}

		// 调用服务
		result, err := service.GetRecommendations(context.Background(), userID.Hex())

		// 验证结果
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "medium", result.HandSize)
		// 由于ratio = 72.0/87.5 ≈ 0.82 < 0.9，握持类型应为"palm"
		assert.Equal(t, "palm", result.GripType)

		// 验证模拟对象被正确调用
		mockDb.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
		mockSingleResult.AssertExpectations(t)
	})

	t.Run("no_measurements", func(t *testing.T) {
		// 创建模拟对象用于GetUserStats
		mockDb := new(MockDatabase)
		mockCollection := new(MockCollection)
		mockSingleResult := new(MockSingleResult)

		// 准备测试数据 - 没有测量记录
		userID := primitive.NewObjectID()
		now := time.Now()

		userStats := models.MeasurementUserStats{
			UserID: userID,
			Averages: models.MeasurementData{
				Unit: "mm",
			},
			HandSize:         "unknown",
			MeasurementCount: 0,
			UpdatedAt:        now,
		}

		// 设置模拟行为
		mockDb.On("Collection", models.MeasurementUserStatsCollection).Return(mockCollection)
		mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockSingleResult)
		mockSingleResult.On("Decode", mock.Anything).Return(func(v interface{}) error {
			stats, ok := v.(*models.MeasurementUserStats)
			if !ok {
				return errors.New("invalid type")
			}
			*stats = userStats
			return nil
		}, nil)

		// 创建服务实例
		service := &service{db: mockDb}

		// 调用服务
		result, err := service.GetRecommendations(context.Background(), userID.Hex())

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "没有足够的测量数据进行推荐")

		// 验证模拟对象被正确调用
		mockDb.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
		mockSingleResult.AssertExpectations(t)
	})

	t.Run("user_stats_not_found", func(t *testing.T) {
		// 创建模拟对象
		mockDb := new(MockDatabase)
		mockCollection := new(MockCollection)
		mockSingleResult := new(MockSingleResult)
		mockCursor := new(MockCursor)

		// 设置模拟行为 - 没有找到用户统计
		mockDb.On("Collection", models.MeasurementUserStatsCollection).Return(mockCollection).Maybe()
		mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(mockSingleResult).Maybe()
		mockSingleResult.On("Decode", mock.Anything).Return(nil, mongo.ErrNoDocuments).Maybe()

		// 为calculateAndSaveUserStats设置模拟行为
		mockDb.On("Collection", models.MeasurementsCollection).Return(mockCollection).Maybe()
		mockCollection.On("Aggregate", mock.Anything, mock.Anything).Return(mockCursor, nil).Maybe()

		// 不使用返回函数
		mockCursor.On("All", mock.Anything, mock.Anything).Return(nil).Maybe()
		mockCursor.On("Close", mock.Anything).Return(nil).Maybe()

		// 为保存统计信息设置模拟行为
		mockDb.On("Collection", models.MeasurementUserStatsCollection).Return(mockCollection).Maybe()
		mockCollection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
			&mongo.UpdateResult{UpsertedCount: 1},
			nil,
		).Maybe()

		// 创建服务实例
		service := &service{db: mockDb}

		userID := primitive.NewObjectID().Hex()

		// 调用服务
		result, err := service.GetRecommendations(context.Background(), userID)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "没有足够的测量数据进行推荐")

		// 验证模拟对象被正确调用
		mockDb.AssertExpectations(t)
		mockCollection.AssertExpectations(t)
		mockSingleResult.AssertExpectations(t)
		mockCursor.AssertExpectations(t)
	})
}
