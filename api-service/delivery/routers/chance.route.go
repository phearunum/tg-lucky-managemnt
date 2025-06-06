package router

import (
	controllers "api-service/delivery/http"
	services "api-service/lib/lucky/service" // Adjust the import path as necessary
	"api-service/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// ChancePointControllerWrapper is a wrapper for the ChancePoint controller
type ChancePointControllerWrapper struct {
	controller *controllers.ChancePointController
}

// NewChancePointControllerWrapper initializes the wrapper for the ChancePoint controller
func NewChancePointControllerWrapper(cc *controllers.ChancePointController) *ChancePointControllerWrapper {
	return &ChancePointControllerWrapper{
		controller: cc,
	}
}

// SetupChancePointRouter initializes and configures the router for ChancePoint-related routes
func SetupChancePointRouter(r *mux.Router, service services.ChancePointServiceInterface) {
	// Initialize controller with ChancePointService
	controller := controllers.NewChancePointController(service)

	// Create a new instance of the wrapper for the controller
	controllerWrapper := NewChancePointControllerWrapper(controller)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	chancePoints := api.PathPrefix("/chance-points").Subrouter()
	chancePoints.Handle("/collection-point", (http.HandlerFunc(controllerWrapper.controller.ChancePointGetByChatIDHandler))).Methods("GET")
	chancePoints.Handle("/list", (http.HandlerFunc(controllerWrapper.controller.GetChancePointListHandler))).Methods("GET")
	chancePoints.Handle("/", utils.HasPermission("chance:point:add")(http.HandlerFunc(controllerWrapper.controller.ChancePointCreateHandler))).Methods("POST")
	chancePoints.Handle("/{id}", utils.HasPermission("chance:point:view")(http.HandlerFunc(controllerWrapper.controller.ChancePointGetByIDHandler))).Methods("GET")
	chancePoints.Handle("/update", (http.HandlerFunc(controllerWrapper.controller.ChancePointUpdateHandler))).Methods("PUT")
}
