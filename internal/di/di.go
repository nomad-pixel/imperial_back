package di

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/config"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	authUsecase "github.com/nomad-pixel/imperial/internal/domain/usecases/auth"
	carUsecase "github.com/nomad-pixel/imperial/internal/domain/usecases/car"
	celebrityUsecase "github.com/nomad-pixel/imperial/internal/domain/usecases/celebrity"
	token "github.com/nomad-pixel/imperial/internal/infrastructure/auth"
	"github.com/nomad-pixel/imperial/internal/infrastructure/email"
	imageSvc "github.com/nomad-pixel/imperial/internal/infrastructure/image"
	"github.com/nomad-pixel/imperial/internal/infrastructure/repositories"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/auth"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car"
	carCategory "github.com/nomad-pixel/imperial/internal/interfaces/http/car/category"
	carImage "github.com/nomad-pixel/imperial/internal/interfaces/http/car/image"
	carMark "github.com/nomad-pixel/imperial/internal/interfaces/http/car/mark"
	carTag "github.com/nomad-pixel/imperial/internal/interfaces/http/car/tag"
	celebrity "github.com/nomad-pixel/imperial/internal/interfaces/http/celebrity"
)

func InitializeApp(ctx context.Context, dbURL string) (*App, error) {
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	userRepo := repositories.NewUserRepositoryImpl(db)
	verifyCodeRepo := repositories.NewVerifyCodeRepositoryImpl(db)
	carCategoryRepo := repositories.NewCarCategoryRepositoryImpl(db)
	carTagRepo := repositories.NewCarTagRepositoryImpl(db)
	carMarkRepo := repositories.NewCarMarkRepositoryImpl(db)
	carRepo := repositories.NewCarRepositoryImpl(db)
	carImageRepo := repositories.NewCarImageRepositoryImpl(db)
	celebrityRepo := repositories.NewCelebrityRepositoryImpl(db)

	emailConfig := config.LoadEmailConfig()
	fmt.Println("emailConfig", emailConfig)
	var emailService ports.EmailService

	if emailConfig.Provider == "smtp" {
		if emailConfig.SMTP.Username == "" {
			log.Println("⚠️  ВНИМАНИЕ: SMTP_USERNAME не установлен!")
		}
		if emailConfig.SMTP.Password == "" {
			log.Println("⚠️  ВНИМАНИЕ: SMTP_PASSWORD не установлен!")
		}

		smtpService, err := email.NewSMTPEmailService(email.SMTPConfig{
			Host:     emailConfig.SMTP.Host,
			Port:     emailConfig.SMTP.Port,
			Username: emailConfig.SMTP.Username,
			Password: emailConfig.SMTP.Password,
			From:     emailConfig.SMTP.From,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize SMTP service: %w", err)
		}
		emailService = smtpService
	} else {
		emailService = email.NewConsoleEmailService()
	}
	tokenSvc := token.NewJWTTokenService()

	imageService, err := imageSvc.NewFileImageService("./uploads", "http://localhost:8080/uploads")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize image service: %w", err)
	}

	signUpUsecase := authUsecase.NewSignUpUsecase(userRepo)
	sendEmailVerificationUsecase := authUsecase.NewSendEmailVerificationUsecase(
		userRepo,
		verifyCodeRepo,
		emailService,
	)
	confirmEmailVerificationUsecase := authUsecase.NewConfirmEmailVerificationUsecase(verifyCodeRepo, userRepo)
	signInUsecase := authUsecase.NewSignInUsecase(userRepo, tokenSvc)
	refreshTokenUsecase := authUsecase.NewRefreshTokenUsecase(tokenSvc)

	authHandler := auth.NewAuthHandler(
		signUpUsecase,
		sendEmailVerificationUsecase,
		confirmEmailVerificationUsecase,
		signInUsecase,
		refreshTokenUsecase,
	)

	createCarUsecase := carUsecase.NewCreateCarUsecase(
		carRepo,
	)
	deleteCarUsecase := carUsecase.NewDeleteCarUsecase(
		carRepo,
	)
	updateCarUsecase := carUsecase.NewUpdateCarUsecase(
		carRepo,
	)
	getCarByIdUsecase := carUsecase.NewGetCarByIdUsecase(
		carRepo,
	)
	getListCarsUsecase := carUsecase.NewGetListCarsUsecase(
		carRepo,
	)

	// CarTag usecases
	createCarTagUsecase := carUsecase.NewCreateCarTagUsecase(carTagRepo)
	getCarTagUsecase := carUsecase.NewGetCarTagUsecase(carTagRepo)
	getCarTagsListUsecase := carUsecase.NewGetCarTagsListUsecase(carTagRepo)
	updateCarTagUsecase := carUsecase.NewUpdateCarTagUsecase(carTagRepo)
	deleteCarTagUsecase := carUsecase.NewDeleteCarTagUsecase(carTagRepo)

	// CarMark usecases
	createCarMarkUsecase := carUsecase.NewCreateCarMarkUsecase(carMarkRepo)
	getCarMarkUsecase := carUsecase.NewGetCarMarkUsecase(carMarkRepo)
	getCarMarksListUsecase := carUsecase.NewGetCarMarksListUsecase(carMarkRepo)
	updateCarMarkUsecase := carUsecase.NewUpdateCarMarkUsecase(carMarkRepo)
	deleteCarMarkUsecase := carUsecase.NewDeleteCarMarkUsecase(carMarkRepo)

	// CarCategory usecases
	createCarCategoryUsecase := carUsecase.NewCreateCarCategoryUsecase(carCategoryRepo)
	getCarCategoryUsecase := carUsecase.NewGetCarCategoryUsecase(carCategoryRepo)
	getCarCategoriesListUsecase := carUsecase.NewGetCarCategoriesListUsecase(carCategoryRepo)
	updateCarCategoryUsecase := carUsecase.NewUpdateCarCategoryUsecase(carCategoryRepo)
	deleteCarCategoryUsecase := carUsecase.NewDeleteCarCategoryUsecase(carCategoryRepo)

	// Car Image usecases
	createCarImageUsecase := carUsecase.NewCreateCarImageUsecase(
		carImageRepo,
		imageService,
	)
	deleteCarImageUsecase := carUsecase.NewDeleteCarImageUsecase(
		carImageRepo,
		imageService,
	)

	getCarImagesListUsecase := carUsecase.NewGetCarImagesListUsecase(
		carImageRepo,
		imageService,
	)

	// Celebrity usecases
	carHandler := car.NewCarHandler(createCarUsecase, deleteCarUsecase, updateCarUsecase, getCarByIdUsecase, getListCarsUsecase)

	carTagHandler := carTag.NewCarTagHandler(
		createCarTagUsecase,
		getCarTagUsecase,
		getCarTagsListUsecase,
		updateCarTagUsecase,
		deleteCarTagUsecase,
	)

	carMarkHandler := carMark.NewCarMarkHandler(
		createCarMarkUsecase,
		getCarMarkUsecase,
		getCarMarksListUsecase,
		updateCarMarkUsecase,
		deleteCarMarkUsecase,
	)

	carCategoryHandler := carCategory.NewCarCategoryHandler(
		createCarCategoryUsecase,
		getCarCategoryUsecase,
		getCarCategoriesListUsecase,
		updateCarCategoryUsecase,
		deleteCarCategoryUsecase,
	)

	carImageHandler := carImage.NewCarImageHandler(
		createCarImageUsecase,
		deleteCarImageUsecase,
		getCarImagesListUsecase,
	)

	celebrityUsecase := celebrityUsecase.NewCreateCelebrityUsecase(celebrityRepo)

	celebrityHandler := celebrity.NewCelebrityHandler(celebrityUsecase)

	app := NewApp(
		db,
		signUpUsecase,
		sendEmailVerificationUsecase,
		confirmEmailVerificationUsecase,
		signInUsecase,
		authHandler,
		tokenSvc,
		celebrityHandler,

		// Car usecases
		createCarUsecase,
		getCarByIdUsecase,
		getListCarsUsecase,
		updateCarUsecase,
		deleteCarUsecase,

		carHandler,
		carImageHandler,
		carTagHandler,
		carMarkHandler,
		carCategoryHandler,

		// CarTag usecases
		createCarTagUsecase,
		getCarTagUsecase,
		getCarTagsListUsecase,
		updateCarTagUsecase,
		deleteCarTagUsecase,
		// CarMark usecases
		createCarMarkUsecase,
		getCarMarkUsecase,
		getCarMarksListUsecase,
		updateCarMarkUsecase,
		deleteCarMarkUsecase,
		// CarCategory usecases
		createCarCategoryUsecase,
		getCarCategoryUsecase,
		getCarCategoriesListUsecase,
		updateCarCategoryUsecase,
		deleteCarCategoryUsecase,

		// Celebrity usecases
		celebrityUsecase,
	)

	return app, nil
}
