package root

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	desc "github.com/drewspitsin/cli-chat/pkg/chat_api_v1"
	"github.com/fatih/color"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func connectChat(ctx context.Context, client desc.ChatV1Client, chatID string, username string, period time.Duration) error {
	stream, err := client.ConnectChat(ctx, &desc.ConnectChatRequest{
		ChatId:   chatID,
		Username: username,
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			message, errRecv := stream.Recv()
			if errRecv == io.EOF {
				return
			}
			if errRecv != nil {
				log.Println("failed to receive message from stream: ", errRecv)
				return
			}

			log.Printf("[%v] - [from: %s]: %s\n",
				color.YellowString(message.GetCreatedAt().AsTime().Format(time.RFC3339)),
				color.BlueString(message.GetFrom()),
				message.GetText(),
			)
		}
	}()

	for {
		time.Sleep(period)

		text := gofakeit.Word()

		_, err = client.SendMessage(ctx, &desc.SendMessageRequest{
			ChatId: chatID,
			Message: &desc.Message{
				From:      username,
				Text:      text,
				CreatedAt: timestamppb.Now(),
			},
		})
		if err != nil {
			log.Println("failed to send message: ", err)
			return err
		}
	}
}
