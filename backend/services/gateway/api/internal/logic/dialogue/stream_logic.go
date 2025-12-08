package dialogue

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/me2/dialogue/rpc/dialogue_client"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type StreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StreamLogic {
	return &StreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (l *StreamLogic) Stream(w http.ResponseWriter, r *http.Request) error {
	userIDVal := l.ctx.Value("user_id")
	if userIDVal == nil {
		return errors.New("未授权")
	}
	userID := userIDVal.(int64)

	avatarIDStr := r.URL.Query().Get("avatar_id")
	avatarID, err := strconv.ParseInt(avatarIDStr, 10, 64)
	if err != nil {
		return err
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var req struct {
			SessionID int64  `json:"session_id"`
			Content   string `json:"content"`
		}
		if err := json.Unmarshal(message, &req); err != nil {
			continue
		}

		stream, err := l.svcCtx.DialogueRpc.ChatStream(l.ctx, &dialogue_client.ChatStreamRequest{
			SessionId: req.SessionID,
			Message:   req.Content,
			UserId:    userID,
			AvatarId:  avatarID,
		})
		if err != nil {
			conn.WriteJSON(map[string]interface{}{"error": err.Error()})
			continue
		}

		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				conn.WriteJSON(map[string]interface{}{"error": err.Error()})
				break
			}

			if err := conn.WriteJSON(map[string]interface{}{
				"content":    chunk.Content,
				"done":       chunk.Done,
				"message_id": chunk.MessageId,
			}); err != nil {
				break
			}
		}
	}

	return nil
}
