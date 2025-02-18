package models

type ResponseJson struct {
	StatusCode string `bson:"statusCode"`
	Status     string `bson:"status"`
	Message    string `bson:"message"`
}
type ResponseData struct {
	StatusCode string                   `bson:"statusCode"`
	Status     string                   `bson:"status"`
	Message    string                   `bson:"message"`
	Data       []map[string]interface{} `bson:"data"`
}
