package runner

import (
	"context"
	"github.com/kuritka/onho.io/common/middleware/rest"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/kuritka/onho.io/common/utils"
)

type RestRunner struct {
	router          chi.Router
	handler         http.Handler
	listenerFactory func() (net.Listener, error)
}

// NewRestRunner doc
func NewRestRunner(port string) *RestRunner {
	router := chi.NewRouter()
	return &RestRunner{
		router:  router,
		handler: router,
		listenerFactory: func() (net.Listener, error) {
			return withAddressListenerFactory(port)
		},
	}
}


// WithWebServerMiddlewares doc
func (r *RestRunner) WithWebServerMiddlewares() *RestRunner {
	r.router.Use(
		rest.LoggingMiddleware(zerologger),
		middleware.DefaultCompress,
		middleware.Recoverer,
		rest.SecurityHeaderMiddleware,
		rest.SetHTTPHeader("Cache-Control",
			"no-cache, no-store, must-revalidate",
			[]string{"/", "index.html"}),
		rest.CacheControlMiddleware(1*time.Hour),
	)
	return r
}

// ExportHealthEndpoint doc
func (r *RestRunner) ExportHealthEndpoint() *RestRunner {
	r.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return r
}



func (r *RestRunner) run(ctx context.Context) error {
	listener, err := r.listenerFactory()
		if err != nil {
		utils.FailOnError(err,"failed to create rest listener")
	}

	zerologger.Info().Msg("successfully started rest listener on port: " + listener.Addr().String())

	//service := contextType.Service.GetFromContext(ctx)
	////for _, registerer := range registerers {
	////	if errR := registerer.RegisterWithRouter(ctx, r.router); errR != nil {
	////		return errors.Wrapf(errR, "cannot register router '%v'", service)
	////	}
	////}
	//zerologger.Info().Msgf("router '%v' started successfully", service)
	// TODO: AllowedOrigins for local development only!!! has to be updated in prod env
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{
			"*",
			"http://localhost",
			"http://localhost:8000",
			"http://localhost:8080"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete},
		AllowCredentials: true,
		AllowedHeaders: []string{
			"Origin",
			"Accept",
			"Accept-Encoding",
			"Content-Type",
			"Content-Length",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Authorization",
			"X-CSRF-Token"},
	})
	svr := &http.Server{Handler: cors.Handler(r.handler)}

	go func() {
		<-ctx.Done()
		if errS := svr.Shutdown(ctx); errS != nil {
			zerologger.Warn().Err(errS).Msg("failed to close http server")
		}
	}()

	err = svr.Serve(listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (r *RestRunner) String() string{
	return "REST-runner"
}