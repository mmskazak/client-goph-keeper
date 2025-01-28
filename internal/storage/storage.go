package storage

import (
	"database/sql"
	"errors"
	"fmt"
)

// Storage структура для хранения базы данных и параметров приложения.
type Storage struct {
	DataBase  *sql.DB
	Token     string
	ServerURL string
}

// connectDB открывает соединение с базой данных SQLite и возвращает его.
func connectDB() (*sql.DB, error) {
	// Открываем базу данных SQLite (если файла нет, он будет создан)
	db, err := sql.Open("sqlite", "./keeper.db")
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging db: %w", err)
	}

	return db, nil
}

// runMigrations выполняет миграции для базы данных, создавая нужные таблицы и вставляя значения по умолчанию.
func runMigrations(db *sql.DB) error {
	// Схема для таблицы параметров приложения
	schema := `
	CREATE TABLE IF NOT EXISTS 'app_params' (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
		'key' TEXT UNIQUE,
		'value' TEXT
	);
	-- Вставка значения для ключа 'jwt', если такой записи нет
	INSERT OR IGNORE INTO app_params (key, value)
	VALUES ('jwt', '');
	       
	-- Вставка значения для ключа 'server_url', если такой записи нет
	INSERT OR IGNORE INTO app_params (key, value)
	VALUES ('server_url', '');
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("ошибка выполнения миграции: %w", err)
	}

	return nil
}

// Init инициализирует соединение с базой данных, запускает миграции и возвращает экземпляр Storage.
func Init() (*Storage, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}

	// Запускаем миграции
	if err := runMigrations(db); err != nil {
		return nil, err
	}

	s := &Storage{DataBase: db}
	err = s.loadAppParams()
	if err != nil {
		return nil, fmt.Errorf("error loading app params: %w", err)
	}
	return s, nil
}

// LoadAppParams загружает параметры приложения из базы данных и устанавливает их в структуре Storage.
func (s *Storage) loadAppParams() error {
	// Извлекаем токен и URL сервера из базы данных
	var jwt, serverURL string
	tokenQuery := `SELECT value FROM app_params WHERE key = 'jwt' LIMIT 1`
	err := s.DataBase.QueryRow(tokenQuery).Scan(&jwt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("токен не найден в базе данных")
		}
		return fmt.Errorf("ошибка получения токена: %w", err)
	}

	serverURLQuery := `SELECT value FROM app_params WHERE key = 'server_url' LIMIT 1`
	err = s.DataBase.QueryRow(serverURLQuery).Scan(&serverURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("адрес сервера не найден в базе данных: %w", err)
		}
		return fmt.Errorf("ошибка получения адреса сервера: %w", err)
	}

	// Устанавливаем значения в структуре
	s.Token = jwt
	s.ServerURL = serverURL

	return nil
}

// RemoveTokenFromDB удаляет токен из базы данных.
func (s *Storage) RemoveTokenFromDB() error {
	// Удаляем токен из таблицы app_params
	deleteQuery := `DELETE FROM app_params WHERE key = ?`
	if _, err := s.DataBase.Exec(deleteQuery, "jwt_token"); err != nil {
		return fmt.Errorf("ошибка удаления jwt токена: %w", err)
	}
	return nil
}
