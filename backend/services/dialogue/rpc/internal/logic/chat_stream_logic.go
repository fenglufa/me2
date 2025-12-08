package logic

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/me2/dialogue/rpc/dialogue"
	"github.com/me2/dialogue/rpc/internal/model"
	"github.com/me2/dialogue/rpc/internal/svc"

	"github.com/me2/ai/rpc/ai_client"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatStreamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatStreamLogic {
	return &ChatStreamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChatStreamLogic) ChatStream(in *dialogue.ChatStreamRequest, stream dialogue.Dialogue_ChatStreamServer) error {
	session, err := l.svcCtx.SessionModel.FindOne(in.SessionId)
	if err != nil {
		return err
	}

	if session.UserId != in.UserId {
		return fmt.Errorf("permission denied")
	}

	builder := NewContextBuilder(l.ctx, l.svcCtx)

	systemPrompt, err := builder.BuildSystemPrompt(in.AvatarId, in.SessionId)
	if err != nil {
		return err
	}

	userMsgID, err := builder.SaveMessage(in.SessionId, "user", in.Message)
	if err != nil {
		return err
	}

	if session.Title == "" {
		session.Title = model.BuildSessionTitle(in.Message)
		session.UpdatedAt = getCurrentTimestamp()
		l.svcCtx.SessionModel.Update(session)
	}

	variables := map[string]string{
		"system_prompt": systemPrompt,
		"user_message":  in.Message,
	}

	l.Infof("Calling AI service: template=%s, user_id=%d, avatar_id=%d", "avatar_chat", in.UserId, in.AvatarId)
	aiStream, err := l.svcCtx.AiRpc.ChatStream(l.ctx, &ai_client.ChatRequest{
		PromptTemplate: "avatar_chat",
		Variables:      variables,
		UserId:         in.UserId,
		AvatarId:       in.AvatarId,
	})
	if err != nil {
		l.Errorf("AI service call failed: %v", err)
		return err
	}
	l.Info("AI service call succeeded, starting to receive stream")

	var fullResponse strings.Builder
	chunkCount := 0

	for {
		l.Infof("Waiting for chunk %d from AI stream...", chunkCount)
		resp, err := aiStream.Recv()
		if err == io.EOF {
			l.Infof("Received EOF from AI stream after %d chunks", chunkCount)
			break
		}
		if err != nil {
			l.Errorf("Error receiving from AI stream: %v", err)
			return err
		}

		chunkCount++
		l.Infof("Received chunk %d: content_len=%d, done=%v", chunkCount, len(resp.Content), resp.Done)

		fullResponse.WriteString(resp.Content)

		if err := stream.Send(&dialogue.ChatStreamResponse{
			Content: resp.Content,
			Done:    resp.Done,
		}); err != nil {
			return err
		}

		if resp.Done {
			l.Infof("Stream marked as done after %d chunks", chunkCount)
			break
		}
	}

	assistantMsgID, err := builder.SaveMessage(in.SessionId, "assistant", fullResponse.String())
	if err != nil {
		return err
	}

	builder.UpdateSessionLastMessage(in.SessionId, fullResponse.String())

	if err := stream.Send(&dialogue.ChatStreamResponse{
		Content:   "",
		Done:      true,
		MessageId: assistantMsgID,
	}); err != nil {
		return err
	}

	l.Infof("Chat completed: user_msg=%d, assistant_msg=%d, session=%d", userMsgID, assistantMsgID, in.SessionId)

	return nil
}
