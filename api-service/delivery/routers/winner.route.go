package router

import (
	controllers "api-service/delivery/http"
	services "api-service/lib/lucky/service" // Adjust the import path as necessary
	"net/http"

	"github.com/gorilla/mux"
)

// WinnerControllerWrapper is a wrapper for the Winner controller
type WinnerControllerWrapper struct {
	controller *controllers.WinnerController
}

// NewWinnerControllerWrapper initializes the wrapper for the Winner controller
func NewWinnerControllerWrapper(cc *controllers.WinnerController) *WinnerControllerWrapper {
	return &WinnerControllerWrapper{
		controller: cc,
	}
}

// SetupWinnerRouter initializes and configures the router for Winner-related routes
func SetupWinnerRouter(r *mux.Router, service services.WinnerServiceInterface) {
	// Initialize controller with WinnerService
	controller := controllers.NewWinnerController(service)

	// Create a new instance of the wrapper for the controller
	controllerWrapper := NewWinnerControllerWrapper(controller)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	Winners := api.PathPrefix("/lucky-winner").Subrouter()
	Winners.Handle("/task-excute-winner", (http.HandlerFunc(controllerWrapper.controller.WinnerTaskExcuteHandler))).Methods("GET")
	Winners.Handle("/task-excute-reset-point", (http.HandlerFunc(controllerWrapper.controller.TaskExcuteResetPointHandler))).Methods("GET")
	Winners.Handle("/winner-send-notif", (http.HandlerFunc(controllerWrapper.controller.TaskSendNotificationWinnerHandler))).Methods("POST")

	Winners.Handle("/collection-point", (http.HandlerFunc(controllerWrapper.controller.WinnerGetByChatIDHandler))).Methods("GET")
	Winners.Handle("/list", (http.HandlerFunc(controllerWrapper.controller.GetWinnerListHandler))).Methods("GET")
	Winners.Handle("/", (http.HandlerFunc(controllerWrapper.controller.WinnerCreateHandler))).Methods("POST")
	Winners.Handle("/{id}", (http.HandlerFunc(controllerWrapper.controller.WinnerGetByIDHandler))).Methods("GET")
	Winners.Handle("/update", (http.HandlerFunc(controllerWrapper.controller.WinnerUpdateHandler))).Methods("PUT")
}
