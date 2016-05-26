package controllers

import (
	"fmt"
	"encoding/json"
	"net/http"
	"../models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/julienschmidt/httprouter"
    "golang.org/x/crypto/bcrypt"
    "github.com/gorilla/sessions"
)

type (
	UserController struct {
		session *mgo.Session
	}
)

type (
    Login struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
)

func NewUserController(s *mgo.Session) *UserController {  
    return &UserController{s}
}

// GetUser retrieves an individual user resource
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {  
    // Get the user's id from the parameters
	id := p.ByName("id")

	// Check to ensure that this is a BSON object
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab the object's id
	oid := bson.ObjectIdHex(id)

    // Stub an example user
    u := models.User{}

    // Fetch the user data
    if err := uc.session.DB("go-deploy").C("users").FindId(oid).One(&u); err != nil {
    	w.WriteHeader(404)
    	return
    }
    
    // Marshal provided interface into JSON structure
    uj, _ := json.Marshal(u)

    // Write content-type, statuscode, payload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", uj)
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {  
    // Stub an user to be populated from the body
    u := models.User{}

    // Populate the user data
    json.NewDecoder(r.Body).Decode(&u)

    // Add an Id
    u.Id = bson.NewObjectId()


    // Set the password
    u.SetPassword()

    // Write the user to mongo
    uc.session.DB("go-deploy").C("users").Insert(u)

    // Marshal provided interface into JSON structure
    uj, _ := json.Marshal(u)

    // Write content-type, statuscode, payload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(201)
    fmt.Fprintf(w, "%s", uj)
}

// RemoveUser removes an existing user resource
func (uc UserController) RemoveUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {  
   	
   	// Get the user's id from the parameters
	id := p.ByName("id")

	// Check to ensure that this is a BSON object
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab the object's id
	oid := bson.ObjectIdHex(id)
	
    // Remove user
    if err := uc.session.DB("go-deploy").C("users").RemoveId(oid); err != nil {
        w.WriteHeader(404)
        return
    }

    // Write Status
    w.WriteHeader(200)
}

//Login validates and returns a user object if they exist in the database.
func (uc UserController) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    // Stub an user to be populated from the body
    u := models.User{}
    l := Login{}

    // Populate the user data
    json.NewDecoder(r.Body).Decode(&l)

    // Finding the user in the database
    err := uc.session.DB("go-deploy").C("users").Find(bson.M{"username": l.Username}).One(&u)
    if err != nil {
        w.WriteHeader(401)
        fmt.Fprintf(w, "User is not found")
        return
    }

    print(l.Password)

    // Compare the user 
    err = bcrypt.CompareHashAndPassword(u.Password_Hash, []byte(l.Password))
    fmt.Println(err)
    if err != nil {
        w.WriteHeader(401)
        fmt.Fprintf(w, "Passwords do not match")
        return
    }

    sess, err := store.Get(req, "go-deploy")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    // Put the user into the session
    sess.Values["logged_in"] = true
    sess.Values["user_id"] = u.Id
    sess.Save(r, w)

    // Marshal provided interface into JSON structure
    uj, _ := json.Marshal(u)

    // Write content-type, statuscode, payload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", uj)
}