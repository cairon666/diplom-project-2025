package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/api/www/auth_router"
	"github.com/cairon666/vkr-backend/internal/api/www/device_router"
	"github.com/cairon666/vkr-backend/internal/api/www/external_apps_router"
	"github.com/cairon666/vkr-backend/internal/api/www/health_router"
	"github.com/cairon666/vkr-backend/internal/api/www/user_router"
	"github.com/cairon666/vkr-backend/internal/config"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/repositories/auth_providers_repo"
	"github.com/cairon666/vkr-backend/internal/repositories/dbqueries"
	"github.com/cairon666/vkr-backend/internal/repositories/device_repo"
	"github.com/cairon666/vkr-backend/internal/repositories/external_apps_repo"
	"github.com/cairon666/vkr-backend/internal/repositories/health_repo"
	"github.com/cairon666/vkr-backend/internal/repositories/roles_repo"
	"github.com/cairon666/vkr-backend/internal/repositories/tempid_repo"
	"github.com/cairon666/vkr-backend/internal/repositories/user_passwords_repo"
	"github.com/cairon666/vkr-backend/internal/repositories/user_repo"
	"github.com/cairon666/vkr-backend/internal/services/auth_service"
	"github.com/cairon666/vkr-backend/internal/services/device_service"
	"github.com/cairon666/vkr-backend/internal/services/external_apps_service"
	"github.com/cairon666/vkr-backend/internal/services/health_service"
	"github.com/cairon666/vkr-backend/internal/services/role_service"
	"github.com/cairon666/vkr-backend/internal/services/telegram_service"
	"github.com/cairon666/vkr-backend/internal/services/user_service"
	"github.com/cairon666/vkr-backend/internal/usecases/auth_usecase"
	"github.com/cairon666/vkr-backend/internal/usecases/device_usecase"
	"github.com/cairon666/vkr-backend/internal/usecases/external_apps_usecase"
	"github.com/cairon666/vkr-backend/internal/usecases/health_usecase"
	"github.com/cairon666/vkr-backend/internal/usecases/user_usecase"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/cairon666/vkr-backend/pkg/password_hasher"
	"github.com/cairon666/vkr-backend/pkg/postgres"
	redis_client "github.com/cairon666/vkr-backend/pkg/redis-client"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func main() {
	flagConfig := flag.String("config", "./config/config.yaml", "path to config file")
	flag.Parse()

	fx.New(
		fx.Provide(func() (*config.Config, error) {
			return config.GetConfig(*flagConfig)
		}),
		fx.Provide(NewLogger),
		fx.Provide(
			fx.Annotate(NewPgxPool,
				fx.As(new(dbqueries.DBTX)),
				fx.As(new(postgres.PostgresClient)),
			),
			dbqueries.New,
		),
		fx.Provide(NewRedisClient),
		fx.Provide(fx.Annotate(password_hasher.NewPasswordHasher,
			fx.As(new(auth_usecase.PasswordHasher)),
			fx.As(new(user_usecase.PasswordHasher)),
		)),
		//	repos
		fx.Provide(fx.Annotate(user_repo.NewUserRepo, fx.As(new(user_service.UserRepo)))),
		fx.Provide(fx.Annotate(roles_repo.NewRolesRepo,
			fx.As(new(role_service.RoleRepos)),
			fx.As(new(auth_service.RoleRepos)),
			fx.As(new(indentity.RolesService)),
		)),
		fx.Provide(fx.Annotate(tempid_repo.NewTempIdRepo, fx.As(new(auth_service.TempIdRepo)))),
		fx.Provide(fx.Annotate(user_passwords_repo.NewUserPasswordsRepo, fx.As(new(auth_service.UserPasswordsRepo)))),
		fx.Provide(fx.Annotate(auth_providers_repo.NewAuthProvidersRepo, fx.As(new(auth_service.AuthProviderRepo)))),
		fx.Provide(fx.Annotate(external_apps_repo.NewExternalAppsRepo,
			fx.As(new(external_apps_service.ExternalAppsRepo)),
		)),
		fx.Provide(fx.Annotate(device_repo.NewDeviceRepo,
			fx.As(new(device_service.DeviceRepo)),
		)),
		fx.Provide(fx.Annotate(health_repo.NewHealthRepo,
			fx.As(new(health_service.HealthRepo)),
		)),
		// services
		fx.Provide(indentity.NewApiKeyService),
		fx.Provide(indentity.NewIdentityService),
		fx.Provide(indentity.NewJWTService),
		fx.Provide(fx.Annotate(indentity.NewJWTService,
			fx.As(new(auth_service.JWTService)),
		)),
		fx.Provide(fx.Annotate(auth_service.NewAuthService,
			fx.As(new(auth_usecase.AuthService)),
			fx.As(new(user_usecase.AuthService)),
		)),
		fx.Provide(fx.Annotate(role_service.NewRoleService,
			fx.As(new(auth_usecase.RolesService)),
			fx.As(new(external_apps_usecase.RolesService)),
		)),
		fx.Provide(fx.Annotate(user_service.NewUserService,
			fx.As(new(auth_usecase.UserService)),
			fx.As(new(user_usecase.UserService)),
		)),
		fx.Provide(fx.Annotate(telegram_service.NewTelegramService,
			fx.As(new(user_usecase.TelegramService)),
			fx.As(new(auth_usecase.TelegramService)),
		)),
		fx.Provide(fx.Annotate(external_apps_service.NewExternalAppsService,
			fx.As(new(indentity.ExternalAppsService)),
			fx.As(new(external_apps_usecase.ExternalAppsService)),
		)),
		fx.Provide(fx.Annotate(device_service.NewDeviceService,
			fx.As(new(device_usecase.DeviceService)),
		)),
		fx.Provide(fx.Annotate(health_service.NewHealthService,
			fx.As(new(health_usecase.HealthService)),
		)),
		// usecases
		fx.Provide(fx.Annotate(auth_usecase.NewAuthUsecase, fx.As(new(auth_router.AuthUsecase)))),
		fx.Provide(fx.Annotate(user_usecase.NewUserUsecase, fx.As(new(user_router.UserUsecase)))),
		fx.Provide(fx.Annotate(external_apps_usecase.NewExternalAppsUsecase, fx.As(new(external_apps_router.ExternalAppsUsecase)))),
		fx.Provide(fx.Annotate(device_usecase.NewDeviceUsecase, fx.As(new(device_router.DeviceUsecase)))),
		fx.Provide(fx.Annotate(health_usecase.NewHealthUsecase, fx.As(new(health_router.HealthUsecase)))),
		// www
		fx.Provide(
			AsRoute(auth_router.NewAuthRouter),
			AsRoute(user_router.NewUserRouter),
			AsRoute(external_apps_router.NewExternalAppsRouter),
			AsRoute(device_router.NewDeviceRouter),
			AsRoute(health_router.NewHealthRouter),
			www.NewHTTPServer,
			fx.Annotate(
				www.NewServeMux,
				fx.ParamTags(`group:"routes"`),
			)),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(www.Router)),
		fx.ResultTags(`group:"routes"`),
	)
}

func NewPgxPool(lc fx.Lifecycle, conf *config.Config) (postgres.PostgresClient, error) {
	pgClient, err := postgres.NewPostgres(context.Background(), conf.PostgresURL)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			if pgClient != nil {
				pgClient.Close()
			}
			return nil
		},
	})

	return pgClient, nil
}

func NewLogger(conf *config.Config) (*logger.Logger, error) {
	if conf.Log.Out == "console" {
		return logger.NewDev()
	}

	return logger.NewProd()
}

func NewRedisClient(lc fx.Lifecycle, conf *config.Config) (*redis.Client, error) {
	redisClient, err := redis_client.NewRedisClient(conf.Redis.Addr, conf.Redis.Password, conf.Redis.DB)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			if redisClient != nil {
				redisClient.Close()
			}
			return nil
		},
	})

	return redisClient, nil
}
