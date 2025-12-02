package di

import (
	"github.com/google/wire"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/auth"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car"
	carCategory "github.com/nomad-pixel/imperial/internal/interfaces/http/car/category"
	carImage "github.com/nomad-pixel/imperial/internal/interfaces/http/car/image"
	carMark "github.com/nomad-pixel/imperial/internal/interfaces/http/car/mark"
	carTag "github.com/nomad-pixel/imperial/internal/interfaces/http/car/tag"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/celebrity"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/driver"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/lead"
)

// HandlerSet provides all HTTP handlers
var HandlerSet = wire.NewSet(
	auth.NewAuthHandler,
	car.NewCarHandler,
	carImage.NewCarImageHandler,
	carTag.NewCarTagHandler,
	carMark.NewCarMarkHandler,
	carCategory.NewCarCategoryHandler,
	celebrity.NewCelebrityHandler,
	lead.NewLeadHandler,
	driver.NewDriverHandler,
)
