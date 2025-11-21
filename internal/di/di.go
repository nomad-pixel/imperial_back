package di

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/config"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/domain/usecases/auth_usecases"
	"github.com/nomad-pixel/imperial/internal/domain/usecases/car_usecases"
	token "github.com/nomad-pixel/imperial/internal/infrastructure/auth"
	"github.com/nomad-pixel/imperial/internal/infrastructure/email"
	imageSvc "github.com/nomad-pixel/imperial/internal/infrastructure/image"
	"github.com/nomad-pixel/imperial/internal/infrastructure/repositories"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/auth"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_category"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_image"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_mark"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car_tag"
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

	imageService, err := imageSvc.NewFileImageService("./uploads", "http://localhost:8080/images")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize image service: %w", err)
	}

	signUpUsecase := auth_usecases.NewSignUpUsecase(userRepo)
	sendEmailVerificationUsecase := auth_usecases.NewSendEmailVerificationUsecase(
		userRepo,
		verifyCodeRepo,
		emailService,
	)
	confirmEmailVerificationUsecase := auth_usecases.NewConfirmEmailVerificationUsecase(verifyCodeRepo, userRepo)
	signInUsecase := auth_usecases.NewSignInUsecase(userRepo, tokenSvc)
	refreshTokenUsecase := auth_usecases.NewRefreshTokenUsecase(tokenSvc)

	authHandler := auth.NewAuthHandler(
		signUpUsecase,
		sendEmailVerificationUsecase,
		confirmEmailVerificationUsecase,
		signInUsecase,
		refreshTokenUsecase,
	)

	createCarUsecase := car_usecases.NewCreateCarUsecase(
		carRepo,
	)
	deleteCarUsecase := car_usecases.NewDeleteCarUsecase(
		carRepo,
	)
	updateCarUsecase := car_usecases.NewUpdateCarUsecase(
		carRepo,
	)
	getCarByIdUsecase := car_usecases.NewGetCarByIdUsecase(
		carRepo,
	)
	getListCarsUsecase := car_usecases.NewGetListCarsUsecase(
		carRepo,
	)

	// CarTag usecases
	createCarTagUsecase := car_usecases.NewCreateCarTagUsecase(carTagRepo)
	getCarTagUsecase := car_usecases.NewGetCarTagUsecase(carTagRepo)
	getCarTagsListUsecase := car_usecases.NewGetCarTagsListUsecase(carTagRepo)
	updateCarTagUsecase := car_usecases.NewUpdateCarTagUsecase(carTagRepo)
	deleteCarTagUsecase := car_usecases.NewDeleteCarTagUsecase(carTagRepo)

	// CarMark usecases
	createCarMarkUsecase := car_usecases.NewCreateCarMarkUsecase(carMarkRepo)
	getCarMarkUsecase := car_usecases.NewGetCarMarkUsecase(carMarkRepo)
	getCarMarksListUsecase := car_usecases.NewGetCarMarksListUsecase(carMarkRepo)
	updateCarMarkUsecase := car_usecases.NewUpdateCarMarkUsecase(carMarkRepo)
	deleteCarMarkUsecase := car_usecases.NewDeleteCarMarkUsecase(carMarkRepo)

	// CarCategory usecases
	createCarCategoryUsecase := car_usecases.NewCreateCarCategoryUsecase(carCategoryRepo)
	getCarCategoryUsecase := car_usecases.NewGetCarCategoryUsecase(carCategoryRepo)
	getCarCategoriesListUsecase := car_usecases.NewGetCarCategoriesListUsecase(carCategoryRepo)
	updateCarCategoryUsecase := car_usecases.NewUpdateCarCategoryUsecase(carCategoryRepo)
	deleteCarCategoryUsecase := car_usecases.NewDeleteCarCategoryUsecase(carCategoryRepo)

	// Car Image usecases
	createCarImageUsecase := car_usecases.NewCreateCarImageUsecase(
		carImageRepo,
		imageService,
	)
	deleteCarImageUsecase := car_usecases.NewDeleteCarImageUsecase(
		carImageRepo,
		imageService,
	)

	getCarImagesListUsecase := car_usecases.NewGetCarImagesListUsecase(
		carImageRepo,
	)

	carHandler := car.NewCarHandler(createCarUsecase, deleteCarUsecase, updateCarUsecase, getCarByIdUsecase, getListCarsUsecase)

	carTagHandler := car_tag.NewCarTagHandler(
		createCarTagUsecase,
		getCarTagUsecase,
		getCarTagsListUsecase,
		updateCarTagUsecase,
		deleteCarTagUsecase,
	)

	carMarkHandler := car_mark.NewCarMarkHandler(
		createCarMarkUsecase,
		getCarMarkUsecase,
		getCarMarksListUsecase,
		updateCarMarkUsecase,
		deleteCarMarkUsecase,
	)

	carCategoryHandler := car_category.NewCarCategoryHandler(
		createCarCategoryUsecase,
		getCarCategoryUsecase,
		getCarCategoriesListUsecase,
		updateCarCategoryUsecase,
		deleteCarCategoryUsecase,
	)

	carImageHandler := car_image.NewCarImageHandler(
		createCarImageUsecase,
		deleteCarImageUsecase,
		getCarImagesListUsecase,
	)

	app := NewApp(
		db,
		signUpUsecase,
		sendEmailVerificationUsecase,
		confirmEmailVerificationUsecase,
		signInUsecase,
		authHandler,
		tokenSvc,

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
	)

	return app, nil
}
