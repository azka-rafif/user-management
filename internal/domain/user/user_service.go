package user

import (
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type UserService interface {
	GetByUserName(userName string) (user User, err error)
	Create(load UserPayload) (user User, err error)
	UpdateName(payload NamePayload, userId uuid.UUID) (user User, err error)
	DeleteByID(userId, userDeleter uuid.UUID) (user User, err error)
}

type UserServiceImpl struct {
	Repo UserRepository
}

func ProvideUserServiceImpl(repo UserRepository) *UserServiceImpl {
	return &UserServiceImpl{Repo: repo}
}

func (s *UserServiceImpl) Create(load UserPayload) (user User, err error) {
	user, err = user.NewFromPayload(load)
	if err != nil {
		return
	}
	err = s.Repo.Create(user)
	if err != nil {
		return
	}
	return
}

func (s *UserServiceImpl) GetByUserName(userName string) (user User, err error) {
	user, err = s.Repo.GetByUserName(userName)

	if err != nil {
		return
	}

	return
}

func (s *UserServiceImpl) UpdateName(payload NamePayload, userId uuid.UUID) (user User, err error) {
	user, err = s.Repo.GetByUserId(userId)
	if err != nil {
		return
	}
	user.UpdateName(payload)
	err = s.Repo.Update(user)
	if err != nil {
		return
	}
	return
}

func (s *UserServiceImpl) DeleteByID(userId, userDeleter uuid.UUID) (user User, err error) {
	exists, err := s.Repo.ExistsByID(userId)
	if err != nil {
		return
	}
	if !exists {
		err = failure.NotFound("user")
		return
	}
	user, err = s.Repo.GetByUserId(userId)
	if err != nil {
		return
	}
	err = user.SoftDelete(userDeleter)
	if err != nil {
		return
	}
	err = s.Repo.Update(user)
	if err != nil {
		return
	}

	return
}