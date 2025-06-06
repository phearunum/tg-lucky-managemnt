package repository

import (
	"api-service/lib/users/dto"
	models "api-service/lib/users/models"
	"api-service/utils"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserList11(page int, limit int, query string) ([]*models.User, int, error) {
	offset := (page - 1) * limit
	var users []*models.User
	var total int64
	log.Print("User Repo Page v%", page)
	log.Print("User Repo Page v%", limit)
	db := r.db // Assuming you have a db connection here
	err := db.Preload("Role").Model(&models.User{}).
		Where("username LIKE ?", "%"+query+"%").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// Then, retrieve the paginated list of users that match the search query
	err = db.Preload("Role").Model(&models.User{}).
		Where("username LIKE ?", "%"+query+"%").
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}
func (r *UserRepository) GetUserList(page int, limit int, query string) ([]*models.User, int, error) {
	offset := (page - 1) * limit
	var users []*models.User
	var total int64
	db := r.db // Assuming you have a db connection here
	// Start building the base query
	baseQuery := db.Model(&models.User{}).Preload("Role")
	if query != "" {
		baseQuery = baseQuery.Where("username LIKE ?", "%"+query+"%")
	}
	err := baseQuery.Count(&total).Error
	if err != nil {
		utils.ErrorLog(err.Error(), "Execute")
		return nil, 0, err
	}
	err = baseQuery.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		utils.ErrorLog(err.Error(), "Execute")
		return nil, 0, err
	}
	utils.InfoLog(nil, "Execute")
	return users, int(total), nil
}

// CreateUser creates a new user using the provided DTO
func (r *UserRepository) CreateUser(dto *dto.UserCreateDTO) (*models.User, error) {

	encryptedPassword, err := utils.EncryptPassword(dto.Password)
	if err != nil {
		log.Printf("create user error encryptedPassword  %v:", err)
		return nil, fmt.Errorf("failed to encrypt password: %v", err)
	}
	log.Printf("create user dto %v:", dto)
	user := &models.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Username:  dto.Username,
		Password:  encryptedPassword,
		RoleID:    uint(dto.RoleID),
		CompanyID: dto.CompanyID,
		Sex:       dto.Sex,
		Phone:     dto.Phone,
		Status:    dto.Status,
		//Token:     dto.Token,
		TGgroup: dto.TGgroup,
	}
	log.Printf("create user %v", user)
	if err := r.db.Create(user).Error; err != nil {
		log.Print("excute user err ", err)
		return nil, err
	}

	return user, nil
}
func (r *UserRepository) ChangeAccountPassword(dto *dto.UserChangePasswordDTO) (*models.User, error) {
	user := &models.User{}
	encryptedPassword, err := utils.EncryptPassword(dto.Password)
	if err != nil {
		log.Printf("create user error encryptedPassword  %v:", err)
		return nil, fmt.Errorf("failed to encrypt password: %v", err)
	}
	// Find the user by ID
	if err := r.db.First(user, dto.ID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	user.ID = uint(dto.ID)
	user.Password = encryptedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
func (r *UserRepository) AssignGroup(dto *dto.AssignGroupDTO) (*models.User, error) {
	user := &models.User{}
	if err := r.db.First(user, dto.ID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	user.ID = uint(dto.ID)
	user.TGgroup = dto.TGgroup
	user.UpdatedAt = time.Now()
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates an existing user using the provided DTO
// UserRepository.go
func (r *UserRepository) UpdateUser(userID uint, dto *dto.UserUpdateDTO) (*models.User, error) {
	user := &models.User{}

	// Find the user by ID
	if err := r.db.First(user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Update the user fields with the new values from the DTO
	if dto.FirstName != "" {
		user.FirstName = dto.FirstName
	}
	if dto.LastName != "" {
		user.LastName = dto.LastName
	}
	if dto.Username != "" {
		user.Username = dto.Username
	}
	if dto.RoleID != 0 {
		user.RoleID = uint(dto.RoleID)
	}
	if dto.CompanyID != 0 {
		user.CompanyID = dto.CompanyID
	}
	if dto.Sex != "" {
		user.Sex = dto.Sex
	}
	if dto.Phone != "" {
		user.Phone = dto.Phone
	}
	if dto.Status != "" {
		user.Status = dto.Status
	}
	if dto.Token != "" {
		user.Token = dto.Token
	}
	if dto.TGgroup != "" {
		user.TGgroup = dto.TGgroup
	}

	// Save the updated user to the database
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) DeleteUserByID(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	return result.Error
}

func (r *UserRepository) FindOneById(userID uint) (*models.User, error) {
	var user models.User
	baseQuery := r.db.Model(&models.User{}).Preload("Role")
	err := baseQuery.Where("id = ?", userID).First(&user).Error

	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}

	utils.LoggerRepository(user, "Execute")
	return &user, nil
}
