package postgres

import (
	"gorm.io/gorm"

	"github.com/kokhno-nikolay/letsgochat/models"
)

type UsersRepo struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) FindById(id int) (models.User, error) {
	user := models.User{}

	err := r.db.First(&user, id).Error
	return user, err
}

func (r *UsersRepo) FindByUsername(username string) (models.User, error) {
	user := models.User{}

	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *UsersRepo) Create(user models.User) error {
	return r.db.Create(&user).Error
}

func (r *UsersRepo) UserExists(username string) (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count, err
}

func (r *UsersRepo) GetAllActive() ([]models.User, error) {
	users := make([]models.User, 0)

	err := r.db.Where("active = true").Find(&users).Error
	return users, err
}

func (r *UsersRepo) SwitchToActive(userID int) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("active", true).Error
}

func (r *UsersRepo) SwitchToInactive(userID int) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("active", false).Error
}
