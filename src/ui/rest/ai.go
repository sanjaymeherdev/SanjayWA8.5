package rest

import (
	domainAI "sanjaywa.com/wa/domains/ai"
	"sanjaywa.com/wa/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type AI struct {
	Service domainAI.IAIUsecase
}

func InitRestAI(app fiber.Router, service domainAI.IAIUsecase) AI {
	handler := AI{Service: service}
	app.Post("/api/ai/chat", handler.Chat)
	return handler
}

func (handler *AI) Chat(c *fiber.Ctx) error {
	var request domainAI.ChatRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseData{
			Code:    "BAD_REQUEST",
			Message: "invalid request body: " + err.Error(),
		})
	}

	if request.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseData{
			Code:    "BAD_REQUEST",
			Message: "message is required",
		})
	}

	response, err := handler.Service.Chat(c.UserContext(), request)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Status:  200,
		Code:    "SUCCESS",
		Message: "AI response generated successfully",
		Results: response,
	})
}
