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

// TelegramSettingWinnerServiceInterface defines the interface for TelegramSettingWinner operations
type TelegramSettingWinnerServiceInterface interface {
	TelegramSettingWinnerServiceGetList(requestDto utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse
	TelegramSettingWinnerServiceCreate(createDto dto.CreateTelegramSettingWinnerDTO) (*dto.ResponseTelegramSettingWinnerDTO, error)
	TelegramSettingWinnerServiceUpdate(updateDto dto.UpdateTelegramSettingWinnerDTO) (*dto.ResponseTelegramSettingWinnerDTO, error)
	TelegramSettingWinnerServiceGetById(requestDto dto.TelegramSettingWinnerFilterDTO) (*dto.ResponseTelegramSettingWinnerDTO, error)
}

// TelegramSettingWinnerService handles TelegramSettingWinner logic
type TelegramSettingWinnerService struct {
	repo *repository.TelegramSettingWinnerRepository
}

// NewTelegramSettingWinnerService creates a new instance of the service
func NewTelegramSettingWinnerService(repo *repository.TelegramSettingWinnerRepository) *TelegramSettingWinnerService {
	return &TelegramSettingWinnerService{repo: repo}
}

// TelegramSettingWinnerServiceGetList retrieves a paginated list of setting winners
func (ts *TelegramSettingWinnerService) TelegramSettingWinnerServiceGetList(requestDto utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse {
	if requestDto.Page <= 0 {
		requestDto.Page = 1
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}
	items, total, err := ts.repo.GetTelegramSettingWinnerList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query, requestDto.TGgroup)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return utils.NewServicePaginationResponse(nil, 0, 0, 0, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel, "TelegramSettingWinnerService [GetList]")
	}
	itemDTOs := make([]*dto.ResponseTelegramSettingWinnerDTO, len(items))
	for i, item := range items {
		dto := mapper.ToResponseTelegramSettingWinnerDTO(*item)
		itemDTOs[i] = &dto
	}
	return utils.NewServicePaginationResponse(itemDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "TelegramSettingWinnerService [GetList]")
}

// TelegramSettingWinnerServiceCreate creates a new setting winner
func (ts *TelegramSettingWinnerService) TelegramSettingWinnerServiceCreate(createDto dto.CreateTelegramSettingWinnerDTO) (*dto.ResponseTelegramSettingWinnerDTO, error) {
	item, err := ts.repo.CreateTelegramSettingWinner(createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	itemDTO := mapper.ToResponseTelegramSettingWinnerDTO(*item)
	return &itemDTO, nil
}

// TelegramSettingWinnerServiceUpdate updates an existing setting winner
func (ts *TelegramSettingWinnerService) TelegramSettingWinnerServiceUpdate(updateDto dto.UpdateTelegramSettingWinnerDTO) (*dto.ResponseTelegramSettingWinnerDTO, error) {
	existing, err := ts.repo.GetTelegramSettingWinnerByID(updateDto.ID)
	if err != nil {
		utils.LoggerRepository(err, "GetTelegramSettingWinnerByID")
		return nil, fmt.Errorf("failed to retrieve Telegram setting winner: %v", err)
	}
	mapper.ApplyUpdateToTelegramSettingWinner(existing, updateDto)
	updated, err := ts.repo.UpdateTelegramSettingWinner(updateDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedDTO := mapper.ToResponseTelegramSettingWinnerDTO(*updated)
	return &updatedDTO, nil
}

// TelegramSettingWinnerServiceGetById retrieves a setting winner by ID
func (ts *TelegramSettingWinnerService) TelegramSettingWinnerServiceGetById(requestDto dto.TelegramSettingWinnerFilterDTO) (*dto.ResponseTelegramSettingWinnerDTO, error) {
	item, err := ts.repo.GetTelegramSettingWinnerByID(requestDto.ID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	itemDTO := mapper.ToResponseTelegramSettingWinnerDTO(*item)
	utils.InfoLog(itemDTO, string(utils.SuccessMessage))
	return &itemDTO, nil
}
