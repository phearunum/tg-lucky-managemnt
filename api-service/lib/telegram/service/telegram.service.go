package service

import (
	dto "api-service/lib/telegram/dto"
	models "api-service/lib/telegram/model"
	repository "api-service/lib/telegram/repository"
	"api-service/utils"
	"encoding/json"
	"log"
	"strconv"

	"github.com/streadway/amqp"
)

type TelegramAccountService struct {
	repo *repository.TelegramAccountRepository
}

func NewTelegramAccountService(repo *repository.TelegramAccountRepository) *TelegramAccountService {
	return &TelegramAccountService{repo: repo}
}
func (us *TelegramAccountService) TelegramAccountListRequest(msgBody []byte, replyTo string, channel *amqp.Channel) {
	var requestData struct {
		Action string `json:"action"`
		Page   string `json:"page"`
		Limit  string `json:"limit"`
		Query  string `json:"query"`
	}
	err := json.Unmarshal(msgBody, &requestData)
	if err != nil {
		utils.SendErrorResponse("Failed to unmarshal JSON", replyTo, channel)
		return
	}
	pageInt, err := strconv.Atoi(requestData.Page)
	if err != nil {
		log.Printf("Failed to convert page to integer: %v", err)
		pageInt = 1 // Default to page 1 if conversion fails
	}

	limitInt, err := strconv.Atoi(requestData.Limit)
	if err != nil {
		log.Printf("Failed to convert limit to integer: %v", err)
		limitInt = 10 // Default to limit 10 if conversion fails
	}
	PermissionsList, total, err := us.repo.GetTelegramAccountList(pageInt, limitInt, requestData.Query)
	if err != nil {
		utils.SendErrorResponse("Failed to get all Permissions", replyTo, channel)
		return
	}
	utils.SendPaginationResponse(PermissionsList, total, pageInt, limitInt, replyTo, channel)
}

// CreateTelegramAccount handles creating a new Telegram account
func (tas *TelegramAccountService) CreateTelegramAccount(msgBody []byte, replyTo string, channel *amqp.Channel) {
	var createDTO dto.TelegramAccountCreateDTO
	if err := json.Unmarshal(msgBody, &createDTO); err != nil {
		utils.SendErrorResponse("Failed to unmarshal JSON", replyTo, channel)
		return
	}

	account := models.TelegramAccount{
		PhoneNumber: createDTO.PhoneNumber,
		AccountName: createDTO.AccountName,
	}

	newAccount, err := tas.repo.Create(&account)
	if err != nil {
		utils.SendErrorResponse(err.Error(), replyTo, channel)
		return
	}

	utils.SendSuccessResponse(newAccount, replyTo, channel)
}

// UpdateTelegramAccount handles updating an existing Telegram account
func (tas *TelegramAccountService) UpdateTelegramAccount(msgBody []byte, replyTo string, channel *amqp.Channel) {
	var updateDTO dto.TelegramAccountUpdateDTO
	if err := json.Unmarshal(msgBody, &updateDTO); err != nil {
		utils.SendErrorResponse("Failed to unmarshal JSON", replyTo, channel)
		return
	}

	account := models.TelegramAccount{
		ID:          updateDTO.ID,
		PhoneNumber: updateDTO.PhoneNumber,
		AccountName: updateDTO.AccountName,
	}

	updatedAccount, err := tas.repo.Update(&account)
	if err != nil {
		utils.SendErrorResponse(err.Error(), replyTo, channel)
		return
	}

	utils.SendSuccessResponse(updatedAccount, replyTo, channel)
}

// DeleteTelegramAccount handles deleting a Telegram account
func (tas *TelegramAccountService) DeleteTelegramAccount(msgBody []byte, replyTo string, channel *amqp.Channel) {
	var deleteDTO dto.TelegramAccountDeleteDTO
	if err := json.Unmarshal(msgBody, &deleteDTO); err != nil {
		utils.SendErrorResponse("Failed to unmarshal JSON", replyTo, channel)
		return
	}
	// Assuming deleteDTO.ID is of type uint
	if err := tas.repo.Delete(deleteDTO.ID); err != nil {
		utils.SendErrorResponse(err.Error(), replyTo, channel)
		return
	}

	utils.SendSuccessResponse(nil, replyTo, channel)
}

func (tas *TelegramAccountService) GetTelegramAccountByID(msgBody []byte, replyTo string, channel *amqp.Channel) {
	log.Print("Call Account By ID ")
	log.Printf("Body: %s", string(msgBody))

	var requestDTO struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(msgBody, &requestDTO); err != nil {
		utils.SendErrorResponse("Failed to unmarshal JSON", replyTo, channel)
		return
	}

	log.Printf("Parsed Request: ID=%s", requestDTO.ID)

	id, err := strconv.ParseUint(requestDTO.ID, 10, 32)
	if err != nil {
		utils.SendErrorResponse("Invalid ID format", replyTo, channel)
		return
	}

	log.Printf("Query account by ID: %d", id)
	account, err := tas.repo.GetByID(uint(id))
	if err != nil {
		utils.SendErrorResponse("Failed to retrieve account", replyTo, channel)
		return
	}

	utils.SendSuccessResponse(account, replyTo, channel)
}

func (us *TelegramAccountService) TelegramTimeBreakList(msgBody []byte, replyTo string, channel *amqp.Channel) {
	var requestData struct {
		Action string `json:"action"`
		Page   string `json:"page"`
		Limit  string `json:"limit"`
		Query  string `json:"query"`
	}
	err := json.Unmarshal(msgBody, &requestData)
	if err != nil {
		utils.SendErrorResponse("Failed to unmarshal JSON", replyTo, channel)
		return
	}
	pageInt, err := strconv.Atoi(requestData.Page)
	if err != nil {
		log.Printf("Failed to convert page to integer: %v", err)
		pageInt = 1 // Default to page 1 if conversion fails
	}

	limitInt, err := strconv.Atoi(requestData.Limit)
	if err != nil {
		log.Printf("Failed to convert limit to integer: %v", err)
		limitInt = 10 // Default to limit 10 if conversion fails
	}
	PermissionsList, total, err := us.repo.GetTelegramSettinglist(pageInt, limitInt, requestData.Query)
	if err != nil {
		utils.SendErrorResponse("Failed to get all Permissions", replyTo, channel)
		return
	}
	utils.SendPaginationResponse(PermissionsList, total, pageInt, limitInt, replyTo, channel)
}
