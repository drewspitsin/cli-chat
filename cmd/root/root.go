package root

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	desc "github.com/drewspitsin/cli-chat/pkg/chat_api_v1"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50052"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chat-app",
	Short: "cli-chat",
}

var testCmd = &cobra.Command{
	Use:   "spam",
	Short: "TestSpam",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("failed to connect to server: %v", err)
		}
		defer conn.Close()

		ctx := context.Background()
		client := desc.NewChatV1Client(conn)

		// Создаем новый чат на сервере
		chatID, err := createChat(ctx, client)
		if err != nil {
			log.Fatalf("failed to create chat: %v", err)
		}

		log.Printf(fmt.Sprintf("%s: %s\n", color.GreenString("Chat created"), color.YellowString(chatID)))

		wg := sync.WaitGroup{}
		wg.Add(2)

		// Подключаемся к чату от имени пользователя oleg
		go func() {
			defer wg.Done()

			err = connectChat(ctx, client, chatID, "oleg", 5*time.Second)
			if err != nil {
				log.Fatalf("failed to connect chat: %v", err)
			}
		}()

		// Подключаемся к чату от имени пользователя ivan
		go func() {
			defer wg.Done()

			err = connectChat(ctx, client, chatID, "ivan", 7*time.Second)
			if err != nil {
				log.Fatalf("failed to connect chat: %v", err)
			}
		}()

		wg.Wait()
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Cоздание",
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Удаление",
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Создает нового пользователя",
	Run: func(cmd *cobra.Command, args []string) {
		usernamesStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("failed to get usernames: %s\n", err.Error())
		}

		log.Printf("user %s created\n", usernamesStr)
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Удаляет пользователя",
	Run: func(cmd *cobra.Command, args []string) {
		usernamesStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("failed to get usernames: %s\n", err.Error())
		}

		log.Printf("user %s deleted\n", usernamesStr)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(testCmd)

	createCmd.AddCommand(createUserCmd)
	deleteCmd.AddCommand(deleteUserCmd)

	createUserCmd.Flags().StringP("username", "u", "", "Имя пользователя")
	err := createUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())
	}

	deleteUserCmd.Flags().StringP("username", "u", "", "Имя пользователя")
	err = deleteUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())
	}
}
