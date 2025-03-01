package models

type UserRegister struct {
	Username        string `bson:"username"`
	Password        string `bson:"password"`
	ConfirmPassword string `bson:"-"`
	Email           string `bson:"email"`
	Role            string `bson:"role"`
	CreateDate      string `bson:"createDate"`
	DeleteDate      string `bson:"deleteDate"`
}

type UserDetail struct {
	Username        string `bson:"username"`
	Email           string `bson:"email"`
	Role            string `bson:"role"`
}

type UserUpdate struct {
	Password string `bson:"password"`
	Email    string `bson:"email"`
	Role     string `bson:"role"`
}
