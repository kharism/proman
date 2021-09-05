package controller

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/eaciit/toolkit"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	db "github.com/kharism/proman/connection"
	"github.com/kharism/proman/model"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func init() {
	token := jwtauth.New("HS256", []byte("secretTest"), nil)
	viper.SetConfigName("api_test")
	viper.SetConfigType("json")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	//start webservice
	authAPI := NewAuth(token)
	promAPI := NewProm(token)
	//itemAPI := NewItem(token)
	//cartAPI := NewCart(token)
	r := chi.NewRouter()

	logger := logrus.New()
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
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
		r.Get("/ping", ping)
		//r.Use(jwtauth.Verifier(token))
		r.Mount("/auth", authAPI.Register())
		r.Mount("/prom", promAPI.Register())
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

	go func() {
		err = http.ListenAndServe(serverAddress, r)
		if err != nil {
			log.Fatal("unable to start web server", err.Error())
		}
	}()
}
func ToStringReader(payload toolkit.M) io.Reader {
	return strings.NewReader(toolkit.JsonString(payload))
}
func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("pong"))
}
func TestPingProm(t *testing.T) {
	userModel := model.User{}

	Convey("Clean up", t, func() {
		client := &http.Client{}
		cli1, err := db.NewClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = cli1.Database(viper.GetString("db")).Collection(userModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)
		payload := toolkit.M{}

		//create basic user
		payload["Username"] = "admin"
		payload["Password"] = "PasswordXX"
		payloadReader := strings.NewReader(toolkit.JsonString(payload))

		resp, err := client.Post("http://localhost:8098/auth/registeruser", "application/json", payloadReader)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode, ShouldEqual, 200)

		//test pong failed
		payloadReader = strings.NewReader("{}")
		resp, err = client.Post("http://localhost:8098/prom/pingprom", "application/json", payloadReader)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode, ShouldNotEqual, 200)
		//respM, err := ProcessResponse(resp)
		//So(err, ShouldBeNil)
		//log.Println(respM)

		//test pong success
		payloadReader = strings.NewReader(toolkit.JsonString(payload))

		resp, err = client.Post("http://localhost:8098/auth", "application/json", payloadReader)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode, ShouldEqual, 200)

		responseJson, err := ProcessResponse(resp)
		So(err, ShouldBeNil)
		So(responseJson, ShouldNotBeNil)
		//t.Log(responseJson)
		token := responseJson["Data"].(map[string]interface{})["Token"].(string)
		payloadReader = strings.NewReader("{}")
		req, err := http.NewRequest("POST", "http://localhost:8098/prom/pingprom", payloadReader)
		So(err, ShouldBeNil)
		req.Header.Add("Authorization", "BEARER "+token)
		//AddRequestBearer(req, token)
		resp, err = client.Do(req)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		//responseJson, err = ProcessResponse(resp)
		//So(err, ShouldBeNil)
		//t.Log(responseJson)
		So(resp.StatusCode, ShouldEqual, 200)
	})
}
func ProcessResponse(resp *http.Response) (toolkit.M, error) {
	responseJson := toolkit.M{}
	contentByte, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return responseJson, err
	}
	err = json.Unmarshal(contentByte, &responseJson)
	if err != nil {
		return responseJson, err
	}
	return responseJson, nil
}
