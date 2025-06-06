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

type ChancePointController struct {
	chancePointService service.ChancePointServiceInterface
}

// NewChancePointController creates a new instance of ChancePointController
func NewChancePointController(service service.ChancePointServiceInterface) *ChancePointController {
	return &ChancePointController{chancePointService: service}
}

// GetChancePointListHandler retrieves a paginated list of ChancePoints
func (cc *ChancePointController) GetChancePointListHandler(w http.ResponseWriter, r *http.Request) {
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

	// Prepare request DTO
	requestDto := utils.PaginationRequestFilterDTO{
		Page:    page,
		Limit:   limit,
		Query:   query,
		TGgroup: tg_group,
	}
	utils.InfoLog(requestDto, "requestDto")

	// Call the service to get the list of ChancePoints
	chancePointListResponse := cc.chancePointService.ChancePointServiceGetList(requestDto)

	// Send the response back to the client
	utils.NewHttpPaginationResponse(w, chancePointListResponse.Data, chancePointListResponse.Meta.Total, chancePointListResponse.Meta.Page, chancePointListResponse.Meta.LastPage, int(chancePointListResponse.Status), chancePointListResponse.Message)
}

// ChancePointCreateHandler creates a new ChancePoint
func (cc *ChancePointController) ChancePointCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createChancePointDTO dto.CreateChancePointDTO
	err := json.NewDecoder(r.Body).Decode(&createChancePointDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}

	response, err := cc.chancePointService.ChancePointServiceCreate(createChancePointDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// ChancePointUpdateHandler updates an existing ChancePoint
func (cc *ChancePointController) ChancePointUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateChancePointDTO dto.UpdateChancePointDTO
	err := json.NewDecoder(r.Body).Decode(&updateChancePointDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}
	utils.InfoLog(updateChancePointDTO, "updateChancePointDTO")
	response, err := cc.chancePointService.ChancePointServiceUpdate(updateChancePointDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// ChancePointGetByIDHandler retrieves a ChancePoint by ID
func (cc *ChancePointController) ChancePointGetByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid ID")
		return
	}
	requestDto := dto.ChancePointFilterDTO{
		ID: uint(id),
	}
	response, err := cc.chancePointService.ChancePointServiceGetById(requestDto)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusNotFound, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// ChancePointGetByIDHandler retrieves a ChancePoint by ID
func (cc *ChancePointController) ChancePointGetByChatIDHandler(w http.ResponseWriter, r *http.Request) {
	chatIdStr := r.URL.Query().Get("chat")

	if chatIdStr == "" {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	requestDto := dto.ChancePointFilterDTO{
		ChatID: chatIdStr,
	}
	response, err := cc.chancePointService.ChancePointServiceGetByChatID(requestDto)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusNotFound, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}
