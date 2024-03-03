package repository

import (
	"adv/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
func (r *Repository) GetBooksByParams(pageNum, pageSize int, sortBy, sortOrder, filterBy, filterValue string) ([]models.Book, error) {
	var orderStr string
	if sortBy != "" && sortOrder != "" {
		orderStr = sortBy + " " + sortOrder
	}
	var filterMap map[string]interface{}
	if filterBy != "" && filterValue != "" {
		filterMap = map[string]interface{}{filterBy: filterValue}
	}
	var books []models.Book
	result := r.db.Model(&models.Book{}).
		Where(filterMap).
		Order(orderStr).
		Limit(pageSize).
		Offset(pageNum * pageSize).
		Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *Repository) GetUser(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *Repository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (r *Repository) CreateUser(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) UpdateUser(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}
func (r *Repository) DeleteUser(id uint) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetBook(id string) (*models.Book, error) {
	var book models.Book
	if err := r.db.First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}
func (r *Repository) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	if err := r.db.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}
func (r *Repository) GetBookByParams() {

}
func (r *Repository) CreateBook(book *models.Book) error {
	if err := r.db.Create(book).Error; err != nil {
		return err
	}
	return nil
}
func (r *Repository) UpdateBook(book *models.Book) error {
	if err := r.db.Save(book).Error; err != nil {
		return err
	}
	return nil
}
func (r *Repository) DeleteBook(id string) error {
	if err := r.db.Delete(&models.Book{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateSubscriber(subscriber *models.Subscribers) error {
	if err := r.db.Create(subscriber).Error; err != nil {
		return err
	}
	return nil
}
func (r *Repository) GetSubscribers() ([]models.Subscribers, error) {
	var subscribers []models.Subscribers
	if err := r.db.Find(&subscribers).Error; err != nil {
		return nil, err
	}
	return subscribers, nil
}
func (r *Repository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *Repository) CheckSubscriptionByEmail(email string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Subscribers{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) GetUserBooks(user models.User) ([]models.Book, error) {
	var userBooks []models.Book
	if err := r.db.Model(&user).Association("Books").Find(&userBooks); err != nil {
		return nil, err
	}
	return userBooks, nil
}

func (r *Repository) AddUserToBook(userID uint, bookID uint) error {
	userBook := models.UserBook{
		UserID: userID,
		BookID: bookID,
	}
	if err := r.db.Create(&userBook).Error; err != nil {
		return err
	}
	return nil
}
