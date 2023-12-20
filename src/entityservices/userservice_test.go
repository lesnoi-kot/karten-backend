package entityservices_test

import (
	"context"
	"text/template"

	"github.com/lesnoi-kot/karten-backend/src/entityservices"
	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun/dbfixture"
)

type userserviceSuite struct {
	suite.Suite
	ContextsContainer *entityservices.ContextsContainer
}

func (suite *userserviceSuite) SetupSuite() {
	suite.ContextsContainer = &entityservices.ContextsContainer{
		Store:       storeService,
		FileStorage: nil,
	}

	fixtures := dbfixture.New(storeService.ORM, dbfixture.WithTruncateTables(), dbfixture.WithTemplateFuncs(template.FuncMap{
		"get_test_user": func() store.User {
			return testUser
		},
	}))
	if err := fixtures.Load(context.Background(), fixturesFS, "userservice_fixtures.yml"); err != nil {
		suite.T().Fatalf("Fixture load error: %s", err)
	}
}

func (suite *userserviceSuite) TestUserService() {
	userContext := suite.ContextsContainer.GetAuthorizedUserContext(entityservices.GetAuthorizedUserContextOptions{
		UserID: testUser.ID,
	})

	suite.Run("Get test user", func() {
		user, err := userContext.UserService.GetUser(&entityservices.GetUserOptions{
			FullInfo:      true,
			IncludeAvatar: true,
		})
		suite.Require().NoError(err)
		suite.Equal(testUser, *user)
		suite.Equal("john_doe", user.Name)
		suite.False(user.IsGuest())
	})

	suite.Run("Test user is valid", func() {
		isValid, err := userContext.UserService.IsValidUser()
		suite.Require().NoError(err)
		suite.True(isValid)
	})

	suite.Run("Edit user", func() {
		s := "Jane Air"
		err := userContext.UserService.EditUser(&entityservices.EditUserOptions{
			Name: &s,
		})
		suite.Require().NoError(err)

		user, err := userContext.UserService.GetUser(&entityservices.GetUserOptions{
			FullInfo: true,
		})
		suite.Require().NoError(err)
		suite.Equal("Jane Air", user.Name)
	})

	suite.Run("Delete user", func() {
		err := userContext.UserService.DeleteUser()
		suite.Require().NoError(err)

		isValid, err := userContext.UserService.IsValidUser()
		suite.Require().NoError(err)
		suite.False(isValid)

		projects, err := userContext.ProjectService.GetProjects(entityservices.GetProjectsOptions{
			IncludeBoards: true,
		})
		suite.Require().NoError(err)
		suite.Len(projects, 0)
	})
}
