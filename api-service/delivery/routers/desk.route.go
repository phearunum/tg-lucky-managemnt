package router

import (
	controllers "api-service/delivery/http"
	services "api-service/lib/desks" // Adjust the import path as necessary

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// DeskSettingControllerWrapper is a wrapper for the desk setting controller
type DeskSettingControllerWrapper struct {
	controller *controllers.DeskSettingController
}

// NewDeskSettingControllerWrapper initializes the wrapper for the desk setting controller
func NewDeskSettingControllerWrapper(dsc *controllers.DeskSettingController) *DeskSettingControllerWrapper {
	return &DeskSettingControllerWrapper{
		controller: dsc,
	}
}

// SetupDeskSettingRouter initializes and configures the router for desk setting-related routes
func SetupDeskSettingRouter(r *mux.Router, service services.DeskSettingServiceInterface) {
	// Initialize controller with DeskSettingService
	controller := controllers.NewDeskSettingController(service)

	// Create a new instance of the wrapper for the controller
	controllerWrapper := NewDeskSettingControllerWrapper(controller)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Desk Setting routes
	deskSettings := api.PathPrefix("/desk-settings").Subrouter()
	deskSettings.HandleFunc("/list", controllerWrapper.controller.DeskSettingListHandler).Methods("GET")
	deskSettings.HandleFunc("/", controllerWrapper.controller.DeskSettingCreateHandler).Methods("POST")
	deskSettings.HandleFunc("/{id}", controllerWrapper.controller.DeskSettingByIDHandler).Methods("GET")
	deskSettings.HandleFunc("/update", controllerWrapper.controller.DeskSettingUpdateHandler).Methods("PUT")
	deskSettings.HandleFunc("/delete/{id}", controllerWrapper.controller.DeskSettingDeleteHandler).Methods("DELETE")
}
