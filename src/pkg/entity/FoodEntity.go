package entity

type Nutriments struct {
	Energy100g        float64 `bson:"energy_100g" json:"energy_100g"`
	EnergyUnit        string  `bson:"energy_unit" json:"energy_unit"`
	EnergyKcal        float64 `bson:"energy_kcal" json:"energy_kcal"`
	EnergyKcal100g    float64 `bson:"energy-kcal_100g" json:"energy-kcal_100g"`
	Carbohydrates100g float64 `bson:"carbohydrates_100g" json:"carbohydrates_100g"`
	Sugars100g        float64 `bson:"sugars_100g" json:"sugars_100g"`
	Fat100g           float64 `bson:"fat_100g" json:"fat_100g"`
	Salt100g          float64 `bson:"salt_100g" json:"salt_100g"`
	Proteins100g      float64 `bson:"proteins_100g" json:"proteins_100g"`
}

type Food struct {
	Id          string     `bson:"id" json:"id"`
	ProductName string     `bson:"product_name" json:"product_name"`
	Nutriments  Nutriments `bson:"nutriments" json:"nutriments"`
}
