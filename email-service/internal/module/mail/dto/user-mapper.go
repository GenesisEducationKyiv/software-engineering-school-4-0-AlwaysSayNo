package dto

import (
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/model"
)

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
