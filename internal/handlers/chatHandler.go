package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/a-h/templ"
	"github.com/ankit-lilly/newsapp/internal/prompts"
	"github.com/ankit-lilly/newsapp/internal/services"

	"github.com/ankit-lilly/newsapp/internal/templates/components/articles"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/ollama/ollama/api"
)

type ChatHandler struct {
	*BaseHandler
	articleService *services.ArticleService
	ws             *melody.Melody
}

func NewChatHandler(articleService *services.ArticleService, ws *melody.Melody) *ChatHandler {
	return &ChatHandler{
		articleService: articleService,
		ws:             ws,
	}
}

type WebsocketMessage struct {
	Chat_mesage string `json:"chat_message"`
	Portal      string `json:"portal"`
	ArticleId   string `json:"articleid"`
}

func (h *ChatHandler) HandleConnect(s *melody.Session) {
	sessionid := s.Request.Header.Get("Sec-WebSocket-Key")
	slog.Debug("Connected to chat, sessionid:", sessionid)
	link, portalName, err := h.extractKeysFromSession(s)
	if err != nil {
		slog.Error(err.Error(), err)
		h.WebSocketResponse(s.Request.Context(), articles.Assistant("assistant", "I'm sorry, I'm having trouble processing your request. Please try again."), s)
		return
	}

	greeting := api.Message{
		Role:    "assistant",
		Content: "How can I help you today?",
	}

	articleDetail, err := h.articleService.GetRawArticle(s.Request.Context(), portalName, link)

	if err != nil {
		slog.Error(err.Error(), err)
		h.WebSocketResponse(s.Request.Context(), articles.Assistant("assistant", "I'm sorry, I'm having trouble processing your request. Please try again."), s)
		return
	}

	s.Keys["history"] = prompts.GetChatPrompt(articleDetail)

	history := s.Keys["history"].([]api.Message)

	if len(history) <= 2 {
		cmp := articles.Assistant(greeting.Role, greeting.Content)
		h.WebSocketResponse(s.Request.Context(), cmp, s)
		history = append(history, greeting)
		s.Keys["history"] = history
	}
}

func (h *ChatHandler) HandleDisconnect(s *melody.Session) {
	delete(s.Keys, "history")
	s.Write([]byte("disconnected"))
}

func (h *ChatHandler) Chat(c echo.Context) error {
	link, portalName, err := h.parseAndValidateIdAndPortal(c)
	if err != nil {
		return err
	}

	chatMessageHandlerMap := make(map[string]any)
	chatMessageHandlerMap["portal"] = portalName
	chatMessageHandlerMap["articleid"] = link

	return h.ws.HandleRequestWithKeys(c.Response().Writer, c.Request(), chatMessageHandlerMap)
}

func (a *ChatHandler) HandleChatMessage(s *melody.Session, msg []byte) {
	var wsMessage WebsocketMessage

	if err := json.Unmarshal(msg, &wsMessage); err != nil {
		s.Write([]byte("Invalid message"))
		return
	}

	slog.Info("Received message:", wsMessage)

	historyRaw, ok := s.Keys["history"]
	if !ok {
		s.Write([]byte("Invalid session"))
		return
	}

	history, ok := historyRaw.([]api.Message)
	if !ok || len(history) == 0 {
		s.Write([]byte("Invalid session"))
		return
	}

	history = append(history, api.Message{
		Role:    "user",
		Content: wsMessage.Chat_mesage,
	})
	s.Keys["history"] = history

	a.WebSocketResponse(s.Request.Context(), articles.User("user", wsMessage.Chat_mesage), s)
	a.WebSocketResponse(s.Request.Context(), articles.AssistantLoader(), s)

	slog.Debug("Chat history:", history)

	resp, err := a.articleService.SendChatRequest(s.Request.Context(), history)

	if err != nil {
		slog.Error(err.Error(), err)
		a.WebSocketResponse(s.Request.Context(), articles.Assistant("assistant", "I'm sorry, I'm having trouble processing your request. Please try again."), s)
		return
	}

	history = append(history, api.Message{
		Role:    resp.Role,
		Content: resp.Content,
	})

	s.Keys["history"] = history

	formatedContent := strings.Split(resp.Content, "\n")
	paragraphWrapped := make([]string, 0)

	for _, content := range formatedContent {
		paragraphWrapped = append(paragraphWrapped, fmt.Sprintf("<p class='text-lg mt-4'>%s</p>", strings.TrimSpace(content)))
	}

	assistantResponse := articles.Assistant(resp.Role, strings.Join(paragraphWrapped, ""))
	a.WebSocketResponse(s.Request.Context(), assistantResponse, s)
}

func (h *ChatHandler) extractKeysFromSession(s *melody.Session) (string, string, error) {
	portal, ok := s.Keys["portal"].(string)
	if !ok {
		return "", "", errors.New("portal is required")
	}

	articleid, ok := s.Keys["articleid"].(string)
	if !ok {
		return "", "", errors.New("articleid is required")
	}

	return articleid, portal, nil
}

func (h *ChatHandler) parseAndValidateIdAndPortal(c echo.Context) (string, string, error) {
	encodedLink := strings.TrimSpace(c.Param("id"))
	portalName := strings.TrimSpace(c.Param("portal"))

	if encodedLink == "" {
		return "", "", errors.New("id is required")
	}
	if portalName == "" {
		return "", "", errors.New("portal is required")
	}

	link, err := url.QueryUnescape(encodedLink)
	if err != nil {
		c.Echo().Logger.Error(err.Error())
		return "", "", errors.New("Invalid link")
	}

	return link, portalName, nil
}

func (a *ChatHandler) WebSocketResponse(ctx context.Context, cmp templ.Component, session *melody.Session) error {
	var buffer bytes.Buffer
	cmp.Render(ctx, &buffer)
	return session.Write(buffer.Bytes())
}
