package commands

import (
	"client-goph-keerper/internal/storage"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// SetGetFileCmd создает команду для получения файла по ID.
func SetGetFileCmd(s *storage.Storage) (*cobra.Command, error) {
	getFileCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a file by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			fileID, err := cmd.Flags().GetString(FileID)
			if err != nil {
				return fmt.Errorf("get file id: %w", err)
			}

			reqURL := path.Join(s.ServerURL, File, "get", fileID)
			req, err := http.NewRequest(http.MethodGet, reqURL, http.NoBody)
			if err != nil {
				return fmt.Errorf("ошибка создания запроса: %w", err)
			}

			// Используем токен из структуры storage
			req.Header.Set(Authorization, s.Token)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf(ErrSendRequest, err)
			}
			defer resp.Body.Close() //nolint:errcheck //опустим проверку

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("не удалось получить файл, статус: %v", resp.Status)
			}

			// Извлечение имени файла из заголовка Content-Disposition
			var fileName string
			contentDisposition := resp.Header.Get("Content-Disposition")
			if contentDisposition != "" {
				// Парсим имя файла из заголовка
				if parts := strings.Split(contentDisposition, "filename="); len(parts) > 1 {
					fileName = strings.Trim(parts[1], "\"")
				}
			}

			// Если имя файла не указано, используем дефолтное имя
			if fileName == "" {
				fileName = "downloaded_file"
			}

			// Создаем файл в корне вызывающего приложения
			outputPath := filepath.Join(".", fileName)
			file, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("ошибка создания файла: %w", err)
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					log.Printf("Error closing file: %v", err)
				}
			}(file)

			// Копируем содержимое ответа в файл
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				return fmt.Errorf("ошибка записи в файл: %w", err)
			}

			fmt.Printf("Файл успешно загружен: %s\n", outputPath)
			return nil
		},
	}

	getFileCmd.Flags().String(FileID, "", "File ID to retrieve")
	err := getFileCmd.MarkFlagRequired(FileID)
	if err != nil {
		return nil, fmt.Errorf("ошибка установки обязательного флага 'file_id': %w", err)
	}

	return getFileCmd, nil
}
