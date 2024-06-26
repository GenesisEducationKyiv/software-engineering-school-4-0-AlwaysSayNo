package user

import (
	"net/http"

	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Saver interface {
	Save(user dto.UserSaveRequestDTO) (*dto.UserResponseDTO, error)
}

type Handler struct {
	saver Saver
}

func NewHandler(saver Saver) *Handler {
	return &Handler{
		saver: saver,
	}
}

func (h *Handler) Add(ctx *gin.Context) {
	email := ctx.PostForm("email")

	var saveDto dto.UserSaveRequestDTO
	saveDto.Email = email

	_, err := h.saver.Save(saveDto)
	if err != nil {
		errors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "E-mail додано")
}

// RegisterRoutes creates an instance of Handler and registers routes for it.
func RegisterRoutes(r *gin.Engine, handler Handler) {
	routes := r.Group("/api/v1/")
	routes.POST("/subscribe", handler.Add)
}
