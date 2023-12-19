package entityservices_test

import (
	"context"
	"text/template"

	"github.com/lesnoi-kot/karten-backend/src/entityservices"
	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun/dbfixture"
)

type projectserviceSuite struct {
	suite.Suite
	ContextsContainer *entityservices.ContextsContainer
}

func (suite *projectserviceSuite) SetupSuite() {
	suite.ContextsContainer = &entityservices.ContextsContainer{
		Store:       storeService,
		FileStorage: nil,
	}

	fixtures := dbfixture.New(storeService.ORM, dbfixture.WithTruncateTables(), dbfixture.WithTemplateFuncs(template.FuncMap{
		"get_test_user": func() store.User {
			return testUser
		},
	}))
	if err := fixtures.Load(context.Background(), fixturesFS, "projectservice_fixtures.yml"); err != nil {
		suite.T().Fatalf("Fixture load error: %s", err)
	}
}

func (suite *projectserviceSuite) TestUserProjects() {
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

	suite.Run("Get a project of a user", func() {
		project, err := userContext.ProjectService.GetProject(entityservices.GetProjectOptions{
			ProjectID:     "2f146153-ee2f-4968-a241-11a4f00bf212",
			IncludeBoards: true,
		})
		suite.Require().NoError(err)
		suite.Equal(testUser.ID, project.UserID)
		suite.Equal("Business", project.Name)
		suite.Equal("2f146153-ee2f-4968-a241-11a4f00bf212", project.ID)

		suite.Len(project.Boards, 2)
		suite.ElementsMatch(
			[]string{project.Boards[0].ID, project.Boards[1].ID},
			[]string{"36461fd5-2eb4-42b9-a921-b4428a448cfa", "ec5da06b-6c0d-4439-9b04-b9da7266ba2f"},
		)
	})

	suite.Run("Get a project that does not exist", func() {
		project, err := userContext.ProjectService.GetProject(entityservices.GetProjectOptions{
			ProjectID:     "6bb81bb6-2eb8-4a99-98fa-6dc6f0d288bb",
			IncludeBoards: false,
		})
		suite.Require().Error(err)
		suite.Nil(project)
	})

	suite.Run("Get a project of another user", func() {
		project, err := userContext.ProjectService.GetProject(entityservices.GetProjectOptions{
			ProjectID:     "6d8065f4-032f-494d-b124-b7bdf3c6480b",
			IncludeBoards: false,
		})
		suite.Require().Error(err)
		suite.Nil(project)
	})

	suite.Run("Add projects", func() {
		project, err := userContext.ProjectService.AddProject(entityservices.AddProjectOptions{
			Name:     "PROJ111",
			AvatarID: nil,
		})
		suite.Require().NoError(err)
		suite.NotEmpty(project.ID)
		suite.Equal("PROJ111", project.Name)
		id := project.ID

		project, err = userContext.ProjectService.GetProject(entityservices.GetProjectOptions{
			ProjectID:     id,
			IncludeBoards: true,
		})
		suite.Require().NoError(err)
		suite.Equal(testUser.ID, project.UserID)
		suite.Equal("PROJ111", project.Name)
		suite.Equal(id, project.ID)
	})

	suite.Run("Delete project", func() {
		err := userContext.ProjectService.DeleteProject(entityservices.DeleteProjectOptions{
			ProjectID: "fd5f451d-fac6-4bc7-a677-34adb39a6701",
		})
		suite.Require().NoError(err)

		_, err = userContext.ProjectService.GetProject(entityservices.GetProjectOptions{
			ProjectID: "fd5f451d-fac6-4bc7-a677-34adb39a6701",
		})
		suite.Require().Error(err)
	})

	suite.Run("Delete not your project", func() {
		err := userContext.ProjectService.DeleteProject(entityservices.DeleteProjectOptions{
			ProjectID: "6d8065f4-032f-494d-b124-b7bdf3c6480b",
		})
		suite.Require().Error(err)
	})
	suite.Run("Delete project that does not exist", func() {
		err := userContext.ProjectService.DeleteProject(entityservices.DeleteProjectOptions{
			ProjectID: "5a931194-ad2f-4c2c-a689-f095f017e65b",
		})
		suite.Require().Error(err)
	})
}
