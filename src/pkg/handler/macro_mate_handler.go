package handler

import (
	"context"
	"errors"
	"github.com/gattma/macro-mate/src/pkg/entity"
	"net/http"
	"strconv"
	"time"

	"github.com/gattma/macro-mate/src/cmd/utils"
	"github.com/gattma/macro-mate/src/pkg/model"
	"github.com/gattma/macro-mate/src/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type FoodHandler interface {
	Healthcheck(*gin.Context)
	Search(*gin.Context)
	SearchKeywords(*gin.Context)
	Add(*gin.Context)
}

type foodHandler struct {
	client         *mongo.Client
	foodRepository repository.FoodRepository
	config         utils.Configuration
}

func NewMacroMateHandler(client *mongo.Client, repo repository.FoodRepository, config utils.Configuration) FoodHandler {
	return &foodHandler{client: client, foodRepository: repo, config: config}
}

func (app *foodHandler) Search(c *gin.Context) {
	ctx, ctxErr := context.WithTimeout(c.Request.Context(), time.Duration(app.config.App.Timeout)*time.Second)
	defer ctxErr()

	var foodsModel []*model.Food

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}
	logrus.Infof("Keyword %s", query)

	result, err := app.foodRepository.Search(query, ctx)
	if !errors.Is(err, mongo.ErrNilCursor) {
		utils.BadRequestError("MacroMate_Handler_Search", err, map[string]interface{}{"Data": query})
	}
	logrus.Infof("Len %d", len(result))

	//convert to entity to model
	for _, item := range result {
		foodsModel = append(foodsModel, mapFoodEntityToModel(item))
	}

	c.IndentedJSON(http.StatusOK, utils.Response("MacroMate_Handler_List", map[string]interface{}{"Data": foodsModel}))
}

func (app *foodHandler) SearchKeywords(c *gin.Context) {
	ctx, ctxErr := context.WithTimeout(c.Request.Context(), time.Duration(app.config.App.Timeout)*time.Second)
	defer ctxErr()

	var foodsModel []*model.Food

	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}
	logrus.Infof("Keyword %s", keyword)

	result, err := app.foodRepository.SearchKeyword(keyword, ctx)
	if !errors.Is(err, mongo.ErrNilCursor) {
		utils.BadRequestError("MacroMate_Handler_Search", err, map[string]interface{}{"Data": keyword})
	}
	logrus.Infof("Len %d", len(result))

	//convert to entity to model
	for _, item := range result {
		foodsModel = append(foodsModel, mapFoodEntityToModel(item))
	}

	c.IndentedJSON(http.StatusOK, utils.Response("MacroMate_Handler_List", map[string]interface{}{"Data": foodsModel}))
}

func (app *foodHandler) Add(c *gin.Context) {

	ctx, ctxErr := context.WithTimeout(c.Request.Context(), time.Duration(app.config.App.Timeout)*time.Second)
	defer ctxErr()

	var appEntity entity.Food
	err := c.BindJSON(&appEntity)
	if err != nil {
		utils.BadRequestError("MacroMate_Handler_Add", err, map[string]interface{}{"Data": appEntity})
		return
	}
	//jsonData, _ := io.ReadAll(c.Request.Body)

	//appEntity = &entity.Food{
	//	ProductName: c.PostForm("product_name"),
	//	Nutriments: entity.Nutriments{
	//		Energy100g:        toFloat64(c.PostForm("energy_100g")),
	//		Carbohydrates100g: toFloat64(c.PostForm("carbohydrates_100g")),
	//		Sugars100g:        toFloat64(c.PostForm("sugars_100g")),
	//		Fat100g:           toFloat64(c.PostForm("fat_100g")),
	//		Salt100g:          toFloat64(c.PostForm("salt_100g")),
	//		Proteins100g:      toFloat64(c.PostForm("proteins_100g")),
	//	},
	//}

	oId, err := app.foodRepository.Add(appEntity, ctx)
	if err != nil {
		utils.BadRequestError("MacroMate_Handler_Add", err, map[string]interface{}{"Data": appEntity})
	}

	c.IndentedJSON(http.StatusCreated, utils.Response("MacroMate_Handler_Add", map[string]interface{}{"OId": oId}))
}

func toFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func (app *foodHandler) Healthcheck(c *gin.Context) {

	ctx, ctxErr := context.WithTimeout(c.Request.Context(), time.Duration(app.config.App.Timeout)*time.Second)
	defer ctxErr()

	if ctxErr != nil {
		logrus.Error("somethig wrong!!!", ctxErr)
	}

	if err := app.client.Ping(ctx, nil); err != nil {
		utils.InternalServerError("Status unhealth", err, map[string]interface{}{"Data": "Please check the Client", "Time": time.Local})
	}

	c.IndentedJSON(http.StatusOK, utils.Response("Pong", map[string]interface{}{"Data": "The MongoDB client is working successfully", "Date": time.Local}))
}

func mapFoodEntityToModel(foodEntity *entity.Food) *model.Food {
	return &model.Food{
		Id:          foodEntity.Id,
		ProductName: foodEntity.ProductName,
		Nutriments: model.Nutriments{
			Energy100g:        foodEntity.Nutriments.Energy100g,
			EnergyUnit:        foodEntity.Nutriments.EnergyUnit,
			EnergyKcal:        foodEntity.Nutriments.EnergyKcal,
			EnergyKcal100g:    foodEntity.Nutriments.EnergyKcal100g,
			Carbohydrates100g: foodEntity.Nutriments.Carbohydrates100g,
			Sugars100g:        foodEntity.Nutriments.Sugars100g,
			Fat100g:           foodEntity.Nutriments.Fat100g,
			Salt100g:          foodEntity.Nutriments.Salt100g,
			Proteins100g:      foodEntity.Nutriments.Proteins100g,
		},
	}
}
