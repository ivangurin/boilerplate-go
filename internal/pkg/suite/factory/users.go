package suite_factory

import (
	"boilerplate/internal/repository"

	"github.com/brianvoe/gofakeit/v7"
)

type UserFactory struct {
	setters []func(*repository.User)
}

func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

func (f *UserFactory) WithAdmin() *UserFactory {
	f.setters = append(f.setters, func(user *repository.User) {
		user.IsAdmin = true
	})
	return f
}

func (f *UserFactory) WithPassword(password string) *UserFactory {
	f.setters = append(f.setters, func(user *repository.User) {
		user.Password = password
	})
	return f
}

func (f *UserFactory) Build() *repository.User {
	return f.generate()
}

func (f *UserFactory) Builds(times int) []*repository.User {
	res := make([]*repository.User, 0, times)
	for i := 0; i < times; i++ {
		res = append(res, f.generate())
	}
	return res
}

func (f *UserFactory) generate() *repository.User {
	user := &repository.User{
		Name:     gofakeit.FirstName() + " " + gofakeit.LastName(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Word(),
	}

	for _, setter := range f.setters {
		setter(user)
	}

	return user
}
