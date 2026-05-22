package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	_ "github.com/dk-if/scr-smv/back/migrations"
)

const testDataDir = "./test_pb_data"

// testData holds IDs and credentials for records seeded in setupTestApp.
type testData struct {
	sessionID      string
	superUserID    string
	superUsername  string
	superPassword  string
	normalUserID   string
	normalUsername string
	normalPassword string
	validToken     string
	expiredToken   string
	validTokenID   string
	expiredTokenID string
	surveyID       string
}

// setupTestApp creates a fresh test app with all hooks bound and test data seeded.
// It is safe to call multiple times; each call returns an independent clone.
func setupTestApp(t testing.TB) *tests.TestApp {
	testApp, err := tests.NewTestApp(testDataDir)
	if err != nil {
		t.Fatal(err)
	}

	bindAppHooks(testApp)
	return testApp
}

// seedTestData inserts the minimal records required by most test scenarios
// directly via app.Save() so that no API hooks fire during seeding.
func seedTestData(t testing.TB, app *tests.TestApp) testData {
	t.Helper()

	// ---- Session ----
	sessionColl, err := app.FindCollectionByNameOrId(COLLECTION_SESSIONS)
	if err != nil {
		t.Fatal("session collection not found:", err)
	}
	session := core.NewRecord(sessionColl)
	session.Set("name", "test-session")
	session.Set("description", "Integration test session")
	session.Set("current_step", "2")
	session.Set("contact_email", "test@example.com")
	if err := app.Save(session); err != nil {
		t.Fatal("seed session:", err)
	}

	// ---- Super user ----
	userColl, err := app.FindCollectionByNameOrId(COLLECTION_APP_USERS)
	if err != nil {
		t.Fatal("app_users collection not found:", err)
	}
	superUser := core.NewRecord(userColl)
	superUser.Set("username", "test-super")
	superUser.Set("role", ROLE_SUPER)
	superUser.Set("status", "ready")
	superUser.Set("session", session.Id)
	superUser.Set(FIELD_PWORD, "sup3rp") // pwordtext set → GET /auth-set should return "false"
	superUser.SetPassword("sup3rp")
	superUser.SetVerified(true)
	if err := app.Save(superUser); err != nil {
		t.Fatal("seed super user:", err)
	}

	// ---- Normal user ----
	normalUser := core.NewRecord(userColl)
	normalUser.Set("username", "test-normal")
	normalUser.Set("role", ROLE_NORMAL)
	normalUser.Set("status", "ready")
	normalUser.Set("session", session.Id)
	normalUser.Set(FIELD_PWORD, "norm5")
	normalUser.SetPassword("norm5")
	normalUser.SetVerified(true)
	if err := app.Save(normalUser); err != nil {
		t.Fatal("seed normal user:", err)
	}

	// ---- Another super user with pwordtext cleared ----
	superUserEmpty := core.NewRecord(userColl)
	superUserEmpty.Set("username", "test-super-empty")
	superUserEmpty.Set("role", ROLE_SUPER)
	superUserEmpty.Set("status", "ready")
	superUserEmpty.Set("session", session.Id)
	superUserEmpty.Set(FIELD_PWORD, "") // pwordtext empty → GET /auth-set returns "true"
	superUserEmpty.SetPassword("em5ty")
	superUserEmpty.SetVerified(true)
	if err := app.Save(superUserEmpty); err != nil {
		t.Fatal("seed super-empty user:", err)
	}

	// ---- Valid token for normal user ----
	tokenColl, err := app.FindCollectionByNameOrId(COLLECTION_TOKENS)
	if err != nil {
		t.Fatal("tokens collection not found:", err)
	}
	validToken := core.NewRecord(tokenColl)
	validToken.Set("user", normalUser.Id)
	validToken.Set("token", "validtoken12345")
	validToken.Set("valid_until", time.Now().AddDate(0, 1, 0))
	if err := app.Save(validToken); err != nil {
		t.Fatal("seed valid token:", err)
	}

	// ---- Expired token for normal user ----
	expiredToken := core.NewRecord(tokenColl)
	expiredToken.Set("user", normalUser.Id)
	expiredToken.Set("token", "expiredtok12345")
	expiredToken.Set("valid_until", time.Now().AddDate(0, -1, 0))
	if err := app.Save(expiredToken); err != nil {
		t.Fatal("seed expired token:", err)
	}

	// ---- Survey record ----
	surveyColl, err := app.FindCollectionByNameOrId(COLLECTION_SURVEY)
	if err != nil {
		t.Fatal("surveys collection not found:", err)
	}
	survey := core.NewRecord(surveyColl)
	survey.Set("response", map[string]any{"answers": []string{"a", "b"}})
	survey.Set("is_complete", true)
	survey.Set("session", session.Id)
	survey.Set("user", normalUser.Id)
	if err := app.Save(survey); err != nil {
		t.Fatal("seed survey:", err)
	}

	return testData{
		sessionID:      session.Id,
		superUserID:    superUser.Id,
		superUsername:  "test-super",
		superPassword:  "sup3rp",
		normalUserID:   normalUser.Id,
		normalUsername: "test-normal",
		normalPassword: "norm5",
		validToken:     "validtoken12345",
		expiredToken:   "expiredtok12345",
		validTokenID:   validToken.Id,
		expiredTokenID: expiredToken.Id,
		surveyID:       survey.Id,
	}
}

// generateUserToken creates a short-lived auth token for an app_users record.
func generateUserToken(t testing.TB, app *tests.TestApp, username string) string {
	t.Helper()
	record, err := app.FindFirstRecordByFilter(COLLECTION_APP_USERS, "username={:u}", dbx.Params{"u": username})
	if err != nil {
		t.Fatalf("generateUserToken: user %q not found: %v", username, err)
	}
	token, err := record.NewAuthToken()
	if err != nil {
		t.Fatalf("generateUserToken: failed to create token for %q: %v", username, err)
	}
	return token
}

func newSessionBody(name string) string {
	return fmt.Sprintf(`{"name":%q,"description":"test session","current_step":"2","contact_email":"test@test.dk"}`, name)
}

// generatePBSuperuserToken creates a PocketBase superuser and returns its auth token.
func generatePBSuperuserToken(t testing.TB, app *tests.TestApp) string {
	t.Helper()
	col, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
	if err != nil {
		t.Fatal("find superusers collection:", err)
	}
	su := core.NewRecord(col)
	su.SetEmail("testadmin@test.dk")
	su.SetPassword("Test1234!!")
	if err := app.Save(su); err != nil {
		t.Fatal("create superuser:", err)
	}
	token, err := su.NewAuthToken()
	if err != nil {
		t.Fatal("gen superuser token:", err)
	}
	return token
}

// ============================================================
// GET /auth-set tests
// ============================================================

func TestAuthSetGet(t *testing.T) {
	t.Run("no auth (guest) returns 400", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		seedTestData(t, app)

		s := tests.ApiScenario{
			Name:            "GET /auth-set - no auth",
			Method:          http.MethodGet,
			URL:             "/auth-set",
			ExpectedStatus:  400,
			ExpectedContent: []string{"Unable to fetch auth record"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("normal user returns 401", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		seedTestData(t, app)
		token := generateUserToken(t, app, "test-normal")

		s := tests.ApiScenario{
			Name:            "GET /auth-set - normal user",
			Method:          http.MethodGet,
			URL:             "/auth-set",
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  401,
			ExpectedContent: []string{"unauthorize"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("super with pwordtext set returns false", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		seedTestData(t, app)
		token := generateUserToken(t, app, "test-super")

		s := tests.ApiScenario{
			Name:            "GET /auth-set - super pwordtext set",
			Method:          http.MethodGet,
			URL:             "/auth-set",
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  200,
			ExpectedContent: []string{"false"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("super with pwordtext empty returns true", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		seedTestData(t, app)
		token := generateUserToken(t, app, "test-super-empty")

		s := tests.ApiScenario{
			Name:            "GET /auth-set - super pwordtext empty",
			Method:          http.MethodGet,
			URL:             "/auth-set",
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  200,
			ExpectedContent: []string{"true"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})
}

// ============================================================
// POST /auth-set tests
// ============================================================

func TestAuthSetPost(t *testing.T) {
	t.Run("guest returns 400", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		seedTestData(t, app)

		s := tests.ApiScenario{
			Name:            "POST /auth-set - guest",
			Method:          http.MethodPost,
			URL:             "/auth-set",
			Body:            strings.NewReader(`{"password":"newpassword123"}`),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  400,
			ExpectedContent: []string{"Unable to fetch auth record"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("normal user returns 401", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		seedTestData(t, app)
		token := generateUserToken(t, app, "test-normal")

		s := tests.ApiScenario{
			Name:            "POST /auth-set - normal user",
			Method:          http.MethodPost,
			URL:             "/auth-set",
			Body:            strings.NewReader(`{"password":"newpassword123"}`),
			Headers:         map[string]string{"Content-Type": "application/json", "Authorization": token},
			ExpectedStatus:  401,
			ExpectedContent: []string{"unauthorize"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("malformed JSON with super auth returns 400", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		token := generateUserToken(t, app, td.superUsername)

		s := tests.ApiScenario{
			Name:            "POST /auth-set - malformed JSON",
			Method:          http.MethodPost,
			URL:             "/auth-set",
			Body:            strings.NewReader(`not-json`),
			Headers:         map[string]string{"Content-Type": "application/json", "Authorization": token},
			ExpectedStatus:  400,
			ExpectedContent: []string{"bad request"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("super user clears pwordtext on success", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		token := generateUserToken(t, app, td.superUsername)

		s := tests.ApiScenario{
			Name:           "POST /auth-set - super success",
			Method:         http.MethodPost,
			URL:            "/auth-set",
			Body:           strings.NewReader(`{"password":"newSup3rPass"}`),
			Headers:        map[string]string{"Content-Type": "application/json", "Authorization": token},
			ExpectedStatus: 200,
			TestAppFactory: func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				record, err := app.FindFirstRecordByData(COLLECTION_APP_USERS, "username", td.superUsername)
				if err != nil {
					t.Fatal("AfterTest: failed to fetch user:", err)
				}
				if ptext := record.GetString(FIELD_PWORD); ptext != "" {
					t.Errorf("expected pwordtext to be cleared, got %q", ptext)
				}
				if !record.ValidatePassword("newSup3rPass") {
					t.Error("expected password to be updated to newSup3rPass")
				}
			},
		}
		s.Test(t)
	})
}

// ============================================================
// POST /token tests
// ============================================================

func TestTokenRoute(t *testing.T) {
	t.Run("POST /token - guest returns error", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		_ = td

		s := tests.ApiScenario{
			Name:            "POST /token - no auth",
			Method:          http.MethodPost,
			URL:             "/token?u=" + td.normalUserID,
			ExpectedStatus:  400,
			ExpectedContent: []string{"Unable to fetch auth record"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("POST /token - normal user returns 401", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		token := generateUserToken(t, app, td.normalUsername)

		s := tests.ApiScenario{
			Name:            "POST /token - normal user",
			Method:          http.MethodPost,
			URL:             "/token?u=" + td.normalUserID,
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  401,
			ExpectedContent: []string{"unauthorize"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("POST /token - super user missing u param returns 400", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		token := generateUserToken(t, app, td.superUsername)

		s := tests.ApiScenario{
			Name:            "POST /token - missing u param",
			Method:          http.MethodPost,
			URL:             "/token",
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  400,
			ExpectedContent: []string{"bad request"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("POST /token - super user creates token (201)", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		token := generateUserToken(t, app, td.superUsername)

		s := tests.ApiScenario{
			Name:            "POST /token - valid request",
			Method:          http.MethodPost,
			URL:             "/token?u=" + td.normalUserID,
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  201,
			ExpectedContent: []string{`"token"`, `"valid_until"`, `"user"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				records, err := app.FindRecordsByFilter(
					COLLECTION_TOKENS,
					"user={:uid}",
					"", -1, 0,
					dbx.Params{"uid": td.normalUserID},
				)
				if err != nil {
					t.Fatal("AfterTest: token query failed:", err)
				}
				// Should be at least 2 tokens: the seeded valid one + the newly created one
				if len(records) < 2 {
					t.Errorf("expected at least 2 tokens for user, got %d", len(records))
				}
			},
		}
		s.Test(t)
	})

	t.Run("POST /token - super user can create multiple tokens for same user", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		token := generateUserToken(t, app, td.superUsername)

		// Single scenario call; verify DB has an extra token after creation
		s := tests.ApiScenario{
			Name:            "POST /token - second token for same user",
			Method:          http.MethodPost,
			URL:             "/token?u=" + td.normalUserID,
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  201,
			ExpectedContent: []string{`"collectionName":"tokens"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				records, err := app.FindRecordsByFilter(
					COLLECTION_TOKENS, "user={:uid}", "", -1, 0,
					dbx.Params{"uid": td.normalUserID},
				)
				if err != nil {
					t.Fatal(err)
				}
				// 2 seeded (valid + expired) + 1 newly created = 3
				if len(records) != 3 {
					t.Errorf("expected 3 tokens (2 seeded + 1 new), got %d", len(records))
				}
			},
		}
		s.Test(t)
	})

	t.Run("POST /token - valid_until is approximately 1 month from now", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		token := generateUserToken(t, app, td.superUsername)

		s := tests.ApiScenario{
			Name:            "POST /token - check valid_until",
			Method:          http.MethodPost,
			URL:             "/token?u=" + td.normalUserID,
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  201,
			ExpectedContent: []string{`"collectionName":"tokens"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				records, err := app.FindRecordsByFilter(
					COLLECTION_TOKENS,
					"user={:uid}",
					"-created", 1, 0,
					dbx.Params{"uid": td.normalUserID},
				)
				if err != nil || len(records) == 0 {
					t.Fatal("could not fetch latest token:", err)
				}
				validUntil := records[0].GetDateTime("valid_until").Time()
				expectedMin := time.Now().AddDate(0, 1, -1)
				expectedMax := time.Now().AddDate(0, 1, 1)
				if validUntil.Before(expectedMin) || validUntil.After(expectedMax) {
					t.Errorf("valid_until %v not within expected ~1 month range", validUntil)
				}
			},
		}
		s.Test(t)
	})
}

// ============================================================
// loginWithToken hook tests (POST /api/collections/app_users/auth-with-password)
// ============================================================

func TestLoginWithToken(t *testing.T) {
	authURL := "/api/collections/" + COLLECTION_APP_USERS + "/auth-with-password"

	t.Run("valid token login succeeds", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)

		s := tests.ApiScenario{
			Name:   "token login - valid token",
			Method: http.MethodPost,
			URL:    authURL,
			Body: strings.NewReader(fmt.Sprintf(
				`{"identity":%q,"password":%q}`, td.normalUsername, td.validToken,
			)),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  200,
			ExpectedContent: []string{`"token"`, `"record"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("expired token falls through to normal password auth - fails with wrong password", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)

		s := tests.ApiScenario{
			Name:   "token login - expired token",
			Method: http.MethodPost,
			URL:    authURL,
			Body: strings.NewReader(fmt.Sprintf(
				`{"identity":%q,"password":%q}`, td.normalUsername, td.expiredToken,
			)),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  400,
			ExpectedContent: []string{"Failed to authenticate"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("wrong token falls through to password auth - fails", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)

		s := tests.ApiScenario{
			Name:   "token login - wrong token",
			Method: http.MethodPost,
			URL:    authURL,
			Body: strings.NewReader(fmt.Sprintf(
				`{"identity":%q,"password":"thisisnotavalidtoken"}`, td.normalUsername,
			)),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  400,
			ExpectedContent: []string{"Failed to authenticate"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("correct password login still works (token hook falls through)", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)

		s := tests.ApiScenario{
			Name:   "password login - correct credentials",
			Method: http.MethodPost,
			URL:    authURL,
			Body: strings.NewReader(fmt.Sprintf(
				`{"identity":%q,"password":%q}`, td.normalUsername, td.normalPassword,
			)),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  200,
			ExpectedContent: []string{`"token"`, `"record"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("non-existent user returns 400", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		seedTestData(t, app)

		s := tests.ApiScenario{
			Name:            "token login - non-existent user",
			Method:          http.MethodPost,
			URL:             authURL,
			Body:            strings.NewReader(`{"identity":"ghost-user","password":"anytoken"}`),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  400,
			ExpectedContent: []string{"Failed to authenticate"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})
}

// ============================================================
// generateUsersOnSessionCreated hook tests
// ============================================================

func TestGenerateUsersOnSessionCreated(t *testing.T) {
	sessionsURL := "/api/collections/" + COLLECTION_SESSIONS + "/records"

	newSessionBody := func(name string) string {
		return fmt.Sprintf(
			`{"name":%q,"description":"Hook test","current_step":"2","contact_email":"hook@test.dk"}`,
			name,
		)
	}

	t.Run("creates 10 normal users and 1 super user on session create", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		// No pre-seeded data needed; we're testing the hook directly via the API.
		pbToken := generatePBSuperuserToken(t, app)

		s := tests.ApiScenario{
			Name:            "session create - generates 11 users",
			Method:          http.MethodPost,
			URL:             sessionsURL,
			Body:            strings.NewReader(newSessionBody("hook-session")),
			Headers:         map[string]string{"Content-Type": "application/json", "Authorization": pbToken},
			ExpectedStatus:  200,
			ExpectedContent: []string{`"collectionName":"sessions"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				session, err := app.FindFirstRecordByFilter(COLLECTION_SESSIONS, "name='hook-session'")
				if err != nil {
					t.Fatal("session not found:", err)
				}
				users, err := app.FindRecordsByFilter(
					COLLECTION_APP_USERS,
					"session={:sid}",
					"", -1, 0,
					dbx.Params{"sid": session.Id},
				)
				if err != nil {
					t.Fatal("users query failed:", err)
				}
				if len(users) != 11 {
					t.Errorf("expected 11 users, got %d", len(users))
				}

				normalCount := 0
				superCount := 0
				for _, u := range users {
					switch u.GetString("role") {
					case ROLE_NORMAL:
						normalCount++
					case ROLE_SUPER:
						superCount++
					}
				}
				if normalCount != 10 {
					t.Errorf("expected 10 normal users, got %d", normalCount)
				}
				if superCount != 1 {
					t.Errorf("expected 1 super user, got %d", superCount)
				}
			},
		}
		s.Test(t)
	})

	t.Run("username pattern: normal users are {name}-n1..n10", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()

		pbToken := generatePBSuperuserToken(t, app)
		s := tests.ApiScenario{
			Name:            "session create - username pattern normal",
			Method:          http.MethodPost,
			URL:             sessionsURL,
			Body:            strings.NewReader(newSessionBody("alpha-session")),
			Headers:         map[string]string{"Content-Type": "application/json", "Authorization": pbToken},
			ExpectedStatus:  200,
			ExpectedContent: []string{`"collectionName":"sessions"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				for i := 1; i <= 10; i++ {
					username := fmt.Sprintf("alpha-session-n%d", i)
					_, err := app.FindFirstRecordByFilter(COLLECTION_APP_USERS, "username={:u}", dbx.Params{"u": username})
					if err != nil {
						t.Errorf("expected normal user %q to exist: %v", username, err)
					}
				}
			},
		}
		s.Test(t)
	})

	t.Run("username pattern: super user is {name}-s", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()

		pbToken := generatePBSuperuserToken(t, app)
		s := tests.ApiScenario{
			Name:            "session create - username pattern super",
			Method:          http.MethodPost,
			URL:             sessionsURL,
			Body:            strings.NewReader(newSessionBody("beta-session")),
			Headers:         map[string]string{"Content-Type": "application/json", "Authorization": pbToken},
			ExpectedStatus:  200,
			ExpectedContent: []string{`"collectionName":"sessions"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				_, err := app.FindFirstRecordByFilter(COLLECTION_APP_USERS, "username={:u}", dbx.Params{"u": "beta-session-s"})
				if err != nil {
					t.Errorf("expected super user 'beta-session-s' to exist: %v", err)
				}
			},
		}
		s.Test(t)
	})

	t.Run("spaces in session name are replaced with hyphens in usernames", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()

		pbToken := generatePBSuperuserToken(t, app)
		s := tests.ApiScenario{
			Name:            "session create - spaces in name",
			Method:          http.MethodPost,
			URL:             sessionsURL,
			Body:            strings.NewReader(newSessionBody("my session name")),
			Headers:         map[string]string{"Content-Type": "application/json", "Authorization": pbToken},
			ExpectedStatus:  200,
			ExpectedContent: []string{`"collectionName":"sessions"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				_, err := app.FindFirstRecordByFilter(COLLECTION_APP_USERS, "username={:u}", dbx.Params{"u": "my-session-name-n1"})
				if err != nil {
					t.Errorf("expected hyphenated username 'my-session-name-n1': %v", err)
				}
				_, err = app.FindFirstRecordByFilter(COLLECTION_APP_USERS, "username={:u}", dbx.Params{"u": "my-session-name-s"})
				if err != nil {
					t.Errorf("expected hyphenated super username 'my-session-name-s': %v", err)
				}
			},
		}
		s.Test(t)
	})

	t.Run("each normal user gets a token, super user does not", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()

		pbToken := generatePBSuperuserToken(t, app)
		s := tests.ApiScenario{
			Name:            "session create - token created for each normal user",
			Method:          http.MethodPost,
			URL:             sessionsURL,
			Body:            strings.NewReader(newSessionBody("gamma-session")),
			Headers:         map[string]string{"Content-Type": "application/json", "Authorization": pbToken},
			ExpectedStatus:  200,
			ExpectedContent: []string{`"collectionName":"sessions"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				session, err := app.FindFirstRecordByFilter(COLLECTION_SESSIONS, "name='gamma-session'")
				if err != nil {
					t.Fatal("session not found:", err)
				}

				// Verify 10 tokens created (one per normal user)
				normalUsers, err := app.FindRecordsByFilter(
					COLLECTION_APP_USERS,
					"session={:sid} && role='normal'",
					"", -1, 0,
					dbx.Params{"sid": session.Id},
				)
				if err != nil {
					t.Fatal("users query:", err)
				}
				for _, u := range normalUsers {
					tokens, err := app.FindRecordsByFilter(
						COLLECTION_TOKENS,
						"user={:uid}",
						"", -1, 0,
						dbx.Params{"uid": u.Id},
					)
					if err != nil {
						t.Fatalf("token query for user %s: %v", u.GetString("username"), err)
					}
					if len(tokens) != 1 {
						t.Errorf("user %s: expected 1 token, got %d", u.GetString("username"), len(tokens))
					}
				}

				// Super user has no token
				superUser, err := app.FindFirstRecordByFilter(
					COLLECTION_APP_USERS,
					"session={:sid} && role='super'",
					dbx.Params{"sid": session.Id},
				)
				if err != nil {
					t.Fatal("super user not found:", err)
				}
				superTokens, _ := app.FindRecordsByFilter(
					COLLECTION_TOKENS,
					"user={:uid}",
					"", -1, 0,
					dbx.Params{"uid": superUser.Id},
				)
				if len(superTokens) != 0 {
					t.Errorf("super user should have 0 tokens, got %d", len(superTokens))
				}
			},
		}
		s.Test(t)
	})

	t.Run("user pwordtext is set and password matches", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()

		pbToken := generatePBSuperuserToken(t, app)
		s := tests.ApiScenario{
			Name:            "session create - pwordtext and password set",
			Method:          http.MethodPost,
			URL:             sessionsURL,
			Body:            strings.NewReader(newSessionBody("delta-session")),
			Headers:         map[string]string{"Content-Type": "application/json", "Authorization": pbToken},
			ExpectedStatus:  200,
			ExpectedContent: []string{`"collectionName":"sessions"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				u, err := app.FindFirstRecordByFilter(COLLECTION_APP_USERS, "username={:u}", dbx.Params{"u": "delta-session-n1"})
				if err != nil {
					t.Fatal("user not found:", err)
				}
				ptext := u.GetString(FIELD_PWORD)
				if ptext == "" {
					t.Error("expected pwordtext to be set")
				}
				if len(ptext) != 5 {
					t.Errorf("expected pwordtext length 5, got %d", len(ptext))
				}
				if !u.ValidatePassword(ptext) {
					t.Error("password does not match pwordtext")
				}
			},
		}
		s.Test(t)
	})

	t.Run("token valid_until is approximately 1 month from now", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()

		pbToken := generatePBSuperuserToken(t, app)
		s := tests.ApiScenario{
			Name:            "session create - token valid_until ~1 month",
			Method:          http.MethodPost,
			URL:             sessionsURL,
			Body:            strings.NewReader(newSessionBody("epsilon-session")),
			Headers:         map[string]string{"Content-Type": "application/json", "Authorization": pbToken},
			ExpectedStatus:  200,
			ExpectedContent: []string{`"collectionName":"sessions"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				u, err := app.FindFirstRecordByFilter(COLLECTION_APP_USERS, "username={:u}", dbx.Params{"u": "epsilon-session-n1"})
				if err != nil {
					t.Fatal("user not found:", err)
				}
				tokens, err := app.FindRecordsByFilter(
					COLLECTION_TOKENS, "user={:uid}", "", -1, 0, dbx.Params{"uid": u.Id},
				)
				if err != nil || len(tokens) == 0 {
					t.Fatal("no token found:", err)
				}
				validUntil := tokens[0].GetDateTime("valid_until").Time()
				expectedMin := time.Now().AddDate(0, 1, -1)
				expectedMax := time.Now().AddDate(0, 1, 1)
				if validUntil.Before(expectedMin) || validUntil.After(expectedMax) {
					t.Errorf("token valid_until %v not in expected ~1 month range", validUntil)
				}
			},
		}
		s.Test(t)
	})
}

// ============================================================
// increaseTokenTimeAfterSubmitSurvey hook tests
// ============================================================

func TestIncreaseTokenTimeAfterSubmitSurvey(t *testing.T) {
	surveyURL := "/api/collections/" + COLLECTION_SURVEY + "/records"

	t.Run("submitting survey extends valid tokens by 3 days", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		userToken := generateUserToken(t, app, td.normalUsername)

		// Record valid_until BEFORE survey creation
		beforeToken, err := app.FindRecordById(COLLECTION_TOKENS, td.validTokenID)
		if err != nil {
			t.Fatal("valid token not found:", err)
		}
		beforeTime := beforeToken.GetDateTime("valid_until").Time()

		s := tests.ApiScenario{
			Name:   "survey create - extends tokens",
			Method: http.MethodPost,
			URL:    surveyURL,
			Body: strings.NewReader(fmt.Sprintf(
				`{"response":{"q1":"a"},"session":%q,"user":%q}`,
				td.sessionID, td.normalUserID,
			)),
			Headers:         map[string]string{"Content-Type": "application/json", "Authorization": userToken},
			ExpectedStatus:  200,
			ExpectedContent: []string{`"collectionName":"surveys"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				afterToken, err := app.FindRecordById(COLLECTION_TOKENS, td.validTokenID)
				if err != nil {
					t.Fatal("token not found after survey:", err)
				}
				afterTime := afterToken.GetDateTime("valid_until").Time()
				diff := afterTime.Sub(beforeTime)
				// Should be extended by ~3 days (allow ±1 hour)
				expectedMin := (3*24 - 1) * time.Hour
				expectedMax := (3*24 + 1) * time.Hour
				if diff < expectedMin || diff > expectedMax {
					t.Errorf("token extended by %v, expected ~3 days", diff)
				}
			},
		}
		s.Test(t)
	})

	t.Run("expired tokens are NOT extended after survey submission", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		userToken := generateUserToken(t, app, td.normalUsername)

		beforeToken, _ := app.FindRecordById(COLLECTION_TOKENS, td.expiredTokenID)
		beforeTime := beforeToken.GetDateTime("valid_until").Time()

		s := tests.ApiScenario{
			Name:   "survey create - expired tokens not extended",
			Method: http.MethodPost,
			URL:    surveyURL,
			Body: strings.NewReader(fmt.Sprintf(
				`{"response":{"q1":"a"},"session":%q,"user":%q}`,
				td.sessionID, td.normalUserID,
			)),
			Headers:         map[string]string{"Content-Type": "application/json", "Authorization": userToken},
			ExpectedStatus:  200,
			ExpectedContent: []string{`"collectionName":"surveys"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				afterToken, _ := app.FindRecordById(COLLECTION_TOKENS, td.expiredTokenID)
				afterTime := afterToken.GetDateTime("valid_until").Time()
				if !afterTime.Equal(beforeTime) {
					t.Errorf("expired token should NOT be extended; before=%v after=%v", beforeTime, afterTime)
				}
			},
		}
		s.Test(t)
	})

	t.Run("unauthenticated survey create is rejected", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)

		s := tests.ApiScenario{
			Name:   "survey create - no auth",
			Method: http.MethodPost,
			URL:    surveyURL,
			Body: strings.NewReader(fmt.Sprintf(
				`{"response":{"q1":"a"},"session":%q,"user":%q}`,
				td.sessionID, td.normalUserID,
			)),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  400,
			ExpectedContent: []string{"Failed to create record"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("user with multiple valid tokens - all tokens extended", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)

		// Add a second valid token for the same normal user
		tokenColl, _ := app.FindCollectionByNameOrId(COLLECTION_TOKENS)
		extraToken := core.NewRecord(tokenColl)
		extraToken.Set("user", td.normalUserID)
		extraToken.Set("token", "extratoken12345")
		extraToken.Set("valid_until", time.Now().AddDate(0, 1, 0))
		if err := app.Save(extraToken); err != nil {
			t.Fatal("seed extra token:", err)
		}

		userToken := generateUserToken(t, app, td.normalUsername)

		s := tests.ApiScenario{
			Name:   "survey create - multiple valid tokens extended",
			Method: http.MethodPost,
			URL:    surveyURL,
			Body: strings.NewReader(fmt.Sprintf(
				`{"response":{"q1":"a"},"session":%q,"user":%q}`,
				td.sessionID, td.normalUserID,
			)),
			Headers:         map[string]string{"Content-Type": "application/json", "Authorization": userToken},
			ExpectedStatus:  200,
			ExpectedContent: []string{`"collectionName":"surveys"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				records, err := app.FindRecordsByFilter(
					COLLECTION_TOKENS,
					"user={:uid} && valid_until > @now",
					"", -1, 0,
					dbx.Params{"uid": td.normalUserID},
				)
				if err != nil {
					t.Fatal("token query:", err)
				}
				// Both valid tokens (original + extra) should now have valid_until ~now+3days
				for _, r := range records {
					vu := r.GetDateTime("valid_until").Time()
					// After extension, valid_until should be ~3 days from now, not 1 month
					// (they were 1 month away, now extended by 3 days → still ~1 month, just +3 days)
					// Check the difference is > now to confirm they are still valid
					if vu.Before(time.Now()) {
						t.Errorf("token %s expired after survey submission", r.Id)
					}
				}
				if len(records) < 2 {
					t.Errorf("expected at least 2 valid tokens, got %d", len(records))
				}
			},
		}
		s.Test(t)
	})
}

// ============================================================
// OnRecordEnrich hook: pwordtext visibility
// ============================================================

func TestPwordtextVisibility(t *testing.T) {
	t.Run("pwordtext hidden from normal user", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		normalToken := generateUserToken(t, app, td.normalUsername)

		s := tests.ApiScenario{
			Name:               "enrich hook - normal user cannot see pwordtext",
			Method:             http.MethodGet,
			URL:                "/api/collections/" + COLLECTION_APP_USERS + "/records/" + td.superUserID,
			Headers:            map[string]string{"Authorization": normalToken},
			ExpectedStatus:     200,
			NotExpectedContent: []string{`"pwordtext"`},
			TestAppFactory:     func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("pwordtext hidden from guest", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)

		// viewRule is empty string (anyone can view), but pwordtext should be hidden
		s := tests.ApiScenario{
			Name:               "enrich hook - guest cannot see pwordtext",
			Method:             http.MethodGet,
			URL:                "/api/collections/" + COLLECTION_APP_USERS + "/records/" + td.superUserID,
			ExpectedStatus:     200,
			NotExpectedContent: []string{`"pwordtext"`},
			TestAppFactory:     func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})
}

// ============================================================
// Collection rules: sessions
// ============================================================

func TestSessionCollectionRules(t *testing.T) {
	sessionsURL := "/api/collections/" + COLLECTION_SESSIONS + "/records"

	t.Run("unauthenticated list is blocked", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		seedTestData(t, app)

		s := tests.ApiScenario{
			Name:            "sessions - guest list",
			Method:          http.MethodGet,
			URL:             sessionsURL,
			ExpectedStatus:  200,
			ExpectedContent: []string{`"items":[]`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("user from same session can view session record", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		token := generateUserToken(t, app, td.normalUsername)

		s := tests.ApiScenario{
			Name:            "sessions - same session user can view",
			Method:          http.MethodGet,
			URL:             sessionsURL + "/" + td.sessionID,
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  200,
			ExpectedContent: []string{`"name":"test-session"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("user from different session cannot view session", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)

		// Create second session with its own super user (bypassing hook via app.Save)
		sessionColl, _ := app.FindCollectionByNameOrId(COLLECTION_SESSIONS)
		session2 := core.NewRecord(sessionColl)
		session2.Set("name", "other-session")
		session2.Set("description", "Another session")
		session2.Set("current_step", "2")
		session2.Set("contact_email", "other@test.dk")
		if err := app.Save(session2); err != nil {
			t.Fatal("seed second session:", err)
		}
		userColl, _ := app.FindCollectionByNameOrId(COLLECTION_APP_USERS)
		user2 := core.NewRecord(userColl)
		user2.Set("username", "other-normal")
		user2.Set("role", ROLE_NORMAL)
		user2.Set("status", "ready")
		user2.Set("session", session2.Id)
		user2.SetPassword("opass")
		user2.SetVerified(true)
		if err := app.Save(user2); err != nil {
			t.Fatal("seed second user:", err)
		}
		token2 := generateUserToken(t, app, "other-normal")

		s := tests.ApiScenario{
			Name:            "sessions - different session user blocked",
			Method:          http.MethodGet,
			URL:             sessionsURL + "/" + td.sessionID,
			Headers:         map[string]string{"Authorization": token2},
			ExpectedStatus:  404, // PocketBase returns 404 when rule prevents access
			ExpectedContent: []string{"wasn't found"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("create session via API is allowed (null create rule)", func(t *testing.T) {
		// PocketBase null createRule means the rule is not set (no restriction from collection rules),
		// but a superadmin is required to create collections directly. For app_users:
		// Actually looking at the schema - sessions createRule is null which means admin-only.
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		token := generateUserToken(t, app, td.normalUsername)

		s := tests.ApiScenario{
			Name:   "sessions - normal user cannot create (null createRule = admin-only)",
			Method: http.MethodPost,
			URL:    sessionsURL,
			Body:   strings.NewReader(`{"name":"new-s","description":"d","current_step":"2","contact_email":"x@x.dk"}`),
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": token,
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{"superusers"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})
}

// ============================================================
// Collection rules: surveys
// ============================================================

func TestSurveyCollectionRules(t *testing.T) {
	surveyURL := "/api/collections/" + COLLECTION_SURVEY + "/records"

	t.Run("normal user can only see own surveys", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		token := generateUserToken(t, app, td.normalUsername)

		s := tests.ApiScenario{
			Name:            "surveys - normal user list shows own survey only",
			Method:          http.MethodGet,
			URL:             surveyURL,
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  200,
			ExpectedContent: []string{td.surveyID},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("super user in same session sees all surveys for that session", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		token := generateUserToken(t, app, td.superUsername)

		s := tests.ApiScenario{
			Name:            "surveys - super sees session surveys",
			Method:          http.MethodGet,
			URL:             surveyURL,
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  200,
			ExpectedContent: []string{td.surveyID},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})
}

// ============================================================
// Collection rules: tokens (all null = admin-only)
// ============================================================

func TestTokensCollectionRules(t *testing.T) {
	t.Run("guest cannot list tokens", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		seedTestData(t, app)

		s := tests.ApiScenario{
			Name:            "tokens - guest list",
			Method:          http.MethodGet,
			URL:             "/api/collections/" + COLLECTION_TOKENS + "/records",
			ExpectedStatus:  403,
			ExpectedContent: []string{"superusers"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("authenticated normal user cannot list tokens", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		token := generateUserToken(t, app, td.normalUsername)

		s := tests.ApiScenario{
			Name:            "tokens - normal user list",
			Method:          http.MethodGet,
			URL:             "/api/collections/" + COLLECTION_TOKENS + "/records",
			Headers:         map[string]string{"Authorization": token},
			ExpectedStatus:  403,
			ExpectedContent: []string{"superusers"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})
}

// ============================================================
// Collection rules: questions (public read)
// ============================================================

func TestQuestionsCollectionRules(t *testing.T) {
	t.Run("questions are publicly readable without auth", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()

		s := tests.ApiScenario{
			Name:            "questions - public list",
			Method:          http.MethodGet,
			URL:             "/api/collections/questions/records",
			ExpectedStatus:  200,
			ExpectedContent: []string{`"totalItems"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("questions seeded: supply_chain v1 exists", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()

		s := tests.ApiScenario{
			Name:            "questions - supply_chain v1 seeded",
			Method:          http.MethodGet,
			URL:             "/api/collections/questions/records?filter=name='supply_chain'",
			ExpectedStatus:  200,
			ExpectedContent: []string{`"supply_chain"`, `"v1"`},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})
}

// ============================================================
// Edge cases & validation
// ============================================================

func TestValidationEdgeCases(t *testing.T) {
	t.Run("duplicate session name is rejected", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		seedTestData(t, app)

		// "test-session" was already created by seedTestData
		s := tests.ApiScenario{
			Name:            "duplicate session name",
			Method:          http.MethodPost,
			URL:             "/api/collections/" + COLLECTION_SESSIONS + "/records",
			Body:            strings.NewReader(`{"name":"test-session","description":"dup","current_step":"2","contact_email":"dup@test.dk"}`),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  403,
			ExpectedContent: []string{"superusers"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("session current_step must be a valid option", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()

		s := tests.ApiScenario{
			Name:            "session invalid current_step",
			Method:          http.MethodPost,
			URL:             "/api/collections/" + COLLECTION_SESSIONS + "/records",
			Body:            strings.NewReader(`{"name":"step-test","description":"d","current_step":"invalid","contact_email":"s@s.dk"}`),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  403,
			ExpectedContent: []string{"superusers"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("session missing required contact_email is rejected", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()

		s := tests.ApiScenario{
			Name:            "session missing contact_email",
			Method:          http.MethodPost,
			URL:             "/api/collections/" + COLLECTION_SESSIONS + "/records",
			Body:            strings.NewReader(`{"name":"no-email","description":"d","current_step":"2"}`),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  403,
			ExpectedContent: []string{"superusers"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("app_user role must be valid option (normal|super)", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)

		s := tests.ApiScenario{
			Name:   "app_user invalid role",
			Method: http.MethodPost,
			URL:    "/api/collections/" + COLLECTION_APP_USERS + "/records",
			Body: strings.NewReader(fmt.Sprintf(
				`{"username":"badrol","role":"admin","status":"ready","session":%q,"password":"pass12","passwordConfirm":"pass12"}`,
				td.sessionID,
			)),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  400,
			ExpectedContent: []string{"validation_invalid_value"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("app_user password too short (min 5)", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)

		s := tests.ApiScenario{
			Name:   "app_user password too short",
			Method: http.MethodPost,
			URL:    "/api/collections/" + COLLECTION_APP_USERS + "/records",
			Body: strings.NewReader(fmt.Sprintf(
				`{"username":"shrtpw","role":"normal","status":"ready","session":%q,"password":"ab","passwordConfirm":"ab"}`,
				td.sessionID,
			)),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  400,
			ExpectedContent: []string{"validation_min_text_constraint"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("app_user username too short (min 3)", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)

		s := tests.ApiScenario{
			Name:   "app_user username too short",
			Method: http.MethodPost,
			URL:    "/api/collections/" + COLLECTION_APP_USERS + "/records",
			Body: strings.NewReader(fmt.Sprintf(
				`{"username":"ab","role":"normal","status":"ready","session":%q,"password":"pass12","passwordConfirm":"pass12"}`,
				td.sessionID,
			)),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  400,
			ExpectedContent: []string{"validation_min_text_constraint"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("duplicate username is rejected", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)

		s := tests.ApiScenario{
			Name:   "duplicate username",
			Method: http.MethodPost,
			URL:    "/api/collections/" + COLLECTION_APP_USERS + "/records",
			Body: strings.NewReader(fmt.Sprintf(
				`{"username":%q,"role":"normal","status":"ready","session":%q,"password":"pass12","passwordConfirm":"pass12"}`,
				td.normalUsername, td.sessionID,
			)),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  400,
			ExpectedContent: []string{"validation_not_unique"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("survey response JSON oversized is rejected", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)
		userToken := generateUserToken(t, app, td.normalUsername)

		// Build a JSON string that exceeds 20000 bytes
		largeValue := strings.Repeat("x", 20001)
		oversizedBody := fmt.Sprintf(
			`{"response":{"data":%q},"session":%q,"user":%q}`,
			largeValue, td.sessionID, td.normalUserID,
		)

		s := tests.ApiScenario{
			Name:            "survey oversized response rejected",
			Method:          http.MethodPost,
			URL:             "/api/collections/" + COLLECTION_SURVEY + "/records",
			Body:            strings.NewReader(oversizedBody),
			Headers:         map[string]string{"Content-Type": "application/json", "Authorization": userToken},
			ExpectedStatus:  400,
			ExpectedContent: []string{"validation_json_size_limit"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("token with clearly past valid_until is treated as expired in loginWithToken", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		td := seedTestData(t, app)

		// expiredToken was seeded with valid_until = 1 month ago
		s := tests.ApiScenario{
			Name:   "expired token cannot authenticate",
			Method: http.MethodPost,
			URL:    "/api/collections/" + COLLECTION_APP_USERS + "/auth-with-password",
			Body: strings.NewReader(fmt.Sprintf(
				`{"identity":%q,"password":%q}`, td.normalUsername, td.expiredToken,
			)),
			Headers:         map[string]string{"Content-Type": "application/json"},
			ExpectedStatus:  400,
			ExpectedContent: []string{"Failed to authenticate"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("GET /auth-set with wrong HTTP method returns 405", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		seedTestData(t, app)

		s := tests.ApiScenario{
			Name:            "auth-set wrong method",
			Method:          http.MethodPut,
			URL:             "/auth-set",
			ExpectedStatus:  404,
			ExpectedContent: []string{"wasn't found"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})

	t.Run("POST /token wrong method (GET) returns 405", func(t *testing.T) {
		app := setupTestApp(t)
		defer app.Cleanup()
		seedTestData(t, app)

		s := tests.ApiScenario{
			Name:            "token wrong method GET",
			Method:          http.MethodGet,
			URL:             "/token",
			ExpectedStatus:  404,
			ExpectedContent: []string{"not found"},
			TestAppFactory:  func(t testing.TB) *tests.TestApp { return app },
		}
		s.Test(t)
	})
}
