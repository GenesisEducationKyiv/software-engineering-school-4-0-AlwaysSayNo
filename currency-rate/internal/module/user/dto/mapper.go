package dto

import (
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/model"
	shareduser "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/user"
)

func ToDTO(entity model.User) shareduser.ResponseDTO {
	return shareduser.ResponseDTO{
		ID:    entity.ID,
		Email: entity.Email,
	}
}

func SaveRequestToModel(dto SaveRequestDTO) model.User {
	return model.User{
		Email: dto.Email,
	}
}
