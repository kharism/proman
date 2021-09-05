package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/kharism/proman/controller"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	debugging bool
	token     *jwtauth.JWTAuth
)
var (
	authAPI controller.IAuthRestAPI
	promAPI controller.IPromRestAPI
)

func init() {
	token = jwtauth.New("HS256", []byte("somethingSecret"), nil)

	viper.SetConfigName("api")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../../config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	debugging = viper.GetBool(`debug`)
	mongourl := os.Getenv("MONGO_URI")
	if mongourl != "" {
		viper.Set("uri", mongourl)
	}
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   debugging,
		FullTimestamp: true,
	})
	authAPI = controller.NewAuth(token)
	promAPI = controller.NewProm(token)
}
func main() {
	r := chi.NewRouter()

	logger := logrus.New()
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:   debugging,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.Info("server api run on DEBUG mode")
	log.SetLevel(log.DebugLevel)

	r.Use(chiMiddleware.RequestID)

	r.Use(chiMiddleware.Recoverer)

	// disable cache control
	r.Use(chiMiddleware.NoCache)

	// apply gzip compression
	r.Use(chiMiddleware.Compress(5, "gzip"))

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)
	r.Group(func(r chi.Router) {
		r.Mount("/auth", authAPI.Register())
		r.Mount("/prom", promAPI.Register())

	})
	r.Get("/metrics", metric)
	staticRoute(r, "/js", "js")
	staticRoute(r, "/img", "img")
	staticRoute(r, "/css", "css")
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte("welcome"))
		content, err := ioutil.ReadFile("./index.html")
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write(content)
	})
	r.Get("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte("welcome"))
		content, err := ioutil.ReadFile("./manifest.json")
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write(content)
	})
	r.Get("/service-worker.js", func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte("welcome"))
		content, err := ioutil.ReadFile("./service-worker.js")
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Header().Add("content-type", "javscript")
		w.Write(content)
	})
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Debugf("[%s] %s", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.Errorf("walk function error : %s\n", err.Error())
	}

	serverAddress := viper.GetString("address")
	log.Infof("server api run at %s", serverAddress)
	err := http.ListenAndServe(serverAddress, r)
	if err != nil {
		log.Fatal("unable to start web server", err.Error())
	}
}
func metric(w http.ResponseWriter, r *http.Request) {
	returnDummy := []byte("my_dummy_metric 100")
	w.Write(returnDummy)
}

// StaticRoute add static route to chi router
func staticRoute(r chi.Router, routePath string, directoryPath string) {
	if strings.ContainsAny(routePath, "{}*") {
		panic("static route does not permit any URL parameters.")
	}

	if routePath != "/" && routePath[len(routePath)-1] != '/' {
		r.Get(routePath, http.RedirectHandler(routePath+"/", 301).ServeHTTP)
		routePath += "/"
	}
	routePath += "*"

	root := http.Dir(directoryPath)
	r.Get(routePath, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
