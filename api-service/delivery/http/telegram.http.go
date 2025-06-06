package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	telegram "api-service/lib/telegram" // Adjust the import path as necessary
	"api-service/utils"

	"github.com/gorilla/mux"
)

type UserRequestController struct {
	userRequestService telegram.UserRequestServiceInterface
}

func (urc *UserRequestController) GetPhoneListHandler(w http.ResponseWriter, r *http.Request) {
	utils.WarnLog(r.URL, "Request List")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	query := r.URL.Query().Get("query")
	status := r.URL.Query().Get("status")
	// Parse page number
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	// Parse limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	// Prepare request DTO without nil date fields
	requestDto := utils.PaginationRequestFilterDTO{
		Page:   page,
		Limit:  limit,
		Query:  query,
		Status: status,
	}
	// Call the service to get the list of user requests
	userRequestListResponse := urc.userRequestService.GetPhoneList(requestDto)
	// Send the response back to the client
	utils.NewHttpPaginationResponse(w, userRequestListResponse.Data, userRequestListResponse.Meta.Total, userRequestListResponse.Meta.Page, userRequestListResponse.Meta.LastPage, int(userRequestListResponse.Status), userRequestListResponse.Message)
}
func (urc *UserRequestController) RequestSavePhoneHandler(w http.ResponseWriter, r *http.Request) {
	var phoneListDTOs []telegram.CreatePhoneListDTO
	err := json.NewDecoder(r.Body).Decode(&phoneListDTOs)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	// Step 3: Call the service to save multiple phones
	responses, err := urc.userRequestService.SavePhone(phoneListDTOs)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, responses, http.StatusOK, string(utils.SuccessMessage))
}

func (urc *UserRequestController) RequestUpdatePhoneHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequestDTO telegram.UpdatePhoneListDTO
	err := json.NewDecoder(r.Body).Decode(&createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}

	response, err := urc.userRequestService.UpdatePhone(createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

func (urc *UserRequestController) RequestBulkDeletePhoneHandler(w http.ResponseWriter, r *http.Request) {
	var deleteRequest telegram.BulkDeleteRequestDTO

	// Decode the request body into the BulkDeleteRequestDTO struct
	err := json.NewDecoder(r.Body).Decode(&deleteRequest)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	// Call the service to delete multiple phones
	err = urc.userRequestService.DeletePhone(deleteRequest)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Failed to delete phones: "+err.Error())
		return
	}

	// Return success response
	utils.HttpSuccessResponse(w, nil, http.StatusOK, "Phones deleted successfully")
}

func (urc *UserRequestController) RequestSaveClockInHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequestDTO telegram.CreateClockTimeDTO
	utils.WarnLog(createUserRequestDTO, "createUserRequestDTO")
	err := json.NewDecoder(r.Body).Decode(&createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}

	response, err := urc.userRequestService.SaveClockTime(createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}
func NewUserRequestController(urs telegram.UserRequestServiceInterface) *UserRequestController {
	return &UserRequestController{userRequestService: urs}
}
func (urc *UserRequestController) GetClockTimeListHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	query := r.URL.Query().Get("query")
	// Parse page number
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	// Parse limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	// Prepare request DTO without nil date fields
	requestDto := utils.PaginationRequestFilterDTO{
		Page:  page,
		Limit: limit,
		Query: query,
	}
	// Call the service to get the list of user requests
	userRequestListResponse := urc.userRequestService.GetClockTimeList(requestDto)
	// Send the response back to the client
	utils.NewHttpPaginationResponse(w, userRequestListResponse.Data, userRequestListResponse.Meta.Total, userRequestListResponse.Meta.Page, userRequestListResponse.Meta.LastPage, int(userRequestListResponse.Status), userRequestListResponse.Message)
}
func (urc *UserRequestController) GetBotLocationListHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	query := r.URL.Query().Get("query")
	// Parse page number
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	// Parse limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	// Prepare request DTO without nil date fields
	requestDto := utils.PaginationRequestFilterDTO{
		Page:  page,
		Limit: limit,
		Query: query,
	}
	// Call the service to get the list of user requests
	userRequestListResponse := urc.userRequestService.GetBotLocationList(requestDto)

	// Send the response back to the client
	utils.NewHttpPaginationResponse(w, userRequestListResponse.Data, userRequestListResponse.Meta.Total, userRequestListResponse.Meta.Page, userRequestListResponse.Meta.LastPage, int(userRequestListResponse.Status), userRequestListResponse.Message)
}
func (urc *UserRequestController) RequestBotLocationCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequestDTO telegram.CreateLocationSettingDTO
	err := json.NewDecoder(r.Body).Decode(&createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}

	response, err := urc.userRequestService.SaveBotLocationSetting(createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}
func (urc *UserRequestController) RequestBotLocationUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequestDTO telegram.UpdateLocationSettingDTO
	err := json.NewDecoder(r.Body).Decode(&createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}

	response, err := urc.userRequestService.UpdateBotLocationSetting(createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

func (urc *UserRequestController) RequestSettingListHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	query := r.URL.Query().Get("query")
	// Parse page number
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	// Parse limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	// Prepare request DTO without nil date fields
	requestDto := utils.PaginationRequestFilterDTO{
		Page:  page,
		Limit: limit,
		Query: query,
	}
	// Call the service to get the list of user requests
	userRequestListResponse := urc.userRequestService.GetSettingList(requestDto)

	// Send the response back to the client
	utils.NewHttpPaginationResponse(w, userRequestListResponse.Data, userRequestListResponse.Meta.Total, userRequestListResponse.Meta.Page, userRequestListResponse.Meta.LastPage, int(userRequestListResponse.Status), userRequestListResponse.Message)
}

func (urc *UserRequestController) RequestSettingCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequestDTO telegram.CreateRequestSettingDTO
	err := json.NewDecoder(r.Body).Decode(&createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}

	response, err := urc.userRequestService.SaveRequestSetting(createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}
func (urc *UserRequestController) RequestSettingupdareHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequestDTO telegram.UpdateRequestSettingDTO
	err := json.NewDecoder(r.Body).Decode(&createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}

	response, err := urc.userRequestService.UpdateRequestSetting(createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

func (urc *UserRequestController) UserRequestListHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	query := r.URL.Query().Get("query")
	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")

	// Parse page number
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	// Parse limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Parse start and end dates
	var startDate, endDate *time.Time
	if startStr != "" {
		start, err := time.Parse("2006-01-02 15:04:05", startStr) // Adjust format as needed
		if err != nil {
			utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid start date format")
			return
		}
		startDate = &start
	}

	if endStr != "" {
		end, err := time.Parse("2006-01-02 15:04:05", endStr) // Adjust format as needed
		if err != nil {
			utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid end date format")
			return
		}
		endDate = &end
	}

	// Prepare request DTO without nil date fields
	requestDto := utils.PaginationRequestFilterDTO{
		Page:  page,
		Limit: limit,
		Query: query,
	}

	// Add date filters only if they are not nil
	if startDate != nil {
		requestDto.StartDate = startDate
	}
	if endDate != nil {
		requestDto.EndDate = endDate
	}

	// Call the service to get the list of user requests
	userRequestListResponse := urc.userRequestService.GetList(requestDto)

	// Send the response back to the client
	utils.NewHttpPaginationResponse(w, userRequestListResponse.Data, userRequestListResponse.Meta.Total, userRequestListResponse.Meta.Page, userRequestListResponse.Meta.LastPage, int(userRequestListResponse.Status), userRequestListResponse.Message)
}

func (urc *UserRequestController) UserRequestByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userRequestID, err := strconv.Atoi(id)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid ID")
		return
	}

	chatID := vars["chat_id"] // Assuming chat_id is a path parameter
	userRequestDTO := telegram.UserRequestFilter{
		ID:     uint(userRequestID),
		ChatID: chatID,
	}

	response, err := urc.userRequestService.GetByID(userRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

func (urc *UserRequestController) UserRequestCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequestDTO telegram.CreateUserRequestDTO
	err := json.NewDecoder(r.Body).Decode(&createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}

	response, err := urc.userRequestService.Create(createUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

func (urc *UserRequestController) UserRequestUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateUserRequestDTO telegram.UpdateUserRequestDTO
	err := json.NewDecoder(r.Body).Decode(&updateUserRequestDTO)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	response, err := urc.userRequestService.Update(updateUserRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

func (urc *UserRequestController) UserRequestDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userRequestID, err := strconv.Atoi(id)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid ID")
		return
	}

	chatID := vars["chat_id"] // Assuming chat_id is a path parameter
	userRequestDTO := telegram.UserRequestFilter{ID: uint(userRequestID), ChatID: chatID}

	deleted, err := urc.userRequestService.Delete(userRequestDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, false, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, deleted, http.StatusOK, string(utils.SuccessMessage))
}
