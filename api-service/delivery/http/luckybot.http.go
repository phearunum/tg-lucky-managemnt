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

type TelegramSettingWinnerController struct {
	telegramSettingWinnerService service.TelegramSettingWinnerServiceInterface
}

// NewTelegramSettingWinnerController creates a new instance of TelegramSettingWinnerController
func NewTelegramSettingWinnerController(service service.TelegramSettingWinnerServiceInterface) *TelegramSettingWinnerController {
	return &TelegramSettingWinnerController{telegramSettingWinnerService: service}
}

// GetTelegramSettingWinnerListHandler retrieves a paginated list of Telegram setting winners
func (uc *TelegramSettingWinnerController) GetTelegramSettingWinnerListHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	query := r.URL.Query().Get("query")
	tg_group := r.Header.Get("tg_group")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	requestDto := utils.PaginationRequestFilterDTO{
		Page:    page,
		Limit:   limit,
		Query:   query,
		TGgroup: tg_group,
	}

	response := uc.telegramSettingWinnerService.TelegramSettingWinnerServiceGetList(requestDto)
	utils.NewHttpPaginationResponse(w, response.Data, response.Meta.Total, response.Meta.Page, response.Meta.LastPage, int(response.Status), response.Message)
}

// TelegramSettingWinnerCreateHandler creates a new Telegram setting winner
func (uc *TelegramSettingWinnerController) TelegramSettingWinnerCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createDto dto.CreateTelegramSettingWinnerDTO
	if err := json.NewDecoder(r.Body).Decode(&createDto); err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}

	result, err := uc.telegramSettingWinnerService.TelegramSettingWinnerServiceCreate(createDto)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, result, http.StatusOK, string(utils.SuccessMessage))
}

// TelegramSettingWinnerUpdateHandler updates an existing Telegram setting winner
func (uc *TelegramSettingWinnerController) TelegramSettingWinnerUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateDto dto.UpdateTelegramSettingWinnerDTO
	if err := json.NewDecoder(r.Body).Decode(&updateDto); err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}

	result, err := uc.telegramSettingWinnerService.TelegramSettingWinnerServiceUpdate(updateDto)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, result, http.StatusOK, string(utils.SuccessMessage))
}

// TelegramSettingWinnerGetByIDHandler retrieves a Telegram setting winner by ID
func (uc *TelegramSettingWinnerController) TelegramSettingWinnerGetByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid ID")
		return
	}

	requestDto := dto.TelegramSettingWinnerFilterDTO{
		ID: uint(id),
	}

	result, err := uc.telegramSettingWinnerService.TelegramSettingWinnerServiceGetById(requestDto)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusNotFound, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, result, http.StatusOK, string(utils.SuccessMessage))
}
