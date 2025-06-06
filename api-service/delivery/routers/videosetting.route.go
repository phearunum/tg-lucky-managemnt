package router

import (
	controllers "api-service/delivery/http"
	services "api-service/lib/videosnap/setting" // Adjust the import path as necessary
	"api-service/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// VideoSnapSettingControllerWrapper is a wrapper for the video snap setting controller
type VideoSnapSettingControllerWrapper struct {
	controller *controllers.VideoSnapSettingController
}

// NewVideoSnapSettingControllerWrapper initializes the wrapper for the video snap setting controller
func NewVideoSnapSettingControllerWrapper(vssc *controllers.VideoSnapSettingController) *VideoSnapSettingControllerWrapper {
	return &VideoSnapSettingControllerWrapper{
		controller: vssc,
	}
}

// SetupVideoSnapSettingRouter initializes and configures the router for video snap setting-related routes
func SetupVideoSnapSettingRouter(r *mux.Router, service services.VideoSnapSettingServiceInterface) {
	// Initialize controller with VideoSnapSettingService
	controller := controllers.NewVideoSnapSettingController(service)

	// Create a new instance of the wrapper for the controller
	controllerWrapper := NewVideoSnapSettingControllerWrapper(controller)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	videoSnapSettings := api.PathPrefix("/videos-settings").Subrouter()
	videoSnapSettings.Handle("/list", utils.HasPermission("setting:video:view")(http.HandlerFunc(controllerWrapper.controller.VideoSnapSettingListHandler))).Methods("GET")
	videoSnapSettings.Handle("/", utils.HasPermission("system:server:add")(http.HandlerFunc(controllerWrapper.controller.VideoSnapSettingCreateHandler))).Methods("POST")
	videoSnapSettings.Handle("/{id}", utils.HasPermission("setting:video:view")(http.HandlerFunc(controllerWrapper.controller.VideoSnapSettingByIDHandler))).Methods("GET")
	videoSnapSettings.Handle("/update", utils.HasPermission("system:server:update")(http.HandlerFunc(controllerWrapper.controller.VideoSnapSettingUpdateHandler))).Methods("PUT")
	videoSnapSettings.Handle("/delete/{id}", utils.HasPermission("system:server:delete")(http.HandlerFunc(controllerWrapper.controller.VideoSnapSettingDeleteHandler))).Methods("DELETE")
	// Video Snap Setting routes
	/*videoSnapSettings := api.PathPrefix("/videos-settings").Subrouter()
	videoSnapSettings.HandleFunc("/list", controllerWrapper.controller.VideoSnapSettingListHandler).Methods("GET")
	videoSnapSettings.HandleFunc("/", controllerWrapper.controller.VideoSnapSettingCreateHandler).Methods("POST")
	videoSnapSettings.HandleFunc("/{id}", controllerWrapper.controller.VideoSnapSettingByIDHandler).Methods("GET")
	videoSnapSettings.HandleFunc("/update", controllerWrapper.controller.VideoSnapSettingUpdateHandler).Methods("PUT")
	videoSnapSettings.HandleFunc("/delete", controllerWrapper.controller.VideoSnapSettingDeleteHandler).Methods("DELETE")
	*/
}
