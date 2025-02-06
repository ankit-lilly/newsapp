package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ankibahuguna/newsapp/internal/articles/views/components"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/ollama/ollama/api"
	"strconv"
)

type WebsocketMessage struct {
	Chat_mesage string `json:"chat_message"`
}

var History map[string][]api.Message = make(map[string][]api.Message)

func (a *ArticleHandler) Chat(c echo.Context) error {
	c.Logger().Info("Connected to chat", c.Param("id"))

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Internal server error")

	}
	chatMessageHandlerMap := make(map[string]any)

	chatMessageHandlerMap["articleid"] = id

	a.m.HandleRequestWithKeys(c.Response().Writer, c.Request(), chatMessageHandlerMap)

	a.m.HandleDisconnect(func(s *melody.Session) {
		sessionid := s.Request.Header.Get("Sec-WebSocket-Key")
		c.Logger().Info("Disconnected from chat", c.Param("id"), sessionid)
		delete(History, sessionid)
	})

	a.m.HandleConnect(func(s *melody.Session) {
		c.Logger().Info("Connected to chat", c.Param("id"))
		cmp := components.Assistant("assistant", "How can I help you today?")
		a.WebSocketResponse(s.Request.Context(), cmp, s)
	})
	return nil

}

func (a *ArticleHandler) HandleChatMessage(s *melody.Session, msg []byte) {

	id, ok := s.Keys["articleid"].(int64)
	sessionid := s.Request.Header.Get("Sec-WebSocket-Key")

	if !ok {
		s.Write([]byte(fmt.Sprintf("Invalid article id %d", id)))
		return
	}

	article, err := a.ArticleService.GetArticleDetail(id)

	if err != nil {
		s.Write([]byte(err.Error()))
		return
	}

	History[sessionid] = append(History[sessionid], api.Message{
		Role:    "system",
		Content: fmt.Sprintf("You are an AI assistant answering questions about this blog post. Your goal is to help users understand the post. This means answering questions about terms, references, or words mentioned in the post. Do not answer any off-topic questions that are not related to post. For instance if the post is about cars and user asks questions abouts sports then say that it's not related to post and I can't answer it.. \n\nPost content:\n %s", article.Body),
	})

	var wsMessage WebsocketMessage

	if err := json.Unmarshal(msg, &wsMessage); err != nil {
		s.Write([]byte("Invalid message"))
		return
	}

	History[sessionid] = append(History[sessionid], api.Message{
		Role:    "user",
		Content: wsMessage.Chat_mesage,
	})

	a.WebSocketResponse(s.Request.Context(), components.User("user", wsMessage.Chat_mesage), s)
	//To show the user loading indicatory while the AI is thinking
	a.WebSocketResponse(s.Request.Context(), components.AssistantLoader(), s)

	ctx := context.Background()
	stream := false
	req := &api.ChatRequest{
		Model:    "llama3.2",
		Messages: History[sessionid],
		Stream:   &stream,
	}

	/*
		fmt.Println("Active sessions")
		for k := range History {
			fmt.Println(k)
		}

		fmt.Println("Active sessions")
	*/
	respFunc := func(resp api.ChatResponse) error {
		if resp.Done {
			currentMessage := resp.Message

			History[sessionid] = append(History[sessionid], api.Message{
				Role:    currentMessage.Role,
				Content: currentMessage.Content,
			})
			a.WebSocketResponse(s.Request.Context(), components.Assistant(currentMessage.Role, resp.Message.Content), s)
		}

		return nil
	}

	err = a.ollama.Chat(ctx, req, respFunc)

	if err != nil {
		fmt.Println(err)
		s.Write([]byte("Invalid message"))
		return
	}
}
