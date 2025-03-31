package service

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"bookstack/internal/repository"
)

type BookService interface {
	//book
	CreateCompleteBook(int, request.CompleteBookCreateRequest) (models.Book, error)
	CreateBook(int, request.BookCreateRequest) (models.Book, error)
	DeleteBook(int) error
	UpdateBook(int, request.BookCreateRequest) (models.Book, error)
	GetAllBook() ([]models.Book, error)
	//shelve
	CreateShelve(int, request.ShelveCreateRequest) (models.Shelve, error)
	GetShelves() ([]models.Shelve, error)
	DeleteShelve(int) error
	//chapter
	CreateChapter(uint, request.BookChapterRequest) (models.Chapter, error)
	GetChaptersOfBook(int) ([]models.Chapter, error)
	DeleteChapter(int) error
	UpdateChapter(int, request.BookChapterRequest) (models.Chapter, error)
	//page
	AddPage(uint, request.PageRequest) (models.Page, error)
	GetPageChapter(int) ([]models.Page, error)
	DeletePage(int) error
	UpdatePage(int, request.PageRequest) (models.Page, error)
}

type BookServiceImpl struct {
	repo repository.BookRepository
}

func NewBookServiceImpl(repository repository.BookRepository) BookService {
	return &BookServiceImpl{
		repo: repository,
	}
}
func (b *BookServiceImpl) UpdateChapter(chapterId int, request request.BookChapterRequest) (models.Chapter, error) {
	return b.repo.UpdateChapter(chapterId, request)
}
func (b *BookServiceImpl) UpdatePage(pageId int, request request.PageRequest) (models.Page, error) {
	return b.repo.UpdatePage(pageId, request)
}
func (b *BookServiceImpl) DeleteChapter(chapterId int) error {
	return b.repo.DeleteChapter(chapterId)
}

func (b *BookServiceImpl) DeletePage(pageId int) error {
	return b.repo.DeletePage(pageId)
}

func (b *BookServiceImpl) DeleteShelve(shelveId int) error {
	return b.repo.DeleteShelve(shelveId)
}

func (b *BookServiceImpl) UpdateBook(bookId int, request request.BookCreateRequest) (models.Book, error) {
	return b.repo.UpdateBook(bookId, request)
}

func (b *BookServiceImpl) DeleteBook(bookId int) error {
	return b.repo.DeleteBook(bookId)
}

func (b *BookServiceImpl) GetShelves() ([]models.Shelve, error) {
	return b.repo.GetShelves()
}
func (b *BookServiceImpl) CreateCompleteBook(userId int, request request.CompleteBookCreateRequest) (models.Book, error) {
	return b.repo.CreateCompleteBook(userId, request)
}

func (b *BookServiceImpl) GetPageChapter(chapterId int) ([]models.Page, error) {
	return b.repo.GetPageChapter(chapterId)
}

func (b *BookServiceImpl) AddPage(chapterId uint, request request.PageRequest) (models.Page, error) {
	return b.repo.AddPage(chapterId, request)
}

func (b *BookServiceImpl) GetChaptersOfBook(bookId int) ([]models.Chapter, error) {
	return b.repo.GetChaptersOfBook(bookId)
}

func (b *BookServiceImpl) CreateChapter(bookId uint, request request.BookChapterRequest) (models.Chapter, error) {
	return b.repo.CreateChapter(bookId, request)
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
