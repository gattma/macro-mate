package model

type Nutriments struct {
	Energy100g        float64 `json:"energy_100g"`
	EnergyUnit        string  `json:"energy_unit"`
	EnergyKcal        float64 `json:"energy_kcal"`
	EnergyKcal100g    float64 `json:"energy-kcal_100g"`
	Carbohydrates100g float64 `json:"carbohydrates_100g"`
	Sugars100g        float64 `json:"sugars_100g"`
	Fat100g           float64 `json:"fat_100g"`
	Salt100g          float64 `json:"salt_100g"`
	Proteins100g      float64 `json:"proteins_100g"`
}
type Food struct {
	Id          string     `json:"id"`
	ProductName string     `json:"product_name"`
	Nutriments  Nutriments `json:"nutriments"`
}
