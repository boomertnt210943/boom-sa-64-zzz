package controller

import (
	"net/http"

	"github.com/boomertnt210943/my-app/entity"
	"github.com/boomertnt210943/my-app/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginPayload login body
type LoginPayload struct {
	StudentId    string `json:"student_id"`
	Password string `json:"password"`
}

// LoginResponse token response
type LoginResponse struct {
	Token string `json:"token"`
	ID    uint   `json:"id"`
	Stdid string `json:"stdid"`
}

// POST /login
func Login(c *gin.Context) {
	var payload LoginPayload
	var studentCouncils entity.StudentCouncil

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// ค้นหา studentCouncils ด้วย รหัสนักศึกษา ที่ผู้ใช้กรอกเข้ามา
	if err := entity.DB().Raw("SELECT * FROM student_councils WHERE student_id = ?", payload.StudentId).Scan(&studentCouncils).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบรหัสผ่าน
	err := bcrypt.CompareHashAndPassword([]byte(studentCouncils.Password), []byte(payload.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid studentCouncils credentials" + payload.StudentId +" : " + payload.Password})
		return
	}

	// กำหนดค่า SecretKey, Issuer และระยะเวลาหมดอายุของ Token สามารถกำหนดเองได้
	// SecretKey ใช้สำหรับการ sign ข้อความเพื่อบอกว่าข้อความมาจากตัวเราแน่นอน
	// Issuer เป็น unique id ที่เอาไว้ระบุตัว client
	// ExpirationHours เป็นเวลาหมดอายุของ token

	jwtWrapper := service.JwtWrapper{
		SecretKey:       "SvNQpBN8y3qlVrsGAYYWoJJk56LtzFHx",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}

	signedToken, err := jwtWrapper.GenerateToken(studentCouncils.Student_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error signing token"})
		return
	}

	tokenResponse := LoginResponse{
		Token: signedToken,
		ID:    studentCouncils.ID,
		Stdid: studentCouncils.Student_id,
	}

	c.JSON(http.StatusOK, gin.H{"data": tokenResponse})
}
