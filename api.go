package servicekit

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oadamia/servicekit/customvalidator"
	"github.com/oadamia/servicekit/nanoid"
)

var e *echo.Echo

const traceKey = "trace_id"

func InitAPI() {
	e = echo.New()
	nanoidgen := nanoid.DefaultGenerator("")
	e.Validator = customvalidator.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		LogMethod:   true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(ctx echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.DebugContext(ctx.Request().Context(), "HTTP request",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("method", v.Method))

			} else {
				slog.ErrorContext(ctx.Request().Context(), "HTTP request failed",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("method", v.Method),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			traceID, err := nanoidgen.Generate()
			if err != nil {
				return err
			}

			ctx := AppendCtx(c.Request().Context(), slog.String(traceKey, traceID))
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})
}

func Start(port string) {
	if e != nil {
		go func() {
			if err := e.Start(port); err != nil {
				slog.Error("shutting down the server: ", "error", err)
			}
		}()
	} else {
		panic("Start. echo does not exist")
	}

	slog.Info("Started On ", "Port ", port)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("Shutting down server...")
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func Echo() *echo.Echo {
	return e
}
