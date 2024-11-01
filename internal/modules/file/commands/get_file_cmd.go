package commands

import (
	"client-goph-keerper/internal/storage"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// SetGetFileCmd создает команду для получения файла по ID
func SetGetFileCmd(s *storage.Storage) (*cobra.Command, error) {
	getFileCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a file by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			fileID, _ := cmd.Flags().GetString("file_id")

			req, err := http.NewRequest("GET", fmt.Sprintf("%s/file/get/%s", s.ServerURL, fileID), nil)
			if err != nil {
				return fmt.Errorf("ошибка создания запроса: %v", err)
			}

			// Используем токен из структуры storage
			req.Header.Set("Authorization", s.Token)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("ошибка отправки запроса: %v", err)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Printf("Error closing response body: %v", err)
				}
			}(resp.Body)

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
				return fmt.Errorf("ошибка создания файла: %v", err)
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
				return fmt.Errorf("ошибка записи в файл: %v", err)
			}

			fmt.Printf("Файл успешно загружен: %s\n", outputPath)
			return nil
		},
	}

	getFileCmd.Flags().String("file_id", "", "File ID to retrieve")
	err := getFileCmd.MarkFlagRequired("file_id")
	if err != nil {
		return nil, fmt.Errorf("ошибка установки обязательного флага 'file_id': %v", err)
	}

	return getFileCmd, nil
}
