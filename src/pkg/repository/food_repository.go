package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gattma/macro-mate/src/cmd/utils"
	"github.com/gattma/macro-mate/src/pkg/entity"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FoodRepository interface {
	Search(query string, ctx context.Context) ([]*entity.Food, error)
	SearchKeyword(keyword string, ctx context.Context) ([]*entity.Food, error)
	Add(appDoc entity.Food, ctx context.Context) (primitive.ObjectID, error)
}

type foodRepository struct {
	client *mongo.Client
	config *utils.Configuration
}

func NewFoodRepository(config *utils.Configuration, client *mongo.Client) FoodRepository {
	return &foodRepository{config: config, client: client}
}

func (app *foodRepository) Search(query string, ctx context.Context) ([]*entity.Food, error) {
	findOptions := options.Find()
	findOptions.SetLimit(int64(50))

	collection := app.client.Database(app.config.Database.DbName).Collection(app.config.Database.Collection)
	filter := bson.M{
		"product_name": bson.M{
			"$regex":   query,
			"$options": "i",
		},
	}
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var Foods []*entity.Food
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		var elem entity.Food
		if err := cursor.Decode(&elem); err != nil {
			logrus.Fatal(err)
			return nil, err
		}

		Foods = append(Foods, &elem)
	}

	cursor.Close(ctx)

	logrus.Infof("Foods Count:", len(Foods))
	return Foods, nil
}

func (app *foodRepository) SearchKeyword(keyword string, ctx context.Context) ([]*entity.Food, error) {
	findOptions := options.Find()
	findOptions.SetLimit(int64(50))

	collection := app.client.Database(app.config.Database.DbName).Collection(app.config.Database.Collection)
	filter := bson.M{
		"_keywords": bson.M{
			"$regex":   keyword,
			"$options": "i",
		},
	}
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var Foods []*entity.Food
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		var elem entity.Food
		if err := cursor.Decode(&elem); err != nil {
			logrus.Fatal(err)
			return nil, err
		}

		Foods = append(Foods, &elem)
	}

	cursor.Close(ctx)

	logrus.Infof("Foods Count:", len(Foods))
	return Foods, nil
}

func (app *foodRepository) Add(appDoc entity.Food, ctx context.Context) (primitive.ObjectID, error) {
	collection := app.client.Database(app.config.Database.DbName).Collection(app.config.Database.Collection)

	insertResult, err := collection.InsertOne(ctx, appDoc)
	if err != mongo.ErrNilCursor {
		return primitive.NilObjectID, err
	}

	if oidResult, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		return oidResult, nil
	} else {
		return primitive.NilObjectID, err
	}
}
