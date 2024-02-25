package routes

import (
	"example.com/rest-apis/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
		return
	}

	err = event.Register(userId)
	if err != nil {
		fmt.Println("Error was ", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Cloud not register the event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Registered successfully"})
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}
	var event models.Event
	event.ID = eventId
	err = event.CancelRegistation(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": " Cloud not cancel the registrations"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "registration canceled successfully"})
}
