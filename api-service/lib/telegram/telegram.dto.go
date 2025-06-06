package telegram

import "time"

// CreateUserRequestDTO represents the Data Transfer Object for creating a UserRequest.
type CreateUserRequestDTO struct {
	ChatID      string  `json:"chat_id"`      // Chat identifier
	AccountName string  `json:"account_name"` // Name of the account
	RequestType string  `json:"request_type"` // Type of request
	StartTime   string  `json:"start_time"`   // Start time of request
	EndTime     string  `json:"end_time"`     // End time of request
	Status      string  `json:"status"`       // Status of request
	AllowTime   float64 `json:"allow_time"`   // Allowed time
	TotalTime   float64 `json:"total_time"`   // Total time consumed
	Message     string  `json:"message"`      // Additional message
	BotName     string  `json:"bot_name"`
}

// UpdateUserRequestDTO represents the Data Transfer Object for updating a UserRequest.
type UpdateUserRequestDTO struct {
	ID          uint    `json:"id"`           // ID of the user request to update
	ChatID      string  `json:"chat_id"`      // Chat identifier
	AccountName string  `json:"account_name"` // Name of the account
	RequestType string  `json:"request_type"` // Type of request
	StartTime   string  `json:"start_time"`   // Start time of request
	EndTime     string  `json:"end_time"`     // End time of request
	TotalTime   float64 `json:"total_time"`   // Total time consumed
	Status      string  `json:"status"`       // Status of request
	AllowTime   float64 `json:"allow_time"`   // Allowed time
	Message     string  `json:"message"`      // Additional message
	BotName     string  `json:"bot_name"`
}

// ResponseUserRequestDTO represents the Data Transfer Object for returning a UserRequest.
type ResponseUserRequestDTO struct {
	ID          uint      `json:"id"`
	ChatID      string    `json:"chat_id"`      // Chat identifier
	AccountName string    `json:"account_name"` // Name of the account
	RequestType string    `json:"request_type"` // Type of request
	StartTime   string    `json:"start_time"`   // Start time of request
	EndTime     string    `json:"end_time"`     // End time of request
	TotalTime   float64   `json:"total_time"`   // Total time consumed
	Status      string    `json:"status"`       // Status of request
	AllowTime   float64   `json:"allow_time"`   // Allowed time
	Message     string    `json:"message"`      // Additional message
	BotName     string    `json:"bot_name"`
	CreatedAt   time.Time `json:"created_at"` // Timestamp when the request was created
	UpdatedAt   time.Time `json:"updated_at"` // Timestamp when the request was last updated
}

type UserRequestFilter struct {
	ID          uint   `json:"id"`
	AccountName string `json:"account_name" binding:"required"`
	ChatID      string `json:"chat_id"`
	Status      string `json:"status"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

// CreateRequestSettingDTO represents the Data Transfer Object for creating a RequestSetting.
type CreateRequestSettingDTO struct {
	Name       string  `json:"name"`        // Name of the request setting
	TimeRest   float64 `json:"time_rest"`   // Time rest value
	ButtonType string  `json:"button_type"` // Type of button
	BotName    string  `json:"bot_name"`    // Type of button
	Status     string  `json:"status"`      // Status of the setting
	OrderNo    int     `json:"order_no"`    // Order number
}

// UpdateRequestSettingDTO represents the Data Transfer Object for updating a RequestSetting.
type UpdateRequestSettingDTO struct {
	ID         uint    `json:"id"`          // ID of the request setting to update
	Name       string  `json:"name"`        // Name of the request setting
	TimeRest   float64 `json:"time_rest"`   // Time rest value
	ButtonType string  `json:"button_type"` // Type of button
	BotName    string  `json:"bot_name"`
	Status     string  `json:"status"`   // Status of the setting
	OrderNo    int     `json:"order_no"` // Order number
}

// ResponseRequestSettingDTO represents the Data Transfer Object for returning a RequestSetting.
type ResponseRequestSettingDTO struct {
	ID         uint      `json:"id"`          // ID of the request setting
	Name       string    `json:"name"`        // Name of the request setting
	TimeRest   float64   `json:"time_rest"`   // Time rest value
	ButtonType string    `json:"button_type"` // Type of button
	BotName    string    `json:"bot_name"`
	Status     string    `json:"status"`     // Status of the setting
	OrderNo    int       `json:"order_no"`   // Order number
	CreatedAt  time.Time `json:"created_at"` // Creation timestamp
	UpdatedAt  time.Time `json:"updated_at"` // Last updated timestamp
}

// CreateLocationSettingDTO represents the Data Transfer Object for creating a LocationSetting.
type CreateLocationSettingDTO struct {
	Name     string  `json:"name"`      // Name of the location setting
	Lat      string  `json:"lat"`       // Latitude of the location
	Long     string  `json:"long"`      // Longitude of the location
	BotToken string  `json:"token_bot"` // Bot token for the location
	Allow    float64 `json:"allow"`     // Allowed distance for the location
	Status   string  `json:"status"`
}

// UpdateLocationSettingDTO represents the Data Transfer Object for updating a LocationSetting.
type UpdateLocationSettingDTO struct {
	ID       uint    `json:"id"`        // ID of the location setting to update
	Name     string  `json:"name"`      // Name of the location setting
	Lat      string  `json:"lat"`       // Latitude of the location
	Long     string  `json:"long"`      // Longitude of the location
	BotToken string  `json:"token_bot"` // Bot token for the location
	Allow    float64 `json:"allow"`     // Allowed distance for the location
	Status   string  `json:"status"`
}

// ResponseLocationSettingDTO represents the Data Transfer Object for returning a LocationSetting.
type ResponseLocationSettingDTO struct {
	ID        uint      `json:"id"`        // ID of the location setting
	Name      string    `json:"name"`      // Name of the location setting
	Lat       string    `json:"lat"`       // Latitude of the location
	Long      string    `json:"long"`      // Longitude of the location
	BotToken  string    `json:"token_bot"` // Bot token for the location
	Allow     float64   `json:"allow"`     // Allowed distance for the location
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"` // Creation timestamp
	UpdatedAt time.Time `json:"updated_at"` // Last updated timestamp
}

// CreateClockTimeDTO represents the Data Transfer Object for creating a ClockTime.
type CreateClockTimeDTO struct {
	ChatID      string `json:"chat_id"`      // Chat ID of the user
	AccountName string `json:"account_name"` // Account name of the user
	RequestType string `json:"request_type"` // Type of the request (e.g., "clock-in", "clock-out")
	StartTime   string `json:"start_time"`   // Start time of the clock
	Message     string `json:"message"`      // Message related to the clock request
	Lat         string `json:"lat"`
	Long        string `json:"long"`
	BotName     string `json:"bot_name"`
}

// UpdateClockTimeDTO represents the Data Transfer Object for updating a ClockTime.
type UpdateClockTimeDTO struct {
	ID          uint   `json:"id"`           // ID of the ClockTime to update
	ChatID      string `json:"chat_id"`      // Chat ID of the user
	AccountName string `json:"account_name"` // Account name of the user
	RequestType string `json:"request_type"` // Type of the request (e.g., "clock-in", "clock-out")
	StartTime   string `json:"start_time"`   // Start time of the clock
	Message     string `json:"message"`      // Message related to the clock request
	Lat         string `json:"lat"`
	Long        string `json:"long"`
	BotName     string `json:"bot_name"`
}

// ResponseClockTimeDTO represents the Data Transfer Object for returning a ClockTime.
type ResponseClockTimeDTO struct {
	ID          uint      `json:"id"`           // ID of the clock time record
	ChatID      string    `json:"chat_id"`      // Chat ID of the user
	AccountName string    `json:"account_name"` // Account name of the user
	RequestType string    `json:"request_type"` // Type of the request (e.g., "clock-in", "clock-out")
	StartTime   string    `json:"start_time"`   // Start time of the clock
	Message     string    `json:"message"`      // Message related to the clock request
	Lat         string    `json:"lat"`
	Long        string    `json:"long"`
	BotName     string    `json:"bot_name"`
	CreatedAt   time.Time `json:"created_at"` // Timestamp of when the record was created
	UpdatedAt   time.Time `json:"updated_at"` // Timestamp of when the record was last updated
}

// CreatePhoneListDTO represents the Data Transfer Object for creating a PhoneLists entry.
type CreatePhoneListDTO struct {
	Phone   string `json:"phone"`    // Phone number
	BotName string `json:"bot_name"` // Name of the bot associated with the phone
	Status  string `json:"status"`   // Status of the phone entry (default: "yes")
}

// UpdatePhoneListDTO represents the Data Transfer Object for updating a PhoneLists entry.
type UpdatePhoneListDTO struct {
	ID      uint   `json:"id"`       // ID of the PhoneLists entry to update
	Phone   string `json:"phone"`    // Phone number
	BotName string `json:"bot_name"` // Name of the bot associated with the phone
	Status  string `json:"status"`   // Status of the phone entry
}

// ResponsePhoneListDTO represents the Data Transfer Object for returning a PhoneLists entry.
type ResponsePhoneListDTO struct {
	ID          uint      `json:"id"`       // ID of the phone list record
	Phone       string    `json:"phone"`    // Phone number
	BotName     string    `json:"bot_name"` // Name of the bot associated with the phone
	Status      string    `json:"status"`   // Status of the phone entry
	Requester   string    `json:"requester"`
	RequestDate string    `json:"request_date"`
	CreatedAt   time.Time `json:"created_at"` // Timestamp of when the record was created
	UpdatedAt   time.Time `json:"updated_at"` // Timestamp of when the record was last updated
}
type BulkDeleteRequestDTO struct {
	ID []uint `json:"id"` // List of IDs to delete
}
