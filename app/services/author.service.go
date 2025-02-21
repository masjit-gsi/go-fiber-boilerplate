package services

import (
	"github.com/fiber-go-template/app/models"
	"github.com/fiber-go-template/app/repository"
	"github.com/fiber-go-template/config/logger"
	"github.com/fiber-go-template/database"
	"github.com/fiber-go-template/helper/pagination"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type AuthorService interface {
	ResolveAll(req models.StandardRequest) (data pagination.Response, err error)
	GetAll() (author []models.Author, err error)
	FindByID(id uuid.UUID) (author models.Author, err error)
	Create(req models.AuthorRequest) (res models.Author, err error)
	Update(id uuid.UUID, req models.AuthorRequest) (res models.Author, err error)
	Delete(id uuid.UUID, userID uuid.UUID) (err error)
}
type AuthorServiceImpl struct {
	DB               database.DBConn
	AuthorRepository repository.AuthorRepository
}

func NewAuthorService(db database.DBConn, author repository.AuthorRepository) *AuthorServiceImpl {
	return &AuthorServiceImpl{
		DB:               db,
		AuthorRepository: author,
	}
}

func (s *AuthorServiceImpl) ResolveAll(req models.StandardRequest) (data pagination.Response, err error) {
	return s.AuthorRepository.ResolveAll(req)
}

func (s *AuthorServiceImpl) GetAll() (res []models.Author, err error) {
	var author models.Author
	err = s.DB.Orm().Table(author.TableName()).
		Select("id", "name", "address").
		Where("coalesce(is_deleted) = false").
		Order("name asc").Scan(&res).Error
	if res == nil {
		return make([]models.Author, 0), nil
	}
	return
}

func (s *AuthorServiceImpl) FindByID(id uuid.UUID) (author models.Author, err error) {
	err = s.DB.Orm().First(&author, "id=?", id).Error
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (s *AuthorServiceImpl) Create(req models.AuthorRequest) (res models.Author, err error) {
	res.BindFromRequest(req)
	err = s.DB.Orm().Create(&res).Error
	if err != nil {
		logger.ErrorWithStack(err)
		return models.Author{}, err
	}

	return
}

func (s *AuthorServiceImpl) Update(id uuid.UUID, req models.AuthorRequest) (res models.Author, err error) {
	var author models.Author
	err = s.DB.Orm().First(&author, "id=?", id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.ErrorWithStack(err)
		return
	}

	author.BindFromRequest(req)
	err = s.DB.Orm().Save(&author).Error
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return author, nil
}

func (s *AuthorServiceImpl) Delete(id uuid.UUID, userID uuid.UUID) (err error) {
	var author models.Author
	err = s.DB.Orm().First(&author, "id=?", id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.ErrorWithStack(err)
		return
	}

	author.SoftDelete(userID)
	err = s.DB.Orm().Save(&author).Error
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return nil
}
