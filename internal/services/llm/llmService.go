package llm

import (
	"context"
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
