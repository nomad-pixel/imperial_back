package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/domain/usecases/auth_usecases"
	"github.com/nomad-pixel/imperial/internal/domain/usecases/car_usecases"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/auth"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_category"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_mark"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_tag"
)

type App struct {
	DB                              *pgxpool.Pool
	AuthHandler                     *auth.AuthHandler
	CarHandler                      *car.CarHandler
	CarTagHandler                   *car_tag.CarTagHandler
	CarMarkHandler                  *car_mark.CarMarkHandler
	CarCategoryHandler              *car_category.CarCategoryHandler
	SignUpUsecase                   auth_usecases.SignUpUsecase
	SignInUsecase                   auth_usecases.SignInUsecase
	SendEmailVerificationUsecase    auth_usecases.SendEmailVerificationUsecase
	ConfirmEmailVerificationUsecase auth_usecases.ConfirmEmailVerificationUsecase
	TokenService                    ports.TokenService

	CreateCarUsecase car_usecases.CreateCarUsecase

	// CarTag usecases
	CreateCarTagUsecase   car_usecases.CreateCarTagUsecase
	GetCarTagUsecase      car_usecases.GetCarTagUsecase
	GetCarTagsListUsecase car_usecases.GetCarTagsListUsecase
	UpdateCarTagUsecase   car_usecases.UpdateCarTagUsecase
	DeleteCarTagUsecase   car_usecases.DeleteCarTagUsecase

	// CarMark usecases
	CreateCarMarkUsecase   car_usecases.CreateCarMarkUsecase
	GetCarMarkUsecase      car_usecases.GetCarMarkUsecase
	GetCarMarksListUsecase car_usecases.GetCarMarksListUsecase
	UpdateCarMarkUsecase   car_usecases.UpdateCarMarkUsecase
	DeleteCarMarkUsecase   car_usecases.DeleteCarMarkUsecase

	// CarCategory usecases
	CreateCarCategoryUsecase    car_usecases.CreateCarCategoryUsecase
	GetCarCategoryUsecase       car_usecases.GetCarCategoryUsecase
	GetCarCategoriesListUsecase car_usecases.GetCarCategoriesListUsecase
	UpdateCarCategoryUsecase    car_usecases.UpdateCarCategoryUsecase
	DeleteCarCategoryUsecase    car_usecases.DeleteCarCategoryUsecase
}

func NewApp(
	db *pgxpool.Pool,
	signUpUsecase auth_usecases.SignUpUsecase,
	sendEmailVerificationUsecase auth_usecases.SendEmailVerificationUsecase,
	confirmEmailVerificationUsecase auth_usecases.ConfirmEmailVerificationUsecase,
	signInUsecase auth_usecases.SignInUsecase,
	authHandler *auth.AuthHandler,
	tokenSvc ports.TokenService,

	createCarUsecase car_usecases.CreateCarUsecase,
	carHandler *car.CarHandler,
	carTagHandler *car_tag.CarTagHandler,
	carMarkHandler *car_mark.CarMarkHandler,
	carCategoryHandler *car_category.CarCategoryHandler,

	createCarTagUsecase car_usecases.CreateCarTagUsecase,
	getCarTagUsecase car_usecases.GetCarTagUsecase,
	getCarTagsListUsecase car_usecases.GetCarTagsListUsecase,
	updateCarTagUsecase car_usecases.UpdateCarTagUsecase,
	deleteCarTagUsecase car_usecases.DeleteCarTagUsecase,

	createCarMarkUsecase car_usecases.CreateCarMarkUsecase,
	getCarMarkUsecase car_usecases.GetCarMarkUsecase,
	getCarMarksListUsecase car_usecases.GetCarMarksListUsecase,
	updateCarMarkUsecase car_usecases.UpdateCarMarkUsecase,
	deleteCarMarkUsecase car_usecases.DeleteCarMarkUsecase,

	createCarCategoryUsecase car_usecases.CreateCarCategoryUsecase,
	getCarCategoryUsecase car_usecases.GetCarCategoryUsecase,
	getCarCategoriesListUsecase car_usecases.GetCarCategoriesListUsecase,
	updateCarCategoryUsecase car_usecases.UpdateCarCategoryUsecase,
	deleteCarCategoryUsecase car_usecases.DeleteCarCategoryUsecase,

) *App {
	return &App{
		DB:                              db,
		AuthHandler:                     authHandler,
		CarHandler:                      carHandler,
		CarTagHandler:                   carTagHandler,
		CarMarkHandler:                  carMarkHandler,
		CarCategoryHandler:              carCategoryHandler,
		SignUpUsecase:                   signUpUsecase,
		ConfirmEmailVerificationUsecase: confirmEmailVerificationUsecase,
		SendEmailVerificationUsecase:    sendEmailVerificationUsecase,
		SignInUsecase:                   signInUsecase,
		TokenService:                    tokenSvc,

		CreateCarUsecase: createCarUsecase,

		CreateCarTagUsecase:   createCarTagUsecase,
		GetCarTagUsecase:      getCarTagUsecase,
		GetCarTagsListUsecase: getCarTagsListUsecase,
		UpdateCarTagUsecase:   updateCarTagUsecase,
		DeleteCarTagUsecase:   deleteCarTagUsecase,

		CreateCarMarkUsecase:   createCarMarkUsecase,
		GetCarMarkUsecase:      getCarMarkUsecase,
		GetCarMarksListUsecase: getCarMarksListUsecase,
		UpdateCarMarkUsecase:   updateCarMarkUsecase,
		DeleteCarMarkUsecase:   deleteCarMarkUsecase,

		CreateCarCategoryUsecase:    createCarCategoryUsecase,
		GetCarCategoryUsecase:       getCarCategoryUsecase,
		GetCarCategoriesListUsecase: getCarCategoriesListUsecase,
		UpdateCarCategoryUsecase:    updateCarCategoryUsecase,
		DeleteCarCategoryUsecase:    deleteCarCategoryUsecase,
	}
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
