package dto

import (
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/broker"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/model"
)

func UserSubscribedDataToUserSaveDTO(data *broker.UserSubscribedData) UserSaveDTO {
	return UserSaveDTO{
		Email: data.Email,
	}
}

func UserSaveDTOToUser(data *UserSaveDTO) model.User {
	return model.User{
		Email: data.Email,
	}
}

func UserToUserResponseDTO(model *model.User) UserResponseDTO {
	return UserResponseDTO{
		ID:           model.ID,
		Email:        model.Email,
		IsSubscribed: model.IsSubscribed,
	}
}
