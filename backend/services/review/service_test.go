package review

import (
	"project/backend/internal/errors"
	"project/backend/models"
	"project/backend/types/review"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 定义MongoDB接口和Mock实现

// 集合接口
type CollectionInterface interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResultInterface
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (CursorInterface, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
}

// 单一结果接口
type SingleResultInterface interface {
	Decode(v interface{}) error
}

// 游标接口
type CursorInterface interface {
	Next(ctx context.Context) bool
	Close(ctx context.Context) error
	Decode(v interface{}) error
	All(ctx context.Context, results interface{}) error
}

// 数据库接口
type DatabaseInterface interface {
	Collection(name string) CollectionInterface
}

// MockCollection 模拟集合实现
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResultInterface {
	args := m.Called(ctx, filter)
	return args.Get(0).(SingleResultInterface)
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (CursorInterface, error) {
	args := m.Called(ctx, filter, opts[0])
	return args.Get(0).(CursorInterface), args.Error(1)
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, replacement)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// MockSingleResult 模拟单一结果实现
type MockSingleResult struct {
	mock.Mock
}

func (m *MockSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	
	// 如果提供了结果，则将其复制到v
	if args.Get(0) != nil {
		src := args.Get(0)
		// 根据具体类型进行转换
		if review, ok := src.(*models.Review); ok {
			// 将review复制到v
			*v.(*models.Review) = *review
		}
	}
	
	return args.Error(1)
}

// MockCursor 模拟游标实现
type MockCursor struct {
	mock.Mock
}

func (m *MockCursor) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockCursor) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCursor) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockCursor) All(ctx context.Context, results interface{}) error {
	args := m.Called(ctx, results)
	
	// 如果提供了结果数组，则将其复制到results
	if args.Get(0) != nil {
		src := args.Get(0)
		// 根据具体类型进行转换
		if reviews, ok := src.([]models.Review); ok && results != nil {
			// 将reviews复制到results
			*results.(*[]models.Review) = reviews
		}
	}
	
	return args.Error(1)
}

// MockDatabase 模拟数据库实现
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Collection(name string) CollectionInterface {
	args := m.Called(name)
	return args.Get(0).(CollectionInterface)
}

// 数据库适配器
type DatabaseAdapter struct {
	mockDB *MockDatabase
}

func (a *DatabaseAdapter) Collection(name string) *mongo.Collection {
	// 这个方法永远不会被调用，因为我们在测试中用Mock拦截了这个调用
	panic("not implemented")
}

func NewDatabaseAdapter(mockDB *MockDatabase) *mongo.Database {
	// 利用unsafe转换，因为我们只需要传递一个符合*mongo.Database类型的值
	// 实际在测试中，所有的Collection调用都被我们的Mock拦截了
	return (*mongo.Database)(mock.Anything)

// ReviewServiceSuite 测试套件
type ReviewServiceSuite struct {
	suite.Suite
	mockDB       *MockDatabase
	mockColl     *MockCollection
	svc          *service
	ctx          context.Context
	
	// 测试数据
	userID       primitive.ObjectID
	reviewID     primitive.ObjectID
	itemID       primitive.ObjectID
	reviewerID   primitive.ObjectID
	testReview   *models.Review
}

// SetupTest 在每个测试前设置环境
func (s *ReviewServiceSuite) SetupTest() {
	s.mockDB = new(MockDatabase)
	s.mockColl = new(MockCollection)
	s.ctx = context.Background()
	
	// 设置数据库和集合
	s.mockDB.On("Collection", models.ReviewsCollection).Return(s.mockColl)
	
	// 由于service结构体的db字段需要*mongo.Database类型
	// 而我们的mockDB是*MockDatabase类型，所以我们需要使用unsafe来修改
	// 这里利用了Go的结构体可比较性
	reviewService := service{}
	// 通过反射来设置私有字段的值
	reviewService.db = (*mongo.Database)(nil)
	
	// 结构体引用
	s.svc = &reviewService
	
	// 替换原生Collection方法，确保服务调用Collection时使用我们的mock
	MongoCollection = func(db *mongo.Database, name string) CollectionInterface {
		return s.mockDB.Collection(name)
	}
	
	// 准备测试数据
	s.userID = primitive.NewObjectID()
	s.reviewID = primitive.NewObjectID()
	s.itemID = primitive.NewObjectID()
	s.reviewerID = primitive.NewObjectID()
	
	// 创建一个测试评测
	now := time.Now()
	s.testReview = &models.Review{
		ID:             s.reviewID,
		ExternalItemID: s.itemID,
		ItemType:       "device",
		UserID:         s.userID,
		Content:        "这是一个测试评测内容，描述一款鼠标的使用体验。内容至少要有50个字符，所以我需要多写一些。",
		Pros:           []string{"轻量化", "无线连接稳定", "电池续航好"},
		Cons:           []string{"价格偏高", "侧键手感一般"},
		Score:          4.5,
		Usage:          "日常办公和FPS游戏",
		Status:         models.ReviewStatusPending,
		Type:           models.ReviewTypeMouse,
		ContentType:    models.ReviewContentSingle,
		ViewCount:      0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// TestCreateReview 测试创建评测
func (s *ReviewServiceSuite) TestCreateReview() {
	// 测试请求数据
	req := review.CreateReviewRequest{
		ExternalItemID: s.itemID.Hex(),
		ItemType:       "device",
		Content:        "这是一个测试评测内容，描述一款鼠标的使用体验。内容至少要有50个字符，所以我需要多写一些。",
		Pros:           []string{"轻量化", "无线连接稳定", "电池续航好"},
		Cons:           []string{"价格偏高", "侧键手感一般"},
		Score:          4.5,
		Usage:          "日常办公和FPS游戏",
		Type:           string(models.ReviewTypeMouse),
		ContentType:    string(models.ReviewContentSingle),
	}
	
	// 设置Mock期望
	s.mockColl.On("InsertOne", s.ctx, mock.Anything).Return(
		&mongo.InsertOneResult{InsertedID: s.reviewID},
		nil,
	)
	
	// 执行测试
	result, err := s.svc.CreateReview(s.ctx, req, s.userID.Hex())
	
	// 验证结果
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.Content, result.Content)
	assert.Equal(s.T(), req.Score, result.Score)
	assert.Equal(s.T(), req.Pros, result.Pros)
	assert.Equal(s.T(), req.Cons, result.Cons)
	assert.Equal(s.T(), string(models.ReviewStatusPending), result.Status)
	
	// 验证Mock
	s.mockColl.AssertExpectations(s.T())
}

// TestCreateReviewWithInvalidID 测试使用无效ID创建评测
func (s *ReviewServiceSuite) TestCreateReviewWithInvalidID() {
	// 测试请求数据
	req := review.CreateReviewRequest{
		ExternalItemID: "invalid-id", // 无效的ObjectID
		ItemType:       "device",
		Content:        "这是一个测试评测内容",
		Pros:           []string{"优点1"},
		Cons:           []string{"缺点1"},
		Score:          4.5,
		Usage:          "日常使用",
		Type:           string(models.ReviewTypeMouse),
		ContentType:    string(models.ReviewContentSingle),
	}
	
	// 执行测试
	result, err := s.svc.CreateReview(s.ctx, req, s.userID.Hex())
	
	// 验证结果
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	
	// 验证错误类型
	appErr, ok := err.(*errors.AppError)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), errors.BadRequest, appErr.Code)
}

// TestGetReviewByID 测试根据ID获取评测
func (s *ReviewServiceSuite) TestGetReviewByID() {
	// 设置Mock期望
	mockResult := new(MockSingleResult)
	mockResult.On("Decode", mock.Anything).Return(s.testReview, nil)
	
	s.mockColl.On("FindOne", s.ctx, bson.M{"_id": s.reviewID}).Return(mockResult)
	
	// 执行测试
	result, err := s.svc.GetReviewByID(s.ctx, s.reviewID.Hex())
	
	// 验证结果
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), s.reviewID.Hex(), result.ID)
	assert.Equal(s.T(), s.testReview.Content, result.Content)
	
	// 验证Mock
	mockResult.AssertExpectations(s.T())
	s.mockColl.AssertExpectations(s.T())
}

// TestGetReviewByIDNotFound 测试获取不存在的评测
func (s *ReviewServiceSuite) TestGetReviewByIDNotFound() {
	// 设置Mock期望
	mockResult := new(MockSingleResult)
	mockResult.On("Decode", mock.Anything).Return(nil, mongo.ErrNoDocuments)
	
	s.mockColl.On("FindOne", s.ctx, bson.M{"_id": s.reviewID}).Return(mockResult)
	
	// 执行测试
	result, err := s.svc.GetReviewByID(s.ctx, s.reviewID.Hex())
	
	// 验证结果
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	
	// 验证错误类型
	appErr, ok := err.(*errors.AppError)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), errors.NotFound, appErr.Code)
	
	// 验证Mock
	mockResult.AssertExpectations(s.T())
	s.mockColl.AssertExpectations(s.T())
}

// TestUpdateReview 测试更新评测
func (s *ReviewServiceSuite) TestUpdateReview() {
	// 准备测试数据
	updatedContent := "这是更新后的评测内容，描述了对鼠标使用体验的更新感受。至少需要50个字符，所以我需要多写一些。"
	updatedScore := 4.0
	
	// 创建更新请求
	req := review.UpdateReviewRequest{
		Content: &updatedContent,
		Score:   &updatedScore,
	}
	
	// 设置用于FindOne的Mock
	mockFindResult := new(MockSingleResult)
	mockFindResult.On("Decode", mock.Anything).Return(s.testReview, nil)
	
	s.mockColl.On("FindOne", s.ctx, mock.Anything).Return(mockFindResult)
	
	// 设置用于ReplaceOne的Mock
	s.mockColl.On("ReplaceOne", s.ctx, mock.Anything, mock.Anything).Return(
		&mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1},
		nil,
	)
	
	// 执行测试
	result, err := s.svc.UpdateReview(s.ctx, req, s.reviewID.Hex(), s.userID.Hex())
	
	// 验证结果
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), updatedContent, result.Content)
	assert.Equal(s.T(), updatedScore, result.Score)
	
	// 验证Mock
	mockFindResult.AssertExpectations(s.T())
	s.mockColl.AssertExpectations(s.T())
}

// TestUpdateReviewNotFound 测试更新不存在的评测
func (s *ReviewServiceSuite) TestUpdateReviewNotFound() {
	// 创建更新请求
	content := "更新内容"
	req := review.UpdateReviewRequest{
		Content: &content,
	}
	
	// 设置用于FindOne的Mock
	mockFindResult := new(MockSingleResult)
	mockFindResult.On("Decode", mock.Anything).Return(nil, mongo.ErrNoDocuments)
	
	s.mockColl.On("FindOne", s.ctx, mock.Anything).Return(mockFindResult)
	
	// 执行测试
	result, err := s.svc.UpdateReview(s.ctx, req, s.reviewID.Hex(), s.userID.Hex())
	
	// 验证结果
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	
	// 验证错误类型
	appErr, ok := err.(*errors.AppError)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), errors.NotFound, appErr.Code)
	
	// 验证Mock
	mockFindResult.AssertExpectations(s.T())
	s.mockColl.AssertExpectations(s.T())
}

// TestApproveReview 测试批准评测
func (s *ReviewServiceSuite) TestApproveReview() {
	// 设置用于FindOne的Mock
	mockFindResult := new(MockSingleResult)
	mockFindResult.On("Decode", mock.Anything).Return(s.testReview, nil)
	
	s.mockColl.On("FindOne", s.ctx, mock.Anything).Return(mockFindResult)
	
	// 设置用于ReplaceOne的Mock
	s.mockColl.On("ReplaceOne", s.ctx, mock.Anything, mock.Anything).Return(
		&mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1},
		nil,
	)
	
	// 执行测试
	result, err := s.svc.ApproveReview(s.ctx, s.reviewID.Hex(), s.reviewerID.Hex(), "审核通过")
	
	// 验证结果
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), string(models.ReviewStatusApproved), result.Status)
	assert.NotNil(s.T(), result.ReviewedAt)
	assert.NotNil(s.T(), result.PublishedAt)
	
	// 验证Mock
	mockFindResult.AssertExpectations(s.T())
	s.mockColl.AssertExpectations(s.T())
}

// TestRejectReview 测试拒绝评测
func (s *ReviewServiceSuite) TestRejectReview() {
	// 拒绝理由
	rejectNotes := "内容质量不够高，请修改后重新提交"
	
	// 设置用于FindOne的Mock
	mockFindResult := new(MockSingleResult)
	mockFindResult.On("Decode", mock.Anything).Return(s.testReview, nil)
	
	s.mockColl.On("FindOne", s.ctx, mock.Anything).Return(mockFindResult)
	
	// 设置用于ReplaceOne的Mock
	s.mockColl.On("ReplaceOne", s.ctx, mock.Anything, mock.Anything).Return(
		&mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1},
		nil,
	)
	
	// 执行测试
	result, err := s.svc.RejectReview(s.ctx, s.reviewID.Hex(), s.reviewerID.Hex(), rejectNotes)
	
	// 验证结果
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), string(models.ReviewStatusRejected), result.Status)
	assert.Equal(s.T(), rejectNotes, result.ReviewerNotes)
	assert.NotNil(s.T(), result.ReviewedAt)
	
	// 验证Mock
	mockFindResult.AssertExpectations(s.T())
	s.mockColl.AssertExpectations(s.T())
}

// TestGetReviewsByUserID 测试获取用户评测列表
func (s *ReviewServiceSuite) TestGetReviewsByUserID() {
	// 准备测试数据
	reviews := []models.Review{*s.testReview}
	
	// 设置用于CountDocuments的Mock
	s.mockColl.On("CountDocuments", s.ctx, mock.Anything).Return(int64(1), nil)
	
	// 设置用于Find的Mock
	mockCursor := new(MockCursor)
	mockCursor.On("All", s.ctx, mock.Anything).Return(reviews, nil)
	mockCursor.On("Close", s.ctx).Return(nil)
	
	s.mockColl.On("Find", s.ctx, mock.Anything, mock.Anything).Return(mockCursor, nil)
	
	// 执行测试
	result, err := s.svc.GetReviewsByUserID(s.ctx, s.userID.Hex(), 1, 10)
	
	// 验证结果
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), 1, result.Total)
	assert.Equal(s.T(), 1, len(result.Reviews))
	assert.Equal(s.T(), s.testReview.Content, result.Reviews[0].Content)
	
	// 验证Mock
	mockCursor.AssertExpectations(s.T())
	s.mockColl.AssertExpectations(s.T())
}

// TestGetPendingReviews 测试获取待审核评测
func (s *ReviewServiceSuite) TestGetPendingReviews() {
	// 准备测试数据
	reviews := []models.Review{*s.testReview}
	
	// 设置用于CountDocuments的Mock
	s.mockColl.On("CountDocuments", s.ctx, mock.Anything).Return(int64(1), nil)
	
	// 设置用于Find的Mock
	mockCursor := new(MockCursor)
	mockCursor.On("All", s.ctx, mock.Anything).Return(reviews, nil)
	mockCursor.On("Close", s.ctx).Return(nil)
	
	s.mockColl.On("Find", s.ctx, mock.Anything, mock.Anything).Return(mockCursor, nil)
	
	// 执行测试
	result, err := s.svc.GetPendingReviews(s.ctx, 1, 10)
	
	// 验证结果
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), 1, result.Total)
	assert.Equal(s.T(), 1, len(result.Reviews))
	assert.Equal(s.T(), string(models.ReviewStatusPending), result.Reviews[0].Status)
	
	// 验证Mock
	mockCursor.AssertExpectations(s.T())
	s.mockColl.AssertExpectations(s.T())
}

// TestAddReviewView 测试增加评测浏览次数
func (s *ReviewServiceSuite) TestAddReviewView() {
	// 设置用于UpdateOne的Mock
	s.mockColl.On("UpdateOne", s.ctx, mock.Anything, mock.Anything).Return(
		&mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1},
		nil,
	)
	
	// 执行测试
	err := s.svc.AddReviewView(s.ctx, s.reviewID.Hex())
	
	// 验证结果
	assert.NoError(s.T(), err)
	
	// 验证Mock
	s.mockColl.AssertExpectations(s.T())
}

// 运行测试套件
func TestReviewServiceSuite(t *testing.T) {
	suite.Run(t, new(ReviewServiceSuite))
}