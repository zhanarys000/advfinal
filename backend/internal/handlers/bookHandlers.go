package handlers

import (
	"adv/internal/config"
	"adv/internal/models"
	"adv/internal/repository"
	"adv/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type BookHandlers struct {
	repo   *repository.Repository
	config config.Config
}

func NewBookHandlers(repo *repository.Repository, config config.Config) *BookHandlers {
	return &BookHandlers{repo: repo, config: config}
}

func (h BookHandlers) GetAllBooks(context *gin.Context) {
	logger.GetLogger().Info("Starting to fetch all books")
	books, err := h.repo.GetAllBooks()
	if err != nil {
		logger.GetLogger().Error("Failed to fetch books:", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	logger.GetLogger().Info("All books fetched successfully")
	context.JSON(http.StatusOK, books)
}
func (h BookHandlers) GetBook(context *gin.Context) {
	logger.GetLogger().Info("Starting book getting")
	bookID := context.Param("id")
	var existingBook *models.Book
	existingBook, err := h.repo.GetBook(bookID)
	if err != nil {
		logger.GetLogger().Error("Failed to find book:", err.Error())
		context.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	logger.GetLogger().Info("Book getting successfully")
	context.JSON(http.StatusOK, gin.H{"message": "Book getting successfully", "book": existingBook})
}

func (h BookHandlers) CreateBook(context *gin.Context) {
	logger.GetLogger().Info("Starting book creation")
	var newBook models.Book
	if err := context.ShouldBind(&newBook); err != nil {
		logger.GetLogger().Error("Failed to bind JSON data for new book:", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.CreateBook(&newBook)
	if err != nil {
		logger.GetLogger().Error("Failed to create book:", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}
	logger.GetLogger().Info("Book created successfully")
	context.JSON(http.StatusCreated, gin.H{"message": "Book created successfully", "book": newBook})
}
func (h BookHandlers) UpdateBook(context *gin.Context) {
	logger.GetLogger().Info("Starting book update")
	bookID := context.Param("id")
	existingBook, err := h.repo.GetBook(bookID)
	if err != nil {
		logger.GetLogger().Error("Failed to find book:", err.Error())
		context.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var updatedBook models.Book
	if err := context.ShouldBindJSON(&updatedBook); err != nil {
		logger.GetLogger().Error("Failed to bind JSON data for updating book:", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if updatedBook.Name != "" {
		existingBook.Name = updatedBook.Name
	}
	if updatedBook.Description != "" {
		existingBook.Description = updatedBook.Description
	}
	if updatedBook.Isbn != "" {
		existingBook.Isbn = updatedBook.Isbn
	}
	if updatedBook.Genre != "" {
		existingBook.Genre = updatedBook.Genre
	}
	if updatedBook.Price != "" {
		existingBook.Price = updatedBook.Price
	}

	err = h.repo.UpdateBook(existingBook)
	if err != nil {
		logger.GetLogger().Error("Failed to update book:", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	logger.GetLogger().Info("Book updated successfully")
	context.JSON(http.StatusOK, gin.H{"message": "Book updated successfully", "book": existingBook})
}
func (h BookHandlers) DeleteBook(context *gin.Context) {
	logger.GetLogger().Info("Starting book delete")
	bookID := context.Param("id")

	err := h.repo.DeleteBook(bookID)
	if err != nil {
		logger.GetLogger().Error("Failed to delete book:", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}
	logger.GetLogger().Info("Book deleted successfully")
	context.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func (h *BookHandlers) GetByParams(context *gin.Context) {
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(context.DefaultQuery("pageSize", "10"))
	sortBy := context.DefaultQuery("sortBy", "")
	sortOrder := context.DefaultQuery("sortOrder", "")
	filterBy := context.DefaultQuery("filterBy", "")
	filterValue := context.DefaultQuery("filterValue", "")

	books, err := h.repo.GetBooksByParams(page-1, pageSize, sortBy, sortOrder, filterBy, filterValue)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}
	context.JSON(http.StatusOK, books)
}

func (h *BookHandlers) BuyBook(context *gin.Context) {
	userEmail, emailExists := context.Get("email")
	if !emailExists {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "User email not found"})
		return
	}

	email, ok := userEmail.(string)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert email to string"})
		return
	}
	user, err := h.repo.GetUserByEmail(email)

	bookID, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	book, err := h.repo.GetBook(strconv.FormatUint(bookID, 10))
	for _, enrolledBook := range user.Books {
		if enrolledBook.ID == uint(bookID) {
			context.JSON(http.StatusBadRequest, gin.H{"error": "User is already enrolled in this book"})
			return
		}
	}
	err = h.repo.AddUserToBook(user.ID, book.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll user to book"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Book purchased successfully", "book": book})
}

func (h *BookHandlers) GetUserBooks(context *gin.Context) {
	userEmail, emailExists := context.Get("email")
	if !emailExists {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "User email not found"})
		return
	}

	email, ok := userEmail.(string)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert email to string"})
		return
	}

	user, err := h.repo.GetUserByEmail(email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	userBooks, err := h.repo.GetUserBooks(*user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user books"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"user": user, "books": userBooks})
}
