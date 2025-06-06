package service

import (
	"api-service/lib/lucky/dto"
	"api-service/lib/lucky/mapper"
	"api-service/lib/lucky/repository"
	"api-service/utils"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// TelegramUserServiceInterface defines the methods for the TelegramUserService
type TelegramUserServiceInterface interface {
	TelegramUserServiceGetList(requestDto utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse
	TelegramUserServiceCreate(createDto dto.CreateTelegramUserDTO) (*dto.ResponseTelegramUserDTO, error)
	TelegramUserServiceUpdate(updateDto dto.UpdateTelegramUserDTO) (*dto.ResponseTelegramUserDTO, error)
	TelegramUserServiceGetById(requestDto dto.TelegramUserFilterDTO) (*dto.ResponseTelegramUserDTO, error)
	TelegramUserServiceGetByChatId(requestDto dto.TelegramUserFilterDTO) (*dto.ResponseTelegramUserDTO, error)
}
type TelegramUserService struct {
	repo *repository.TelegramUserRepository
}

func NewTelegramUserService(repo *repository.TelegramUserRepository) *TelegramUserService {
	return &TelegramUserService{repo: repo}
}

// TelegramUserServiceGetList retrieves a paginated list of Telegram users
func (ts *TelegramUserService) TelegramUserServiceGetList(requestDto utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse {
	if requestDto.Page <= 0 {
		requestDto.Page = 1
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}
	utils.InfoLog(requestDto, "TelegramUserServiceGetList")
	users, total, err := ts.repo.GetTelegramUserList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query, requestDto.TGgroup)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		utils.ErrorLog(err.Error(), err.Error())
		return utils.NewServicePaginationResponse(nil, 0, 0, 0, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel, "TelegramUserService [GetList]")
	}
	userDTOs := make([]*dto.ResponseTelegramUserDTO, len(users))
	for i, user := range users {
		userDTO := mapper.ToResponseTelegramUserDTO(*user)
		userDTOs[i] = &userDTO
	}
	return utils.NewServicePaginationResponse(userDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "TelegramUserService [GetList]")
}

// TelegramUserServiceCreate creates a new Telegram user
func (ts *TelegramUserService) TelegramUserServiceCreate(createDto dto.CreateTelegramUserDTO) (*dto.ResponseTelegramUserDTO, error) {
	createdUser, err := ts.repo.CreateTelegramUser(createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	userDTO := mapper.ToResponseTelegramUserDTO(*createdUser)
	return &userDTO, nil
}

// TelegramUserServiceUpdate updates an existing Telegram user
func (ts *TelegramUserService) TelegramUserServiceUpdate(updateDto dto.UpdateTelegramUserDTO) (*dto.ResponseTelegramUserDTO, error) {
	user, err := ts.repo.GetTelegramUserByID(updateDto.ID)
	if err != nil {
		utils.LoggerRepository(err, "GetTelegramUserByID")
		return nil, fmt.Errorf("failed to retrieve Telegram user: %v", err)
	}
	mapper.ApplyUpdateToTelegramUser(user, updateDto)
	updatedUser, err := ts.repo.UpdateTelegramUser(updateDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedDto := mapper.ToResponseTelegramUserDTO(*updatedUser)
	return &updatedDto, nil
}

// TelegramUserServiceGetById retrieves a Telegram user by ID
func (ts *TelegramUserService) TelegramUserServiceGetById(requestDto dto.TelegramUserFilterDTO) (*dto.ResponseTelegramUserDTO, error) {
	user, err := ts.repo.GetTelegramUserByID(requestDto.ID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	userDTO := mapper.ToResponseTelegramUserDTO(*user)
	utils.InfoLog(userDTO, string(utils.SuccessMessage))
	return &userDTO, nil
}

func (ts *TelegramUserService) TelegramUserServiceGetByChatId(requestDto dto.TelegramUserFilterDTO) (*dto.ResponseTelegramUserDTO, error) {
	user, err := ts.repo.GetTelegramUserByChatID(*requestDto.ChatID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	userDTO := mapper.ToResponseTelegramUserDTO(*user)
	utils.InfoLog(userDTO, string(utils.SuccessMessage))
	return &userDTO, nil
}
