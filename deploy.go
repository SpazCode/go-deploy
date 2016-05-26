package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"./controllers"
)

// Database handlers
func getSession() *mgo.Session {  
    // Connect to our local mongo
    s, err := mgo.Dial("mongodb://localhost")

    // Check if connection error, is mongo running?
    if err != nil {
        panic(err)
    }
    return s
}

func rootHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "Welcome to the Deployment APIs")
}

func main() {
	// Create a new Router
	r := httprouter.New()
	r.GET("/", rootHandler)

	// User Routes
	uc := controllers.NewUserController(getSession())

	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.RemoveUser) 
	r.POST("/login", uc.Login)
	// Listen for the port with this router
	http.ListenAndServe(":9000", r)
}
