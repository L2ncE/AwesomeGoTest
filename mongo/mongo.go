package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongo/mongo.go
// Mongo定义一个mongodb的数据访问对象
type Mongo struct {
	col *mongo.Collection
}

// 使用NewMongo来初始化一个mongodb的数据访问对象
func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("test"),
	}
}

// 将test_id解析为ObjID
func (m *Mongo) ResolveObjID(c context.Context, testID string) (string, error) {
	res := m.col.FindOneAndUpdate(c, bson.M{
		"test_id": testID,
	}, bson.M{
		"$set": bson.M{
			"test_id": testID,
		},
	}, options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After))

	if err := res.Err(); err != nil {
		return "", fmt.Errorf("cannot findOneAndUpdate: %v", err)
	}

	var row struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err := res.Decode(&row)
	if err != nil {
		return "", fmt.Errorf("cannot decode result: %v", err)
	}
	return row.ID.Hex(), nil
}
