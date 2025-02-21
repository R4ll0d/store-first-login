package models

type UserRegister struct {
	Username   string `bson:"username"`
	Password   string `bson:"password"`
	Email      string `bson:"email"`
	Role       string `bson:"role"`
	CreateDate string `bson:"createDate"`
	DeleteDate string `bson:"deleteDate"`
}

type UserUpdate struct {
	Password string `bson:"password"`
	Email    string `bson:"email"`
	Role     string `bson:"role"`
}
