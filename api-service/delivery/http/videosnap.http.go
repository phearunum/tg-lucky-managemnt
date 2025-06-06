package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	video "api-service/lib/videosnap"
	"api-service/utils"

	"github.com/gorilla/mux"
)

type VideoSnapController struct {
	videoSnapService video.VideosnapServiceInterface
}

func NewVideoSnapController(vs video.VideosnapServiceInterface) *VideoSnapController {
	return &VideoSnapController{videoSnapService: vs}
}

// VideoSnapListHandler handles requests to list video snapshots
// @Summary List video snapshots
// @Description Get a list of video snapshots
// @Tags video snapshots
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param query query string false "Search query"
// @Success 200 {object} utils.PaginationResponse
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /videosnaps/list [get]
func (vc *VideoSnapController) VideoSnapListHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	query := r.URL.Query().Get("query")

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

	videoSnapListResponse := vc.videoSnapService.VideosnapServiceGetList(requestDto)
	utils.NewHttpPaginationResponse(w, videoSnapListResponse.Data, videoSnapListResponse.Meta.Total, videoSnapListResponse.Meta.Page, videoSnapListResponse.Meta.LastPage, int(videoSnapListResponse.Status), videoSnapListResponse.Message)
}

// VideoSnapByIDHandler handles requests to get a video snapshot by ID
// @Summary Get video snapshot by ID
// @Description Get details of a video snapshot by ID
// @Tags video snapshots
// @Accept  json
// @Produce  json
// @Param id path int true "Video Snapshot ID"
// @Success 200 {object} video.VideoSnapDTO
// @Failure 400 {string} string "Invalid video snapshot ID"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /videosnaps/{id} [get]
func (vc *VideoSnapController) VideoSnapByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	videoSnapID, err := strconv.Atoi(id)

	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid video snapshot ID")
		return
	}

	videoSnapDTO := &video.VideosnapResponseDTO{ID: uint(videoSnapID)}
	data, err := json.Marshal(videoSnapDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusInternalServerError, "Failed to marshal request data")
		return
	}

	response, err := vc.videoSnapService.VideosnapServiceGetById(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// VideoSnapCreateHandler handles requests to create a new video snapshot
// @Summary Create a new video snapshot
// @Description Create a new video snapshot
// @Tags video snapshots
// @Accept  json
// @Produce  json
// @Param videoSnap body video.VideoSnapDTO true "Video Snapshot data"
// @Success 200 {object} video.VideoSnapDTO
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /videosnaps/create [post]
func (vc *VideoSnapController) VideoSnapCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createVideoSnapDTO video.VideosnapCreateDTO

	// Detect content type
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		// Handle JSON request body
		err := json.NewDecoder(r.Body).Decode(&createVideoSnapDTO)
		if err != nil {
			utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid JSON payload")
			return
		}

	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		// Handle form-data
		err := r.ParseMultipartForm(30 << 50) // Limit to 10MB
		if err != nil {
			utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid form-data payload")
			return
		}
		// Extract form-data values
		createVideoSnapDTO.GameNo = r.FormValue("period")
		createVideoSnapDTO.Streamkey = r.FormValue("rtmpurl")
	} else if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		// Handle application/x-www-form-urlencoded
		err := r.ParseForm() // Parse the form data
		if err != nil {
			utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid form data payload")
			return
		}

		// Extract form values
		createVideoSnapDTO.GameNo = r.FormValue("period")     // Adjust field names accordingly
		createVideoSnapDTO.Streamkey = r.FormValue("rtmpurl") // Adjust field names accordingly

		//fmt.Println("Received Form Data: Period =", createVideoSnapDTO.Period, ", RTMP URL =", createVideoSnapDTO.Rtmpurl) // Log form data
	} else {
		utils.HttpSuccessResponse(w, nil, http.StatusUnsupportedMediaType, "Unsupported Content-Type")
		return
	}
	// Validate the required parameters
	if createVideoSnapDTO.GameNo == "" || createVideoSnapDTO.Streamkey == "" {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "period and rtmpurl are required")
		return
	}

	// Call the service method
	_, err := vc.videoSnapService.VideosnapServiceCreate(createVideoSnapDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}
	//utils.InfoLog(response, "Success")
	responseTask := video.VideoSnapTaskResponseDTO{
		Status: http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseTask)

}

// VideoSnapUpdateHandler handles requests to update an existing video snapshot
// @Summary Update an existing video snapshot
// @Description Update an existing video snapshot
// @Tags video snapshots
// @Accept  json
// @Produce  json
// @Param videoSnap body video.VideoSnapDTO true "Video Snapshot data"
// @Success 200 {object} video.VideoSnapDTO
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to marshal request data"
// @Router /videosnaps/update [put]
func (vc *VideoSnapController) VideoSnapUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateVideoSnapDTO video.VideosnapStartDTO
	// Detect content type
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		// Handle JSON request body
		err := json.NewDecoder(r.Body).Decode(&updateVideoSnapDTO)
		if err != nil {
			utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid JSON payload")
			return
		}
	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		// Handle form-data
		err := r.ParseMultipartForm(10 << 20) // Limit to 10MB
		if err != nil {
			utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid form-data payload")
			return
		}

		// Extract form-data values
		updateVideoSnapDTO.GameNo = r.FormValue("period")
		updateVideoSnapDTO.Streamkey = r.FormValue("rtmpurl")
	} else if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		// Handle application/x-www-form-urlencoded
		err := r.ParseForm() // Parse the form data
		if err != nil {
			utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid form data payload")
			return
		}

		// Extract form values
		updateVideoSnapDTO.GameNo = r.FormValue("period")     // Adjust field names accordingly
		updateVideoSnapDTO.Streamkey = r.FormValue("rtmpurl") // Adjust field names accordingly

		//fmt.Println("Received Form Data: Period =", createVideoSnapDTO.Period, ", RTMP URL =", createVideoSnapDTO.Rtmpurl) // Log form data
	} else {
		utils.HttpSuccessResponse(w, nil, http.StatusUnsupportedMediaType, "Unsupported Content-Type")
		return
	}
	// Validate the required parameters
	if updateVideoSnapDTO.GameNo == "" || updateVideoSnapDTO.Streamkey == "" {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "period and rtmpurl are required")
		return
	}
	// Marshal DTO to JSON for service handling
	data, err := json.Marshal(updateVideoSnapDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusInternalServerError, "Failed to process data")
		return
	}

	serviceResponse, err := vc.videoSnapService.VideosnapServiceUpdate(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}
	responseTask := video.VideoSnapTaskResponseDTO{
		ImageURL: serviceResponse.ImageURL,
		VideoURL: serviceResponse.VideoURL,
		Status:   http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseTask)
	//utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// VideoSnapDeleteHandler handles requests to delete a video snapshot
// @Summary Delete a video snapshot
// @Description Delete a video snapshot
// @Tags video snapshots
// @Accept  json
// @Produce  json
// @Param id query int true "Video Snapshot ID"
// @Success 200 {string} string "Video snapshot deleted successfully"
// @Failure 400 {string} string "ID parameter is required"
// @Failure 500 {string} string "Failed to encode JSON"
// @Router /videosnaps/delete [delete]
func (vc *VideoSnapController) VideoSnapDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	videoSnapID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid video snapshot ID", http.StatusBadRequest)
		return
	}

	videoSnapDTO := &video.VideosnapResponseDTO{ID: uint(videoSnapID)}
	data, err := json.Marshal(videoSnapDTO)
	if err != nil {
		http.Error(w, "Failed to marshal request data", http.StatusInternalServerError)
		return
	}

	deleted, err := vc.videoSnapService.VideosnapServiceDelete(data)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, deleted, http.StatusOK, string(utils.SuccessMessage))
}
