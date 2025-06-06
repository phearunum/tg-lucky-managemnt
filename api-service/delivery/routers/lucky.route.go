package router

import (
	controllers "api-service/delivery/http"
	services "api-service/lib/lucky/service" // Adjust the import path as necessary
	"api-service/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// TelegramUserControllerWrapper is a wrapper for the Telegram user controller
type TelegramUserControllerWrapper struct {
	controller *controllers.TelegramUserController
}

// NewTelegramUserControllerWrapper initializes the wrapper for the Telegram user controller
func NewTelegramUserControllerWrapper(tuc *controllers.TelegramUserController) *TelegramUserControllerWrapper {
	return &TelegramUserControllerWrapper{
		controller: tuc,
	}
}

// SetupTelegramUserRouter initializes and configures the router for Telegram user-related routes
func SetupTelegramUserRouter(r *mux.Router, service services.TelegramUserServiceInterface) {
	// Initialize controller with TelegramUserService
	controller := controllers.NewTelegramUserController(service)

	// Create a new instance of the wrapper for the controller
	controllerWrapper := NewTelegramUserControllerWrapper(controller)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()
	telegramUsers := api.PathPrefix("/telegram-lucky").Subrouter()
	telegramUsers.Handle("/member", (http.HandlerFunc(controllerWrapper.controller.TelegramUserGetByChatIDHandler))).Methods("GET")
	telegramUsers.Handle("/list", (http.HandlerFunc(controllerWrapper.controller.GetTelegramUserListHandler))).Methods("GET")
	telegramUsers.Handle("/", utils.HasPermission("telegram:user:add")(http.HandlerFunc(controllerWrapper.controller.TelegramUserCreateHandler))).Methods("POST")
	telegramUsers.Handle("/{id}", utils.HasPermission("telegram:user:view")(http.HandlerFunc(controllerWrapper.controller.TelegramUserGetByIDHandler))).Methods("GET")
	telegramUsers.Handle("/update", utils.HasPermission("telegram:user:update")(http.HandlerFunc(controllerWrapper.controller.TelegramUserUpdateHandler))).Methods("PUT")
}
