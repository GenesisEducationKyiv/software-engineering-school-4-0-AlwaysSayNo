package dto

import "genesis-currency-api/internal/module/user/model"
import shareduser "genesis-currency-api/internal/shared/dto/user"

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
