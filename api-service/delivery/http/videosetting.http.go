package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	video "api-service/lib/videosnap/setting"
	"api-service/utils"

	"github.com/gorilla/mux"
)

// VideoSnapSettingController handles requests related to video snap settings
type VideoSnapSettingController struct {
	videoSnapSettingService video.VideoSnapSettingServiceInterface
}

// NewVideoSnapSettingController initializes a new VideoSnapSettingController
func NewVideoSnapSettingController(vs video.VideoSnapSettingServiceInterface) *VideoSnapSettingController {
	return &VideoSnapSettingController{videoSnapSettingService: vs}
}

// VideoSnapSettingListHandler handles requests to list video snap settings
// @Summary List video snap settings
// @Description Get a list of video snap settings
// @Tags video snap settings
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param query query string false "Search query"
// @Success 200 {object} utils.PaginationResponse
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /videosnapsettings/list [get]
func (vsc *VideoSnapSettingController) VideoSnapSettingListHandler(w http.ResponseWriter, r *http.Request) {
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

	// Call service to get list of settings
	videoSnapSettingListResponse := vsc.videoSnapSettingService.VideoSnapSettingServiceGetList(requestDto)
	utils.NewHttpPaginationResponse(w, videoSnapSettingListResponse.Data, videoSnapSettingListResponse.Meta.Total, videoSnapSettingListResponse.Meta.Page, videoSnapSettingListResponse.Meta.LastPage, int(videoSnapSettingListResponse.Status), videoSnapSettingListResponse.Message)
}

// VideoSnapSettingByIDHandler handles requests to get a video snap setting by ID
// @Summary Get video snap setting by ID
// @Description Get details of a video snap setting by ID
// @Tags video snap settings
// @Accept json
// @Produce json
// @Param id path int true "Video Snap Setting ID"
// @Success 200 {object} video.VideoSnapSettingDTO
// @Failure 400 {string} string "Invalid video snap setting ID"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /videosnapsettings/{id} [get]
func (vsc *VideoSnapSettingController) VideoSnapSettingByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	videoSnapSettingID, err := strconv.Atoi(id)

	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid video snap setting ID")
		return
	}

	videoSnapSettingDTO := &video.VideoSnapSettingResponseDTO{ID: uint(videoSnapSettingID)}
	data, err := json.Marshal(videoSnapSettingDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusInternalServerError, "Failed to marshal request data")
		return
	}

	response, err := vsc.videoSnapSettingService.VideoSnapSettingServiceGetById(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// VideoSnapSettingCreateHandler handles requests to create a new video snap setting
// @Summary Create a new video snap setting
// @Description Create a new video snap setting
// @Tags video snap settings
// @Accept json
// @Produce json
// @Param videoSnapSetting body video.VideoSnapSettingDTO true "Video Snap Setting data"
// @Success 200 {object} video.VideoSnapSettingDTO
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /videosnapsettings/create [post]
func (vsc *VideoSnapSettingController) VideoSnapSettingCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createVideoSnapSettingDTO video.VideoSnapSettingResponseDTO
	if err := json.NewDecoder(r.Body).Decode(&createVideoSnapSettingDTO); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(createVideoSnapSettingDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}

	response, err := vsc.videoSnapSettingService.VideoSnapSettingServiceCreate(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// VideoSnapSettingUpdateHandler handles requests to update an existing video snap setting
// @Summary Update an existing video snap setting
// @Description Update an existing video snap setting
// @Tags video snap settings
// @Accept json
// @Produce json
// @Param videoSnapSetting body video.VideoSnapSettingDTO true "Video Snap Setting data"
// @Success 200 {object} video.VideoSnapSettingDTO
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /videosnapsettings/update [put]
func (vsc *VideoSnapSettingController) VideoSnapSettingUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateVideoSnapSettingDTO video.VideoSnapSettingResponseDTO
	if err := json.NewDecoder(r.Body).Decode(&updateVideoSnapSettingDTO); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(updateVideoSnapSettingDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}

	response, err := vsc.videoSnapSettingService.VideoSnapSettingServiceUpdate(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// VideoSnapSettingDeleteHandler handles requests to delete a video snap setting
// @Summary Delete a video snap setting
// @Description Delete a video snap setting
// @Tags video snap settings
// @Accept json
// @Produce json
// @Param id query int true "Video Snap Setting ID"
// @Success 200 {string} string "Video snap setting deleted successfully"
// @Failure 400 {string} string "ID parameter is required"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /videosnapsettings/delete [delete]
func (vsc *VideoSnapSettingController) VideoSnapSettingDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	videoSnapSettingID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid video snap setting ID", http.StatusBadRequest)
		return
	}
	if videoSnapSettingID == 1 {
		http.Error(w, "Primary Server can't delete it", http.StatusInternalServerError)
		return
	}
	if videoSnapSettingID <= 1 {
		http.Error(w, "Primary Server can't delete it", http.StatusInternalServerError)
		return
	}
	videoSnapSettingDTO := &video.VideoSnapSettingResponseDTO{ID: uint(videoSnapSettingID)}
	data, err := json.Marshal(videoSnapSettingDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}

	deleted, err := vsc.videoSnapSettingService.VideoSnapSettingServiceDelete(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, deleted, http.StatusOK, string(utils.SuccessMessage))
}
