package application

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/storage/postgres"
	httpapi "github.com/tetafro/nott-backend-go/internal/transport/http"
)

// Application is an entity that gets things done.
type Application struct {
	addr   string
	router *chi.Mux
	log    logrus.FieldLogger
}

// New creates main application instance that handles all requests.
func New(db *gorm.DB, addr string, log logrus.FieldLogger) (*Application, error) {
	app := &Application{addr: addr, log: log}

	foldersRepo := postgres.NewFoldersRepo(db)
	foldersController := httpapi.NewFoldersController(foldersRepo, log)

	notepadsRepo := postgres.NewNotepadsRepo(db)
	notepadsController := httpapi.NewNotepadsController(notepadsRepo, log)

	notesRepo := postgres.NewNotesRepo(db)
	notesController := httpapi.NewNotesController(notesRepo, log)

	usersRepo := postgres.NewUsersRepo(db)

	tokensRepo := postgres.NewTokensRepo(db)

	authController := httpapi.NewAuthController(usersRepo, tokensRepo, log)

	mwAuth := httpapi.NewAuthMiddleware(usersRepo, log)
	mwLog := middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log})

	// Auth router
	ra := chi.NewRouter()
	ra.MethodFunc(http.MethodPost, "/register", authController.Register)
	ra.MethodFunc(http.MethodPost, "/login", authController.Login)
	ra.MethodFunc(http.MethodPost, "/logout", authController.Logout)

	// Application router
	r := chi.NewRouter()
	// Users
	r.MethodFunc(http.MethodGet, "/profile", authController.GetProfile)
	r.MethodFunc(http.MethodPut, "/profile", authController.UpdateProfile)
	// Folders
	r.MethodFunc(http.MethodGet, "/folders", foldersController.GetList)
	r.MethodFunc(http.MethodPost, "/folders", foldersController.Create)
	r.MethodFunc(http.MethodGet, "/folders/{id}", foldersController.GetOne)
	r.MethodFunc(http.MethodPut, "/folders/{id}", foldersController.Update)
	r.MethodFunc(http.MethodDelete, "/folders/{id}", foldersController.Delete)
	// Notepads
	r.MethodFunc(http.MethodGet, "/notepads", notepadsController.GetList)
	r.MethodFunc(http.MethodPost, "/notepads", notepadsController.Create)
	r.MethodFunc(http.MethodGet, "/notepads/{id}", notepadsController.GetOne)
	r.MethodFunc(http.MethodPut, "/notepads/{id}", notepadsController.Update)
	r.MethodFunc(http.MethodDelete, "/notepads/{id}", notepadsController.Delete)
	// Notes
	r.MethodFunc(http.MethodGet, "/notes", notesController.GetList)
	r.MethodFunc(http.MethodPost, "/notes", notesController.Create)
	r.MethodFunc(http.MethodGet, "/notes/{id}", notesController.GetOne)
	r.MethodFunc(http.MethodPut, "/notes/{id}", notesController.Update)
	r.MethodFunc(http.MethodDelete, "/notes/{id}", notesController.Delete)

	app.router = chi.NewRouter()
	app.router.MethodFunc(http.MethodGet, "/healthz", healthz)
	app.router.Mount("/api/v1/auth", mwLog(ra))
	app.router.Mount("/api/v1", mwLog(mwAuth(r)))

	return app, nil
}

// Run starts application.
func (app *Application) Run() error {
	app.log.Infof("Start listening at %s", app.addr)
	if err := http.ListenAndServe(app.addr, app.router); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	return nil
}

func healthz(w http.ResponseWriter, r *http.Request) {
	// nolint: errcheck,gosec
	w.Write([]byte("Ok"))
}
