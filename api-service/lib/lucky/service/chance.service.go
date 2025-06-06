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

// ChancePointServiceInterface defines the methods for the ChancePoint service
type ChancePointServiceInterface interface {
	ChancePointServiceGetList(requestDto utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse
	ChancePointServiceCreate(createDto dto.CreateChancePointDTO) (*dto.ResponseChancePointDTO, error)
	ChancePointServiceUpdate(updateDto dto.UpdateChancePointDTO) (*dto.ResponseChancePointDTO, error)
	ChancePointServiceGetById(requestDto dto.ChancePointFilterDTO) (*dto.ResponseChancePointDTO, error)
	ChancePointServiceGetByChatID(requestDto dto.ChancePointFilterDTO) ([]*dto.ResponseChancePointDTO, error)
}

// ChancePointService provides methods to manage chance points
type ChancePointService struct {
	repo *repository.ChancePointRepository
}

// NewChancePointService initializes a new ChancePointService
func NewChancePointService(repo *repository.ChancePointRepository) *ChancePointService {
	return &ChancePointService{repo: repo}
}

// ChancePointServiceGetList retrieves a paginated list of chance points
func (cs *ChancePointService) ChancePointServiceGetList(requestDto utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse {
	if requestDto.Page <= 0 {
		requestDto.Page = 1
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}
	utils.InfoLog(requestDto, "ChancePointServiceGetList")
	points, total, err := cs.repo.GetChancePointList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query, requestDto.TGgroup)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return utils.NewServicePaginationResponse(nil, 0, 0, 0, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel, "ChancePointService [GetList]")
	}
	pointDTOs := make([]*dto.ResponseChancePointDTO, len(points))
	for i, point := range points {
		pointDTO := mapper.ToResponseChancePointDTO(*point)
		pointDTOs[i] = &pointDTO
	}
	return utils.NewServicePaginationResponse(pointDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "ChancePointService [GetList]")
}

// ChancePointServiceCreate creates a new chance point
func (cs *ChancePointService) ChancePointServiceCreate(createDto dto.CreateChancePointDTO) (*dto.ResponseChancePointDTO, error) {
	createdPoint, err := cs.repo.CreateChancePoint(createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	pointDTO := mapper.ToResponseChancePointDTO(*createdPoint)
	return &pointDTO, nil
}

// ChancePointServiceUpdate updates an existing chance point
func (cs *ChancePointService) ChancePointServiceUpdate(updateDto dto.UpdateChancePointDTO) (*dto.ResponseChancePointDTO, error) {
	point, err := cs.repo.GetChancePointByID(updateDto.ID)
	if err != nil {
		utils.LoggerRepository(err, "GetChancePointByID")
		return nil, fmt.Errorf("failed to retrieve chance point: %v", err)
	}
	utils.InfoLog(updateDto, "updateDto")
	mapper.ApplyUpdateToChancePoint(point, updateDto)
	updatedPoint, err := cs.repo.UpdateChancePoint(updateDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedDto := mapper.ToResponseChancePointDTO(*updatedPoint)
	return &updatedDto, nil
}

// ChancePointServiceGetById retrieves a chance point by ID
func (cs *ChancePointService) ChancePointServiceGetById(requestDto dto.ChancePointFilterDTO) (*dto.ResponseChancePointDTO, error) {
	point, err := cs.repo.GetChancePointByID(requestDto.ID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	pointDTO := mapper.ToResponseChancePointDTO(*point)
	utils.InfoLog(pointDTO, string(utils.SuccessMessage))
	return &pointDTO, nil
}

// GetByChatID retrieves multiple chance points by ChatID
func (cs *ChancePointService) ChancePointServiceGetByChatID(requestDto dto.ChancePointFilterDTO) ([]*dto.ResponseChancePointDTO, error) {
	points, err := cs.repo.GetChanceCollection(requestDto.ChatID)
	if err != nil {
		utils.ErrorLog(nil, "Failed to retrieve chance points: "+err.Error())
		return nil, err
	}
	utils.InfoLog(requestDto, "Check Point list")
	pointDTOs := make([]*dto.ResponseChancePointDTO, len(points))
	for i, point := range points {
		pointDTO := mapper.ToResponseChancePointDTO(*point)
		pointDTOs[i] = &pointDTO
	}
	utils.InfoLog(pointDTOs, "Successfully retrieved chance points")
	return pointDTOs, nil
}
