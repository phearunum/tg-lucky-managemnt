// repository/telegram_account_repository.go
package repository

import (
	luckyModel "api-service/lib/lucky/models"
	"api-service/lib/telegram/dto"
	models "api-service/lib/telegram/model"
	"api-service/utils"
	"errors"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TelegramAccountRepository struct {
	db *gorm.DB
}

func NewTelegramAccountRepository(db *gorm.DB) *TelegramAccountRepository {
	return &TelegramAccountRepository{db: db}
}
func (r *TelegramAccountRepository) GetTelegramAccountList(page int, limit int, query string) ([]*models.TelegramAccount, int, error) {
	offset := (page - 1) * limit
	var users []*models.TelegramAccount
	var total int64
	db := r.db
	baseQuery := db.Model(&models.TelegramAccount{}).Order("created_at DESC")
	if query != "" {
		baseQuery = baseQuery.Where("phone_number LIKE ?", "%"+query+"%")
		log.Printf("Generated SQL Query: %v", baseQuery.Statement.SQL.String())
	}
	// Count the total number of records matching the query
	err := baseQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Then, retrieve the paginated list
	err = baseQuery.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}
func (repo *TelegramAccountRepository) Create(account *models.TelegramAccount) (*models.TelegramAccount, error) {
	if err := repo.db.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (repo *TelegramAccountRepository) Update(account *models.TelegramAccount) (*models.TelegramAccount, error) {
	if err := repo.db.Save(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (repo *TelegramAccountRepository) FindByPhoneNumber(phoneNumber string) (*models.TelegramAccount, error) {
	var account models.TelegramAccount
	if err := repo.db.Where("phone_number = ?", phoneNumber).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *TelegramAccountRepository) Delete(id uint) error {
	if err := repo.db.Delete(&models.TelegramAccount{}, id).Error; err != nil {
		return err
	}
	return nil
}

/*
func (r *TelegramAccountRepository) GetByID(accountID uint) (*models.TelegramAccount, error) {
	var account models.TelegramAccount
	baseQuery := r.db.Model(&models.TelegramAccount{})
	err := baseQuery.Where("id = ?", accountID).First(&account).Error
	log.Printf("FindById SQL Query: %v", baseQuery.Statement.SQL.String())
	if err != nil {
		return nil, err
	}
	return &account, nil
}
*/

func (repo *TelegramAccountRepository) GetByID(id uint) (*models.TelegramAccount, error) {
	var account models.TelegramAccount
	if err := repo.db.First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *TelegramAccountRepository) GetTelegramAccountByID(id int) (*models.TelegramAccount, error) {
	var account models.TelegramAccount
	if err := repo.db.First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *TelegramAccountRepository) FilterRequest(chat_id int64, requestType string) (*models.UserRequest, error) {
	var request models.UserRequest
	if err := repo.db.
		Where("chat_id = ? AND request_type = ? AND status = ?", chat_id, requestType, "pending"). // Fixed "pedding" typo
		Order("id DESC").                                                                          // Sort by ID in descending order to get the latest record
		Limit(1).                                                                                  // Ensure only one record is returned
		First(&request).                                                                           // Fetch the first record of the ordered result (effectively the last one)
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or handle as needed if no record is found
		}
		return nil, err
	}
	return &request, nil
}

func (repo *TelegramAccountRepository) SaveRequest(request *models.UserRequest) (*models.UserRequest, error) {
	if err := repo.db.Save(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (repo *TelegramAccountRepository) FilterRequestOne(chat_id int64) (*models.UserRequest, error) {
	var request models.UserRequest
	if err := repo.db.
		Where("chat_id = ?  AND status = ?", chat_id, "pending"). // Fixed "pedding" typo
		Order("id DESC").                                         // Sort by ID in descending order to get the latest record
		Limit(1).                                                 // Ensure only one record is returned
		First(&request).                                          // Fetch the first record of the ordered result (effectively the last one)
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or handle as needed if no record is found
		}
		return nil, err
	}
	return &request, nil
}
func (repo *TelegramAccountRepository) FilterRequestSetting(RequestName string) (*models.RequestSetting, error) {
	var request models.RequestSetting

	if err := repo.db.
		Where("name = ?  AND status = ?", RequestName, "yes"). // Fixed "pedding" typo
		Order("id DESC").                                      // Sort by ID in descending order to get the latest record
		Limit(1).                                              // Ensure only one record is returned
		First(&request).                                       // Fetch the first record of the ordered result (effectively the last one)
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or handle as needed if no record is found
		}
		return nil, err
	}
	return &request, nil
}

func (repo *TelegramAccountRepository) FilterRequestSettingAll() (*[]models.RequestSetting, error) {
	var request []models.RequestSetting
	if err := repo.db.
		Where(" status = ?", "yes"). // Fixed "pedding" typo
		Order("order_no ASC").       // Sort by ID in descending order to get the latest record                                            // Ensure only one record is returned
		Find(&request).              // Fetch the first record of the ordered result (effectively the last one)
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or handle as needed if no record is found
		}
		return nil, err
	}
	return &request, nil
}

func (r *TelegramAccountRepository) GetTelegramSettinglist(page int, limit int, query string) ([]*models.RequestSetting, int, error) {
	offset := (page - 1) * limit
	var users []*models.RequestSetting
	var total int64
	db := r.db
	baseQuery := db.Model(&models.RequestSetting{})
	if query != "" {
		baseQuery = baseQuery.Where("name LIKE ?", "%"+query+"%")
	}
	// Count the total number of records matching the query
	err := baseQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Then, retrieve the paginated list
	err = baseQuery.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

func (repo *TelegramAccountRepository) FilterRequestOneWithType(chat_id int64, btnType string) (*models.UserRequest, error) {
	var request models.UserRequest
	if err := repo.db.
		Where("chat_id = ? AND  request_type =? AND status = ?", chat_id, btnType, "pending"). // Fixed "pedding" typo
		Order("id DESC").                                                                      // Sort by ID in descending order to get the latest record
		Limit(1).                                                                              // Ensure only one record is returned
		First(&request).                                                                       // Fetch the first record of the ordered result (effectively the last one)
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or handle as needed if no record is found
		}
		return nil, err
	}
	return &request, nil
}

func (repo *TelegramAccountRepository) FilterLocationSetting(TokenName string) (*models.LocationSetting, error) {
	var request models.LocationSetting
	if err := repo.db.
		Where("token_bot = ?", TokenName). // Fixed "pedding" typo
		Order("id DESC").                  // Sort by ID in descending order to get the latest record
		Limit(1).                          // Ensure only one record is returned
		First(&request).                   // Fetch the first record of the ordered result (effectively the last one)
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or handle as needed if no record is found
		}
		return nil, err
	}
	return &request, nil
}
func (repo *TelegramAccountRepository) GetLocationSettingAll() ([]*models.LocationSetting, error) {
	var request []*models.LocationSetting
	if err := repo.db. // Fixed "pedding" typo
				Order("id DESC"). // Sort by ID in descending order to get the latest record                         // Ensure only one record is returned
				Find(&request).   // Fetch the first record of the ordered result (effectively the last one)
				Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or handle as needed if no record is found
		}
		return nil, err
	}
	return request, nil
}

func (repo *TelegramAccountRepository) SaveClockTimeRequest(request *models.ClockTime) (*models.ClockTime, error) {
	if err := repo.db.Save(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (repo *TelegramAccountRepository) FilterRequestGetBreakOne(chat_id int64, status string) (*models.UserRequest, error) {
	var request models.UserRequest
	if err := repo.db.
		Where("chat_id = ? AND status = ? AND request_type <>'Location'", chat_id, status). // Fixed "pedding" typo
		Order("id DESC").                                                                   // Sort by ID in descending order to get the latest record
		Limit(1).                                                                           // Ensure only one record is returned
		First(&request).                                                                    // Fetch the first record of the ordered result (effectively the last one)
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or handle as needed if no record is found
		}
		return nil, err
	}
	return &request, nil
}

func (repo *TelegramAccountRepository) UpdateRequestStatus(requestId int, status string, message string) bool {
	// Attempt to update the record in the database
	if err := repo.db.
		Model(&models.UserRequest{}). // Specify the model for the table you're updating
		Where("id = ?", requestId).
		Updates(models.UserRequest{Status: status, Message: message}). // Use Updates with a struct
		Error; err != nil {

		// Check for specific error types
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false // No record found, handle as needed
		}
		// Log or handle unexpected errors
		return false
	}

	return true // Update was successful
}
func (s *TelegramAccountRepository) RemoveLast5Phones(bot_name string) error {
	// Step 1: Get the last 5 phone numbers
	phoneLists, err := s.FilterGet5Numbers(bot_name)
	if err != nil {
		return err // Handle the error as necessary
	}

	// Step 2: Collect IDs to delete
	var ids []uint
	for _, phoneList := range phoneLists {
		ids = append(ids, phoneList.ID) // Assuming ID is of type uint
	}

	// Step 3: Delete the selected records
	if err := s.BulkDelete(ids); err != nil {
		return err // Handle the error as necessary
	}

	return nil // Return nil if deletion was successful
}

func (repo *TelegramAccountRepository) FilterGet5Numbers(bot_name string) ([]models.PhoneLists, error) {
	var phoneLists []models.PhoneLists
	if err := repo.db.
		Where("bot_name = ? and status ='yes'", bot_name). // Filter by chat_id and status
		Order("id DESC").                                  // Sort by ID in descending order
		Limit(5).                                          // Limit the result to 5 records
		Find(&phoneLists).                                 // Retrieve the records
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No records found, return nil
		}
		return nil, err // Return the error
	}
	return phoneLists, nil // Return the slice of PhoneLists
}
func (repo *TelegramAccountRepository) BulkDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil // No IDs to delete, no action needed
	}

	// Use GORM to delete records by IDs
	if err := repo.db.Where("id IN ?", ids).Delete(&models.PhoneLists{}).Error; err != nil {
		return err // Return the error for further handling
	}
	return nil // Return nil if deletion was successful
}

func (repo *TelegramAccountRepository) BulkUpdatePhone(dtoFilter dto.BulkRequestDTO) bool {
	if err := repo.db.
		Model(&models.PhoneLists{}). // Specify the model for the table you're updating
		Where("id IN (?)", dtoFilter.ID).
		Updates(models.PhoneLists{Status: "used", Requester: dtoFilter.Requester, RequestDate: dtoFilter.RequesterDate}). // Use Updates with a struct
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false // No record found, handle as needed
		}
		return false
	}
	return true // Update was successful
}

func (repo *TelegramAccountRepository) SaveRegister(request *luckyModel.TelegramUsers) (*luckyModel.TelegramUsers, error) {
	if err := repo.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chat_id"}},
		UpdateAll: true,
	}).Create(&request).Error; err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	return request, nil
}
