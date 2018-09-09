package application

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/auth"
	"github.com/tetafro/nott-backend-go/internal/folders"
	"github.com/tetafro/nott-backend-go/internal/notepads"
	"github.com/tetafro/nott-backend-go/internal/notes"
)

const version = "v1"

// Application is an entity that gets things done.
type Application struct {
	addr   string
	router *chi.Mux
	log    logrus.FieldLogger
}

// New creates main application instance that handles all requests.
func New(db *gorm.DB, addr string, log logrus.FieldLogger) (*Application, error) {
	app := &Application{addr: addr, log: log}

	foldersRepo := folders.NewPostgresRepo(db)
	foldersController := folders.NewController(foldersRepo, log)

	notepadsRepo := notepads.NewPostgresRepo(db)
	notepadsController := notepads.NewController(notepadsRepo, log)

	notesRepo := notes.NewPostgresRepo(db)
	notesController := notes.NewController(notesRepo, log)

	usersRepo := auth.NewUsersPostgresRepo(db)

	tokensRepo := auth.NewTokensPostgresRepo(db)

	usersController := auth.NewController(usersRepo, tokensRepo, log)

	mwAuth := auth.NewAuthMiddleware(usersRepo, log)
	mwLog := middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log})

	r := chi.NewRouter()
	// Auth
	r.Method(http.MethodPost, "/login",
		http.HandlerFunc(usersController.Login))
	r.Method(http.MethodPost, "/logout",
		http.HandlerFunc(usersController.Logout))
	r.Method(http.MethodGet, "/profile",
		mwAuth(http.HandlerFunc(usersController.GetProfile)))
	r.Method(http.MethodPut, "/profile",
		mwAuth(http.HandlerFunc(usersController.UpdateProfile)))
	// Folders
	r.Method(http.MethodGet, "/folders",
		mwAuth(http.HandlerFunc(foldersController.GetList)))
	r.Method(http.MethodPost, "/folders",
		mwAuth(http.HandlerFunc(foldersController.Create)))
	r.Method(http.MethodGet, "/folders/{id}",
		mwAuth(http.HandlerFunc(foldersController.GetOne)))
	r.Method(http.MethodPut, "/folders/{id}",
		mwAuth(http.HandlerFunc(foldersController.Update)))
	r.Method(http.MethodDelete, "/folders/{id}",
		mwAuth(http.HandlerFunc(foldersController.Delete)))
	// Notepads
	r.Method(http.MethodGet, "/notepads",
		mwAuth(http.HandlerFunc(notepadsController.GetList)))
	r.Method(http.MethodPost, "/notepads",
		mwAuth(http.HandlerFunc(notepadsController.Create)))
	r.Method(http.MethodGet, "/notepads/{id}",
		mwAuth(http.HandlerFunc(notepadsController.GetOne)))
	r.Method(http.MethodPut, "/notepads/{id}",
		mwAuth(http.HandlerFunc(notepadsController.Update)))
	r.Method(http.MethodDelete, "/notepads/{id}",
		mwAuth(http.HandlerFunc(notepadsController.Delete)))
	// Notes
	r.Method(http.MethodGet, "/notes",
		mwAuth(http.HandlerFunc(notesController.GetList)))
	r.Method(http.MethodPost, "/notes",
		mwAuth(http.HandlerFunc(notesController.Create)))
	r.Method(http.MethodGet, "/notes/{id}",
		mwAuth(http.HandlerFunc(notesController.GetOne)))
	r.Method(http.MethodPut, "/notes/{id}",
		mwAuth(http.HandlerFunc(notesController.Update)))
	r.Method(http.MethodDelete, "/notes/{id}",
		mwAuth(http.HandlerFunc(notesController.Delete)))

	app.router = chi.NewRouter()
	app.router.Method(http.MethodGet, "/healthz",
		http.HandlerFunc(healthz))
	app.router.Mount("/api/"+version, mwLog(r))

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
