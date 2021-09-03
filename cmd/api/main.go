package main

import(
	"net/http"
	"os"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/kharism/microservice_simple/controller"
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
}
func main() {

}
