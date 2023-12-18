package entityservices_test

import (
	"context"
	"os"
	"testing"
	"text/template"
	"time"

	"github.com/lesnoi-kot/karten-backend/src/entityservices"
	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/samber/lo"
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
)

type userserviceSuite struct {
	suite.Suite
	ContextsContainer *entityservices.ContextsContainer
}

func TestIntegrationUserserviceSuite(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip("entityservices intergation tests are skipped")
	}

	suite.Run(t, new(userserviceSuite))
}

func (suite *userserviceSuite) SetupSuite() {
	storeService := store.NewStore(store.StoreConfig{
		DSN:    os.Getenv("STORE_DSN"),
		Logger: zap.NewNop().Sugar(),
		Debug:  false,
	})

	suite.ContextsContainer = &entityservices.ContextsContainer{
		Store:       storeService,
		FileStorage: nil,
	}

	fixtures := dbfixture.New(storeService.ORM, dbfixture.WithTruncateTables(), dbfixture.WithTemplateFuncs(template.FuncMap{
		"get_test_user": func() store.User {
			return testUser
		},
	}))
	if err := fixtures.Load(context.Background(), os.DirFS("./testdata"), "userservice_fixtures.yml"); err != nil {
		suite.T().Fatalf("Fixture load error: %s", err)
	}
}

func (suite *userserviceSuite) TestUserProjects() {
	userContext := suite.ContextsContainer.GetAuthorizedUserContext(entityservices.GetAuthorizedUserContextOptions{
		UserID: testUser.ID,
	})

	suite.Run("Get all projects of a user", func() {
		all, err := userContext.ProjectService.GetProjects(entityservices.GetProjectsOptions{
			IncludeBoards: false,
		})
		suite.Require().NoError(err)
		suite.Len(all, 4)

		suite.ElementsMatch(
			lo.Map[*store.Project, string](all, func(p *store.Project, _ int) string {
				return p.ID
			}),
			[]string{
				"fd5f451d-fac6-4bc7-a677-34adb39a6701",
				"2f146153-ee2f-4968-a241-11a4f00bf212",
				"2d2712eb-266d-4626-b017-697a67907e28",
				"1f894df2-f233-4885-81ef-e21aee62e2cd",
			},
		)

		for _, p := range all {
			suite.Equal(testUser.ID, p.UserID)
		}
	})
}
