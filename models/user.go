package models

import (
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type (  
    User struct {
        Username        string `json:"username" bson:"username"`
        Password        string `json:"-" bson:"-"`
        Password_Hash   []byte `json:"-" bson:"password_hash"`
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
  u.Password_Hash = hpass
}