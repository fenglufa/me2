package dialogue

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/me2/dialogue/rpc/dialogue_client"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type StreamMessage struct {
	SessionID int64  `json:"session_id"`
	Message   string `json:"message"`
}

type StreamResponse struct {
	Content   string `json:"content"`
	Done      bool   `json:"done"`
	MessageID int64  `json:"message_id,omitempty"`
}

func StreamHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(int64)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		avatarIDStr := r.URL.Query().Get("avatar_id")
		if avatarIDStr == "" {
			http.Error(w, "avatar_id required", http.StatusBadRequest)
			return
		}

		avatarID, err := strconv.ParseInt(avatarIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid avatar_id", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Errorf("WebSocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		logx.Infof("WebSocket connected: user_id=%d, avatar_id=%d", userID, avatarID)

		for {
			var msg StreamMessage
			if err := conn.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logx.Errorf("WebSocket read error: %v", err)
				}
				break
			}

			if msg.SessionID == 0 || msg.Message == "" {
				conn.WriteJSON(map[string]string{"error": "invalid message"})
				continue
			}

			ctx := context.Background()
			stream, err := svcCtx.DialogueRpc.ChatStream(ctx, &dialogue_client.ChatStreamRequest{
				SessionId: msg.SessionID,
				Message:   msg.Message,
				UserId:    userID,
				AvatarId:  avatarID,
			})
			if err != nil {
				logx.Errorf("ChatStream RPC failed: %v", err)
				conn.WriteJSON(map[string]string{"error": err.Error()})
				continue
			}

			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					logx.Errorf("Stream receive error: %v", err)
					conn.WriteJSON(map[string]string{"error": err.Error()})
					break
				}

				response := StreamResponse{
					Content:   resp.Content,
					Done:      resp.Done,
					MessageID: resp.MessageId,
				}

				if err := conn.WriteJSON(response); err != nil {
					logx.Errorf("WebSocket write error: %v", err)
					return
				}

				if resp.Done {
					break
				}
			}
		}

		logx.Infof("WebSocket disconnected: user_id=%d", userID)
	}
}

func StreamHandlerWithAuth(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			token = r.Header.Get("Authorization")
		}

		if token == "" {
			http.Error(w, "token required", http.StatusUnauthorized)
			return
		}

		claims, err := parseToken(token, svcCtx.Config.Auth.AccessSecret)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		r = r.WithContext(ctx)

		StreamHandler(svcCtx)(w, r)
	}
}

type Claims struct {
	UserID int64 `json:"user_id"`
}

func parseToken(tokenString, secret string) (*Claims, error) {
	var claims Claims
	if err := json.Unmarshal([]byte(tokenString), &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}
