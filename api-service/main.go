package main

import (
	"api-service/config"
	"api-service/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	// Import the RabbitMQ package
	auth "api-service/auth"
	roleModel "api-service/lib/roles/models"
	roleRepo "api-service/lib/roles/repository"
	roleService "api-service/lib/roles/services"
	userModel "api-service/lib/users/models"
	"api-service/lib/users/repository"
	userService "api-service/lib/users/services"

	//----

	permissionRepo "api-service/lib/permission/repository"
	permissionService "api-service/lib/permission/service"

	// ------
	employeeRepo "api-service/lib/employee"
	employeeService "api-service/lib/employee"

	// ------
	department "api-service/lib/department"

	//--------
	menusService "api-service/lib/menus"
	menusRepo "api-service/lib/menus/repository"

	// telegram

	telegrambot "api-service/lib/telegram/bot"
	telegramModel "api-service/lib/telegram/model"
	telegramRepo "api-service/lib/telegram/repository"

	telegram "api-service/lib/telegram"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"

	"api-service/delivery/middleware"
	router "api-service/delivery/routers"

	//luckyModel "api-service/lib/lucky/models"
	luckyRepo "api-service/lib/lucky/repository"
	luckyService "api-service/lib/lucky/service"

	"github.com/gorilla/mux"
)

var whitelist = map[string]bool{

	"/api/telegram-lucky/update":                true,
	"/api/telegram-lucky/member":                true,
	"/api/lucky-winner/task-excute-reset-point": true,
	"/api/lucky-winner/task-excute-winner":      true,
	"/api/chance-points/collection-point":       true,
	//"/api/chance-points/list":                   true,
	"/api/chance-points/update":                true,
	"/api/user-requests/time-records":          true,
	"/api/user-requests/bot-location":          true,
	"/api/auth/login":                          true,
	"/api/auth/logout":                         true,
	"/api/videosnap/":                          true,
	"/api/videostop/":                          true,
	"/api/videosnap":                           true,
	"/api/videostop":                           true,
	"/api/video/videosnap":                     true,
	"/api/video/videostop":                     true,
	"/api/v1/items/list":                       true,
	"/api/swagger/index.html":                  true,
	"/swagger/index.html":                      true,
	"/swagger/swagger-ui-bundle.js":            true,
	"/swagger/swagger-ui.css":                  true,
	"/swagger/swagger-ui-standalone-preset.js": true,
	"/swagger/doc.json":                        true,
	"/swagger/favicon-32x32.png":               true,
	"/swagger/favicon-16x16.png":               true,
}

func main() {

	// Initialize configuration and database
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	cfg, db, err := initConfigAndDatabase("config/config.yml")
	if err != nil {
		logger.Errorf("Initialization error: %v", err)
	}

	err = config.InitRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		utils.ErrorLog(err, "Failed to initialize Redis")
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	// Container Timezoe
	loc, errTime := time.LoadLocation(cfg.Service.TimeZone)
	if errTime != nil {
		fmt.Println("Error loading location:", errTime)
		return
	}
	now := time.Now().In(loc)
	fmt.Println("Current Time in", cfg.Service.TimeZone, ":", now)
	utcNow := time.Now().UTC()
	fmt.Println("Current Time in UTC:", utcNow)
	fmt.Println("Current Time Zone:", loc.String(), "Offset:", now.Format("-07:00"))

	// Initialize Telegram repository
	telegramRepo := telegramRepo.NewTelegramAccountRepository(db)
	//telegramRepo := repository.NewTelegramAccountRepository(db)
	telegramModel.MigrateTelegramAccount(db)
	telegramModel.MigrateTelegramMessage(db)
	//config.InitTelegram(cfg.Telegram, telegramRepo)
	for _, botConfig := range cfg.Telegram.Bots {
		if !botConfig.Status {
			log.Printf("Skipping bot %s (status: false)", botConfig.Name)
			continue
		}

		log.Printf("Initializing Telegram bot: %s", botConfig.Name)
		// Initialize and start the bot
		bot := telegrambot.NewBot(botConfig.Name, botConfig.Token, botConfig.Debug, botConfig.GroupID, botConfig.BackBtn, telegramRepo)
		go bot.Start()

		log.Printf("Bot %s started successfully", botConfig.Name)
	}
	utils.InfoLog(cfg.Telegram, " Init Telegram")
	utils.InitializeLogger(cfg.Service.LogPtah)

	// HTTP handler route
	r := mux.NewRouter()
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3300", "http://localhost:3000", "http://192.168.50.102:8080", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := corsHandler.Handler(r)

	userRepository := auth.NewUserRepository(db)
	authService := auth.NewAuthService(userRepository)

	userRepo := repository.NewUserRepository(db)
	userService := userService.NewUserService(userRepo)
	userModel.MigrateUsers(db)

	roleRepo_ := roleRepo.NewRoleRepository(db)
	roleService := roleService.NewRoleService(roleRepo_)
	roleModel.MigrateRoles(db)

	permissionRepo_ := permissionRepo.NewPermissionRepository(db)
	//permissionModel.MigrateRoleAccesses(db)
	permissionService := permissionService.NewPermissionService(permissionRepo_)

	employeeRepo_ := employeeRepo.NewEmployeeRepository(db)
	employeeService_ := employeeService.NewEmployeeService(employeeRepo_)

	departmentRepo_ := department.NewDepartmentRepository(db)
	deprtmentService := department.NewDepartmentService(departmentRepo_)

	menusRepo_ := menusRepo.NewMenuRepository(db)
	//menusmodel.MigrateMenu(db)
	menusService_ := menusService.NewMenuService(menusRepo_)
	tgRepo := telegram.NewUserRequestRepository(db)
	tsService := telegram.NewUserRequestService(tgRepo)

	luckyRepos := luckyRepo.NewTelegramUserRepository(db)
	luckService := luckyService.NewTelegramUserService(luckyRepos)
	//luckyModel.MigrateDB(db)

	chanceRepo := luckyRepo.NewChancePointRepository(db)
	chanceService := luckyService.NewChancePointService(chanceRepo)
	winnerRepo := luckyRepo.NewWinnerRepository(db)
	winnerService := luckyService.NewWinnerService(winnerRepo)
	luckySettingRepo := luckyRepo.NewTelegramSettingWinnerRepository(db)
	luckyService := luckyService.NewTelegramSettingWinnerService(luckySettingRepo)

	r.Use(middleware.LoggingMiddleware)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.Use(utils.AuthMiddlewareWithWhiteList(whitelist))

	// Enable HTTP Handler Direct
	router.SetupAuthRouter(r, authService)
	router.SetupRouter(r, userService)
	router.SetupRoleRouter(r, roleService)
	router.SetupPermissionRouter(r, permissionService)
	router.SetupMenuRouter(r, menusService_)
	router.SetupEmployeeRouter(r, employeeService_)
	router.SetupDepartmentRouter(r, deprtmentService)
	router.SetupUserRequestRouter(r, tsService)
	router.SetupTelegramUserRouter(r, luckService)
	router.SetupChancePointRouter(r, chanceService)
	router.SetupWinnerRouter(r, winnerService)
	router.SetupTelegramSettingWinnerRouter(r, luckyService)

	log.Println("      \n                   *       \n.         *        *        *\n.        ***     **=**     ***\n        *\"\"\"*   *|***|*   *\"\"\"*\n       *|***|*  *|*+*|*  *|***|*\n**********\"\"\"*___*//+\\\\*___*\"\"\"*********\n@@@@@@@@@@@@@@@@//   \\\\@@@@@@@@@@@@@@@@@\n###############||ព្រះពុទ្ធ||#################\nTTTTTTTTTTTTTTT||ព្រះធម័||TTTTTTTTTTTTTTTTT\n------------- -//ព្រះសង្ឃ\\\\----------------\n៚ សូមប្រោសប្រទានពរឱ្យប្រតិប័ត្តិការណ៍ប្រព្រឹត្តទៅជាធម្មតា ៚ \n ៚ ជោគជ័យ      ៚សិរីសួរស្តី       ៚សុវត្តិភាព \n_________________________________________\n៚  Application Service is Running Port: " + cfg.Service.Port)

	log.Fatal(http.ListenAndServe(":"+cfg.Service.Port, handler))
}

func initConfigAndDatabase(configPath string) (config.Config, *gorm.DB, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return cfg, nil, err
	}

	db := config.InitDatabase(configPath)
	return cfg, db, nil
}
