package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/security"

	_ "github.com/dk-if/scr-smv/back/migrations"
)

const ROLE_NORMAL = "normal"
const ROLE_SUPER = "super"
const COLLECTION_APP_USERS = "app_users"
const COLLECTION_SURVEY = "surveys"
const COLLECTION_TOKENS = "tokens"
const COLLECTION_SESSIONS = "sessions"

const FIELD_PWORD = "pwordtext"

var migratecmdConfig migratecmd.Config

func init() {
	migratecmdConfig = migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: false,
	}
}

func serveStaticFile(e *core.ServeEvent) error {
	e.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), true))
	return nil
}

type AuthSetPostBody struct {
	Password string `json:"password"`
}

func serveSuperSetPassword(se *core.ServeEvent) {
	se.Router.GET("/auth-set", func(e *core.RequestEvent) error {
		record := e.Auth

		if record == nil {
			return apis.NewBadRequestError("unable to fetch auth record", nil)
		}

		if role := record.GetString("role"); role != "super" {
			return e.HTML(http.StatusUnauthorized, "unauthorize")
		}

		if ptext := record.GetString(FIELD_PWORD); ptext == "" {
			return e.HTML(http.StatusOK, "true")
		} else {
			return e.HTML(http.StatusOK, "false")
		}
	})

	se.Router.POST("/auth-set", func(e *core.RequestEvent) error {
		record := e.Auth

		if record == nil {
			return apis.NewBadRequestError("unable to fetch auth record", nil)
		}

		if role := record.GetString("role"); role != "super" {
			return e.HTML(http.StatusUnauthorized, "unauthorize")
		}

		body := new(AuthSetPostBody)
		if err := e.BindBody(&body); err != nil {
			return e.HTML(http.StatusBadRequest, "bad request")
		}

		record.SetPassword(body.Password)

		record.Set(FIELD_PWORD, "")
		if err := e.App.Save(record); err != nil {
			return err
		}
		return e.Next()
	})
}

func serveTokenRoutes(se *core.ServeEvent) {
	se.Router.POST("/token", func(e *core.RequestEvent) error {
		record := e.Auth
		if record == nil {
			return apis.NewBadRequestError("unable to fetch auth record", nil)
		}

		if role := record.GetString("role"); role != "super" {
			return e.HTML(http.StatusUnauthorized, "unauthorize")
		}

		user := e.Request.URL.Query().Get("u")
		if user == "" {
			return e.HTML(http.StatusBadRequest, "bad request")
		}
		tokensCollection, err := e.App.FindCollectionByNameOrId(COLLECTION_TOKENS)
		if err != nil {
			log.Fatalln("missing '" + COLLECTION_TOKENS + "' collection")
		}
		tokenRecord := core.NewRecord(tokensCollection)
		tokenRecord.Load(newTokenData(user))
		if err := e.App.Save(tokenRecord); err != nil {
			return err
		}
		return e.JSON(http.StatusCreated, tokenRecord)
	})
}

func loginWithToken(app core.App) func(e *core.RecordAuthWithPasswordRequestEvent) error {
	return func(e *core.RecordAuthWithPasswordRequestEvent) error {
		r, err := app.FindFirstRecordByFilter(
			COLLECTION_TOKENS,
			"user.username={:user} && token={:token} && valid_until > @now",
			dbx.Params{"user": e.Identity, "token": e.Password},
		)
		if err != nil {
			log.Printf("[FindFirstRecordByFilter] failed to fetch: %v", err)
			return e.Next()
		}
		if r != nil {
			if errs := app.ExpandRecord(r, []string{"user"}, nil); len(errs) > 0 {
				return fmt.Errorf("[FindFirstRecordByFilter] failed to expand: %v", errs)
			}
			if err := apis.RecordAuthResponse(e.RequestEvent, r.ExpandedOne("user"), "password", nil); err != nil {
				return err
			}
		}
		return e.Next()
	}
}

func generateUsersOnSessionCreated(app core.App) func(e *core.RecordRequestEvent) error {
	return func(e *core.RecordRequestEvent) error {

		if err := e.Next(); err != nil {
			return err
		}

		userCollection, err := app.FindCollectionByNameOrId(COLLECTION_APP_USERS)
		if err != nil {
			log.Fatalln("missing '" + COLLECTION_APP_USERS + "' collection")
		}
		tokensCollection, err := app.FindCollectionByNameOrId(COLLECTION_TOKENS)
		if err != nil {
			log.Fatalln("missing '" + COLLECTION_TOKENS + "' collection")
		}

		sessionId := e.Record.Id
		sessionName := e.Record.GetString("name")

		for i := 0; i < 10; i++ {
			record := core.NewRecord(userCollection)
			record.Load(newUserData(sessionId, ROLE_NORMAL))
			record.Set("username", strings.ReplaceAll(fmt.Sprintf("%s-n%d", sessionName, i+1), " ", "-"))
			password := security.RandomString(5)
			record.Set(FIELD_PWORD, password)
			record.SetPassword(password)
			record.SetVerified(true)
			if err := app.Save(record); err != nil {
				return err
			}

			tokenRecord := core.NewRecord(tokensCollection)
			tokenRecord.Load(newTokenData(record.Id))
			if err := app.Save(tokenRecord); err != nil {
				return err
			}
		}
		record := core.NewRecord(userCollection)
		record.Load(newUserData(sessionId, ROLE_SUPER))
		record.Set("username", strings.ReplaceAll(sessionName+"-s", " ", "-"))
		password := security.RandomString(5)
		record.Set(FIELD_PWORD, password)
		record.SetPassword(password)
		record.SetVerified(true)
		if err := app.Save(record); err != nil {
			return err
		}
		return nil
	}
}

func increaseTokenTimeAfterSubmitSurvey(app core.App) func(e *core.RecordRequestEvent) error {
	return func(e *core.RecordRequestEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		record := e.Auth
		username := record.GetString("username")
		if username == "" {
			return fmt.Errorf("unable to fetch auth record")
		}

		rs, err := app.FindRecordsByFilter(
			COLLECTION_TOKENS,
			"user.username={:user} && valid_until > @now",
			"", // sort
			-1, 0,
			dbx.Params{"user": username},
		)

		if err != nil {
			log.Printf("[FindRecordsByFilter] failed to fetch: %v", err)
			return nil
		}

		for _, r := range rs {
			r.Set("valid_until", r.GetDateTime("valid_until").Time().AddDate(0, 0, 3))
			if err := app.Save(r); err != nil {
				log.Printf("[Save] failed to save: %v", err)
			}
		}

		return nil
	}
}

func logRequest(e *core.RequestEvent) error {
	start := time.Now()

	// Get request details
	method := e.Request.Method
	url := e.Request.URL.String()
	userAgent := e.Request.UserAgent()
	realIP := e.RealIP() // Use PocketBase's built-in RealIP method

	// Get auth information if available
	var authInfo string
	if e.Auth != nil {
		username := e.Auth.GetString("username")
		role := e.Auth.GetString("role")
		authInfo = fmt.Sprintf("user:%s,role:%s", username, role)
	} else {
		authInfo = "guest"
	}

	// Log the incoming request
	e.App.Logger().Info("Incoming request",
		"method", method,
		"url", url,
		"real_ip", realIP,
		"user_agent", userAgent,
		"auth", authInfo,
		"timestamp", start.Format(time.RFC3339),
	)

	// Process the request
	err := e.Next()

	// Log request completion with duration and status
	duration := time.Since(start)
	status := "success"
	if err != nil {
		status = fmt.Sprintf("error: %s", err.Error())
	}

	e.App.Logger().Info("Request completed",
		"method", method,
		"url", url,
		"real_ip", realIP,
		"auth", authInfo,
		"duration_ms", duration.Milliseconds(),
		"status", status,
		"timestamp", time.Now().Format(time.RFC3339),
	)

	return err
}
func bindAppHooks(app core.App) {
	app.OnRecordAuthWithPasswordRequest(COLLECTION_APP_USERS).BindFunc(loginWithToken(app))

	app.OnRecordCreateRequest(COLLECTION_SESSIONS).BindFunc(generateUsersOnSessionCreated(app))
	app.OnRecordCreateRequest(COLLECTION_SURVEY).BindFunc(increaseTokenTimeAfterSubmitSurvey(app))

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.BindFunc(logRequest)

		serveSuperSetPassword(se)
		serveTokenRoutes(se)

		// serves static files from the provided public dir (if exists)
		serveStaticFile(se)
		return se.Next()
	})

	app.OnRecordEnrich(COLLECTION_APP_USERS).BindFunc(func(e *core.RecordEnrichEvent) error {
		if e.RequestInfo.Auth == nil || !e.RequestInfo.Auth.IsSuperuser() {
			e.Record.Hide(FIELD_PWORD, "")
		}
		return nil
	})
}

func main() {
	app := pocketbase.New()

	migratecmd.MustRegister(app, app.RootCmd, migratecmdConfig)

	bindAppHooks(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func newUserData(sessionId string, role string) map[string]any {
	return map[string]any{
		"role":    role,
		"status":  "ready",
		"session": sessionId,
	}
}

func newTokenData(userId string) map[string]any {
	now := time.Now()
	nextWeek := now.AddDate(0, 1, 0)
	return map[string]any{
		"user":        userId,
		"token":       security.RandomString(15),
		"valid_until": nextWeek,
	}
}
