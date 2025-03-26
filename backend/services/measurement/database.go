package measurement

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseInterface 定义了服务层需要的数据库操作
type DatabaseInterface interface {
	Collection(name string) CollectionInterface
}

// CollectionInterface 定义了集合操作
type CollectionInterface interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResultInterface
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (CursorInterface, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (CursorInterface, error)
}

// SingleResultInterface 定义了单条结果操作
type SingleResultInterface interface {
	Decode(v interface{}) error
}

// CursorInterface 定义了游标操作
type CursorInterface interface {
	Next(ctx context.Context) bool
	Close(ctx context.Context) error
	Decode(v interface{}) error
	All(ctx context.Context, results interface{}) error
}

// 适配器将*mongo.Database转换为DatabaseInterface
type DatabaseAdapter struct {
	DB *mongo.Database
}

func (a *DatabaseAdapter) Collection(name string) CollectionInterface {
	return &CollectionAdapter{Collection: a.DB.Collection(name)}
}

// 适配器将*mongo.Collection转换为CollectionInterface
type CollectionAdapter struct {
	Collection *mongo.Collection
}

func (a *CollectionAdapter) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return a.Collection.InsertOne(ctx, document, opts...)
}

func (a *CollectionAdapter) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResultInterface {
	return &SingleResultAdapter{SingleResult: a.Collection.FindOne(ctx, filter, opts...)}
}

func (a *CollectionAdapter) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (CursorInterface, error) {
	cursor, err := a.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return &CursorAdapter{Cursor: cursor}, nil
}

func (a *CollectionAdapter) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return a.Collection.UpdateOne(ctx, filter, update, opts...)
}

func (a *CollectionAdapter) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return a.Collection.CountDocuments(ctx, filter, opts...)
}

func (a *CollectionAdapter) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (CursorInterface, error) {
	cursor, err := a.Collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	return &CursorAdapter{Cursor: cursor}, nil
}

// 适配器将*mongo.SingleResult转换为SingleResultInterface
type SingleResultAdapter struct {
	SingleResult *mongo.SingleResult
}

func (a *SingleResultAdapter) Decode(v interface{}) error {
	return a.SingleResult.Decode(v)
}

// 适配器将*mongo.Cursor转换为CursorInterface
type CursorAdapter struct {
	Cursor *mongo.Cursor
}

func (a *CursorAdapter) Next(ctx context.Context) bool {
	return a.Cursor.Next(ctx)
}

func (a *CursorAdapter) Close(ctx context.Context) error {
	return a.Cursor.Close(ctx)
}

func (a *CursorAdapter) Decode(v interface{}) error {
	return a.Cursor.Decode(v)
}

func (a *CursorAdapter) All(ctx context.Context, results interface{}) error {
	return a.Cursor.All(ctx, results)
}
