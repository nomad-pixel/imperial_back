//go:build wireinject
// +build wireinject

package di

import (
	"context"

	"github.com/google/wire"
)

// InitializeApp initializes the application with all dependencies using Wire
func InitializeApp(ctx context.Context) (*App, error) {
	wire.Build(ProviderSet)
	return nil, nil
}
