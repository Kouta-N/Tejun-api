package routes

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
	"golang.org/x/crypto/bcrypt"
)

var sessionKey string = os.Getenv("SESSION_KEY")

func Signup(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	u := new(models.User)
	if errs := binding.Bind(req, u); errs != nil {
		log.Println(errs.Error())
		return
	}
	u.Email = strings.TrimSpace(u.Email)

	// locale := req.Context().Value("locale").(string)

	err := u.GetByEmail(u.Email)

	if err == nil {
		return
	}

	// u.SetLocale(locale)

	user := models.CreateUser(*u)

	session := sessions.GetSession(req)
	err = handlers.AuthenticateSession(session, &user)
	if err != nil {
		helpers.WriteErr(w, http.StatusUnauthorized, err)
		return
	}
}

func executeLogin(w http.ResponseWriter, req *http.Request) {
	var email, password string

	if req.Method == http.MethodPost {
		lf := new(models.LoginPostForm)
		if errs := binding.Bind(req, lf); errs != nil {
			log.Println(errs.Error())
			return
		}

		email = lf.Email
		password = lf.Password
	} else {
		qs := req.URL.Query()
		email = qs.Get("email")
		password = qs.Get("password")
	}

	// check empty for avoiding unnecessary DB query

	user := new(models.User)
	err := user.LoginByEmail(email, password)

	session := sessions.GetSession(req)
	err = handlers.AuthenticateSession(session, user)
	if err != nil {
		helpers.WriteErr(w, http.StatusUnauthorized, err)
		return
	}
}

func Logout(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	session := sessions.GetSession(req)
	user := req.Context().Value("user").(handlers.User)
	handlers.Logout(session, user)
	session.Delete("user")
}

func CreateUser(user User) User {
	now := time.Now()
	user.CreatedAt = &now
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)

	if user.Name == "" {
		// let name be the local part of email
		s := strings.Split(user.Email, "@")
		user.Name = s[0]
	}

	err := DbMap(true).Insert(&user)
	if err != nil {
		if helpers.HasDuplicateEntry(err) {
			// Return the user that already exists
			// if adding a user violates the unique constraint.
			user.GetByEmail(user.Email, true)
			return user
		}
		helpers.CheckErr(err, "Insert user failed")
	}

	helpers.ExecuteInBackground(func() {
		SubscribeToMailChimp(user)
	})

	return user
}

func AuthenticateSession(s sessions.Session, user entity.AuthorizedUser) {
	s.Set(sessionKey, user)
}