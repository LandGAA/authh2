package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/LandGAA/authh2/internal/entity"
	"github.com/LandGAA/authh2/pkg/logger"
	"go.uber.org/zap"
)

type Repository interface {
	GetAll() ([]entity.User, error)
	GetByID(id int) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
	Delete(id int) error
	Create(user entity.User) error
	UpdatePassword(user entity.User) error
}

type UserRepository struct {
	db *sql.DB
}

func NewRep(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (u *UserRepository) GetAll() ([]entity.User, error) {
	query := `SELECT * FROM users`
	logger.Logger.Debug("Получение всех пользователей")
	rows, err := u.db.Query(query)
	if err != nil {
		logger.Logger.Error("Ошибка получения пользователей",
			zap.Error(err),
			zap.String("rep", "GetAll"))
		return nil, fmt.Errorf("Ошибка получения пользователей: %w", err)
	}

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreateAt); err != nil {
			msg := fmt.Errorf("Ошибка поиска пользователей: %w", err)

			if errors.Is(err, sql.ErrNoRows) {
				logger.Logger.Error("Ошибка поиска",
					zap.Error(msg),
					zap.String("rep", "GetAll"))
				return nil, msg
			}

			logger.Logger.Error("Ошибка поиска",
				zap.Error(msg),
				zap.String("rep", "GetAll"))
			return nil, msg
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UserRepository) GetByID(id int) (entity.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	row := u.db.QueryRow(query, id)

	var user entity.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreateAt); err != nil {
		msg := fmt.Errorf("Ошибка получения пользователя по ID = %d -> %w", id, err)

		logger.Logger.Error("Ошибка поиска пользователя",
			zap.Error(msg),
			zap.String("rep", "GetByID"))
		return entity.User{}, msg
	}
	return user, nil
}

func (u *UserRepository) GetByEmail(email string) (entity.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	row := u.db.QueryRow(query, email)

	var user entity.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreateAt); err != nil {
		msg := fmt.Errorf("Ошибка получения пользователя по email = %s -> %w", email, err)
		logger.Logger.Error("Ошибка поиска пользователя",
			zap.Error(msg),
			zap.String("rep", "GetByEmail"))
		return entity.User{}, msg
	}
	return user, nil
}

func (u *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	rowEd, err := u.db.Exec(query, id)
	if err != nil {
		msg := fmt.Errorf("Ошибка запроса на удаление пользователя с ID = %d", id)
		logger.Logger.Error("Ошибка запроса на удаление",
			zap.Error(msg),
			zap.String("rep", "Delete"))
		return msg
	}

	affected, err := rowEd.RowsAffected()
	if err != nil {
		msg := fmt.Errorf("Ошибка получения измененных строк при удалении пользователя с ID = %d", id)
		logger.Logger.Error("Ошибка удаления пользователя",
			zap.Error(msg),
			zap.String("rep", "Delete"))
		return msg
	}

	if affected == 0 {
		msg := fmt.Errorf("Ошибка удаления пользователя с ID = %d, пользователя не существует, тк 0 строк", id)
		logger.Logger.Error("Ошибка удаления пользователя",
			zap.Error(msg),
			zap.String("rep", "Delete"))
		return msg
	}
	return nil
}

func (u *UserRepository) Create(user entity.User) error {
	query := `INSERT INTO users (name, email, password, role, create_at)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id`
	err := u.db.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.CreateAt).Scan(&user.ID)
	if err != nil {
		msg := fmt.Errorf("Ошибка при создании пользователя: %w", err)
		logger.Logger.Error("Ошибка создания пользователя",
			zap.Error(msg),
			zap.Any("user", user),
			zap.String("rep", "Create"))
		return msg
	}
	return nil
}

func (u *UserRepository) UpdatePassword(user entity.User) error {
	query := `UPDATE users
			  SET password = $2
		      WHERE id = $1`

	exec, err := u.db.Exec(
		query,
		user.ID,
		user.Password)

	if err != nil {
		msg := fmt.Errorf("Ошибка отправки запроса на обновления данных, %w", err)
		logger.Logger.Error("Ошибка отправки запроса на обновления данных",
			zap.Error(err),
			zap.String("rep", "UpdatePassword"))
		return msg
	}

	affected, err := exec.RowsAffected()
	if err != nil {
		msg := fmt.Errorf("Ошибка получения количества измененных строк при обновлении данных, %w", err)
		logger.Logger.Error("Ошибка получения количества измененных строк",
			zap.Error(err),
			zap.String("rep", "UpdatePassword"))
		return msg
	}

	if affected == 0 {
		msg := fmt.Errorf("0 измененных строк")
		logger.Logger.Error("0 измененных строк",
			zap.String("rep", "UpdatePassword"))
		return msg
	}
	return nil
}
