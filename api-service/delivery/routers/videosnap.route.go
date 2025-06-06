package router

import (
	controllers "api-service/delivery/http"
	services "api-service/lib/videosnap" // Adjust the import path as necessary

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// VideoSnapControllerWrapper is a wrapper for the video snapshot controller
type VideoSnapControllerWrapper struct {
	controller *controllers.VideoSnapController
}

// NewVideoSnapControllerWrapper initializes the wrapper for the video snapshot controller
func NewVideoSnapControllerWrapper(vsc *controllers.VideoSnapController) *VideoSnapControllerWrapper {
	return &VideoSnapControllerWrapper{
		controller: vsc,
	}
}

// SetupVideoSnapRouter initializes and configures the router for video snapshot-related routes
func SetupVideoSnapRouter(r *mux.Router, service services.VideosnapServiceInterface) {
	// Initialize controller with VideoSnapService
	controller := controllers.NewVideoSnapController(service)

	// Create a new instance of the wrapper for the controller
	controllerWrapper := NewVideoSnapControllerWrapper(controller)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/videosnap", controllerWrapper.controller.VideoSnapCreateHandler).Methods("POST")
	api.HandleFunc("/videostop", controllerWrapper.controller.VideoSnapUpdateHandler).Methods("POST") // ✅ /api/videostop
	api.HandleFunc("/videosnap/", controllerWrapper.controller.VideoSnapCreateHandler).Methods("POST")
	api.HandleFunc("/videostop/", controllerWrapper.controller.VideoSnapUpdateHandler).Methods("POST") // ✅ /api/videostop

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Video Snapshot routes (under /api/video/)
	videoSnap := api.PathPrefix("/video").Subrouter()
	videoSnap.HandleFunc("/list", controllerWrapper.controller.VideoSnapListHandler).Methods("GET")
	videoSnap.HandleFunc("/videosnap", controllerWrapper.controller.VideoSnapCreateHandler).Methods("POST")
	videoSnap.HandleFunc("/videostop", controllerWrapper.controller.VideoSnapUpdateHandler).Methods("POST") // ✅ /api/video/videostop
	videoSnap.HandleFunc("/delete", controllerWrapper.controller.VideoSnapDeleteHandler).Methods("DELETE")
	videoSnap.HandleFunc("/{id}", controllerWrapper.controller.VideoSnapByIDHandler).Methods("GET")
}
