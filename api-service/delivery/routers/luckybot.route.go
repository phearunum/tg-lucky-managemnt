package router

import (
	controllers "api-service/delivery/http"
	services "api-service/lib/lucky/service"
	"net/http"

	"github.com/gorilla/mux"
)

// TelegramSettingWinnerControllerWrapper wraps the TelegramSettingWinnerController
type TelegramSettingWinnerControllerWrapper struct {
	controller *controllers.TelegramSettingWinnerController
}

// NewTelegramSettingWinnerControllerWrapper initializes the wrapper for TelegramSettingWinnerController
func NewTelegramSettingWinnerControllerWrapper(tswc *controllers.TelegramSettingWinnerController) *TelegramSettingWinnerControllerWrapper {
	return &TelegramSettingWinnerControllerWrapper{
		controller: tswc,
	}
}

// SetupTelegramSettingWinnerRouter configures the routes for Telegram setting winner endpoints
func SetupTelegramSettingWinnerRouter(r *mux.Router, service services.TelegramSettingWinnerServiceInterface) {
	controller := controllers.NewTelegramSettingWinnerController(service)
	controllerWrapper := NewTelegramSettingWinnerControllerWrapper(controller)

	api := r.PathPrefix("/api").Subrouter()
	telegramSetting := api.PathPrefix("/lucky-setting-winner").Subrouter()

	telegramSetting.Handle("/list", (http.HandlerFunc(controllerWrapper.controller.GetTelegramSettingWinnerListHandler))).Methods("GET")
	telegramSetting.Handle("/", (http.HandlerFunc(controllerWrapper.controller.TelegramSettingWinnerCreateHandler))).Methods("POST")
	telegramSetting.Handle("/update", (http.HandlerFunc(controllerWrapper.controller.TelegramSettingWinnerUpdateHandler))).Methods("PUT")
	telegramSetting.Handle("/{id}", (http.HandlerFunc(controllerWrapper.controller.TelegramSettingWinnerGetByIDHandler))).Methods("GET")
}
