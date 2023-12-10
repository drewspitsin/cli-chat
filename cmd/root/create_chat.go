package root

import (
	"context"

	desc "github.com/drewspitsin/cli-chat/pkg/chat_api_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func createChat(ctx context.Context, client desc.ChatV1Client) (string, error) {
	res, err := client.CreateChat(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}

	return res.GetChatId(), nil
}
