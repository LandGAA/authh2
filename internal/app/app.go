package app

import (
	"github.com/LandGAA/authh2/internal/delivery"
	"github.com/LandGAA/authh2/internal/repository"
	"github.com/LandGAA/authh2/internal/usecase"
	"github.com/LandGAA/authh2/pkg/database"
)

var GlobalUseCase usecase.UseCase

func Run() {
	db := database.Connect()
	rep := repository.NewRep(db)
	u := usecase.NewUserUseCase(&rep)
	GlobalUseCase = u
	router := delivery.SetupRouters(u)
	err := router.Run(":8081")
	defer db.Close()
	if err != nil {
		return
	}
}
