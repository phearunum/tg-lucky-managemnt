package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	desk "api-service/lib/desks"
	"api-service/utils"

	"github.com/gorilla/mux"
)

// DeskSettingController handles requests related to desk settings
type DeskSettingController struct {
	deskSettingService desk.DeskSettingServiceInterface
}

// NewDeskSettingController initializes a new DeskSettingController
func NewDeskSettingController(ds desk.DeskSettingServiceInterface) *DeskSettingController {
	return &DeskSettingController{deskSettingService: ds}
}

// DeskSettingListHandler handles requests to list desk settings
// @Summary List desk settings
// @Description Get a list of desk settings
// @Tags desk settings
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param query query string false "Search query"
// @Success 200 {object} utils.PaginationResponse
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /desksettings/list [get]
func (dsc *DeskSettingController) DeskSettingListHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	query := r.URL.Query().Get("query")

	// Default values for pagination
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil || pageInt <= 0 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil || limitInt <= 0 {
		limitInt = 10
	}

	requestDto := utils.PaginationRequestDTO{
		Page:  pageInt,
		Limit: limitInt,
		Query: query,
	}

	deskSettingListResponse := dsc.deskSettingService.DeskSettingServiceGetList(requestDto)
	utils.NewHttpPaginationResponse(w, deskSettingListResponse.Data, deskSettingListResponse.Meta.Total, deskSettingListResponse.Meta.Page, deskSettingListResponse.Meta.LastPage, int(deskSettingListResponse.Status), deskSettingListResponse.Message)
}

// DeskSettingByIDHandler handles requests to get a desk setting by ID
// @Summary Get desk setting by ID
// @Description Get details of a desk setting by ID
// @Tags desk settings
// @Accept json
// @Produce json
// @Param id path int true "Desk Setting ID"
// @Success 200 {object} desk.DeskSettingDTO
// @Failure 400 {string} string "Invalid desk setting ID"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /desksettings/{id} [get]
func (dsc *DeskSettingController) DeskSettingByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	deskSettingID, err := strconv.Atoi(id)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid desk setting ID")
		return
	}

	deskSettingDTO := &desk.DeskSettingResponseDTO{ID: uint(deskSettingID)}
	data, err := json.Marshal(deskSettingDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusInternalServerError, "Failed to marshal request data")
		return
	}

	response, err := dsc.deskSettingService.DeskSettingServiceGetById(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// DeskSettingCreateHandler handles requests to create a new desk setting
// @Summary Create a new desk setting
// @Description Create a new desk setting
// @Tags desk settings
// @Accept json
// @Produce json
// @Param deskSetting body desk.DeskSettingDTO true "Desk Setting data"
// @Success 200 {object} desk.DeskSettingDTO
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /desksettings/create [post]
func (dsc *DeskSettingController) DeskSettingCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createDeskSettingDTO desk.DeskSettingResponseDTO
	if err := json.NewDecoder(r.Body).Decode(&createDeskSettingDTO); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(createDeskSettingDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}

	response, err := dsc.deskSettingService.DeskSettingServiceCreate(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// DeskSettingUpdateHandler handles requests to update an existing desk setting
// @Summary Update an existing desk setting
// @Description Update an existing desk setting
// @Tags desk settings
// @Accept json
// @Produce json
// @Param deskSetting body desk.DeskSettingDTO true "Desk Setting data"
// @Success 200 {object} desk.DeskSettingDTO
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /desksettings/update [put]
func (dsc *DeskSettingController) DeskSettingUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateDeskSettingDTO desk.DeskSettingResponseDTO
	if err := json.NewDecoder(r.Body).Decode(&updateDeskSettingDTO); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(updateDeskSettingDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}

	response, err := dsc.deskSettingService.DeskSettingServiceUpdate(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// DeskSettingDeleteHandler handles requests to delete a desk setting
// @Summary Delete a desk setting
// @Description Delete a desk setting
// @Tags desk settings
// @Accept json
// @Produce json
// @Param id query int true "Desk Setting ID"
// @Success 200 {string} string "Desk setting deleted successfully"
// @Failure 400 {string} string "ID parameter is required"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /desksettings/delete [delete]
func (dsc *DeskSettingController) DeskSettingDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	deskSettingID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid desk setting ID", http.StatusBadRequest)
		return
	}

	deskSettingDTO := &desk.DeskSettingResponseDTO{ID: uint(deskSettingID)}
	data, err := json.Marshal(deskSettingDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}

	deleted, err := dsc.deskSettingService.DeskSettingServiceDelete(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, deleted, http.StatusOK, string(utils.SuccessMessage))
}
