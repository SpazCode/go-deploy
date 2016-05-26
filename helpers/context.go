import (
	"github.com/gorilla/sessions"
)

type Context struct {
	Database *mgo.Database
	Session *sessions.Session
}

func NewContext(req *http.Request, session *mgo.Database) (*Context, error) {
	sess, err := store.Get(req, "gostbook")
	return &Context{
		Database: session.Clone().DB(database),
		Session:  sess,
	}, err
}