package jasmine

import (
	"git.sr.ht/~sircmpwn/tokidoki/storage"
	"github.com/1f349/cardcaldav"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"net/http"
	"path/filepath"
	"time"
)

type Conf struct {
	Listen string `json:"listen"`
	DB     string `json:"db"`
}

type jasmineHandler struct {
	auth    AuthProvider
	backend caldav.Backend
}

func (j *jasmineHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	principlePath, err := j.auth.CurrentUserPrincipal(req.Context())
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var homeSets []webdav.BackendSuppliedHomeSet
	path, err := j.backend.CalendarHomeSetPath(req.Context())
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	homeSets = append(homeSets, caldav.NewCalendarHomeSet(path))

	if req.URL.Path == principlePath {
		opts := webdav.ServePrincipalOptions{
			CurrentUserPrincipalPath: principlePath,
			HomeSets:                 homeSets,
			Capabilities: []webdav.Capability{
				caldav.CapabilityCalendar,
			},
		}
		webdav.ServePrincipal(rw, req, &opts)
		return
	}

	if req.URL.Path == "/" {
		http.Error(rw, "Jasmine API Endpoint", http.StatusOK)
		return
	}

	http.NotFound(rw, req)
}

type AuthProvider interface {
	cardcaldav.ProviderMiddleware
	webdav.UserPrincipalBackend
}

func NewHttpServer(conf Conf, wd string) *http.Server {
	cardcaldav.SetupLogger(Logger)
	principle := NewAuth(conf.DB, Logger)

	calStorage, _, err := storage.NewFilesystem(filepath.Join(wd, "storage"), "/calendar/", "/contacts/", principle)
	if err != nil {
		Logger.Fatal("Failed to load storage backend", "err", err)
	}
	calHandler := &caldav.Handler{Backend: calStorage}

	handler := &jasmineHandler{
		auth:    principle,
		backend: calStorage,
	}

	r := http.NewServeMux()
	r.Handle("/", principle.Middleware(handler))
	r.Handle("GET /health", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		http.Error(rw, "Health OK", http.StatusOK)
	}))
	r.Handle("/.well-known/caldav", principle.Middleware(calHandler))
	r.Handle("/{user}/calendar/", principle.Middleware(calHandler))

	r2 := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t := time.Now()
		r.ServeHTTP(rw, req)
		td := time.Since(t)
		Logger.Debug("Request", "method", req.Method, "url", req.URL.String(), "remote", req.RemoteAddr, "dur", td.String())
	})

	return &http.Server{
		Addr:              conf.Listen,
		Handler:           r2,
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    2500,
	}
}
