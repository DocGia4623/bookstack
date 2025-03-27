package service

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"bookstack/internal/repository"
)

type BookService interface {
	CreateBook(int, request.BookCreateRequest) (models.Book, error)
	CreateShelve(int, request.ShelveCreateRequest) (models.Shelve, error)
	GetAllBook() ([]models.Book, error)
	CreateChapter(request.BookChapterRequest) (models.Chapter, error)
	GetChaptersOfBook(int) ([]models.Chapter, error)
	AddPage(request.PageRequest) (models.Page, error)
}

type BookServiceImpl struct {
	repo repository.BookRepository
}

func NewBookServiceImpl(repository repository.BookRepository) BookService {
	return &BookServiceImpl{
		repo: repository,
	}
}

func (b *BookServiceImpl) AddPage(request request.PageRequest) (models.Page, error) {
	return b.repo.AddPage(request)
}
func (b *BookServiceImpl) GetChaptersOfBook(bookId int) ([]models.Chapter, error) {
	return b.repo.GetChaptersOfBook(bookId)
}

func (b *BookServiceImpl) CreateChapter(request request.BookChapterRequest) (models.Chapter, error) {
	return b.repo.CreateChapter(request)
}

func (b *BookServiceImpl) GetAllBook() ([]models.Book, error) {
	return b.repo.GetAllBook()
}

func (b *BookServiceImpl) CreateShelve(userId int, request request.ShelveCreateRequest) (models.Shelve, error) {
	return b.repo.CreateShelve(userId, request)
}

func (b *BookServiceImpl) CreateBook(userId int, request request.BookCreateRequest) (models.Book, error) {
	return b.repo.CreateBook(userId, request)
}
