package di

import (
	"github.com/google/wire"
	authUsecase "github.com/nomad-pixel/imperial/internal/domain/usecases/auth"
	carUsecase "github.com/nomad-pixel/imperial/internal/domain/usecases/car"
	celebrityUsecase "github.com/nomad-pixel/imperial/internal/domain/usecases/celebrity"
)

// AuthUsecaseSet provides all auth-related use cases
var AuthUsecaseSet = wire.NewSet(
	authUsecase.NewSignUpUsecase,
	authUsecase.NewSendEmailVerificationUsecase,
	authUsecase.NewConfirmEmailVerificationUsecase,
	authUsecase.NewSignInUsecase,
	authUsecase.NewRefreshTokenUsecase,
)

// CarUsecaseSet provides all car-related use cases
var CarUsecaseSet = wire.NewSet(
	// Car CRUD
	carUsecase.NewCreateCarUsecase,
	carUsecase.NewDeleteCarUsecase,
	carUsecase.NewUpdateCarUsecase,
	carUsecase.NewGetCarByIdUsecase,
	carUsecase.NewGetListCarsUsecase,

	// Car Tag
	carUsecase.NewCreateCarTagUsecase,
	carUsecase.NewGetCarTagUsecase,
	carUsecase.NewGetCarTagsListUsecase,
	carUsecase.NewUpdateCarTagUsecase,
	carUsecase.NewDeleteCarTagUsecase,

	// Car Mark
	carUsecase.NewCreateCarMarkUsecase,
	carUsecase.NewGetCarMarkUsecase,
	carUsecase.NewGetCarMarksListUsecase,
	carUsecase.NewUpdateCarMarkUsecase,
	carUsecase.NewDeleteCarMarkUsecase,

	// Car Category
	carUsecase.NewCreateCarCategoryUsecase,
	carUsecase.NewGetCarCategoryUsecase,
	carUsecase.NewGetCarCategoriesListUsecase,
	carUsecase.NewUpdateCarCategoryUsecase,
	carUsecase.NewDeleteCarCategoryUsecase,

	// Car Image
	carUsecase.NewCreateCarImageUsecase,
	carUsecase.NewDeleteCarImageUsecase,
	carUsecase.NewGetCarImagesListUsecase,
)

// CelebrityUsecaseSet provides all celebrity-related use cases
var CelebrityUsecaseSet = wire.NewSet(
	celebrityUsecase.NewCreateCelebrityUsecase,
	celebrityUsecase.NewUploadCelebrityImageUsecase,
	celebrityUsecase.NewGetCelebrityByIdUsecase,
	celebrityUsecase.NewListCelebritiesUsecase,
	celebrityUsecase.NewUpdateCelebrityUsecase,
	celebrityUsecase.NewDeleteCelebrityUsecase,
)
