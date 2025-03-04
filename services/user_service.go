package services

import (
	"encoding/json"
	"fmt"
	"os"
	"store-first-login/errs"
	"store-first-login/logs"
	"store-first-login/models"
	"store-first-login/repositories"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pquerna/otp/totp"
)

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// InsertUser implements UserService.
func (s userService) InsertUser(user models.UserRegister) error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		logs.Error("Invalid hashPassword")
		return errs.NewValidationError("Invalid hashPassword")
	}
	user.Password = hashedPassword
	user.CreateDate = time.Now().UTC().Local().Format("2006-01-02T15:04:05.999-0700")
	insertedID, err := s.userRepo.Insert(user)
	if err != nil {
		if insertedID == nil {
			logs.Error("Username Already Exits")
			return errs.NewAlreadyExits("Already Exits")
		}
		logs.Error("Unexpected error")
		return errs.NewUnexpectedError()
	}
	logs.Info(fmt.Sprintf("User %s Created", user.Username))
	return nil
}

// UpdateUser implements UserService.
func (s userService) UpdateUser(username string, user models.UserUpdate) error {
	if user.Password != "" {
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			logs.Error("Invalid hashPassword")
			return errs.NewValidationError("Invalid hashPassword")
		}
		user.Password = hashedPassword
	}
	result, err := s.userRepo.Update(username, user)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	if result.MatchedCount == 0 {
		logs.Error("User Not Found")
		return errs.NewNotFoundError("User Not Found")
	}
	if result.ModifiedCount == 0 {
		logs.Error(errs.NewUnexpectedError())
		return errs.NewUnexpectedError()
	}
	logs.Info(fmt.Sprintf("User %s Updated with details: %+v", username, user))
	return nil
}

// DeleteUser implements UserService.
func (s userService) DeleteUser(username string) error {
	err := s.userRepo.Delete(username)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	logs.Info(fmt.Sprintf("User %s Deleted", username))
	return nil
}

// GetUser implements UserService.
func (s userService) GetUser(username string) (models.UserDetail, error) {
	// result เป็น map[string]interface{}
	result, err := s.userRepo.GetOne(username)
	if err != nil {
		logs.Error(err.Error())
		return models.UserDetail{}, err // คืนค่า struct เปล่าแทน nil
	}

	// แปลง map[string]interface{} เป็น JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		logs.Error(err.Error())
		return models.UserDetail{}, fmt.Errorf("failed to marshal result: %v", err)
	}

	// แปลง JSON เป็น models.UserDetail
	var userDetails models.UserDetail
	if err := json.Unmarshal(jsonData, &userDetails); err != nil {
		logs.Error(err.Error())
		return models.UserDetail{}, fmt.Errorf("failed to unmarshal to models.UserDetail: %v", err)
	}
	logs.Info(fmt.Sprintf("User %s Geted with details: %+v", username, userDetails))
	return userDetails, nil
}

// LoginUser implements UserService.
func (s userService) LoginUser(input models.UserLogin) (string, error) {
	// ดึงข้อมูลผู้ใช้จากฐานข้อมูล
	user, err := s.userRepo.GetOne(input.Username)
	if err != nil {
		logs.Error(fmt.Errorf("user not found"))
		return "", err
	}
	// แปลง map[string]interface{} เป็น JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		logs.Error(fmt.Errorf("failed to marshal result: %v", err))
		return "", err
	}

	// แปลง JSON เป็น models.UserDetail
	var userDetails models.UserRegister
	if err := json.Unmarshal(jsonData, &userDetails); err != nil {
		logs.Error(fmt.Errorf("failed to unmarshal to models.UserDetail: %v", err))
		return "", err
	}

	// ตรวจสอบรหัสผ่าน
	if !CheckPasswordHash(input.Password, userDetails.Password) {
		logs.Error("Invalid password for user: " + input.Username)
		return "", err
	}

	// สร้าง JWT Token
	tokenString, err := s.GenerateJWT(userDetails)
	if err != nil {
		logs.Error(fmt.Errorf("failed to generate JWT token: %v", err))
		return "", err
	}
	logs.Info(fmt.Sprintf("User %s Login successfully", userDetails.Username))
	return tokenString, nil // ล็อกอินสำเร็จ
}

// ฟังก์ชันสำหรับสร้าง JWT Token
func (s userService) GenerateJWT(userDetails models.UserRegister) (string, error) {
	// สร้าง JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// สร้าง claims (ข้อมูลที่เก็บใน token)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userDetails.Username
	claims["role"] = userDetails.Role
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix() // ตั้งเวลาให้หมดอายุใน 24 ชั่วโมง

	// เซ็น JWT ด้วยคีย์ลับ
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		logs.Error("Failed to sign JWT token: " + err.Error())
		return "", fmt.Errorf("failed to create token")
	}

	return tokenString, nil
}

// ฟังก์ชันสำหรับสร้าง OTP และส่งให้ผู้ใช้
func (s userService) SendOTP(user models.UserRegister) error {
	// สร้าง OTP Secret
	otpKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "store-first-login",
		AccountName: user.Username,
	})
	if err != nil {
		logs.Error(fmt.Errorf("failed to generate OTP key: %v", err))
		return err
	}

	// สร้าง OTP code
	otpCode, err := totp.GenerateCode(otpKey.Secret(), time.Now())
	if err != nil {
		logs.Error(fmt.Errorf("failed to generate OTP code: %v", err))
		return err
	}

	// ส่ง OTP ผ่านช่องทางที่เหมาะสม (แทน fmt.Println)
	err = s.SendOTPToUser(user, otpCode)
	if err != nil {
		logs.Error(fmt.Errorf("failed to send OTP: %v", err))
		return err
	}

	return nil
}

// SendOTPToUser ส่ง OTP ไปยังหมายเลขโทรศัพท์ของผู้ใช้
func (s *userService) SendOTPToUser(user models.UserRegister, otpCode string) error {
	fmt.Printf("Successfully sent OTP '%s' to number '%s'\n", otpCode, user.PhoneNumber)
	return nil
}

// func (u userService) VerifyOTP(otp models.OTPRequest) error {
// 	// ดึงข้อมูลผู้ใช้จากฐานข้อมูล
// 	res, err := u.userRepo.FindByUsername(&user.Username)
// 	if err != nil {
// 		return fmt.Errorf("user not found: %v", err)
// 	}
// 	// ตรวจสอบว่า OTP ที่กรอกมาถูกต้องหรือไม่
// 	if !totp.Validate(user.OTP, res.OTPKey) {
// 		return fmt.Errorf("invalid OTP code")
// 	}

// 	return nil
// }
