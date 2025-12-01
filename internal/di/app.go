package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	authUsecase "github.com/nomad-pixel/imperial/internal/domain/usecases/auth"
	carUsecase "github.com/nomad-pixel/imperial/internal/domain/usecases/car"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/auth"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_category"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_image"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_mark"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_tag"
)

type App struct {
	DB                              *pgxpool.Pool
	AuthHandler                     *auth.AuthHandler
	CarHandler                      *car.CarHandler
	CarImageHandler                 *car_image.CarImageHandler
	CarTagHandler                   *car_tag.CarTagHandler
	CarMarkHandler                  *car_mark.CarMarkHandler
	CarCategoryHandler              *car_category.CarCategoryHandler
	SignUpUsecase                   authUsecase.SignUpUsecase
	SignInUsecase                   authUsecase.SignInUsecase
	SendEmailVerificationUsecase    authUsecase.SendEmailVerificationUsecase
	ConfirmEmailVerificationUsecase authUsecase.ConfirmEmailVerificationUsecase
	TokenService                    ports.TokenService

	// Car usecases
	CreateCarUsecase   carUsecase.CreateCarUsecase
	GetCarByIdUsecase  carUsecase.GetCarByIdUsecase
	GetListCarsUsecase carUsecase.GetListCarsUsecase
	UpdateCarUsecase   carUsecase.UpdateCarUsecase
	DeleteCarUsecase   carUsecase.DeleteCarUsecase

	// CarTag usecases
	CreateCarTagUsecase   carUsecase.CreateCarTagUsecase
	GetCarTagUsecase      carUsecase.GetCarTagUsecase
	GetCarTagsListUsecase carUsecase.GetCarTagsListUsecase
	UpdateCarTagUsecase   carUsecase.UpdateCarTagUsecase
	DeleteCarTagUsecase   carUsecase.DeleteCarTagUsecase

	// CarMark usecases
	CreateCarMarkUsecase   carUsecase.CreateCarMarkUsecase
	GetCarMarkUsecase      carUsecase.GetCarMarkUsecase
	GetCarMarksListUsecase carUsecase.GetCarMarksListUsecase
	UpdateCarMarkUsecase   carUsecase.UpdateCarMarkUsecase
	DeleteCarMarkUsecase   carUsecase.DeleteCarMarkUsecase

	// CarCategory usecases
	CreateCarCategoryUsecase    carUsecase.CreateCarCategoryUsecase
	GetCarCategoryUsecase       carUsecase.GetCarCategoryUsecase
	GetCarCategoriesListUsecase carUsecase.GetCarCategoriesListUsecase
	UpdateCarCategoryUsecase    carUsecase.UpdateCarCategoryUsecase
	DeleteCarCategoryUsecase    carUsecase.DeleteCarCategoryUsecase
}

func NewApp(
	db *pgxpool.Pool,
	signUpUsecase authUsecase.SignUpUsecase,
	sendEmailVerificationUsecase authUsecase.SendEmailVerificationUsecase,
	confirmEmailVerificationUsecase authUsecase.ConfirmEmailVerificationUsecase,
	signInUsecase authUsecase.SignInUsecase,
	authHandler *auth.AuthHandler,
	tokenSvc ports.TokenService,

	createCarUsecase carUsecase.CreateCarUsecase,
	getCarByIdUsecase carUsecase.GetCarByIdUsecase,
	getListCarsUsecase carUsecase.GetListCarsUsecase,
	updateCarUsecase carUsecase.UpdateCarUsecase,
	deleteCarUsecase carUsecase.DeleteCarUsecase,
	carHandler *car.CarHandler,
	carImageHandler *car_image.CarImageHandler,
	carTagHandler *car_tag.CarTagHandler,
	carMarkHandler *car_mark.CarMarkHandler,
	carCategoryHandler *car_category.CarCategoryHandler,

	createCarTagUsecase carUsecase.CreateCarTagUsecase,
	getCarTagUsecase carUsecase.GetCarTagUsecase,
	getCarTagsListUsecase carUsecase.GetCarTagsListUsecase,
	updateCarTagUsecase carUsecase.UpdateCarTagUsecase,
	deleteCarTagUsecase carUsecase.DeleteCarTagUsecase,

	createCarMarkUsecase carUsecase.CreateCarMarkUsecase,
	getCarMarkUsecase carUsecase.GetCarMarkUsecase,
	getCarMarksListUsecase carUsecase.GetCarMarksListUsecase,
	updateCarMarkUsecase carUsecase.UpdateCarMarkUsecase,
	deleteCarMarkUsecase carUsecase.DeleteCarMarkUsecase,

	createCarCategoryUsecase carUsecase.CreateCarCategoryUsecase,
	getCarCategoryUsecase carUsecase.GetCarCategoryUsecase,
	getCarCategoriesListUsecase carUsecase.GetCarCategoriesListUsecase,
	updateCarCategoryUsecase carUsecase.UpdateCarCategoryUsecase,
	deleteCarCategoryUsecase carUsecase.DeleteCarCategoryUsecase,

) *App {
	return &App{
		DB:                              db,
		AuthHandler:                     authHandler,
		CarHandler:                      carHandler,
		CarImageHandler:                 carImageHandler,
		CarTagHandler:                   carTagHandler,
		CarMarkHandler:                  carMarkHandler,
		CarCategoryHandler:              carCategoryHandler,
		SignUpUsecase:                   signUpUsecase,
		ConfirmEmailVerificationUsecase: confirmEmailVerificationUsecase,
		SendEmailVerificationUsecase:    sendEmailVerificationUsecase,
		SignInUsecase:                   signInUsecase,
		TokenService:                    tokenSvc,

		CreateCarUsecase:   createCarUsecase,
		GetCarByIdUsecase:  getCarByIdUsecase,
		GetListCarsUsecase: getListCarsUsecase,
		UpdateCarUsecase:   updateCarUsecase,
		DeleteCarUsecase:   deleteCarUsecase,

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
