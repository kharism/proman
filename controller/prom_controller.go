package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"sync"

	"github.com/eaciit/toolkit"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/kharism/proman/util"
	"github.com/spf13/viper"
)

// IAuthRestAPI user controller interface
// RegisterUser -> the register user request handler
// VerifyLogin  -> this is the one handling login request
type IPromRestAPI interface {
	GetYaml(w http.ResponseWriter, r *http.Request)
	SaveYaml(w http.ResponseWriter, r *http.Request)
	Register() http.Handler
}

// authentication controller struct
// simple implementation with mongodb backend as prototyping backend
type promRestController struct {
	tokenAuth *jwtauth.JWTAuth
	lock      *sync.Mutex
	//auth func() service.IAuth
	//rkas      func() service.IRKAS
}

// create new auth controller
func NewProm(token *jwtauth.JWTAuth) IPromRestAPI {
	return &promRestController{
		tokenAuth: token,
		lock:      new(sync.Mutex),
	}
}

type Yaml struct {
	Content string
}

func (c *promRestController) GetYaml(w http.ResponseWriter, r *http.Request) {
	yamlPath := viper.GetString("prometheus_yaml")

	content, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	yy := Yaml{}
	yy.Content = string(content)
	util.WriteJSONData(w, yy)
}

type ymlPayload struct {
	Content string
}

// handler for login. On success login return the userdata with JWT Token.
// Probably need to remove unneccessary field due to security reason
// payload :
// {
//	"Username":"string",
//  "Password":"string",
// }
// response: json representation of the User struct with Token field filled with JWTAuth token
// Some service in the microserve will require JWT Token as authentication method
// use the token as HTTP Header "Authorization" with value "BEARER "+token
// for more example refer to TestItems functions on auth_controller_test.go
func (c *promRestController) SaveYaml(w http.ResponseWriter, r *http.Request) {

	data := ymlPayload{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		util.WriteJSONError(w, err)
		return
	}
	//log.Println(data)
	c.lock.Lock()
	defer c.lock.Unlock()
	yamlPath := viper.GetString("prometheus_yaml")
	content, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	err = ioutil.WriteFile(yamlPath, []byte(data.Content), 0644)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	//try reload data
	prometheusReloadUrl := viper.GetString("prometheus_path")
	resp, err := http.Post(prometheusReloadUrl, "application/json", strings.NewReader("{}"))
	if err != nil || resp.StatusCode != 200 {
		//rollback
		ioutil.WriteFile(yamlPath, []byte(content), 0644)
	}
	yy := toolkit.M{}
	util.WriteJSONData(w, yy)
}

func (c *promRestController) Register() http.Handler {
	r := chi.NewRouter()
	log.Println("TokenAuth", c.tokenAuth)
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(c.tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/getyaml", c.GetYaml)
		r.Post("/pingprom", pingProm)
		r.Post("/saveyaml", c.SaveYaml)
		//r.Post("/", c.SaveItem)
		//r.Put("/{id}", c.SaveItem)
		//r.Delete("/{id}", c.HideItem)
	})

	return r
}
func pingProm(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(200)
	w.Write([]byte("pongprom"))
}
