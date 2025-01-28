package commands

import (
	"bytes"
	"client-goph-keerper/internal/storage"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/spf13/cobra"
)

// RegisterCommand создаёт команду регистрации пользователя.
func RegisterCommand(s *storage.Storage) (*cobra.Command, error) {
	registerCmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Получаем логин и пароль из флагов
			username, err := cmd.Flags().GetString(Username)
			if err != nil {
				return fmt.Errorf("get username: %w", err)
			}
			password, err := cmd.Flags().GetString(Password)
			if err != nil {
				return fmt.Errorf("get password: %w", err)
			}

			// Создаем JSON объект для передачи на сервер
			data := map[string]string{
				Username: username,
				Password: password,
			}

			body, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("ошибка кодирования JSON: %w", err)
			}

			registrationURL := path.Join(s.ServerURL, "registration")
			// Создаём HTTP запрос для регистрации
			req, err := http.NewRequest(http.MethodPost, registrationURL, bytes.NewBuffer(body))
			if err != nil {
				return fmt.Errorf("error creating http request: %w", err)
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("http request: %w", err)
			}
			defer resp.Body.Close() //nolint:errcheck // опустим здесь проверку

			// Читаем ответ от сервера
			all, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("error reading response body: %w", err)
			}

			if resp.StatusCode != http.StatusCreated {
				fmt.Printf(Response, resp.Status)
				fmt.Println(string(all))
				return nil
			}

			// Декодируем JSON ответ, чтобы получить JWT токен
			var result map[string]string
			if err := json.Unmarshal(all, &result); err != nil {
				return fmt.Errorf("ошибка декодирования JSON ответа: %w", err)
			}

			token, ok := result["jwt"]
			if !ok || token == "" {
				return errors.New("токен не найден в ответе")
			}

			// Сохраняем токен в базе данных
			err = saveTokenToDB(s, token)
			if err != nil {
				return fmt.Errorf("ошибка сохранения токена: %w", err)
			}
			fmt.Println("JWT токен:", token)

			fmt.Printf(Response, resp.Status)
			return nil
		},
	}

	registerCmd.Flags().String(Username, "", "Username for the new user")
	registerCmd.Flags().String(Password, "", "Password for the new user")
	err := registerCmd.MarkFlagRequired(Username)
	if err != nil {
		return nil, fmt.Errorf("mark username: %w", err)
	}
	err = registerCmd.MarkFlagRequired(Password)
	if err != nil {
		return nil, fmt.Errorf("mark password: %w", err)
	}

	return registerCmd, nil
}

// saveTokenToDB сохраняет токен в базе данных, используя переданный объект *storage.Storage.
func saveTokenToDB(s *storage.Storage, jwt string) error {
	insertQuery := `INSERT INTO users (jwt) VALUES (?)`
	if _, err := s.DataBase.Exec(insertQuery, jwt); err != nil {
		return fmt.Errorf("ошибка вставки токена в базу данных: %w", err)
	}

	fmt.Println("Токен успешно сохранен в базу данных.")
	return nil
}

// checkTokenExists проверяет наличие токена в базе данных.
func checkTokenExists(s *storage.Storage) (bool, error) {
	query := `SELECT EXISTS(SELECT * FROM users WHERE 1)`
	var exists bool
	err := s.DataBase.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}

	return exists, nil
}
