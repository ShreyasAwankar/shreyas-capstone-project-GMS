package models

type ProductRequestBody struct {
	ProductName         string `form:"ProductName" binding:"required"`
	Category            string `form:"Category" binding:"required"`
	Price               int64  `form:"Price" binding:"required,gte=1"`
	Weight              string `form:"Weight" binding:"required"`
	Brand               string `form:"Brand" binding:"required"`
	Vegetarian          bool   `form:"Vegeterian"`
	PackageInformation  string `form:"PackageInformation" binding:"required"`
	ItemPackageQuantity int    `form:"ItemPackageQuantity" binding:"required,gte=1"`
	CountryofOrigin     string `form:"CountryofOrigin" binding:"required"`
	Manufacturer        string `form:"Manufacturer" binding:"required"`
}

type Product struct {
	ProductID           string `form:"ProductID"`
	ProductName         string `form:"ProductName" binding:"required"`
	Category            string `form:"Category" binding:"required"`
	Price               int64  `form:"Price" binding:"required,gte=1"`
	Weight              string `form:"Weight" binding:"required"`
	ImageURL            string `json:"ImageURL"`
	ThumbnailURL        string `json:"ThumbnailURL"`
	Brand               string `form:"Brand" binding:"required"`
	Vegetarian          bool   `form:"Vegeterian"`
	PackageInformation  string `form:"PackageInformation" binding:"required"`
	ItemPackageQuantity int    `form:"ItemPackageQuantity" binding:"required,gte=1"`
	CountryofOrigin     string `form:"CountryofOrigin" binding:"required"`
	Manufacturer        string `form:"Manufacturer" binding:"required"`
}
