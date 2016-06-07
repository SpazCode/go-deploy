package models

import (
  "time"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type (  
    User struct {
        Created         time.Time `json:"created_at" bson:"created_at"`
        Updated         time.Time `json:"updated_at" bson:"updated_at"`
        LastLogin       time.Time `json:"last_login" bson:"last_login"`
        Username        string `json:"username" bson:"username"`
        Password        string `json:"password,omitempty" bson:"-"`
        Password_Hash   string `json:"-" bson:"password_hash"`
        Email           string `json:"email" bson:"email"`
        Id              bson.ObjectId `json:"id" bson:"_id"`
    }
)

//SetPassword takes a plaintext password and hashes it with bcrypt and sets the
//password field to the hash.
func (u *User) SetPassword() {
  hpass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
  if err != nil {
      panic(err) //this is a panic because bcrypt errors on invalid costs
  }
  u.Password = ""
  u.Password_Hash = string(hpass)
}