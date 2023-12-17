package entityservices

import (
	"errors"

	"github.com/lesnoi-kot/karten-backend/src/filestorage"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

var (
	ErrPermissionDenied error = errors.New("Permission denied")
	ErrResourceNotFound       = errors.New("resource not found")
)

type Actor struct {
	UserID store.UserID
}

type StoreInjector interface {
	GetStore() *store.Store
}

type ActorInjector interface {
	GetActor() Actor
}

type UserServiceInjector interface {
	GetUserService() IUserService
}

type FileServiceInjector interface {
	GetFileService() IFileService
}

type GetAuthorizedUserContextOptions struct {
	UserID store.UserID
}

type IContextsContainer interface {
	GetAuthorizedUserContext(args GetAuthorizedUserContextOptions) *AuthorizedUserContext
}

type ContextsContainer struct {
	Store       *store.Store
	FileStorage filestorage.FileStorage
}

func (container ContextsContainer) GetAuthorizedUserContext(args GetAuthorizedUserContextOptions) *AuthorizedUserContext {
	authorizedContext := &AuthorizedUserContext{
		Actor:       Actor{UserID: args.UserID},
		Store:       container.Store,
		FileStorage: container.FileStorage,
	}
	authorizedContext.UserService = &UserService{authorizedContext}
	authorizedContext.ProjectService = &ProjectService{authorizedContext}
	authorizedContext.BoardService = &BoardService{authorizedContext}
	authorizedContext.TaskListService = &TaskListService{authorizedContext}
	authorizedContext.TaskService = &TaskService{authorizedContext}
	authorizedContext.CommentService = &CommentService{authorizedContext}
	authorizedContext.FileService = &FileService{authorizedContext, container.FileStorage}
	return authorizedContext
}

type AuthorizedUserContext struct {
	Actor
	Store           *store.Store
	FileStorage     filestorage.FileStorage
	UserService     IUserService
	ProjectService  IProjectService
	BoardService    IBoardService
	TaskListService ITaskListService
	TaskService     ITaskService
	CommentService  ICommentService
	FileService     IFileService
}

func (auc AuthorizedUserContext) GetStore() *store.Store {
	return auc.Store
}

func (auc AuthorizedUserContext) GetActor() Actor {
	return auc.Actor
}

func (auc AuthorizedUserContext) GetFileService() IFileService {
	return auc.FileService
}

func (auc AuthorizedUserContext) GetUserService() IUserService {
	return auc.UserService
}
