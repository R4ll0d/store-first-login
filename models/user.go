package models

type UserRegister struct {
	Username        string `bson:"username"`
	Password        string `bson:"password"`
	ConfirmPassword string `bson:"-"`
	PhoneNumber     string `bson:"phoneNumber"`
	Email           string `bson:"email"`
	Role            string `bson:"role"`
	Verify          bool   `bson:"verify"`
	CreateDate      string `bson:"createDate"`
	DeleteDate      string `bson:"deleteDate"`
}

type UserLogin struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type UserDetail struct {
	Username string `bson:"username"`
	Email    string `bson:"email"`
	Role     string `bson:"role"`
}

type UserUpdate struct {
	Password string `bson:"password"`
	Email    string `bson:"email"`
	Role     string `bson:"role"`
}

type OTPRequest struct {
	OTP string `bson:"otp"`
}
