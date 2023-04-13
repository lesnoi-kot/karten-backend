package authservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/lesnoi-kot/karten-backend/src/authservice/oauth"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

type AuthService struct {
	Store *store.Store
}

type RegisterUserOptions struct{}

func (service AuthService) Authenticate(userInfo *oauth.UserInfo) (*store.User, error) {
	db_social_id := fmt.Sprintf("%s_%s", userInfo.AuthProvider, userInfo.ID)

	db_user, err := service.Store.Users.GetBySocialID(context.Background(), db_social_id)
	if errors.Is(err, store.ErrNotFound) {
		db_user = &store.User{
			SocialID: db_social_id,
			Name:     userInfo.Name,
			Login:    userInfo.Login,
			Email:    userInfo.Email,
			URL:      userInfo.URL,
		}
		// TODO: avatar

		if err := service.Store.Users.Add(context.Background(), db_user); err != nil {
			return nil, err
		}

		service.OnRegister(db_user)
	} else if err != nil {
		return nil, err
	}

	return db_user, nil
}

func (service AuthService) OnRegister(*store.User) error {
	// Add a tutorial project and board
	return nil
}
