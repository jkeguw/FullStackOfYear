// tests/integration/measurement_test.go

package integration

import (
	measurementHandler "FullStackOfYear/backend/handlers/measurement"
	"FullStackOfYear/backend/models"
	measurementTypes "FullStackOfYear/backend/types/measurement"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// 创建模拟的服务
type MockMeasurementService struct {
	mock.Mock
}

func (m *MockMeasurementService) CreateMeasurement(ctx context.Context, userID string, request measurementTypes.CreateMeasurementRequest) (*models.Measurement, error) {
	args := m.Called(ctx, userID, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Measurement), args.Error(1)
}

func (m *MockMeasurementService) GetMeasurement(ctx context.Context, userID, measurementID string) (*models.Measurement, error) {
	args := m.Called(ctx, userID, measurementID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Measurement), args.Error(1)
}

func (m *MockMeasurementService) UpdateMeasurement(ctx context.Context, userID, measurementID string, request measurementTypes.UpdateMeasurementRequest) (*models.Measurement, error) {
	args := m.Called(ctx, userID, measurementID, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Measurement), args.Error(1)
}

func (m *MockMeasurementService) DeleteMeasurement(ctx context.Context, userID, measurementID string) error {
	args := m.Called(ctx, userID, measurementID)
	return args.Error(0)
}

func (m *MockMeasurementService) ListMeasurements(ctx context.Context, userID string, request measurementTypes.MeasurementListRequest) (*measurementTypes.MeasurementListResponse, error) {
	args := m.Called(ctx, userID, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*measurementTypes.MeasurementListResponse), args.Error(1)
}

func (m *MockMeasurementService) GetUserStats(ctx context.Context, userID string) (*models.MeasurementUserStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MeasurementUserStats), args.Error(1)
}

func (m *MockMeasurementService) GetRecommendations(ctx context.Context, userID string) (*measurementTypes.MeasurementRecommendationResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*measurementTypes.MeasurementRecommendationResponse), args.Error(1)
}

func setupTestRouter(mockService *MockMeasurementService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := measurementHandler.NewHandler(mockService)

	// 测试路由
	group := r.Group("/api/v1/measurements")
	{
		group.POST("", func(c *gin.Context) {
			// 设置模拟的用户ID
			c.Set("userId", "5f8d0c5b7cb6c50c84b8a5f1")
			handler.CreateMeasurement(c)
		})

		group.GET("/:id", func(c *gin.Context) {
			c.Set("userId", "5f8d0c5b7cb6c50c84b8a5f1")
			handler.GetMeasurement(c)
		})

		group.PUT("/:id", func(c *gin.Context) {
			c.Set("userId", "5f8d0c5b7cb6c50c84b8a5f1")
			handler.UpdateMeasurement(c)
		})

		group.DELETE("/:id", func(c *gin.Context) {
			c.Set("userId", "5f8d0c5b7cb6c50c84b8a5f1")
			handler.DeleteMeasurement(c)
		})

		group.GET("", func(c *gin.Context) {
			c.Set("userId", "5f8d0c5b7cb6c50c84b8a5f1")
			handler.ListMeasurements(c)
		})

		group.GET("/stats", func(c *gin.Context) {
			c.Set("userId", "5f8d0c5b7cb6c50c84b8a5f1")
			handler.GetUserStats(c)
		})

		group.GET("/recommend", func(c *gin.Context) {
			c.Set("userId", "5f8d0c5b7cb6c50c84b8a5f1")
			handler.GetRecommendations(c)
		})
	}

	return r
}

func TestCreateMeasurementAPI(t *testing.T) {
	mockService := new(MockMeasurementService)
	router := setupTestRouter(mockService)

	// 准备测试数据
	now := time.Now()
	userID := "5f8d0c5b7cb6c50c84b8a5f1"
	userObjID, _ := primitive.ObjectIDFromHex(userID)
	measurementID := primitive.NewObjectID()

	// 创建请求
	request := measurementTypes.CreateMeasurementRequest{
		Palm:       85.5,
		Length:     70.2,
		Unit:       "mm",
		Device:     "test-device",
		Calibrated: true,
	}

	// 创建期望的响应数据
	expectedMeasurement := &models.Measurement{
		ID:     measurementID,
		UserID: userObjID,
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

	// 设置模拟服务预期
	mockService.On("CreateMeasurement", mock.Anything, userID, request).Return(expectedMeasurement, nil)

	// 序列化请求
	jsonRequest, _ := json.Marshal(request)

	// 创建HTTP请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/measurements", bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	router.ServeHTTP(w, req)

	// 检查响应
	assert.Equal(t, http.StatusCreated, w.Code)

	var response measurementTypes.MeasurementResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, measurementID.Hex(), response.ID)
	assert.Equal(t, 85.5, response.Palm)
	assert.Equal(t, 70.2, response.Length)
	assert.Equal(t, "mm", response.Unit)
	assert.Equal(t, 85, response.Quality.Score)

	// 验证模拟服务被正确调用
	mockService.AssertExpectations(t)
}

func TestGetMeasurementAPI(t *testing.T) {
	mockService := new(MockMeasurementService)
	router := setupTestRouter(mockService)

	// 准备测试数据
	now := time.Now()
	userID := "5f8d0c5b7cb6c50c84b8a5f1"
	userObjID, _ := primitive.ObjectIDFromHex(userID)
	measurementID := primitive.NewObjectID()

	// 创建期望的响应数据
	expectedMeasurement := &models.Measurement{
		ID:     measurementID,
		UserID: userObjID,
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

	// 设置模拟服务预期
	mockService.On("GetMeasurement", mock.Anything, userID, measurementID.Hex()).Return(expectedMeasurement, nil)

	// 创建HTTP请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/measurements/"+measurementID.Hex(), nil)

	// 发送请求
	router.ServeHTTP(w, req)

	// 检查响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response measurementTypes.MeasurementResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, measurementID.Hex(), response.ID)
	assert.Equal(t, 85.5, response.Palm)
	assert.Equal(t, 70.2, response.Length)

	// 验证模拟服务被正确调用
	mockService.AssertExpectations(t)
}
