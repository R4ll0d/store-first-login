package services
import "golang.org/x/crypto/bcrypt"

// HashPassword ทำการเข้ารหัสรหัสผ่านด้วย bcrypt
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// CheckPasswordHash ตรวจสอบว่ารหัสผ่านที่ป้อนเข้ามาตรงกับที่แฮชไว้หรือไม่
func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
