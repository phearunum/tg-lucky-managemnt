package telegram

import (
	"api-service/utils"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRequestRepository struct {
	db *gorm.DB
}

// NewUserRequestRepository creates a new instance of UserRequestRepository
func NewUserRequestRepository(db *gorm.DB) *UserRequestRepository {
	return &UserRequestRepository{db: db}
}

// GetUserRequestList retrieves a paginated list of user requests
func (r *UserRequestRepository) GetUserRequestList(page, limit int, query string, createdAtStart, createdAtEnd *time.Time) ([]*UserRequest, int, error) {
	var userRequests []*UserRequest
	var total int64
	offset := (page - 1) * limit

	// Build base query
	baseQuery := r.db.Model(&UserRequest{}).Order("created_at DESC")
	if query != "" {
		baseQuery = baseQuery.Where("(account_name LIKE ? OR request_type LIKE ?)", "%"+query+"%", "%"+query+"%")
	}
	// Filter by created_at if provided
	if createdAtStart != nil {
		baseQuery = baseQuery.Where("created_at >= ?", *createdAtStart)
	}
	if createdAtEnd != nil {
		baseQuery = baseQuery.Where("created_at <= ?", *createdAtEnd)
	}

	// Use transactions for better performance when counting and fetching
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := baseQuery.Offset(offset).Limit(limit).Find(&userRequests).Error; err != nil {
		return nil, 0, err
	}

	return userRequests, int(total), nil
}

// CreateUserRequest adds a new user request
func (r *UserRequestRepository) CreateUserRequest(createDTO CreateUserRequestDTO) (*UserRequest, error) {
	utils.InfoLog(createDTO, "createDTO Repo")
	userRequest := ToUserRequestModelFromCreateDTO(createDTO)

	// Insert record (pass pointer)
	if err := r.db.Create(&userRequest).Error; err != nil {
		utils.LoggerRepository(err, "CreateUserRequest: Insert failed")
		return nil, err
	}
	return &userRequest, nil
}

// UpdateUserRequest modifies an existing user request
func (r *UserRequestRepository) UpdateUserRequest(updateDTO UpdateUserRequestDTO) (*UserRequest, error) {
	existingUserRequest, err := r.GetUserRequestByID(updateDTO.ID, updateDTO.ChatID)
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	updatedUserRequest := ToUserRequestModelFromUpdateDTO(updateDTO)
	if err := r.db.Model(existingUserRequest).Updates(updatedUserRequest).Error; err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}

	return existingUserRequest, nil
}

// DeleteUserRequest removes a user request by ID
func (r *UserRequestRepository) DeleteUserRequest(id uint, chatID string) (bool, error) {
	// Check if the user request exists
	if _, err := r.GetUserRequestByID(id, chatID); err != nil {
		utils.LoggerRepository(err, "Execute")
		return false, err
	}
	// Delete the user request from the database
	if err := r.db.Where("id = ? AND chat_id = ?", id, chatID).Delete(&UserRequest{}).Error; err != nil {
		utils.LoggerRepository(err, "Execute")
		return false, err
	}
	return true, nil
}

// GetUserRequestByID retrieves a user request by ID and chat ID
func (r *UserRequestRepository) GetUserRequestByID(id uint, chatID string) (*UserRequest, error) {
	var userRequest UserRequest
	err := r.db.Where("id = ? AND chat_id = ?", id, chatID).First(&userRequest).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	return &userRequest, nil
}

func (r *UserRequestRepository) GetTelegramSettingList(page, limit int, query string) ([]*RequestSetting, int, error) {
	var settingRequests []*RequestSetting
	var total int64
	offset := (page - 1) * limit
	// Build base query
	baseQuery := r.db.Model(&RequestSetting{}).Order("created_at DESC")
	if query != "" {
		baseQuery = baseQuery.Where("(name LIKE ? OR button_type LIKE ?)", "%"+query+"%", "%"+query+"%")
	}

	// Use transactions for better performance when counting and fetching
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := baseQuery.Offset(offset).Limit(limit).Find(&settingRequests).Error; err != nil {
		return nil, 0, err
	}

	return settingRequests, int(total), nil
}

func (r *UserRequestRepository) SaveRequestSetting(createDTO CreateRequestSettingDTO) (*RequestSetting, error) {
	userRequest := ToRequestSettingModelFromCreateDTO(createDTO)
	// Insert record (pass pointer)
	if err := r.db.Model(&RequestSetting{}).Save(&userRequest).Error; err != nil {
		utils.LoggerRepository(err, "Create Request: Insert failed")
		return nil, err
	}
	return &userRequest, nil
}

func (r *UserRequestRepository) UpdateRequestSetting(updateDTO UpdateRequestSettingDTO) (*RequestSetting, error) {
	userRequest := ToRequestSettingModelFromUpdateDTO(updateDTO)

	// Check if record exists
	var existing RequestSetting
	if err := r.db.First(&existing, updateDTO.ID).Error; err != nil {
		utils.LoggerRepository(err, "Request: Record not found")
		return nil, err
	}

	// Update record
	if err := r.db.Model(&existing).Updates(userRequest).Error; err != nil {
		utils.LoggerRepository(err, "Request: Update failed")
		return nil, err
	}

	return &existing, nil
}

func (repo *UserRequestRepository) FilterRequestSettingAll() (*[]RequestSetting, error) {
	var request []RequestSetting
	if err := repo.db.
		Where(" status = ? AND bot_name=?", "yes"). // Fixed "pedding" typo
		Order("order_no ASC").                      // Sort by ID in descending order to get the latest record                                            // Ensure only one record is returned
		Find(&request).                             // Fetch the first record of the ordered result (effectively the last one)
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or handle as needed if no record is found
		}
		return nil, err
	}
	return &request, nil
}

func (repo *UserRequestRepository) FilterRequestSetting(RequestName string) (*RequestSetting, error) {
	var request RequestSetting
	if err := repo.db.
		Where("name = ?  ", RequestName, "yes"). // Fixed "pedding" typo
		Order("id DESC").                        // Sort by ID in descending order to get the latest record
		Limit(1).                                // Ensure only one record is returned
		First(&request).                         // Fetch the first record of the ordered result (effectively the last one)
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or handle as needed if no record is found
		}
		return nil, err
	}
	return &request, nil
}

func (r *UserRequestRepository) GetLocationSettingList(page, limit int, query string) ([]*LocationSetting, int, error) {
	var settingRequests []*LocationSetting
	var total int64
	offset := (page - 1) * limit
	// Build base query
	baseQuery := r.db.Model(&LocationSetting{}).Order("created_at DESC")
	if query != "" {
		baseQuery = baseQuery.Where("(name LIKE ?)", "%"+query+"%")
	}

	// Use transactions for better performance when counting and fetching
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := baseQuery.Offset(offset).Limit(limit).Find(&settingRequests).Error; err != nil {
		return nil, 0, err
	}

	return settingRequests, int(total), nil
}

func (repo *UserRequestRepository) FilterLocationSetting(TokenName string) (*LocationSetting, error) {
	var request LocationSetting
	if err := repo.db.
		Where("token_bot = ? AND status='yes'", TokenName). // Fixed "pedding" typo
		Order("id DESC").                                   // Sort by ID in descending order to get the latest record
		Limit(1).                                           // Ensure only one record is returned
		First(&request).                                    // Fetch the first record of the ordered result (effectively the last one)
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or handle as needed if no record is found
		}
		return nil, err
	}
	return &request, nil
}

func (r *UserRequestRepository) SaveLocationSetting(createDTO CreateLocationSettingDTO) (*LocationSetting, error) {
	userRequest := ToLocationSettingModelFromCreateDTO(createDTO)
	// Insert record (pass pointer)
	if err := r.db.Model(&LocationSetting{}).Save(&userRequest).Error; err != nil {
		utils.LoggerRepository(err, "Create Request: Insert failed")
		return nil, err
	}
	return &userRequest, nil
}

func (r *UserRequestRepository) UpdateLocationSetting(updateDTO UpdateLocationSettingDTO) (*LocationSetting, error) {
	userRequest := ToLocationSettingModelFromUpdateDTO(updateDTO)

	// Check if record exists
	var existing LocationSetting
	if err := r.db.First(&existing, updateDTO.ID).Error; err != nil {
		utils.LoggerRepository(err, "Request: Record not found")
		return nil, err
	}
	// Update record
	if err := r.db.Model(&existing).Updates(userRequest).Error; err != nil {
		utils.LoggerRepository(err, "Request: Update failed")
		return nil, err
	}

	return &existing, nil
}

func (repo *UserRequestRepository) FilterBotLocationSettingAll() (*[]LocationSetting, error) {
	var request []LocationSetting
	if err := repo.db.
		Order("order_no ASC"). // Sort by ID in descending order to get the latest record                                            // Ensure only one record is returned
		Find(&request).        // Fetch the first record of the ordered result (effectively the last one)
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or handle as needed if no record is found
		}
		return nil, err
	}
	return &request, nil
}

func (r *UserRequestRepository) GetClockList(page, limit int, query string) ([]*ClockTime, int, error) {
	var settingRequests []*ClockTime
	var total int64
	offset := (page - 1) * limit
	// Build base query
	baseQuery := r.db.Model(&ClockTime{}).Order("created_at DESC")
	if query != "" {
		baseQuery = baseQuery.Where("account_name LIKE ? OR chat_id LIKE ? ", "%"+query+"%", "%"+query+"%")
	}

	// Use transactions for better performance when counting and fetching
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := baseQuery.Offset(offset).Limit(limit).Find(&settingRequests).Error; err != nil {
		return nil, 0, err
	}

	return settingRequests, int(total), nil
}

func (r *UserRequestRepository) SaveClockIn(createDTO CreateClockTimeDTO) (*ClockTime, error) {
	utils.WarnLog(createDTO, "SaveClockIn")
	userRequest := ToClockTimeModelFromCreateDTO(createDTO)
	// Insert record (pass pointer)
	if err := r.db.Model(&ClockTime{}).Save(&userRequest).Error; err != nil {
		utils.LoggerRepository(err, "Create Request: Insert failed")
		return nil, err
	}
	return &userRequest, nil
}

func (r *UserRequestRepository) GetPhoneList(page, limit int, query string, status string) ([]*PhoneLists, int, error) {
	var data []*PhoneLists
	var total int64
	offset := (page - 1) * limit
	// Build base query
	baseQuery := r.db.Model(&PhoneLists{}).Where("status = ?", status).Order("created_at DESC")
	if query != "" {
		baseQuery = baseQuery.Where("(phone LIKE ? OR bot_name LIKE ?)", "%"+query+"%", "%"+query+"%")
	}

	// Use transactions for better performance when counting and fetching
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := baseQuery.Offset(offset).Limit(limit).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, int(total), nil
}

func (r *UserRequestRepository) SavePhones_(createDTOs []CreatePhoneListDTO) ([]PhoneLists, error) {
	var phoneLists []PhoneLists

	// Convert CreatePhoneListDTO to PhoneLists models
	for _, createDTO := range createDTOs {
		utils.WarnLog(createDTO, "SavePhones")
		userRequest := ToPhoneListModelFromCreateDTO(createDTO)
		phoneLists = append(phoneLists, userRequest)
	}

	// Insert records in bulk
	if err := r.db.Create(&phoneLists).Error; err != nil {
		utils.LoggerRepository(err, "Create Request: Bulk insert failed")
		return nil, err
	}

	return phoneLists, nil
}

func (r *UserRequestRepository) SavePhones(createDTOs []CreatePhoneListDTO) ([]PhoneLists, error) {
	var phoneLists []PhoneLists

	// Convert CreatePhoneListDTO to PhoneLists models
	for _, createDTO := range createDTOs {
		utils.WarnLog(createDTO, "SavePhones")
		userRequest := ToPhoneListModelFromCreateDTO(createDTO)
		phoneLists = append(phoneLists, userRequest)
	}
	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "phone"}},
		UpdateAll: true,
	}).Create(&phoneLists).Error; err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}

	return phoneLists, nil
}

func (r *UserRequestRepository) UpdatePhone(updateDTO UpdatePhoneListDTO) (*PhoneLists, error) {
	userRequest := ToPhoneListModelFromUpdateDTO(updateDTO)
	var existing PhoneLists
	if err := r.db.First(&existing, updateDTO.ID).Error; err != nil {
		utils.LoggerRepository(err, "Request: Record not found")
		return nil, err
	}
	if err := r.db.Model(&existing).Updates(userRequest).Error; err != nil {
		utils.LoggerRepository(err, "Request: Update failed")
		return nil, err
	}
	return &existing, nil
}

func (r *UserRequestRepository) BulkPhoneDelete(ids []uint) error {
	// Check if there are any IDs to delete
	if len(ids) == 0 {
		return nil // No IDs to delete, consider this a no-op
	}

	// Use GORM to delete records by IDs
	if err := r.db.Where("id IN ?", ids).Delete(&PhoneLists{}).Error; err != nil {
		// Log the error if necessary (optional)
		return err // Return the error for further handling
	}
	return nil // Return nil if deletion was successful
}
