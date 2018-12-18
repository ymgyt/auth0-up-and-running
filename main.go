package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/negroni"
	"github.com/ymgyt/appkit/middlewares"
	"go.uber.org/zap"

	"github.com/julienschmidt/httprouter"

	"github.com/ymgyt/appkit/envvar"
	"github.com/ymgyt/appkit/handlers"
	"github.com/ymgyt/appkit/logging"
	"github.com/ymgyt/appkit/server"
	"github.com/ymgyt/appkit/services"
)

type config struct {
	Root          string `envvar:"APP_ROOT,required"`
	Host          string `envvar:"APP_HOST,required,default=localhost"`
	Port          string `envvar:"APP_PORT,required,default=443"`
	GCPProjectID  string `envvar:"GCP_PROJECT_ID,required"`
	GCPCredential string `envvar:"GCP_CREDENTIAL,required"`
	Auth0Domain   string `envvar:"AUTH0_DOMAIN,required"`
	Auth0ClientID string `envvar:"AUTH0_CLIENT_ID,required"`
}

func router(ctx context.Context, cfg *config, logger *zap.Logger) http.Handler {
	r := httprouter.New()

	static := handlers.MustStatic(cfg.Root+"/static", "/static")
	ts := handlers.MustTemplateSet(&handlers.TemplateSetConfig{Root: cfg.Root + "/templates", AlwaysReload: true})
	auth0 := &Auth0{ts: ts, logger: logger, cfg: cfg}

	r.GET("/static/*filepath", static.ServeStatic)
	r.GET("/login", auth0.RenderLogin)

	common := negroni.New(middlewares.MustLogging(&middlewares.LoggingConfig{Logger: logger, Console: true}))
	common.UseHandler(r)

	return common
}

func main() {

	ctx := context.Background()
	cfg := configFromEnv()
	logger := logging.Must(&logging.Config{Out: os.Stdout, Level: "debug"})
	router := router(ctx, cfg, logger)
	server := server.Must(&server.Config{
		Addr:            ":" + cfg.Port,
		DisableHTTPS:    true,
		Handler:         router,
		DatastoreClient: services.MustDatastore(services.NewDatastoreFromFile(ctx, cfg.GCPProjectID, cfg.GCPCredential)),
	})

	logger.Sugar().Info("running on ", cfg.Port)
	logger.Sugar().Info(server.Run())
}

func configFromEnv() *config {
	cfg := &config{}
	if err := envvar.Inject(cfg); err != nil {
		fail(err)
	}
	return cfg
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
