package main

import (
	"client-goph-keerper/internal/modules/auth"
	"client-goph-keerper/internal/modules/connecttoserver"
	"client-goph-keerper/internal/modules/file"
	"client-goph-keerper/internal/modules/passwords"
	"client-goph-keerper/internal/storage"
	"log"

	"github.com/spf13/cobra"
)

func main() {
	// Инициализация хранилища
	s, err := storage.Init()
	if err != nil {
		log.Fatalf("Ошибка инициализации хранилища: %v", err)
	}

	// Команды для первоначальной настройки
	installingCmd, err := connecttoserver.StartsCmd(s)
	if err != nil {
		log.Fatalf("Ошибка установки начальных команд: %v", err)
	}

	// Команды для аутентификации
	authCmd, err := auth.InitAuthCmd(s)
	if err != nil {
		log.Fatalf("Ошибка инициализации команд аутентификации: %v", err)
	}

	// Команды для управления файлами
	fileCmd, err := file.InitFileCmd(s)
	if err != nil {
		log.Fatalf("Ошибка инициализации команд управления файлами: %v", err)
	}

	// Команды для управления паролями
	pwdCmd, err := passwords.InitPwdCmd(s)
	if err != nil {
		log.Fatalf("Ошибка инициализации команд управления паролями: %v", err)
	}

	var rootCmd = &cobra.Command{Use: "app"}

	// Добавляем команды
	rootCmd.AddCommand(installingCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(fileCmd)
	rootCmd.AddCommand(pwdCmd)

	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
