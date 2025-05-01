package cart

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"project/backend/internal/errors"
	"project/backend/models"
)

// Service 购物车服务接口
type Service interface {
	GetCart(ctx context.Context, userID string) (*models.Cart, error)
	AddToCart(ctx context.Context, userID string, item models.CartItem) error
	UpdateQuantity(ctx context.Context, userID string, productID string, quantity int) error
	RemoveFromCart(ctx context.Context, userID string, productID string) error
	ClearCart(ctx context.Context, userID string) error
}

// MongoService 实现了购物车服务接口
type MongoService struct {
	collection *mongo.Collection
}

// NewService 创建新的购物车服务
func NewService(db *mongo.Database) Service {
	if db == nil {
		// 安全检查，如果数据库为nil，返回Mock服务
		return &MockService{}
	}
	return &MongoService{
		collection: db.Collection("carts"),
	}
}

// GetCart 获取用户的购物车
func (s *MongoService) GetCart(ctx context.Context, userID string) (*models.Cart, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("无效的用户ID")
	}

	filter := bson.M{"user_id": objectID}
	cart := &models.Cart{}

	err = s.collection.FindOne(ctx, filter).Decode(cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 若用户尚无购物车，创建新购物车
			cart = &models.Cart{
				UserID:    objectID,
				Items:     []models.CartItem{},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			result, err := s.collection.InsertOne(ctx, cart)
			if err != nil {
				return nil, errors.NewInternalServerError("创建购物车失败")
			}
			cart.ID = result.InsertedID.(primitive.ObjectID)
		} else {
			return nil, errors.NewInternalServerError("获取购物车失败")
		}
	}

	return cart, nil
}

// AddToCart 向购物车添加商品
func (s *MongoService) AddToCart(ctx context.Context, userID string, item models.CartItem) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.NewBadRequestError("无效的用户ID")
	}

	// 获取用户购物车
	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	// 检查商品是否已在购物车中
	itemIndex := -1
	for i, cartItem := range cart.Items {
		if cartItem.ProductID == item.ProductID {
			itemIndex = i
			break
		}
	}

	if itemIndex >= 0 {
		// 若商品已存在，更新数量
		cart.Items[itemIndex].Quantity += item.Quantity
		update := bson.M{
			"$set": bson.M{
				"items":      cart.Items,
				"updated_at": time.Now(),
			},
		}
		_, err = s.collection.UpdateOne(ctx, bson.M{"user_id": objectID}, update)
	} else {
		// 若商品不存在，添加到购物车
		update := bson.M{
			"$push": bson.M{"items": item},
			"$set":  bson.M{"updated_at": time.Now()},
		}
		_, err = s.collection.UpdateOne(ctx, bson.M{"user_id": objectID}, update)
	}

	if err != nil {
		return errors.NewInternalServerError("添加商品到购物车失败")
	}
	return nil
}

// UpdateQuantity 更新购物车商品数量
func (s *MongoService) UpdateQuantity(ctx context.Context, userID string, productID string, quantity int) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.NewBadRequestError("无效的用户ID")
	}

	productObjID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.NewBadRequestError("无效的商品ID")
	}

	if quantity <= 0 {
		// 若数量为0，从购物车中移除
		return s.RemoveFromCart(ctx, userID, productID)
	}

	update := bson.M{
		"$set": bson.M{
			"items.$[elem].quantity": quantity,
			"updated_at":             time.Now(),
		},
	}

	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem.product_id": productObjID}},
	})

	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"user_id": objectID},
		update,
		opts,
	)
	if err != nil {
		return errors.NewInternalServerError("更新商品数量失败")
	}
	return nil
}

// RemoveFromCart 从购物车移除商品
func (s *MongoService) RemoveFromCart(ctx context.Context, userID string, productID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.NewBadRequestError("无效的用户ID")
	}

	productObjID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.NewBadRequestError("无效的商品ID")
	}

	update := bson.M{
		"$pull": bson.M{
			"items": bson.M{"product_id": productObjID},
		},
		"$set": bson.M{"updated_at": time.Now()},
	}

	_, err = s.collection.UpdateOne(ctx, bson.M{"user_id": objectID}, update)
	if err != nil {
		return errors.NewInternalServerError("从购物车移除商品失败")
	}
	return nil
}

// ClearCart 清空购物车
func (s *MongoService) ClearCart(ctx context.Context, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.NewBadRequestError("无效的用户ID")
	}

	update := bson.M{
		"$set": bson.M{
			"items":      []models.CartItem{},
			"updated_at": time.Now(),
		},
	}

	_, err = s.collection.UpdateOne(ctx, bson.M{"user_id": objectID}, update)
	if err != nil {
		return errors.NewInternalServerError("清空购物车失败")
	}
	return nil
}
