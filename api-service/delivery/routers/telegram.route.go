package router

import (
	controllers "api-service/delivery/http"
	services "api-service/lib/telegram" // Adjust the import path as necessary

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// UserRequestControllerWrapper is a wrapper for the user request controller
type UserRequestControllerWrapper struct {
	controller *controllers.UserRequestController
}

// NewUserRequestControllerWrapper initializes the wrapper for the user request controller
func NewUserRequestControllerWrapper(urc *controllers.UserRequestController) *UserRequestControllerWrapper {
	return &UserRequestControllerWrapper{
		controller: urc,
	}
}

// SetupUserRequestRouter initializes and configures the router for user request-related routes
func SetupUserRequestRouter(r *mux.Router, service services.UserRequestServiceInterface) {
	// Initialize controller with UserRequestService
	controller := controllers.NewUserRequestController(service)

	// Create a new instance of the wrapper for the controller
	controllerWrapper := NewUserRequestControllerWrapper(controller)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// User Request routes

	userRequests := api.PathPrefix("/user-requests").Subrouter()
	userRequests.HandleFunc("/time-records", controllerWrapper.controller.GetClockTimeListHandler).Methods("GET")
	userRequests.HandleFunc("/time-records", controllerWrapper.controller.RequestSaveClockInHandler).Methods("POST")
	userRequests.HandleFunc("/bot-phone", controllerWrapper.controller.GetPhoneListHandler).Methods("GET")
	userRequests.HandleFunc("/bot-phone", controllerWrapper.controller.RequestSavePhoneHandler).Methods("POST")
	userRequests.HandleFunc("/bot-phone/delete", controllerWrapper.controller.RequestBulkDeletePhoneHandler).Methods("POST")
	userRequests.HandleFunc("/bot-phone", controllerWrapper.controller.RequestUpdatePhoneHandler).Methods("PUT")

	userRequests.HandleFunc("/bot-location", controllerWrapper.controller.GetBotLocationListHandler).Methods("GET")
	userRequests.HandleFunc("/bot-location", controllerWrapper.controller.RequestBotLocationCreateHandler).Methods("POST")
	userRequests.HandleFunc("/bot-location", controllerWrapper.controller.RequestBotLocationUpdateHandler).Methods("PUT")

	userRequests.HandleFunc("/setting", controllerWrapper.controller.RequestSettingListHandler).Methods("GET")
	userRequests.HandleFunc("/setting", controllerWrapper.controller.RequestSettingCreateHandler).Methods("POST")
	userRequests.HandleFunc("/setting", controllerWrapper.controller.RequestSettingupdareHandler).Methods("PUT")

	userRequests.HandleFunc("/list", controllerWrapper.controller.UserRequestListHandler).Methods("GET")
	userRequests.HandleFunc("/", controllerWrapper.controller.UserRequestCreateHandler).Methods("POST")
	userRequests.HandleFunc("/{id}", controllerWrapper.controller.UserRequestByIDHandler).Methods("GET")
	userRequests.HandleFunc("/update", controllerWrapper.controller.UserRequestUpdateHandler).Methods("PUT")
	userRequests.HandleFunc("/delete/{id}", controllerWrapper.controller.UserRequestDeleteHandler).Methods("DELETE")

}
