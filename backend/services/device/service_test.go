package device

//
//import (
//	"project/backend/internal/errors"
//	"project/backend/models"
//	"project/backend/types/device"
//	"context"
//	"testing"
//	"time"
//
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"github.com/stretchr/testify/suite"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//)
//
//// 定义MongoDB接口和Mock实现
//
//// 集合接口
//type CollectionInterface interface {
//	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
//	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResultInterface
//	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (CursorInterface, error)
//	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
//	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
//	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
//	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
//}
//
//// 单一结果接口
//type SingleResultInterface interface {
//	Decode(v interface{}) error
//}
//
//// 游标接口
//type CursorInterface interface {
//	Next(ctx context.Context) bool
//	Close(ctx context.Context) error
//	Decode(v interface{}) error
//	All(ctx context.Context, results interface{}) error
//}
//
//// 数据库接口
//type DatabaseInterface interface {
//	Collection(name string) CollectionInterface
//}
//
//// MockCollection 模拟集合实现
//type MockCollection struct {
//	mock.Mock
//}
//
//func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
//	args := m.Called(ctx, document)
//	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
//}
//
//func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResultInterface {
//	args := m.Called(ctx, filter)
//	return args.Get(0).(SingleResultInterface)
//}
//
//func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (CursorInterface, error) {
//	args := m.Called(ctx, filter, opts[0])
//	return args.Get(0).(CursorInterface), args.Error(1)
//}
//
//func (m *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
//	args := m.Called(ctx, filter, update)
//	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
//}
//
//func (m *MockCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
//	args := m.Called(ctx, filter, replacement)
//	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
//}
//
//func (m *MockCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
//	args := m.Called(ctx, filter)
//	return args.Get(0).(int64), args.Error(1)
//}
//
//func (m *MockCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
//	args := m.Called(ctx, filter)
//	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
//}
//
//// MockSingleResult 模拟单一结果实现
//type MockSingleResult struct {
//	mock.Mock
//}
//
//func (m *MockSingleResult) Decode(v interface{}) error {
//	args := m.Called(v)
//
//	// 如果提供了结果，则将其复制到v
//	if args.Get(0) != nil {
//		src := args.Get(0)
//		// 根据具体类型进行转换
//		switch src := src.(type) {
//		case *models.HardwareDevice:
//			*v.(*models.HardwareDevice) = *src
//		case *models.MouseDevice:
//			*v.(*models.MouseDevice) = *src
//		case *models.UserDevice:
//			*v.(*models.UserDevice) = *src
//		}
//	}
//
//	return args.Error(1)
//}
//
//// MockCursor 模拟游标实现
//type MockCursor struct {
//	mock.Mock
//}
//
//func (m *MockCursor) Next(ctx context.Context) bool {
//	args := m.Called(ctx)
//	return args.Bool(0)
//}
//
//func (m *MockCursor) Close(ctx context.Context) error {
//	args := m.Called(ctx)
//	return args.Error(0)
//}
//
//func (m *MockCursor) Decode(v interface{}) error {
//	args := m.Called(v)
//	return args.Error(0)
//}
//
//func (m *MockCursor) All(ctx context.Context, results interface{}) error {
//	args := m.Called(ctx, results)
//
//	// 如果提供了结果数组，则将其复制到results
//	if args.Get(0) != nil {
//		src := args.Get(0)
//		// 根据具体类型进行转换
//		switch src := src.(type) {
//		case []models.HardwareDevice:
//			*results.(*[]models.HardwareDevice) = src
//		case []models.UserDevice:
//			*results.(*[]models.UserDevice) = src
//		}
//	}
//
//	return args.Error(1)
//}
//
//// MockDatabase 模拟数据库实现
//type MockDatabase struct {
//	mock.Mock
//}
//
//func (m *MockDatabase) Collection(name string) CollectionInterface {
//	args := m.Called(name)
//	return args.Get(0).(CollectionInterface)
//}
//
//// 用于拦截mongo.Database的Collection调用的函数变量
//var MongoCollection = func(db *mongo.Database, name string) CollectionInterface {
//	// 默认实现，应该被在测试设置中覆盖
//	panic("MongoCollection not mocked")
//}
//
//// DeviceServiceSuite 测试套件
//type DeviceServiceSuite struct {
//	suite.Suite
//	mockDB           *MockDatabase
//	mockDeviceColl   *MockCollection
//	mockUserDeviceColl *MockCollection
//	svc              *service
//	ctx              context.Context
//
//	// 测试数据
//	deviceID        primitive.ObjectID
//	userID          primitive.ObjectID
//	userDeviceID    primitive.ObjectID
//	testMouseDevice *models.MouseDevice
//	testUserDevice  *models.UserDevice
//}
//
//// SetupTest 在每个测试前设置环境
//func (s *DeviceServiceSuite) SetupTest() {
//	s.mockDB = new(MockDatabase)
//	s.mockDeviceColl = new(MockCollection)
//	s.mockUserDeviceColl = new(MockCollection)
//	s.ctx = context.Background()
//
//	// 设置数据库和集合
//	s.mockDB.On("Collection", models.DevicesCollection).Return(s.mockDeviceColl)
//	s.mockDB.On("Collection", models.UserDevicesCollection).Return(s.mockUserDeviceColl)
//
//	// 由于service结构体的db字段需要*mongo.Database类型
//	// 而我们的mockDB是*MockDatabase类型，所以我们需要使用unsafe来修改
//	// 这里利用了Go的结构体可比较性
//	deviceService := service{}
//	// 通过反射来设置私有字段的值
//	deviceService.db = (*mongo.Database)(nil)
//
//	// 结构体引用
//	s.svc = &deviceService
//
//	// 替换原生Collection方法，确保服务调用Collection时使用我们的mock
//	MongoCollection = func(db *mongo.Database, name string) CollectionInterface {
//		return s.mockDB.Collection(name)
//	}
//
//	// 准备测试数据
//	s.deviceID = primitive.NewObjectID()
//	s.userID = primitive.NewObjectID()
//	s.userDeviceID = primitive.NewObjectID()
//
//	// 创建一个测试鼠标设备
//	now := time.Now()
//	s.testMouseDevice = &models.MouseDevice{
//		HardwareDevice: models.HardwareDevice{
//			ID:          s.deviceID,
//			Name:        "Test Mouse",
//			Brand:       "TestBrand",
//			Type:        models.DeviceTypeMouse,
//			ImageURL:    "https://example.com/image.jpg",
//			Description: "This is a test mouse",
//			CreatedAt:   now,
//			UpdatedAt:   now,
//		},
//		Dimensions: models.MouseDimensions{
//			Length: 120,
//			Width:  65,
//			Height: 40,
//			Weight: 70,
//		},
//		Shape: models.MouseShape{
//			Type:              "ergonomic",
//			HumpPlacement:     "center",
//			FrontFlare:        "medium",
//			SideCurvature:     "curved",
//			HandCompatibility: "medium",
//		},
//		Technical: models.MouseTechnical{
//			Connectivity: []string{"wireless", "bluetooth"},
//			Sensor:       "PixArt 3370",
//			MaxDPI:       20000,
//			PollingRate:  1000,
//			SideButtons:  2,
//		},
//		Recommended: models.MouseRecommended{
//			GameTypes:    []string{"FPS", "MOBA"},
//			GripStyles:   []string{"claw", "palm"},
//			HandSizes:    []string{"medium", "large"},
//			DailyUse:     true,
//			Professional: true,
//		},
//	}
//
//	// 创建一个测试用户设备配置
//	s.testUserDevice = &models.UserDevice{
//		ID:          s.userDeviceID,
//		UserID:      s.userID,
//		Name:        "My Gaming Setup",
//		Description: "Setup for FPS games",
//		Devices: []models.UserDeviceSettings{
//			{
//				DeviceID:   s.deviceID,
//				DeviceType: models.DeviceTypeMouse,
//				Settings: map[string]any{
//					"dpi":         1600,
//					"pollingRate": 1000,
//				},
//			},
//		},
//		IsPublic:  true,
//		CreatedAt: now,
//		UpdatedAt: now,
//	}
//}
//
//// TestCreateMouseDevice 测试创建鼠标设备
//func (s *DeviceServiceSuite) TestCreateMouseDevice() {
//	// 测试请求数据
//	req := device.CreateMouseRequest{
//		Name:        "Test Mouse",
//		Brand:       "TestBrand",
//		ImageURL:    "https://example.com/image.jpg",
//		Description: "This is a test mouse",
//		Dimensions: models.MouseDimensions{
//			Length: 120,
//			Width:  65,
//			Height: 40,
//			Weight: 70,
//		},
//		Shape: models.MouseShape{
//			Type:              "ergonomic",
//			HumpPlacement:     "center",
//			FrontFlare:        "medium",
//			SideCurvature:     "curved",
//			HandCompatibility: "medium",
//		},
//		Technical: models.MouseTechnical{
//			Connectivity: []string{"wireless", "bluetooth"},
//			Sensor:       "PixArt 3370",
//			MaxDPI:       20000,
//			PollingRate:  1000,
//			SideButtons:  2,
//		},
//		Recommended: models.MouseRecommended{
//			GameTypes:    []string{"FPS", "MOBA"},
//			GripStyles:   []string{"claw", "palm"},
//			HandSizes:    []string{"medium", "large"},
//			DailyUse:     true,
//			Professional: true,
//		},
//	}
//
//	// 设置Mock期望
//	s.mockDeviceColl.On("InsertOne", s.ctx, mock.Anything).Return(
//		&mongo.InsertOneResult{InsertedID: s.deviceID},
//		nil,
//	)
//
//	// 执行测试
//	result, err := s.svc.CreateMouseDevice(s.ctx, req)
//
//	// 验证结果
//	assert.NoError(s.T(), err)
//	assert.NotNil(s.T(), result)
//	assert.Equal(s.T(), req.Name, result.Name)
//	assert.Equal(s.T(), req.Brand, result.Brand)
//	assert.Equal(s.T(), models.DeviceTypeMouse, result.Type)
//	assert.Equal(s.T(), req.Dimensions, result.Dimensions)
//	assert.Equal(s.T(), req.Technical.Sensor, result.Technical.Sensor)
//
//	// 验证Mock
//	s.mockDeviceColl.AssertExpectations(s.T())
//}
//
//// TestGetDeviceByID 测试根据ID获取设备
//func (s *DeviceServiceSuite) TestGetDeviceByID() {
//	// 设置Mock期望
//	mockResult := new(MockSingleResult)
//	mockResult.On("Decode", mock.Anything).Return(&s.testMouseDevice.HardwareDevice, nil)
//
//	s.mockDeviceColl.On("FindOne", s.ctx, bson.M{"_id": s.deviceID}).Return(mockResult)
//
//	// 执行测试
//	result, err := s.svc.GetDeviceByID(s.ctx, s.deviceID.Hex())
//
//	// 验证结果
//	assert.NoError(s.T(), err)
//	assert.NotNil(s.T(), result)
//	assert.Equal(s.T(), s.deviceID, result.ID)
//	assert.Equal(s.T(), s.testMouseDevice.Name, result.Name)
//	assert.Equal(s.T(), s.testMouseDevice.Brand, result.Brand)
//
//	// 验证Mock
//	mockResult.AssertExpectations(s.T())
//	s.mockDeviceColl.AssertExpectations(s.T())
//}
//
//// TestGetDeviceByIDNotFound 测试获取不存在的设备
//func (s *DeviceServiceSuite) TestGetDeviceByIDNotFound() {
//	// 设置Mock期望
//	mockResult := new(MockSingleResult)
//	mockResult.On("Decode", mock.Anything).Return(nil, mongo.ErrNoDocuments)
//
//	s.mockDeviceColl.On("FindOne", s.ctx, bson.M{"_id": s.deviceID}).Return(mockResult)
//
//	// 执行测试
//	result, err := s.svc.GetDeviceByID(s.ctx, s.deviceID.Hex())
//
//	// 验证结果
//	assert.Error(s.T(), err)
//	assert.Nil(s.T(), result)
//
//	// 验证错误类型
//	appErr, ok := err.(*errors.AppError)
//	assert.True(s.T(), ok)
//	assert.Equal(s.T(), errors.NotFound, appErr.Code)
//
//	// 验证Mock
//	mockResult.AssertExpectations(s.T())
//	s.mockDeviceColl.AssertExpectations(s.T())
//}
//
//// TestGetDeviceList 测试获取设备列表
//func (s *DeviceServiceSuite) TestGetDeviceList() {
//	// 准备测试数据
//	devices := []models.HardwareDevice{s.testMouseDevice.HardwareDevice}
//
//	// 设置用于CountDocuments的Mock
//	s.mockDeviceColl.On("CountDocuments", s.ctx, mock.Anything).Return(int64(1), nil)
//
//	// 设置用于Find的Mock
//	mockCursor := new(MockCursor)
//	mockCursor.On("All", s.ctx, mock.Anything).Return(devices, nil)
//	mockCursor.On("Close", s.ctx).Return(nil)
//
//	s.mockDeviceColl.On("Find", s.ctx, mock.Anything, mock.Anything).Return(mockCursor, nil)
//
//	// 执行测试
//	result, err := s.svc.GetDeviceList(s.ctx, string(models.DeviceTypeMouse), 1, 10)
//
//	// 验证结果
//	assert.NoError(s.T(), err)
//	assert.NotNil(s.T(), result)
//	assert.Equal(s.T(), 1, result.Total)
//	assert.Equal(s.T(), 1, len(result.Devices))
//	assert.Equal(s.T(), s.deviceID.Hex(), result.Devices[0].ID)
//	assert.Equal(s.T(), s.testMouseDevice.Name, result.Devices[0].Name)
//
//	// 验证Mock
//	mockCursor.AssertExpectations(s.T())
//	s.mockDeviceColl.AssertExpectations(s.T())
//}
//
//// TestCreateUserDevice 测试创建用户设备配置
//func (s *DeviceServiceSuite) TestCreateUserDevice() {
//	// 测试请求数据
//	req := device.CreateUserDeviceRequest{
//		Name:        "My Gaming Setup",
//		Description: "Setup for FPS games",
//		Devices: []device.UserDeviceSettingsRequest{
//			{
//				DeviceID:   s.deviceID.Hex(),
//				DeviceType: string(models.DeviceTypeMouse),
//				Settings: map[string]any{
//					"dpi":         1600,
//					"pollingRate": 1000,
//				},
//			},
//		},
//		IsPublic: true,
//	}
//
//	// 设置Mock期望
//	s.mockUserDeviceColl.On("InsertOne", s.ctx, mock.Anything).Return(
//		&mongo.InsertOneResult{InsertedID: s.userDeviceID},
//		nil,
//	)
//
//	// 执行测试
//	result, err := s.svc.CreateUserDevice(s.ctx, s.userID.Hex(), req)
//
//	// 验证结果
//	assert.NoError(s.T(), err)
//	assert.NotNil(s.T(), result)
//	assert.Equal(s.T(), req.Name, result.Name)
//	assert.Equal(s.T(), req.Description, result.Description)
//	assert.Equal(s.T(), req.IsPublic, result.IsPublic)
//	assert.Len(s.T(), result.Devices, 1)
//	assert.Equal(s.T(), s.deviceID, result.Devices[0].DeviceID)
//
//	// 验证Mock
//	s.mockUserDeviceColl.AssertExpectations(s.T())
//}
//
//// TestUpdateUserDevice 测试更新用户设备配置
//func (s *DeviceServiceSuite) TestUpdateUserDevice() {
//	// 准备测试数据
//	updatedName := "Updated Gaming Setup"
//	updatedDesc := "Updated Setup for FPS and MOBA games"
//
//	// 创建更新请求
//	req := device.UpdateUserDeviceRequest{
//		Name:        updatedName,
//		Description: updatedDesc,
//		IsPublic:    true,
//		Devices: []device.UserDeviceSettingsRequest{
//			{
//				DeviceID:   s.deviceID.Hex(),
//				DeviceType: string(models.DeviceTypeMouse),
//				Settings: map[string]any{
//					"dpi":         1800,
//					"pollingRate": 1000,
//				},
//			},
//		},
//	}
//
//	// 设置用于FindOne的Mock
//	mockFindResult := new(MockSingleResult)
//	mockFindResult.On("Decode", mock.Anything).Return(s.testUserDevice, nil)
//
//	s.mockUserDeviceColl.On("FindOne", s.ctx, mock.Anything).Return(mockFindResult)
//
//	// 设置用于ReplaceOne的Mock
//	s.mockUserDeviceColl.On("ReplaceOne", s.ctx, mock.Anything, mock.Anything).Return(
//		&mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1},
//		nil,
//	)
//
//	// 执行测试
//	result, err := s.svc.UpdateUserDevice(s.ctx, s.userID.Hex(), s.userDeviceID.Hex(), req)
//
//	// 验证结果
//	assert.NoError(s.T(), err)
//	assert.NotNil(s.T(), result)
//	assert.Equal(s.T(), updatedName, result.Name)
//	assert.Equal(s.T(), updatedDesc, result.Description)
//	assert.Len(s.T(), result.Devices, 1)
//
//	// 验证Mock
//	mockFindResult.AssertExpectations(s.T())
//	s.mockUserDeviceColl.AssertExpectations(s.T())
//}
//
//// TestGetUserDeviceByID 测试获取用户设备配置
//func (s *DeviceServiceSuite) TestGetUserDeviceByID() {
//	// 设置Mock期望
//	mockResult := new(MockSingleResult)
//	mockResult.On("Decode", mock.Anything).Return(s.testUserDevice, nil)
//
//	s.mockUserDeviceColl.On("FindOne", s.ctx, mock.Anything).Return(mockResult)
//
//	// 执行测试
//	result, err := s.svc.GetUserDeviceByID(s.ctx, s.userID.Hex(), s.userDeviceID.Hex())
//
//	// 验证结果
//	assert.NoError(s.T(), err)
//	assert.NotNil(s.T(), result)
//	assert.Equal(s.T(), s.userDeviceID, result.ID)
//	assert.Equal(s.T(), s.userID, result.UserID)
//	assert.Equal(s.T(), s.testUserDevice.Name, result.Name)
//
//	// 验证Mock
//	mockResult.AssertExpectations(s.T())
//	s.mockUserDeviceColl.AssertExpectations(s.T())
//}
//
//// TestGetUserDevices 测试获取用户设备列表
//func (s *DeviceServiceSuite) TestGetUserDevices() {
//	// 准备测试数据
//	userDevices := []models.UserDevice{*s.testUserDevice}
//
//	// 设置用于CountDocuments的Mock
//	s.mockUserDeviceColl.On("CountDocuments", s.ctx, mock.Anything).Return(int64(1), nil)
//
//	// 设置用于Find的Mock
//	mockCursor := new(MockCursor)
//	mockCursor.On("All", s.ctx, mock.Anything).Return(userDevices, nil)
//	mockCursor.On("Close", s.ctx).Return(nil)
//
//	s.mockUserDeviceColl.On("Find", s.ctx, mock.Anything, mock.Anything).Return(mockCursor, nil)
//
//	// 执行测试
//	result, err := s.svc.GetUserDevices(s.ctx, s.userID.Hex(), 1, 10)
//
//	// 验证结果
//	assert.NoError(s.T(), err)
//	assert.NotNil(s.T(), result)
//	assert.Equal(s.T(), 1, result.Total)
//	assert.Equal(s.T(), 1, len(result.UserDevices))
//	assert.Equal(s.T(), s.userDeviceID.Hex(), result.UserDevices[0].ID)
//	assert.Equal(s.T(), s.userID.Hex(), result.UserDevices[0].UserID)
//	assert.Equal(s.T(), s.testUserDevice.Name, result.UserDevices[0].Name)
//
//	// 验证Mock
//	mockCursor.AssertExpectations(s.T())
//	s.mockUserDeviceColl.AssertExpectations(s.T())
//}
//
//// TestDeleteUserDevice 测试删除用户设备配置
//func (s *DeviceServiceSuite) TestDeleteUserDevice() {
//	// 设置Mock期望
//	s.mockUserDeviceColl.On("DeleteOne", s.ctx, mock.Anything).Return(
//		&mongo.DeleteResult{DeletedCount: 1},
//		nil,
//	)
//
//	// 执行测试
//	err := s.svc.DeleteUserDevice(s.ctx, s.userID.Hex(), s.userDeviceID.Hex())
//
//	// 验证结果
//	assert.NoError(s.T(), err)
//
//	// 验证Mock
//	s.mockUserDeviceColl.AssertExpectations(s.T())
//}
//
//// TestDeleteUserDeviceNotFound 测试删除不存在的用户设备配置
//func (s *DeviceServiceSuite) TestDeleteUserDeviceNotFound() {
//	// 设置Mock期望
//	s.mockUserDeviceColl.On("DeleteOne", s.ctx, mock.Anything).Return(
//		&mongo.DeleteResult{DeletedCount: 0},
//		nil,
//	)
//
//	// 执行测试
//	err := s.svc.DeleteUserDevice(s.ctx, s.userID.Hex(), s.userDeviceID.Hex())
//
//	// 验证结果
//	assert.Error(s.T(), err)
//
//	// 验证错误类型
//	appErr, ok := err.(*errors.AppError)
//	assert.True(s.T(), ok)
//	assert.Equal(s.T(), errors.NotFound, appErr.Code)
//
//	// 验证Mock
//	s.mockUserDeviceColl.AssertExpectations(s.T())
//}
//
//// TestGetPublicUserDevices 测试获取公开的用户设备列表
//func (s *DeviceServiceSuite) TestGetPublicUserDevices() {
//	// 准备测试数据
//	userDevices := []models.UserDevice{*s.testUserDevice}
//
//	// 设置用于CountDocuments的Mock
//	s.mockUserDeviceColl.On("CountDocuments", s.ctx, bson.M{"isPublic": true}).Return(int64(1), nil)
//
//	// 设置用于Find的Mock
//	mockCursor := new(MockCursor)
//	mockCursor.On("All", s.ctx, mock.Anything).Return(userDevices, nil)
//	mockCursor.On("Close", s.ctx).Return(nil)
//
//	s.mockUserDeviceColl.On("Find", s.ctx, bson.M{"isPublic": true}, mock.Anything).Return(mockCursor, nil)
//
//	// 执行测试
//	result, err := s.svc.GetPublicUserDevices(s.ctx, 1, 10)
//
//	// 验证结果
//	assert.NoError(s.T(), err)
//	assert.NotNil(s.T(), result)
//	assert.Equal(s.T(), 1, result.Total)
//	assert.Equal(s.T(), 1, len(result.UserDevices))
//	assert.Equal(s.T(), s.userDeviceID.Hex(), result.UserDevices[0].ID)
//	assert.Equal(s.T(), s.userID.Hex(), result.UserDevices[0].UserID)
//	assert.Equal(s.T(), s.testUserDevice.Name, result.UserDevices[0].Name)
//	assert.True(s.T(), result.UserDevices[0].IsPublic)
//
//	// 验证Mock
//	mockCursor.AssertExpectations(s.T())
//	s.mockUserDeviceColl.AssertExpectations(s.T())
//}
//
//// 运行测试套件
//func TestDeviceServiceSuite(t *testing.T) {
//	suite.Run(t, new(DeviceServiceSuite))
//}
