package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"tmhub/admin/helpers"
	"tmhub/admin/models/entity"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

var (
	sessionKey string = os.Getenv("SESSION_KEY")

	ServiceAccountFilePath = "./tmhub-management-portal-ce3bd21320d9.json"
)

func AuthenticateSession(s sessions.Session, user entity.AuthorizedUser) {
	s.Set(sessionKey, user)
}

func Logout(s sessions.Session) {
	s.Delete(sessionKey)
}

func LoginRequired(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

		s := sessions.GetSession(req)
		sessionUser := s.Get(sessionKey)

		if sessionUser != nil {
			ctx := req.Context()
			ctx = context.WithValue(ctx, "sessionUser", sessionUser)
			h(w, req.WithContext(ctx), ps)
		} else {
			t, _ := template.ParseFiles("templates/logout.html")
			t.Execute(w, map[string]interface{}{})
		}
	}
}

// Check if sessionUser is a member of tmhub-console group
func MemberRequired(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		sessionUser := req.Context().Value("sessionUser").(entity.AuthorizedUser)
		memberKey := sessionUser.Email
		groupKey := os.Getenv("GROUP_KEY")

		// To access the Admin SDK Directory API,
		// the service account needs to impersonate one of the users who can access the Admin API.
		adminEmail := os.Getenv("ADMIN_EMAIL")
		membersService, err := CreateMembersService(adminEmail)
		if err != nil {
			http.Error(w, "Failed to create directory service: "+err.Error(), http.StatusForbidden)
			return
		}

		getMember := membersService.Get(groupKey, memberKey)
		member, err := getMember.Do()

		if err != nil {
			r := helpers.NewRender(req)
			r.HTML(w, http.StatusNotFound, "error", map[string]interface{}{
				"Error": err,
			})
		} else {
			sessionUser.Role = member.Role

			ctx := req.Context()
			ctx = context.WithValue(ctx, "sessionUser", sessionUser)
			h(w, req.WithContext(ctx), ps)
		}
	}
}

// Instantiate an Admin SDK Directory service object.
func CreateMembersService(userEmail string) (*admin.MembersService, error) {
	ctx := context.Background()

	jsonCredentials, err := ioutil.ReadFile(ServiceAccountFilePath)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: %v", err)
	}

	config, err := google.JWTConfigFromJSON(jsonCredentials, admin.AdminDirectoryGroupReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("JWTConfigFromJSON: %v", err)
	}
	config.Subject = userEmail

	ts := config.TokenSource(ctx)

	service, err := admin.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("NewService: %v", err)
	}

	return service.Members, nil
}
