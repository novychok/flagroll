package platformapiv1

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"
	"github.com/novychok/flagroll/platform/internal/service"
	platformapiv1 "github.com/novychok/flagroll/platform/pkg/api/platform/v1"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
)

type Config struct {
	Port              int `mapstructure:"PLATFORM_API_V1_PORT"`
	ReadHeaderTimeout int `mapstructure:"READ_HEADER_TIMEOUT"`
	QuietDownPeriod   int `mapstructure:"QUIET_DOWN_PERIOD"`
	CorsMaxAge        int `mapstructure:"CORS_MAX_AGE"`
}

type Server struct {
	l   *slog.Logger
	cfg *Config

	authorizationService service.Authorization

	h platformapiv1.ServerInterface
}

func (s *Server) Run(ctx context.Context) error {
	logger := httplog.NewLogger("platformApi", httplog.Options{
		JSON:             true,
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		QuietDownRoutes: []string{
			"/health",
		},
		QuietDownPeriod: 10 * time.Second,
	})

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8010/",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Heartbeat("/health"))
	r.Use(ContextMiddleware())

	swagger, err := platformapiv1.GetSwagger()
	if err != nil {
		return err
	}

	swagger.Servers = nil

	options := oapimiddleware.Options{
		ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			b, _ := json.Marshal(platformapiv1.Error{
				Message: message,
			})
			_, _ = w.Write(b)
		},
		Options: openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
				req := input.RequestValidationInput.Request
				customCtx := ContextFromRequest(req)

				switch input.SecuritySchemeName {
				case "bearerAuth":
					customCtx.Set("authRequired", true)
				case "keyAuth":
					customCtx.Set("apiKeyRequired", true)
				}

				return nil
			},
			ExcludeRequestBody:    true,
			ExcludeResponseBody:   false,
			IncludeResponseStatus: true,
			MultiError:            false,
		},
	}

	r.Use(oapimiddleware.OapiRequestValidatorWithOptions(swagger, &options))

	r.Use(s.auth)
	r.Use(s.keyAuth)

	platformapiv1.HandlerFromMux(s.h, r)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", s.cfg.Port),
	}

	s.l.Info("platform api server is running", slog.Int("port", s.cfg.Port))

	err = srv.ListenAndServe()
	if err != nil {
		s.l.Error("platform server error", slog.Any("error", err))
	}

	return nil
}

func NewServer(
	l *slog.Logger,
	cfg *Config,

	authorizationService service.Authorization,

	h platformapiv1.ServerInterface,
) *Server {
	return &Server{
		l:   l,
		cfg: cfg,

		authorizationService: authorizationService,

		h: h,
	}
}
