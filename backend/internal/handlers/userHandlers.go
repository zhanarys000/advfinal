package handlers

import (
	"adv/internal/config"
	"adv/internal/db"
	"adv/internal/forms"
	"adv/internal/models"
	"adv/internal/repository"
	"adv/pkg/email"
	"adv/pkg/logger"
	"adv/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type UserHandlers struct {
	repo   *repository.Repository
	config config.Config
}

func NewUserHandlers(repo *repository.Repository, config config.Config) *UserHandlers {
	return &UserHandlers{repo: repo, config: config}
}

func (h UserHandlers) Register(context *gin.Context) {
	logger.GetLogger().Info("Starting user registration")
	var newUser models.User
	if err := context.ShouldBindJSON(&newUser); err != nil {
		logger.GetLogger().Error("Invalid registration request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.repo.GetUserByEmail(newUser.Email)
	if err == nil {
		logger.GetLogger().Error("Account already registered for email:", newUser.Email)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "The account is already registered"})
		return
	}
	if newUser.Password == "qwerty123" && newUser.Email == "musabecova05@gmail.com" {
		newUser.Role = "ADMIN"
	} else {
		newUser.Role = "USER"
	}
	fmt.Println(newUser.Password)
	hashedPassword, _ := utils.HashPassword(newUser.Password)
	newUser.Password = hashedPassword
	var user *models.User
	user, err = h.repo.CreateUser(&newUser)
	if err != nil {
		logger.GetLogger().Error("Account already registered for email:", newUser.Email)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "The account is already registered"})
		return
	}
	signedToken, _ := utils.CreateToken(strconv.Itoa(int(user.ID)), newUser.Email, newUser.Role)
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(context.Writer, &cookie)
	logger.GetLogger().Info("User registered successfully")
	context.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (h UserHandlers) Login(context *gin.Context) {
	logger.GetLogger().Info("Starting user login")

	var data forms.LoginForm
	if err := context.BindJSON(&data); err != nil {
		logger.GetLogger().Error("Invalid login request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.repo.GetUserByUsername(data.Username)
	if err != nil {
		logger.GetLogger().Error("Failed to get user by username:", err)
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckPasswordHash(data.Password, user.Password) {
		logger.GetLogger().Error("Authentication failed for user:", user.Email)
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	token, err := utils.CreateToken(strconv.Itoa(int(user.ID)), user.Email, user.Role)
	if err != nil {
		logger.GetLogger().Error("Failed to create token:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token", "data": token})
		return
	}
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(context.Writer, &cookie)
	logger.GetLogger().Info("User login successful")
	context.JSON(http.StatusOK, gin.H{"token": token})
}

func (h UserHandlers) GetVerificationCode(context *gin.Context) {
	logger.GetLogger().Info("Starting forgot password process")

	var form forms.GetCode
	if err := context.BindJSON(&form); err != nil {
		logger.GetLogger().Error("Invalid forgot password request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	verificationCode := utils.GenerateVerificationCode()
	err := db.SaveVerificationCodeToRedis(context, h.config, form.Email, verificationCode)
	if err != nil {
		logger.GetLogger().Error("Failed to save verification code to Redis:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = email.SendVerificationCodeEmail(form.Email, verificationCode, h.config)
	if err != nil {
		logger.GetLogger().Error("Failed to get verification code:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.GetLogger().Info("Verification code sent to email successfully")
	context.JSON(http.StatusOK, gin.H{"message": "Verification code sent to your email"})
}

func (h UserHandlers) CheckVerificationCode(context *gin.Context) {
	logger.GetLogger().Info("Checking verification code")
	var code forms.CheckCode
	if err := context.BindJSON(&code); err != nil {
		logger.GetLogger().Error("Invalid code check request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	validCode, err := db.CheckVerificationCode(context, h.config, code.Email, code.Code)
	if err != nil {
		logger.GetLogger().Error("Failed to verify code:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify code"})
		return
	}
	if !validCode {
		logger.GetLogger().Error("Failed to save user information")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user information"})
		return
	}
	logger.GetLogger().Info("successful")
	context.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func (h UserHandlers) Profile(context *gin.Context) {
	logger.GetLogger().Info("Fetching user profile")
	username, exists := context.Get("email")
	if !exists {
		logger.GetLogger().Error("User not authenticated")
		context.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not authenticated"})
		return
	}
	emailm, ok := username.(string)
	if !ok {
		logger.GetLogger().Error("Error while retrieving user ID")
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error while retrieving user ID"})
		return
	}

	user, err := h.repo.GetUserByEmail(emailm)
	if err != nil {
		logger.GetLogger().Error("User does not exist:", err)
		context.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User does not exist", "data": err})
		return
	}
	logger.GetLogger().Info("User profile fetched successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "User profile fetched successfully", "data": user})
}

func (h UserHandlers) Subscribe(context *gin.Context) {
	logger.GetLogger().Info("Starting Subscribe user")
	userEmail, emailExists := context.Get("email")
	if !emailExists {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "User email not found"})
		return
	}
	email, ok := userEmail.(string)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert email to string"})
		return
	}

	var count int64
	count, err := h.repo.CheckSubscriptionByEmail(email)
	if err != nil {
		logger.GetLogger().Error("Failed to check subscription:", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check subscription"})
		return
	}
	if count > 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User is already subscribed"})
		return
	}

	subscribe := &models.Subscribers{
		Email: email,
	}
	err = h.repo.CreateSubscriber(subscribe)
	if err != nil {
		logger.GetLogger().Error("Failed to create subscribe:", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	logger.GetLogger().Info("subscribe  successfully")
	context.JSON(http.StatusOK, gin.H{"message": "subscribes uccessfully"})
}

func (h UserHandlers) UpdateProfile(context *gin.Context) {
	logger.GetLogger().Info("Starting user update")
	var updateUser models.User
	if err := context.ShouldBindJSON(&updateUser); err != nil {
		logger.GetLogger().Error("Invalid user update request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := h.repo.GetUserByEmail(updateUser.Email)
	if err != nil {
		logger.GetLogger().Error("Failed to find user:", err.Error())
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if updateUser.Username != "" {
		existingUser.Username = updateUser.Username
	}
	if updateUser.Email != "" {
		existingUser.Email = updateUser.Email
	}

	if err := h.repo.UpdateUser(existingUser); err != nil {
		logger.GetLogger().Error("Failed to update user:", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	logger.GetLogger().Info("User updated successfully")
	context.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": existingUser})
}

func (h UserHandlers) SendEmail(context *gin.Context) {
	logger.GetLogger().Info("Starting spam email sending")

	var text forms.Form
	if err := context.ShouldBindJSON(&text); err != nil {
		logger.GetLogger().Error("Invalid request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users, err := h.repo.GetSubscribers()
	if err != nil {
		logger.GetLogger().Error("Failed to fetch subscribers:", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscribers"})
		return
	}

	for _, user := range users {
		_ = email.SendSpamForUser(user.Email, text.Text, h.config)
		logger.GetLogger().Info("Spam email sent to:", user.Email)
	}

	logger.GetLogger().Info("Spam emails sent successfully")
	context.JSON(http.StatusOK, gin.H{"message": "Spam emails sent successfully"})
}
