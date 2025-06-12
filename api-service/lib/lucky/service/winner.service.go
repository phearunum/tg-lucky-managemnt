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

// WinnerServiceInterface defines the methods for the Winner service
type WinnerServiceInterface interface {
	WinnerServiceGetList(requestDto utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse
	WinnerServiceCreate(createDto dto.CreateWinnerDTO) (*dto.ResponseWinnerDTO, error)
	WinnerServiceUpdate(updateDto dto.UpdateWinnerDTO) (*dto.ResponseWinnerDTO, error)
	WinnerServiceGetById(requestDto dto.WinnerFilterDTO) (*dto.ResponseWinnerDTO, error)
	WinnerServiceGetByChatID(requestDto dto.WinnerFilterDTO) ([]*dto.ResponseWinnerDTO, error)
	WinnerServiceTaskExcuteGetWiner() bool
	ServiceTaskExcuteResetPoint() bool
	ServiceGetWinner(filters dto.WinnerFilterIDsDTO) ([]*dto.ResponseWinnerDTO, error)
	LuckySettingService(requestDto dto.TelegramSettingWinnerFilterDTO) (*dto.ResponseTelegramSettingWinnerDTO, error)
}

// WinnerService provides methods to manage chance points
type WinnerService struct {
	repo *repository.WinnerRepository
}

// NewWinnerService initializes a new WinnerService
func NewWinnerService(repo *repository.WinnerRepository) *WinnerService {
	return &WinnerService{repo: repo}
}

// WinnerServiceGetList retrieves a paginated list of chance points
func (cs *WinnerService) WinnerServiceGetList(requestDto utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse {
	if requestDto.Page <= 0 {
		requestDto.Page = 1
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}
	points, total, err := cs.repo.GetWinnerList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query, requestDto.Status, requestDto.TGgroup)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return utils.NewServicePaginationResponse(nil, 0, 0, 0, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel, "WinnerService [GetList]")
	}
	pointDTOs := make([]*dto.ResponseWinnerDTO, len(points))
	for i, point := range points {
		pointDTO := mapper.ToResponseWinnerDTO(*point)
		pointDTOs[i] = &pointDTO
	}
	return utils.NewServicePaginationResponse(pointDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "WinnerService [GetList]")
}

// WinnerServiceCreate creates a new chance point
func (cs *WinnerService) WinnerServiceCreate(createDto dto.CreateWinnerDTO) (*dto.ResponseWinnerDTO, error) {
	createdPoint, err := cs.repo.CreateWinner(createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	pointDTO := mapper.ToResponseWinnerDTO(*createdPoint)
	return &pointDTO, nil
}

// WinnerServiceUpdate updates an existing chance point
func (cs *WinnerService) WinnerServiceUpdate(updateDto dto.UpdateWinnerDTO) (*dto.ResponseWinnerDTO, error) {
	point, err := cs.repo.GetWinnerByID(updateDto.ID)
	if err != nil {
		utils.LoggerRepository(err, "GetWinnerByID")
		return nil, fmt.Errorf("failed to retrieve chance point: %v", err)
	}
	mapper.ApplyUpdateToWinner(point, updateDto)
	updatedPoint, err := cs.repo.UpdateWinner(updateDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedDto := mapper.ToResponseWinnerDTO(*updatedPoint)
	return &updatedDto, nil
}

// WinnerServiceGetById retrieves a chance point by ID
func (cs *WinnerService) WinnerServiceGetById(requestDto dto.WinnerFilterDTO) (*dto.ResponseWinnerDTO, error) {
	point, err := cs.repo.GetWinnerByID(requestDto.ID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	pointDTO := mapper.ToResponseWinnerDTO(*point)
	utils.InfoLog(pointDTO, string(utils.SuccessMessage))
	return &pointDTO, nil
}

// GetByChatID retrieves multiple chance points by ChatID
func (cs *WinnerService) WinnerServiceGetByChatID(requestDto dto.WinnerFilterDTO) ([]*dto.ResponseWinnerDTO, error) {
	points, err := cs.repo.GetChanceCollection(requestDto.ChatID)
	if err != nil {
		utils.ErrorLog(nil, "Failed to retrieve chance points: "+err.Error())
		return nil, err
	}
	pointDTOs := make([]*dto.ResponseWinnerDTO, len(points))
	for i, point := range points {
		pointDTO := mapper.ToResponseWinnerDTO(*point)
		pointDTOs[i] = &pointDTO
	}
	utils.InfoLog(pointDTOs, "Successfully retrieved chance points")
	return pointDTOs, nil
}
func (cs *WinnerService) WinnerServiceTaskExcuteGetWiner() bool {
	//	success := cs.repo.TaskExcuteWinner()
	success := cs.repo.TaskExcuteWinnerTask()
	if !success {
		utils.ErrorLog(nil, "Failed to execute winner task")
		return false
	}
	return true
}
func (cs *WinnerService) ServiceTaskExcuteResetPoint() bool {
	success := cs.repo.TaskExcutePointReset()
	if !success {
		utils.ErrorLog(nil, "Failed to execute winner task")
		return false
	}
	return true
}

type TelegramMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func (cs *WinnerService) ServiceGetWinner(filters dto.WinnerFilterIDsDTO) ([]*dto.ResponseWinnerDTO, error) {
	utils.WarnLog(filters.ID, "filters.ID")
	winners, err := cs.repo.GetWinnerToday(filters.ID, filters.Limit)
	if err != nil {
		utils.ErrorLog(nil, "Failed to retrieve chance points: "+err.Error())
		return nil, err
	}
	pointDTOs := make([]*dto.ResponseWinnerDTO, len(winners))
	for i, win := range winners {
		pointDTO := mapper.ToResponseWinnerDTO(*win)
		pointDTOs[i] = &pointDTO
	}
	utils.InfoLog(pointDTOs, "Successfully retrieved winner")
	return pointDTOs, nil

}

func (cs *WinnerService) LuckySettingService(requestDto dto.TelegramSettingWinnerFilterDTO) (*dto.ResponseTelegramSettingWinnerDTO, error) {
	point, err := cs.repo.GetLuckyWinnerSetting(*requestDto.GroupID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	pointDTO := mapper.ToResponseTelegramSettingWinnerDTO(*point)
	utils.InfoLog(pointDTO, string(utils.SuccessMessage))
	return &pointDTO, nil
}
