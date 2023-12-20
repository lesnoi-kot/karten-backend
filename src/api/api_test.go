package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"text/template"
	"time"

	"github.com/lesnoi-kot/karten-backend/src/api"
	"github.com/lesnoi-kot/karten-backend/src/entityservices"
	"github.com/lesnoi-kot/karten-backend/src/filestorage"
	"github.com/lesnoi-kot/karten-backend/src/settings"
	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun/dbfixture"
	"go.uber.org/zap"
)

var (
	testDate = time.Unix(0, 0).UTC()
	testUser = store.User{
		ID:       33,
		SocialID: "ae4f869a3f59db15cd8173b451c97025",
		Name:     "john_doe",
		Login:    "John Doe",
		Email:    "john_doe@hotmail.com",
		URL:      "www.google.com",
	}
	fixturesFS = os.DirFS("./testdata")
)

type APITestSuite struct {
	suite.Suite
	store         *store.Store
	fileStorage   filestorage.FileStorage
	api           *api.APIService
	fixtureLoader *dbfixture.Fixture
}

func TestIntegrationAPI(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip("api intergation tests are skipped")
	}

	suite.Run(t, new(APITestSuite))
}

func (suite *APITestSuite) SetupSuite() {
	settings.AppConfig.SessionsSecretKey = "test"
	settings.AppConfig.SessionsStorePath = "" // = os.TempDir()
	settings.AppConfig.EnableGuest = true
	settings.AppConfig.MediaURL = "http://nginx:80/"

	suite.fileStorage, _ = filestorage.NewFileSystemStorage(os.TempDir())
	suite.store = store.NewStore(store.StoreConfig{
		DSN:    os.Getenv("STORE_DSN"),
		Logger: zap.NewNop().Sugar(),
		Debug:  false,
	})

	suite.fixtureLoader = dbfixture.New(
		suite.store.ORM,
		dbfixture.WithTruncateTables(),
		dbfixture.WithTemplateFuncs(template.FuncMap{
			"get_test_user": func() store.User {
				return testUser
			},
		}))

	if err := suite.fixtureLoader.Load(context.Background(), fixturesFS, "fixtures.yml"); err != nil {
		suite.T().Fatalf("Fixture load error: %s", err)
	}

	suite.api = api.NewAPI(api.APIConfig{
		Store:       suite.store,
		Logger:      zap.Must(zap.NewDevelopment()).Sugar(),
		FileStorage: nil,
		ContextsContainer: entityservices.ContextsContainer{
			Store:       suite.store,
			FileStorage: nil,
		},
		FrontendURL:  "",
		APIPrefix:    "",
		CookieDomain: "",
		Debug:        false,
	})
}

func (suite *APITestSuite) authorizeRequest(r *http.Request, userID store.UserID) error {
	sess, err := suite.api.SessionStore.Get(r, api.USER_SESSION_KEY)
	sess.Values[api.SESSION_KEY_USER_ID] = userID
	if err != nil {
		return err
	}

	rec := httptest.NewRecorder()
	if err = sess.Save(r, rec); err != nil {
		return err
	}

	for _, c := range rec.Result().Cookies() {
		r.AddCookie(c)
	}

	return nil
}

func (suite *APITestSuite) setCSRF(outer *http.Request) string {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	suite.api.Server().Handler.ServeHTTP(rec, req)

	for _, c := range rec.Result().Cookies() {
		if c.Name == "_csrf" {
			outer.AddCookie(c)
			outer.Header.Add("X-CSRF-Token", c.Value)
			return c.Value
		}
	}

	return ""
}
