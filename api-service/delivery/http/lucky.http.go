package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	dto "api-service/lib/lucky/dto"
	service "api-service/lib/lucky/service"
	"api-service/utils"

	"github.com/gorilla/mux"
)

type TelegramUserController struct {
	telegramUserService service.TelegramUserServiceInterface
}

// NewTelegramUserController creates a new instance of TelegramUserController
func NewTelegramUserController(service service.TelegramUserServiceInterface) *TelegramUserController {
	return &TelegramUserController{telegramUserService: service}
}

// GetTelegramUserListHandler retrieves a paginated list of Telegram users
func (uc *TelegramUserController) GetTelegramUserListHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	query := r.URL.Query().Get("query")
	tg_group := r.Header.Get("tg_group")
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
		Page:    page,
		Limit:   limit,
		Query:   query,
		TGgroup: tg_group,
	}
	// Call the service to get the list of Telegram users
	telegramUserListResponse := uc.telegramUserService.TelegramUserServiceGetList(requestDto)

	// Send the response back to the client
	utils.NewHttpPaginationResponse(w, telegramUserListResponse.Data, telegramUserListResponse.Meta.Total, telegramUserListResponse.Meta.Page, telegramUserListResponse.Meta.LastPage, int(telegramUserListResponse.Status), telegramUserListResponse.Message)
}

// TelegramUserCreateHandler creates a new Telegram user
func (uc *TelegramUserController) TelegramUserCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createUserDTO dto.CreateTelegramUserDTO
	err := json.NewDecoder(r.Body).Decode(&createUserDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}

	response, err := uc.telegramUserService.TelegramUserServiceCreate(createUserDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// TelegramUserUpdateHandler updates an existing Telegram user
func (uc *TelegramUserController) TelegramUserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateUserDTO dto.UpdateTelegramUserDTO
	err := json.NewDecoder(r.Body).Decode(&updateUserDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}
	response, err := uc.telegramUserService.TelegramUserServiceUpdate(updateUserDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// TelegramUserGetByIDHandler retrieves a Telegram user by ID
func (uc *TelegramUserController) TelegramUserGetByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid user ID")
		return
	}
	requestDto := dto.TelegramUserFilterDTO{
		ID: uint(id),
	}
	response, err := uc.telegramUserService.TelegramUserServiceGetById(requestDto)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusNotFound, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}
func (uc *TelegramUserController) TelegramUserGetByChatIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	requestDto := dto.TelegramUserFilterDTO{
		ChatID: &idStr,
	}
	response, err := uc.telegramUserService.TelegramUserServiceGetByChatId(requestDto)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusNotFound, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}
