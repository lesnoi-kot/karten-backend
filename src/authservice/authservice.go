package authservice

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/lesnoi-kot/karten-backend/src/authservice/oauth"
	"github.com/lesnoi-kot/karten-backend/src/fileservice"
	"github.com/lesnoi-kot/karten-backend/src/modules/images"
	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/lesnoi-kot/karten-backend/src/userservice"
)

type AuthService struct {
	Store       *store.Store
	FileService *fileservice.FileService
}

func (service AuthService) generateSocialID(userInfo *oauth.UserInfo) string {
	return fmt.Sprintf("%s_%s", userInfo.AuthProvider, userInfo.ID)
}

func (service AuthService) Authenticate(ctx context.Context, userInfo *oauth.UserInfo) (*store.User, error) {
	db_social_id := service.generateSocialID(userInfo)
	db_user, err := service.GetUserBySocialID(ctx, db_social_id)

	if errors.Is(err, sql.ErrNoRows) {
		// Register user if not found in the db.
		db_user = &store.User{
			SocialID: db_social_id,
			Name:     userInfo.Name,
			Login:    userInfo.Login,
			Email:    userInfo.Email,
			URL:      userInfo.URL,
		}

		if avatarID, err := service.copyAvatar(ctx, userInfo.AvatarURL); err == nil {
			db_user.AvatarID = avatarID
		}

		if err := service.RegisterUser(ctx, db_user); err != nil {
			return nil, err
		}

		service.onRegister(ctx, db_user)
	} else if err != nil {
		return nil, err
	}

	return db_user, nil
}

func (service AuthService) GetUserBySocialID(ctx context.Context, socialID string) (*store.User, error) {
	user := new(store.User)

	err := service.Store.ORM.
		NewSelect().
		Model(user).
		Where("social_id = ?", socialID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service AuthService) RegisterUser(ctx context.Context, db_user *store.User) error {
	_, err := service.Store.ORM.NewInsert().
		Model(db_user).
		Column("social_id", "avatar_id", "name", "login", "email", "url").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (service AuthService) onRegister(ctx context.Context, user *store.User) error {
	userService := userservice.UserService{
		Context: ctx,
		UserID:  user.ID,
		Store:   service.Store,
	}

	project, err := userService.AddProject(&userservice.AddProjectOptions{
		Name: user.Name,
	})
	if err != nil {
		return err
	}

	board, err := userService.AddBoard(&userservice.AddBoardOptions{
		ProjectID: project.ID,
		Name:      "Tutorial board",
		Color:     0x0094ae,
	})
	if err != nil {
		return err
	}

	{
		list, err := userService.AddTaskList(&userservice.AddTaskListOptions{
			BoardID:  board.ID,
			Name:     "Stuff to try (this is a list)",
			Position: 0,
		})
		if err != nil {
			return err
		}

		_, err = userService.AddTask(&userservice.AddTaskOptions{
			TaskListID: list.ID,
			Name:       "This is a card. Drag it to the \"Tried It\" List to show it's done. →",
			Position:   0,
		})
		if err != nil {
			return err
		}
	}

	{
		list, err := userService.AddTaskList(&userservice.AddTaskListOptions{
			BoardID:  board.ID,
			Name:     "Tried it",
			Position: 10000,
		})
		if err != nil {
			return err
		}

		_, err = userService.AddTask(&userservice.AddTaskOptions{
			TaskListID: list.ID,
			Name:       "Lets go",
			Position:   0,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Upload user avatar and add to the database.
func (service AuthService) copyAvatar(ctx context.Context, avatarURL string) (store.FileID, error) {
	if avatarURL == "" {
		return "", nil
	}

	resp, err := http.Get(avatarURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	imgInfo, err := images.ParseImage(bytes.NewReader(body))
	fileExtension := strings.TrimPrefix(imgInfo.MIMEType, "image/")

	avatar, err := service.FileService.AddImage(ctx, fileservice.AddFileOptions{
		Name:     fmt.Sprintf("avatar.%s", fileExtension),
		MIMEType: imgInfo.MIMEType,
		Data:     bytes.NewReader(body),
	})
	if err != nil {
		return "", err
	}

	return avatar.ID, nil
}
