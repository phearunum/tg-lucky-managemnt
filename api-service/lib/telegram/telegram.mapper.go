package telegram

// ToUserRequestModelFromCreateDTO converts CreateUserRequestDTO to UserRequest model
func ToUserRequestModelFromCreateDTO(createDTO CreateUserRequestDTO) UserRequest {
	return UserRequest{
		ChatID:      createDTO.ChatID,
		AccountName: createDTO.AccountName,
		RequestType: createDTO.RequestType,
		StartTime:   createDTO.StartTime,
		EndTime:     createDTO.EndTime,
		Status:      createDTO.Status,
		AllowTime:   createDTO.AllowTime,
		TotalTime:   createDTO.TotalTime,
		Message:     createDTO.Message,
		BotName:     createDTO.BotName,
	}
}

// ToUserRequestModelFromUpdateDTO converts UpdateUserRequestDTO to UserRequest model
func ToUserRequestModelFromUpdateDTO(updateDTO UpdateUserRequestDTO) UserRequest {
	return UserRequest{
		ID:          updateDTO.ID,
		ChatID:      updateDTO.ChatID,
		AccountName: updateDTO.AccountName,
		RequestType: updateDTO.RequestType,
		StartTime:   updateDTO.StartTime,
		EndTime:     updateDTO.EndTime,
		TotalTime:   updateDTO.TotalTime,
		Status:      updateDTO.Status,
		AllowTime:   updateDTO.AllowTime,
		Message:     updateDTO.Message,
		BotName:     updateDTO.BotName,
	}
}

// ToUserRequestResponseDTO converts UserRequest model to ResponseUserRequestDTO
// ToUserRequestResponseDTO converts UserRequest model to ResponseUserRequestDTO
func ToUserRequestResponseDTO(userRequest UserRequest) *ResponseUserRequestDTO {
	return &ResponseUserRequestDTO{
		ID:          userRequest.ID,
		ChatID:      userRequest.ChatID,
		AccountName: userRequest.AccountName,
		RequestType: userRequest.RequestType,
		StartTime:   userRequest.StartTime,
		EndTime:     userRequest.EndTime,
		TotalTime:   userRequest.TotalTime,
		Status:      userRequest.Status,
		AllowTime:   userRequest.AllowTime,
		Message:     userRequest.Message,
		BotName:     userRequest.BotName,
		CreatedAt:   userRequest.CreatedAt,
		UpdatedAt:   userRequest.UpdatedAt,
	}
}

// ToRequestSettingModelFromCreateDTO converts CreateRequestSettingDTO to RequestSetting model
func ToRequestSettingModelFromCreateDTO(createDTO CreateRequestSettingDTO) RequestSetting {
	return RequestSetting{
		Name:       createDTO.Name,
		TimeRest:   createDTO.TimeRest,
		ButtonType: createDTO.ButtonType,
		BotName:    createDTO.BotName,
		Status:     createDTO.Status,
		OrderNo:    createDTO.OrderNo,
	}
}

// ToRequestSettingModelFromUpdateDTO converts UpdateRequestSettingDTO to RequestSetting model
func ToRequestSettingModelFromUpdateDTO(updateDTO UpdateRequestSettingDTO) RequestSetting {
	return RequestSetting{
		ID:         updateDTO.ID,
		Name:       updateDTO.Name,
		TimeRest:   updateDTO.TimeRest,
		ButtonType: updateDTO.ButtonType,
		BotName:    updateDTO.BotName,
		Status:     updateDTO.Status,
		OrderNo:    updateDTO.OrderNo,
	}
}

// ToRequestSettingResponseDTO converts RequestSetting model to ResponseRequestSettingDTO
func ToRequestSettingResponseDTO(requestSetting RequestSetting) *ResponseRequestSettingDTO {
	return &ResponseRequestSettingDTO{
		ID:         requestSetting.ID,
		Name:       requestSetting.Name,
		TimeRest:   requestSetting.TimeRest,
		ButtonType: requestSetting.ButtonType,
		BotName:    requestSetting.BotName,
		Status:     requestSetting.Status,
		OrderNo:    requestSetting.OrderNo,
		CreatedAt:  requestSetting.CreatedAt,
		UpdatedAt:  requestSetting.UpdatedAt,
	}
}

// ToLocationSettingModelFromCreateDTO converts CreateLocationSettingDTO to LocationSetting model
func ToLocationSettingModelFromCreateDTO(createDTO CreateLocationSettingDTO) LocationSetting {
	return LocationSetting{
		Name:     createDTO.Name,
		Lat:      createDTO.Lat,
		Long:     createDTO.Long,
		BotToken: createDTO.BotToken,
		Allow:    createDTO.Allow,
		Status:   createDTO.Status,
	}
}

// ToLocationSettingModelFromUpdateDTO converts UpdateLocationSettingDTO to LocationSetting model
func ToLocationSettingModelFromUpdateDTO(updateDTO UpdateLocationSettingDTO) LocationSetting {
	return LocationSetting{
		ID:       updateDTO.ID,
		Name:     updateDTO.Name,
		Lat:      updateDTO.Lat,
		Long:     updateDTO.Long,
		BotToken: updateDTO.BotToken,
		Allow:    updateDTO.Allow,
		Status:   updateDTO.Status,
	}
}

// ToLocationSettingResponseDTO converts LocationSetting model to ResponseLocationSettingDTO
func ToLocationSettingResponseDTO(locationSetting LocationSetting) *ResponseLocationSettingDTO {
	return &ResponseLocationSettingDTO{
		ID:        locationSetting.ID,
		Name:      locationSetting.Name,
		Lat:       locationSetting.Lat,
		Long:      locationSetting.Long,
		BotToken:  locationSetting.BotToken,
		Allow:     locationSetting.Allow,
		Status:    locationSetting.Status,
		CreatedAt: locationSetting.CreatedAt,
		UpdatedAt: locationSetting.UpdatedAt,
	}
}

// ToClockTimeModelFromCreateDTO converts CreateClockTimeDTO to ClockTime model
func ToClockTimeModelFromCreateDTO(createDTO CreateClockTimeDTO) ClockTime {
	return ClockTime{
		ChatID:      createDTO.ChatID,
		AccountName: createDTO.AccountName,
		RequestType: createDTO.RequestType,
		StartTime:   createDTO.StartTime,
		Message:     createDTO.Message,
		BotName:     createDTO.BotName,
		Lat:         createDTO.Lat,
		Long:        createDTO.Long,
	}
}

// ToClockTimeModelFromUpdateDTO converts UpdateClockTimeDTO to ClockTime model
func ToClockTimeModelFromUpdateDTO(updateDTO UpdateClockTimeDTO) ClockTime {
	return ClockTime{
		ID:          updateDTO.ID,
		ChatID:      updateDTO.ChatID,
		AccountName: updateDTO.AccountName,
		RequestType: updateDTO.RequestType,
		StartTime:   updateDTO.StartTime,
		Message:     updateDTO.Message,
		BotName:     updateDTO.BotName,
		Lat:         updateDTO.Lat,
		Long:        updateDTO.Long,
	}
}

// ToClockTimeResponseDTO converts ClockTime model to ResponseClockTimeDTO
func ToClockTimeResponseDTO(clockTime ClockTime) *ResponseClockTimeDTO {
	return &ResponseClockTimeDTO{
		ID:          clockTime.ID,
		ChatID:      clockTime.ChatID,
		AccountName: clockTime.AccountName,
		RequestType: clockTime.RequestType,
		StartTime:   clockTime.StartTime,
		Message:     clockTime.Message,
		Lat:         clockTime.Lat,
		Long:        clockTime.Long,
		CreatedAt:   clockTime.CreatedAt,
		UpdatedAt:   clockTime.UpdatedAt,
		BotName:     clockTime.BotName,
	}
}

// ToPhoneListModelFromCreateDTO converts CreatePhoneListDTO to PhoneLists model.
func ToPhoneListModelFromCreateDTO(createDTO CreatePhoneListDTO) PhoneLists {
	return PhoneLists{
		Phone:   createDTO.Phone,
		BotName: createDTO.BotName,
		Status:  createDTO.Status, // Use default value if not set
	}
}

// ToPhoneListModelFromUpdateDTO converts UpdatePhoneListDTO to PhoneLists model.
func ToPhoneListModelFromUpdateDTO(updateDTO UpdatePhoneListDTO) PhoneLists {
	return PhoneLists{
		ID:      updateDTO.ID,
		Phone:   updateDTO.Phone,
		BotName: updateDTO.BotName,
		Status:  updateDTO.Status,
	}
}

// ToPhoneListResponseDTO converts PhoneLists model to ResponsePhoneListDTO.
func ToPhoneListResponseDTO(phoneList PhoneLists) *ResponsePhoneListDTO {
	return &ResponsePhoneListDTO{
		ID:          phoneList.ID,
		Phone:       phoneList.Phone,
		BotName:     phoneList.BotName,
		Status:      phoneList.Status,
		Requester:   phoneList.Requester,
		RequestDate: phoneList.RequestDate,
		CreatedAt:   phoneList.CreatedAt,
		UpdatedAt:   phoneList.UpdatedAt,
	}
}
