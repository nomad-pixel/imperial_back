package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/config"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/auth"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car"
	carCategory "github.com/nomad-pixel/imperial/internal/interfaces/http/car/category"
	carImage "github.com/nomad-pixel/imperial/internal/interfaces/http/car/image"
	carMark "github.com/nomad-pixel/imperial/internal/interfaces/http/car/mark"
	carTag "github.com/nomad-pixel/imperial/internal/interfaces/http/car/tag"
	celebrity "github.com/nomad-pixel/imperial/internal/interfaces/http/celebrity"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/driver"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/lead"
)

// App contains all application dependencies
type App struct {
	Config             *config.Config
	DB                 *pgxpool.Pool
	TokenService       ports.TokenService
	AuthHandler        *auth.AuthHandler
	CarHandler         *car.CarHandler
	CarImageHandler    *carImage.CarImageHandler
	CarTagHandler      *carTag.CarTagHandler
	CarMarkHandler     *carMark.CarMarkHandler
	CarCategoryHandler *carCategory.CarCategoryHandler
	CelebrityHandler   *celebrity.CelebrityHandler
	LeadHandler        *lead.LeadHandler
	DriverHandler      *driver.DriverHandler
}

// NewApp creates a new App instance with all dependencies injected
func NewApp(
	cfg *config.Config,
	db *pgxpool.Pool,
	tokenSvc ports.TokenService,
	authHandler *auth.AuthHandler,
	carHandler *car.CarHandler,
	carImageHandler *carImage.CarImageHandler,
	carTagHandler *carTag.CarTagHandler,
	carMarkHandler *carMark.CarMarkHandler,
	carCategoryHandler *carCategory.CarCategoryHandler,
	celebrityHandler *celebrity.CelebrityHandler,
	leadHandler *lead.LeadHandler,
	driverHandler *driver.DriverHandler,
) *App {
	return &App{
		Config:             cfg,
		DB:                 db,
		TokenService:       tokenSvc,
		AuthHandler:        authHandler,
		CarHandler:         carHandler,
		CarImageHandler:    carImageHandler,
		CarTagHandler:      carTagHandler,
		CarMarkHandler:     carMarkHandler,
		CarCategoryHandler: carCategoryHandler,
		CelebrityHandler:   celebrityHandler,
		LeadHandler:        leadHandler,
		DriverHandler:      driverHandler,
	}
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
