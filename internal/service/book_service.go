package service

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"bookstack/internal/repository"
)

type BookService interface {
	CreateBook(int, request.BookCreateRequest) (models.Book, error)
	CreateShelve(int, request.ShelveCreateRequest) (models.Shelve, error)
}

type BookServiceImpl struct {
	repo repository.BookRepository
}

func NewBookServiceImpl(repository repository.BookRepository) BookService {
	return &BookServiceImpl{
		repo: repository,
	}
}

func (b *BookServiceImpl) CreateShelve(userId int, request request.ShelveCreateRequest) (models.Shelve, error) {
	return b.repo.CreateShelve(userId, request)
}

func (b *BookServiceImpl) CreateBook(userId int, request request.BookCreateRequest) (models.Book, error) {
	return b.repo.CreateBook(userId, request)
}
