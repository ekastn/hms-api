package handlers

import (
	"log"

	"github.com/ekastn/hms-api/internal/service"
	"github.com/ekastn/hms-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type ActivityHandler struct {
	activityService *service.ActivityService
}

func NewActivityHandler(activityService *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService}
}

// HandleGetAllActivities handles the request to get all activities.
//
//	@Summary		Get all activities
//	@Description	Get a list of all recorded activities. Admin access required.
//	@Tags			Activities
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	utils.SuccessResponse{data=[]domain.Activity}	"Activities retrieved successfully"
//	@Failure		500	{object}	utils.ErrorResponse								"Failed to retrieve activities"
//	@Router			/activities [get]
func (h *ActivityHandler) HandleGetAllActivities(c *fiber.Ctx) error {
	activities, err := h.activityService.GetAllActivities(c.Context())
	if err != nil {
		log.Printf("Error getting activities: %v", err)
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to retrieve activities", err.Error())
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "Activities retrieved successfully", activities)
}
