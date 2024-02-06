package models

type CreationSuccessResponse struct {
	DocumentID  string `json:"Product created with id"`
	ProductData Product
}

type UpdateSuccessResponse struct {
	DocumentID  string `json:"Product updated with id"`
	ProductData Product
}

type ErrorResponse struct {
	Error string `json:"Error"`
}

type SuccessResponse struct {
	Success string `json:"Success"`
}

type ListGrocerySuccessResponse struct {
	Message          string      `json:"message"`
	Page             int         `json:"page"`
	NumberOfProducts int         `json:"number of products"`
	Products         []Product   `json:"products"`
	NextPage         interface{} `json:"next page"`
}
