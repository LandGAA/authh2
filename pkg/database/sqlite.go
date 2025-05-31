package database

import (
	"database/sql"
	"github.com/LandGAA/authh2/pkg/jwt"
	"github.com/LandGAA/authh2/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"log"
	_ "modernc.org/sqlite"
)

func Connect() *sql.DB {
	logger.Logger.Debug("SQLite подключается...")
	jwt.Init()

	db, err := sql.Open("sqlite", "../database.db")
	if err != nil {
		logger.Logger.Fatal("Ошибка подключения к базе данных!",
			zap.Error(err),
			zap.String("db", "connect"))
		return nil
	}

	if err := db.Ping(); err != nil {
		logger.Logger.Fatal("Нет отклика от базы данных",
			zap.Error(err),
			zap.String("db", "connect"))
		db.Close()
	}
	logger.Logger.Info("Успешное подключение к SQLite!")

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Ошибка создания драйвера миграции: %v\n", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3", driver,
	)
	if err != nil {
		log.Fatalf("ошибка создания миграции: %v\n", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("ошибка получения версии: %v\n", err)
	}

	if dirty {
		log.Println("Обнаружена грязная база данных")
		err = m.Force(int(version))
		if err != nil {
			log.Fatalf("не удалось исправить грязную версию: %v\n", err)
		}
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("нет новых миграций для выполнения")
		} else {
			log.Fatalf("ошибка выполнения миграции: %v\n", err)
		}
	} else {
		log.Println("миграция выполнена успешно")
	}

	return db
}
