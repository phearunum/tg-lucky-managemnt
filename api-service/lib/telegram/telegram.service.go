package telegram

import (
	redis "api-service/config"
	"api-service/utils"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// UserRequestServiceInterface defines the methods for the UserRequestService
type UserRequestServiceInterface interface {
	GetList(requestDto utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse
	Create(dtoCreate CreateUserRequestDTO) (*ResponseUserRequestDTO, error)
	Update(updateDTO UpdateUserRequestDTO) (*ResponseUserRequestDTO, error)
	GetByID(requestDTO UserRequestFilter) (*ResponseUserRequestDTO, error)
	Delete(requestDTO UserRequestFilter) (bool, error)
	GetSettingList(q utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse
	SaveRequestSetting(createDTO CreateRequestSettingDTO) (*ResponseRequestSettingDTO, error)
	UpdateRequestSetting(createDTO UpdateRequestSettingDTO) (*ResponseRequestSettingDTO, error)

	GetBotLocationList(q utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse
	SaveBotLocationSetting(createDTO CreateLocationSettingDTO) (*ResponseLocationSettingDTO, error)
	UpdateBotLocationSetting(createDTO UpdateLocationSettingDTO) (*ResponseLocationSettingDTO, error)

	GetClockTimeList(q utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse
	SaveClockTime(createDTO CreateClockTimeDTO) (*ResponseClockTimeDTO, error)
	GetPhoneList(q utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse
	SavePhone(createDTOs []CreatePhoneListDTO) ([]ResponsePhoneListDTO, error)
	UpdatePhone(createDTO UpdatePhoneListDTO) (*ResponsePhoneListDTO, error)
	DeletePhone(requestDTO BulkDeleteRequestDTO) error
}

type UserRequestService struct {
	repo *UserRequestRepository
}

func NewUserRequestService(repo *UserRequestRepository) *UserRequestService {
	return &UserRequestService{repo: repo}
}

// logError logs the error and returns it
func (s *UserRequestService) logError(err error) error {
	utils.ErrorLog(err.Error(), "internal server error")
	return fmt.Errorf("internal server error: %v", err)
}

// GetClockList
// GetList retrieves a paginated list of user requests
func (s *UserRequestService) GetClockTimeList(q utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse {
	// Validate page number and set default if invalid
	if q.Page <= 0 {
		q.Page = 1
	}
	// Validate limit and set default if invalid
	if q.Limit <= 0 {
		q.Limit = 10
	}
	// Fetch user requests from the repository
	requests, total, err := s.repo.GetClockList(int(q.Page), int(q.Limit), q.Query)
	if err != nil {
		// Return an error response if fetching requests fails
		return utils.NewServicePaginationResponse(nil, 0, q.Page, q.Limit, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel, "UserRequestService [GetList]")
	}
	// Map the retrieved requests to DTOs
	responseDTOs := make([]*ResponseClockTimeDTO, len(requests))
	for i, request := range requests {
		responseDTOs[i] = ToClockTimeResponseDTO(*request)
	}
	// Construct and return a successful response
	return utils.NewServicePaginationResponse(responseDTOs, total, q.Page, q.Limit, http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "UserRequestService [GetList]")
}
func (s *UserRequestService) SaveClockTime(createDTO CreateClockTimeDTO) (*ResponseClockTimeDTO, error) {
	createdRequest, err := s.repo.SaveClockIn(createDTO)
	if err != nil {
		return nil, s.logError(err)
	}

	responseDTO := ToClockTimeResponseDTO(*createdRequest)
	return responseDTO, nil
}

// GetList retrieves a paginated list of user requests
func (s *UserRequestService) GetSettingList(q utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse {
	// Validate page number and set default if invalid
	if q.Page <= 0 {
		q.Page = 1
	}
	// Validate limit and set default if invalid
	if q.Limit <= 0 {
		q.Limit = 10
	}
	// Fetch user requests from the repository
	requests, total, err := s.repo.GetTelegramSettingList(int(q.Page), int(q.Limit), q.Query)
	if err != nil {
		// Return an error response if fetching requests fails
		return utils.NewServicePaginationResponse(nil, 0, q.Page, q.Limit, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel, "UserRequestService [GetList]")
	}
	// Map the retrieved requests to DTOs
	responseDTOs := make([]*ResponseRequestSettingDTO, len(requests))
	for i, request := range requests {
		responseDTOs[i] = ToRequestSettingResponseDTO(*request)
	}
	// Construct and return a successful response
	return utils.NewServicePaginationResponse(responseDTOs, total, q.Page, q.Limit, http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "UserRequestService [GetList]")
}

func (s *UserRequestService) SaveRequestSetting(createDTO CreateRequestSettingDTO) (*ResponseRequestSettingDTO, error) {
	createdRequest, err := s.repo.SaveRequestSetting(createDTO)
	if err != nil {
		return nil, s.logError(err)
	}
	// Init
	keybords, errf := s.repo.FilterRequestSettingAll()
	if errf != nil {
		utils.ErrorLog(errf.Error(), "Failed init redis")

	}
	redisKey := fmt.Sprintf("request:bot:%s", "setting")
	if err := redis.SetWithoutExpired(redisKey, keybords); err != nil {
		utils.WarnLog(err.Error(), "Error Redis Init")
	}
	responseDTO := ToRequestSettingResponseDTO(*createdRequest)
	return responseDTO, nil
}
func (s *UserRequestService) UpdateRequestSetting(createDTO UpdateRequestSettingDTO) (*ResponseRequestSettingDTO, error) {
	createdRequest, err := s.repo.UpdateRequestSetting(createDTO)
	if err != nil {
		return nil, s.logError(err)
	}
	// Init
	keybords, errf := s.repo.FilterRequestSettingAll()
	if errf != nil {
		utils.ErrorLog(errf.Error(), "Failed init redis")

	}
	redisKey := fmt.Sprintf("request:bot:%s", "setting")
	if err := redis.SetWithoutExpired(redisKey, keybords); err != nil {
		utils.WarnLog(err.Error(), "Error Redis Init")
	}
	responseDTO := ToRequestSettingResponseDTO(*createdRequest)
	return responseDTO, nil
}

// GetList retrieves a paginated list of user requests
func (s *UserRequestService) GetList(q utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse {
	// Validate page number and set default if invalid
	if q.Page <= 0 {
		q.Page = 1
	}
	// Validate limit and set default if invalid
	if q.Limit <= 0 {
		q.Limit = 10
	}
	// Fetch user requests from the repository
	requests, total, err := s.repo.GetUserRequestList(int(q.Page), int(q.Limit), q.Query, q.StartDate, q.EndDate)
	if err != nil {
		// Return an error response if fetching requests fails
		return utils.NewServicePaginationResponse(nil, 0, q.Page, q.Limit, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel, "UserRequestService [GetList]")
	}
	// Map the retrieved requests to DTOs
	responseDTOs := make([]*ResponseUserRequestDTO, len(requests))
	for i, request := range requests {
		responseDTOs[i] = ToUserRequestResponseDTO(*request)
	}
	// Construct and return a successful response
	return utils.NewServicePaginationResponse(responseDTOs, total, q.Page, q.Limit, http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "UserRequestService [GetList]")
}

// Create creates a new user request
func (s *UserRequestService) Create(createDTO CreateUserRequestDTO) (*ResponseUserRequestDTO, error) {
	utils.InfoLog(createDTO, "createDTO")
	createdRequest, err := s.repo.CreateUserRequest(createDTO)
	if err != nil {
		return nil, s.logError(err)
	}
	responseDTO := ToUserRequestResponseDTO(*createdRequest)
	return responseDTO, nil // Return the created user request response DTO
}

// Update modifies an existing user request
func (s *UserRequestService) Update(updateDTO UpdateUserRequestDTO) (*ResponseUserRequestDTO, error) {
	updatedRequest, err := s.repo.UpdateUserRequest(updateDTO)
	if err != nil {
		return nil, s.logError(err)
	}
	updatedRequestDTO := ToUserRequestResponseDTO(*updatedRequest)
	return updatedRequestDTO, nil
}

// GetByID retrieves a user request by ID
func (s *UserRequestService) GetByID(requestDTO UserRequestFilter) (*ResponseUserRequestDTO, error) {
	request, err := s.repo.GetUserRequestByID(requestDTO.ID, requestDTO.ChatID)
	if err != nil {
		return nil, s.logError(err)
	}
	responseDTO := ToUserRequestResponseDTO(*request)
	return responseDTO, nil
}

// Delete removes a user request by ID
func (s *UserRequestService) Delete(requestDTO UserRequestFilter) (bool, error) {
	success, err := s.repo.DeleteUserRequest(requestDTO.ID, requestDTO.ChatID)
	if err != nil {
		return false, s.logError(err)
	}

	return success, nil
}

// Location

// GetList retrieves a paginated list of user requests
func (s *UserRequestService) GetBotLocationList(q utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse {
	// Validate page number and set default if invalid
	if q.Page <= 0 {
		q.Page = 1
	}
	// Validate limit and set default if invalid
	if q.Limit <= 0 {
		q.Limit = 10
	}
	// Fetch user requests from the repository
	requests, total, err := s.repo.GetLocationSettingList(int(q.Page), int(q.Limit), q.Query)
	if err != nil {
		// Return an error response if fetching requests fails
		return utils.NewServicePaginationResponse(nil, 0, q.Page, q.Limit, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel, "UserRequestService [GetList]")
	}
	// Map the retrieved requests to DTOs
	responseDTOs := make([]*ResponseLocationSettingDTO, len(requests))
	for i, request := range requests {
		responseDTOs[i] = ToLocationSettingResponseDTO(*request)
	}
	// Construct and return a successful response
	return utils.NewServicePaginationResponse(responseDTOs, total, q.Page, q.Limit, http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "UserRequestService [GetList]")
}

func (s *UserRequestService) SaveBotLocationSetting(createDTO CreateLocationSettingDTO) (*ResponseLocationSettingDTO, error) {
	createdRequest, err := s.repo.SaveLocationSetting(createDTO)
	if err != nil {
		return nil, s.logError(err)
	}
	// Init
	BotInfo, errf := s.repo.FilterLocationSetting(createDTO.BotToken)
	if errf != nil {
		utils.ErrorLog(errf.Error(), "Failed init redis")

	}
	redisKey := fmt.Sprintf("location:bot:%s", createDTO.Name)
	if err := redis.SetWithoutExpired(redisKey, BotInfo); err != nil {
		utils.WarnLog(err.Error(), "Error Redis Init")
	}
	responseDTO := ToLocationSettingResponseDTO(*createdRequest)
	return responseDTO, nil
}
func (s *UserRequestService) UpdateBotLocationSetting(createDTO UpdateLocationSettingDTO) (*ResponseLocationSettingDTO, error) {
	createdRequest, err := s.repo.UpdateLocationSetting(createDTO)
	if err != nil {
		return nil, s.logError(err)
	}
	// Init
	redisKey := fmt.Sprintf("location:bot:%s", createDTO.Name)
	// Remove

	if errr := redis.RemoveRedisKey(redisKey); errr != nil {
		utils.WarnLog(errr.Error(), "Error Redis Init")
	}

	BotInfo, errf := s.repo.FilterLocationSetting(createDTO.BotToken)
	if errf != nil {
		utils.ErrorLog(errf.Error(), "Failed init redis")
	}

	if err := redis.SetWithoutExpired(redisKey, BotInfo); err != nil {
		utils.WarnLog(err.Error(), "Error Redis Init")
	}
	responseDTO := ToLocationSettingResponseDTO(*createdRequest)
	return responseDTO, nil
}

func (s *UserRequestService) GetPhoneList(q utils.PaginationRequestFilterDTO) utils.ServicePaginationResponse {
	utils.InfoLog(q, "PaginationRequestFilterDTO")
	// Validate page number and set default if invalid
	if q.Page <= 0 {
		q.Page = 1
	}
	// Validate limit and set default if invalid
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Status == "" {
		q.Status = "yes"
	}
	// Fetch user requests from the repository
	requests, total, err := s.repo.GetPhoneList(int(q.Page), int(q.Limit), q.Query, q.Status)
	if err != nil {
		// Return an error response if fetching requests fails
		return utils.NewServicePaginationResponse(nil, 0, q.Page, q.Limit, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel, "GetPhoneList [GetList]")
	}
	// Map the retrieved requests to DTOs
	responseDTOs := make([]*ResponsePhoneListDTO, len(requests))
	for i, request := range requests {
		responseDTOs[i] = ToPhoneListResponseDTO(*request)
	}
	// Construct and return a successful response
	return utils.NewServicePaginationResponse(responseDTOs, total, q.Page, q.Limit, http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "GetPhoneList [GetList]")
}

func (s *UserRequestService) SavePhone(createDTOs []CreatePhoneListDTO) ([]ResponsePhoneListDTO, error) {
	createdPhoneLists, err := s.repo.SavePhones(createDTOs)
	if err != nil {
		return nil, s.logError(err)
	}

	// Convert to response DTOs
	var responseDTOs []ResponsePhoneListDTO
	for _, phoneList := range createdPhoneLists {
		responseDTOs = append(responseDTOs, *ToPhoneListResponseDTO(phoneList))
	}

	return responseDTOs, nil
}

func (s *UserRequestService) UpdatePhone(createDTO UpdatePhoneListDTO) (*ResponsePhoneListDTO, error) {
	createdRequest, err := s.repo.UpdatePhone(createDTO)
	if err != nil {
		return nil, s.logError(err)
	}
	responseDTO := ToPhoneListResponseDTO(*createdRequest)
	return responseDTO, nil
}

func (s *UserRequestService) DeletePhone(requestDTO BulkDeleteRequestDTO) error {
	return s.repo.BulkPhoneDelete(requestDTO.ID)
}
