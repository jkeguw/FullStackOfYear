package device

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	appErrors "project/backend/internal/errors"
	"project/backend/models"
	"project/backend/services/similarity"
)

// 比较鼠标
func (s *ServiceImpl) CompareMice(ctx context.Context, mouseIDs []string) (*similarity.ComparisonResult, error) {
	if len(mouseIDs) < 1 {
		return nil, appErrors.NewBadRequestError("至少需要一个鼠标ID进行比较")
	}

	// 将字符串ID转换为ObjectID
	objectIDs := make([]primitive.ObjectID, 0, len(mouseIDs))
	for _, id := range mouseIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, appErrors.NewBadRequestError("无效的鼠标ID: " + id)
		}
		objectIDs = append(objectIDs, objID)
	}

	// 查询所有鼠标
	filter := bson.M{
		"_id": bson.M{"$in": objectIDs},
		"type": models.DeviceTypeMouse,
	}

	cursor, err := s.db.Collection(models.DevicesCollection).Find(ctx, filter)
	if err != nil {
		return nil, appErrors.NewInternalServerError("查询鼠标失败: " + err.Error())
	}
	defer cursor.Close(ctx)

	// 解析鼠标数据
	var mice []models.MouseDevice
	if err = cursor.All(ctx, &mice); err != nil {
		return nil, appErrors.NewInternalServerError("解析鼠标数据失败: " + err.Error())
	}

	if len(mice) == 0 {
		return nil, appErrors.NewNotFoundError("未找到指定的鼠标")
	}

	if len(mice) < len(mouseIDs) {
		return nil, appErrors.NewNotFoundError("部分鼠标ID无效")
	}

	// 创建相似度服务
	similarityService := similarity.NewSimilarityService()

	// 计算比较结果
	result, err := similarityService.CompareMice(mice)
	if err != nil {
		return nil, appErrors.NewInternalServerError("计算鼠标比较结果失败: " + err.Error())
	}

	return result, nil
}

// 查找相似鼠标
func (s *ServiceImpl) FindSimilarMice(ctx context.Context, mouseID string, limit int) ([]models.MouseDevice, error) {
	if mouseID == "" {
		return nil, appErrors.NewBadRequestError("鼠标ID不能为空")
	}

	if limit <= 0 {
		limit = 5 // 默认限制
	}

	// 将鼠标ID转换为ObjectID
	targetID, err := primitive.ObjectIDFromHex(mouseID)
	if err != nil {
		return nil, appErrors.NewBadRequestError("无效的鼠标ID")
	}

	// 先获取目标鼠标
	var targetMouse models.MouseDevice
	err = s.db.Collection(models.DevicesCollection).FindOne(ctx, bson.M{
		"_id": targetID,
		"type": models.DeviceTypeMouse,
	}).Decode(&targetMouse)

	if err != nil {
		return nil, appErrors.NewNotFoundError("未找到目标鼠标")
	}

	// 获取所有其他鼠标
	// 限制返回的鼠标数量，因为我们需要计算相似度
	maxMice := 100 // 最多处理100个鼠标
	limitInt64 := int64(maxMice)

	cursor, err := s.db.Collection(models.DevicesCollection).Find(ctx, 
		bson.M{
			"_id": bson.M{"$ne": targetID},
			"type": models.DeviceTypeMouse,
		},
		&options.FindOptions{
			Limit: &limitInt64,
		})
	
	if err != nil {
		return nil, appErrors.NewInternalServerError("查询鼠标列表失败: " + err.Error())
	}
	defer cursor.Close(ctx)

	// 解析鼠标数据
	var allMice []models.MouseDevice
	if err = cursor.All(ctx, &allMice); err != nil {
		return nil, appErrors.NewInternalServerError("解析鼠标数据失败: " + err.Error())
	}

	// 添加目标鼠标
	allMice = append(allMice, targetMouse)

	// 创建相似度服务
	similarityService := similarity.NewSimilarityService()

	// 寻找相似鼠标
	similarMice, err := similarityService.FindSimilarMice(targetID, allMice, limit)
	if err != nil {
		return nil, errors.Wrap(err, "寻找相似鼠标失败")
	}

	return similarMice, nil
}