package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ankit-lilly/newsapp/internal/articles/views/components"
	"github.com/ankit-lilly/newsapp/pkg/llm"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/ollama/ollama/api"
	"strconv"
	"time"
)

type WebsocketMessage struct {
	Chat_mesage string `json:"chat_message"`
}

type Score int
type BiasScore struct {
	Score     Score    `json:"score"`
	Reasoning string   `json:"reasoning"`
	Keywords  []string `json:"keywords"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *Score) UnmarshalJSON(b []byte) error {
	// First, try unmarshalling the value as an integer.
	var intVal int
	if err := json.Unmarshal(b, &intVal); err == nil {
		*s = Score(intVal)
		return nil
	}

	// If that fails, try unmarshalling as a string.
	var strVal string
	if err := json.Unmarshal(b, &strVal); err == nil {
		num, err := strconv.Atoi(strVal)
		if err != nil {
			return fmt.Errorf("Score: cannot convert string %q to int: %v", strVal, err)
		}
		*s = Score(num)
		return nil
	}

	return fmt.Errorf("Score: unable to unmarshal %s", string(b))
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
		sessionid := s.Request.Header.Get("Sec-WebSocket-Key")
		fmt.Println("Connected to chat", sessionid)
		id, ok := s.Keys["articleid"].(int64)
		if !ok {
			s.Write([]byte(fmt.Sprintf("Invalid article id %d", id)))
			return
		}
		cmp := components.Assistant("assistant", "How can I help you today?")

		article, err := a.ArticleService.GetArticleDetail(id)
		if err != nil {
			c.Logger().Error(err.Error())
			s.Write([]byte(err.Error()))
			return
		}

		History[sessionid] = append(History[sessionid], api.Message{
			Role:    "system",
			Content: fmt.Sprintf("You are an AI assistant answering questions about this blog post. Your goal is to help users understand the post. This means answering questions about terms, references, or words mentioned in the post. Do not answer any off-topic questions that are not related to post. For instance if the post is about cars and user asks questions abouts sports then say that it's not related to post and I can't answer it.. \n\nPost content:\n %s", article.Body),
		})

		go func() {
			a.GetBiasScore(article.Body, s)
		}()
		a.WebSocketResponse(s.Request.Context(), cmp, s)
	})
	return nil
}

func (a *ArticleHandler) HandleChatMessage(s *melody.Session, msg []byte) {

	var wsMessage WebsocketMessage
	sessionid := s.Request.Header.Get("Sec-WebSocket-Key")

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
	l := llm.New(a.ollama, "llama3.2")

	errorChan, messageChat := l.ChatRequest(ctx, History[sessionid])
	for {
		select {
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
				continue
			}
			if err != nil {
				fmt.Println("Error", err)
				s.Write([]byte("Invalid message"))
				return // Exit on error
			}

		case resp, ok := <-messageChat:
			if !ok {
				messageChat = nil
				continue
			}
			History[sessionid] = append(History[sessionid], api.Message{
				Role:    resp.Role,
				Content: resp.Content,
			})
			a.WebSocketResponse(s.Request.Context(), components.Assistant(resp.Role, resp.Content), s)
			return
		}

		// Check if both channels are closed
		if messageChat == nil && errorChan == nil {
			return
		}
	}

}

func (a *ArticleHandler) GetBiasScore(article string, s *melody.Session) error {
	// Define the system prompt to instruct the model on how to analyze bias.
	const biasPrompt = `Your job is to analyse news articles for biasness. Consider:
- Language choice and tone
- Source selection and citations
- Context inclusion/omission
- Emotional appeals vs factual reporting
- Use the knowledge you have of the world together with the text
- Balance of perspectives
- Any other factors you think are relevant

You only respond in JSON format with the following fields:
- score: a number between -100 (far left) to 100 (far right) indicating the bias of the article
- reasoning: a brief explanation of the rating
- keywords: an array of key biased terms found in the article

The outout should look like following:
{
  "score": <number between -100 (far left) to 100 (far right)>,
  "reasoning": "<brief explanation of rating>",
  "keywords": ["key", "biased", "terms", "found"]
}

Output only the JSON object. Do not include any other text or fomatting.
`
	prompt := fmt.Sprintf("Analyze this article:  %s", article)

	ctx := s.Request.Context()

	llmHandler := llm.New(a.ollama, "olmo2")

	outputChan, errChan := llmHandler.GenerateRequest(ctx, biasPrompt, prompt, false)

	var response string
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case resp, ok := <-outputChan:
			if !ok {
				outputChan = nil
				continue

			}
			response = resp
			break
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			if err != nil {
				s.Write([]byte(err.Error()))
				return err
			}
			break

		case <-ticker.C:
			s.Write([]byte("Still processing..."))

		}

		if outputChan == nil && errChan == nil {
			break
		}

	}

	var biasScore BiasScore

	if err := json.Unmarshal([]byte(response), &biasScore); err != nil {
		s.Write([]byte(err.Error()))
		return err
	}
	return a.WebSocketResponse(ctx, components.NeutralityIndicator(int(biasScore.Score), biasScore.Reasoning, biasScore.Keywords), s)
}
