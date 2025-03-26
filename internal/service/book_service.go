package service

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"bookstack/internal/repository"
)

type BookService interface {
	CreateBook(request.BookCreateRequest) (models.Book, error)
}

type BookServiceImpl struct {
	repo repository.BookRepository
}

func NewBookServiceImpl(repository repository.BookRepository) BookService {
	return &BookServiceImpl{
		repo: repository,
	}
}

func (b *BookServiceImpl) CreateBook(request request.BookCreateRequest) (models.Book, error) {
	return b.repo.CreateBook(request)
}
