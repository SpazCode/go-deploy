package controllers

import (
	"fmt"
    "time"
    "bytes"
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
        store *sessions.CookieStore
	}
)

type (
    Login struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
)

func NewUserController(s *mgo.Session, st *sessions.CookieStore) *UserController {
    return &UserController{s, st}
}

// GetUsers retrieves all users
func (uc UserController) GetUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    // Create user array
    us := []models.User{}
    // Fetch users
    if err := uc.session.DB("go-deploy").C("users").Find(nil).All(&us)
    err != nil {
        w.WriteHeader(404)
        return
    }

    // Marshal provided interface into JSON structure
    uj, _ := json.Marshal(us)

    // Write content-type, statuscode, payload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", uj)
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

    err := uc.session.DB("go-deploy").C("users").Find(bson.M{"username": u.Username}).One(&u)
    if err == nil {
        w.WriteHeader(409)
        fmt.Fprintf(w, "This username already exists")
        return
    }

    err = uc.session.DB("go-deploy").C("users").Find(bson.M{"email": u.Email}).One(&u)
    if err == nil {
        w.WriteHeader(409)
        fmt.Fprintf(w, "This email already exists")
        return
    }

    // Add an Id
    u.Id = bson.NewObjectId()

    // Set the password
    u.SetPassword()
    // Set the Created and Updated values
    u.Created = time.Now()
    u.Updated = time.Now()
    u.LastLogin = time.Now()

    // Write the user to mongo
    uc.session.DB("go-deploy").C("users").Insert(u)

    // Marshal provided interface into JSON structure    
    uj, _ := json.Marshal(u)

    // Write content-type, statuscode, payload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(201)
    fmt.Fprintf(w, "%s", uj)
}

func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {  
    // User Object for the body
    u := models.User{}

    // Get the user's id from the parameters
    id := p.ByName("id")

    // Populae user data that we want to update
    json.NewDecoder(r.Body).Decode(&u)

    if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        fmt.Fprintf(w, "User ID is invalid")
        return
    }

    // Update the updated timestamp
    u.Updated = time.Now()

    if len(u.Password) > 0 {
        u.SetPassword()
    }

    // Save the updated user object to the db
    err := uc.session.DB("go-deploy").C("users").UpdateId(id, u)
    if err != nil {
        w.WriteHeader(400)
        fmt.Fprintf(w, "Error updating the User")
        return
    }

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

    // Compare the user 
    err = bcrypt.CompareHashAndPassword([]byte(u.Password_Hash), []byte(l.Password))
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(401)
        fmt.Fprintf(w, "Passwords do not match")
        return
    }

    sess, err := uc.store.Get(r, "GODEPLOYSESSION")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    // Put the user into the session
    sess.Values["logged_in"] = true
    sess.Values["user_id"] = u.Id.Hex()
    sess.Save(r, w)

    // Marshal provided interface into JSON structure
    uj, _ := json.Marshal(u)

    // Write content-type, statuscode, payload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", uj)
    var buffer bytes.Buffer
    buffer.WriteString(u.Id.String())
    buffer.WriteString(" Logged In")
    fmt.Println(buffer.String())
}


// Logout clears the user session 
func (uc UserController) Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    sess, err := uc.store.Get(r, "GODEPLOYSESSION")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    // Put the user into the session
    sess.Values["logged_in"] = false
    sess.Values["user_id"] = ""
    sess.Save(r, w)
}