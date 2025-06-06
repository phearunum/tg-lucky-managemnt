package auth

import (
	redis "api-service/config"
	users "api-service/lib/users/models"
	utils "api-service/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type AuthController struct {
	authService *AuthService
}

func NewAuthController(authService *AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	log.Print("Login Action")
	var credentials LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := c.authService.Login(&credentials)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	encryptedPassword, err := utils.EncryptPassword(credentials.Password)
	if err != nil {
		log.Printf("create user error encryptedPassword  %v:", err)

	}
	utils.InfoLog(encryptedPassword, "Token")
	if err := utils.ComparePassword(user.Password, credentials.Password); err != nil {

		utils.ErrorLog(err, "ComparePassword")
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := c.authService.GenerateToken(user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	response := struct {
		Token string      `json:"token"`
		User  *users.User `json:"user"`
	}{
		Token: token,
		User:  user,
	}
	log.Print("Start Setting Redis")
	expiration := time.Minute * 1 // 1 nimutes
	key := fmt.Sprintf("user:%d", user.ID)
	user.Token = token
	log.Printf("response %v ", response)
	log.Print("Redis Key ", key)
	log.Print("Redis expiration", expiration)
	if err := redis.SetWithExpiration(key, response, expiration); err != nil {
		utils.ErrorLog(err, "Failed Store User Infor in Redis")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to store user information in Redis")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("id")
	if id == "" {
		utils.RespondWithJSON(w, http.StatusOK, "logout")
		//http.Error(w, "Role ID not found in headers", http.StatusBadRequest)
		return
	}

	roleId, err := strconv.Atoi(id)
	if err != nil || roleId <= 0 {
		//utils.RespondWithError(w, http.StatusInternalServerError, "Failed to logout")
		//return
	}
	key := fmt.Sprintf("user:%d", roleId)
	if err := redis.RemoveRedisKey(key); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to store user information in Redis")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "logout")
}
