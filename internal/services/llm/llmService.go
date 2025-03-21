package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ankit-lilly/newsapp/internal/prompts"
	"github.com/ollama/ollama/api"
)

type LLMHandler struct {
	ollama *api.Client
	model  string
}

func New(ollamaClient *api.Client, model string) *LLMHandler {
	return &LLMHandler{
		ollama: ollamaClient,
		model:  model,
	}
}

func (llm *LLMHandler) UpdateModel(model string) {
	llm.model = model
}

func (llm *LLMHandler) GenerateRequest(ctx context.Context, system, prompt string, stream bool) (<-chan string, <-chan error) {
	output := make(chan string)
	errChan := make(chan error, 1)

	req := &api.GenerateRequest{
		System: system,
		Model:  llm.model,
		Prompt: prompt,
		Stream: &stream,
	}

	var callback func(api.GenerateResponse) error

	if stream {
		callback = func(resp api.GenerateResponse) error {
			output <- resp.Response
			if resp.Done {
				close(output)
				return nil
			}
			return nil
		}
	} else {
		callback = func(resp api.GenerateResponse) error {
			if resp.Done {
				defer close(output)
				output <- resp.Response
			}
			return nil
		}
	}

	// Call the API in a separate goroutine to avoid blocking.
	go func() {
		if err := llm.ollama.Generate(ctx, req, callback); err != nil {
			errChan <- err
			close(errChan)
		} else {
			close(errChan)
		}
	}()

	return output, errChan
}

type ArticleQuality struct {
	Content string
	Rating  int
}

func (llm *LLMHandler) GetRating(ctx context.Context, content string) (chan api.Message, chan error) {

	format := json.RawMessage(`{
    "type": "object",
    "properties": {
      "summary": {
        "type": "string"
      },
      "rating": {
        "type": "number"
      },
      "keywords": {
        "type": "array",
        "items": {
          "type": "string"
        }
      }
    },
    "required": [
      "summary",
      "rating", 
      "keywords"
    ]

  }`)

	stream := false
	req := &api.ChatRequest{
		Model:  llm.model,
		Format: format,
		Stream: &stream,
		Messages: []api.Message{
			{
				Role:    "system",
				Content: prompts.ARTICLE_QUALITY,
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Evaluate the quality of the following article: %q", content),
			},
		},
	}

	outputChan := make(chan api.Message)
	errorsChan := make(chan error, 1)

	chatResponseCallback := func(resp api.ChatResponse) error {
		if resp.Done {
			defer close(outputChan)
			outputChan <- resp.Message
		}
		return nil
	}
	err := llm.ollama.Chat(ctx, req, chatResponseCallback)
	if err != nil {
		defer close(errorsChan)
		errorsChan <- err
	}

	return outputChan, errorsChan
}

func (llm *LLMHandler) ChatRequest(ctx context.Context, messages []api.Message) (<-chan error, <-chan api.Message) {
	stream := false
	outputChan := make(chan api.Message)
	errorsChan := make(chan error, 1)

	req := &api.ChatRequest{
		Model:    llm.model,
		Messages: messages,
		Stream:   &stream,
	}

	chatResponseCallback := func(resp api.ChatResponse) error {
		if resp.Done {
			defer close(outputChan)
			outputChan <- resp.Message
		}
		return nil
	}

	go func() {
		err := llm.ollama.Chat(ctx, req, chatResponseCallback)
		if err != nil {
			errorsChan <- err
			close(errorsChan)
		} else {
			close(errorsChan)
		}
	}()

	return errorsChan, outputChan
}
