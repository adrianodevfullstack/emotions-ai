package emotions

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/joho/godotenv"
)

const (
	defaultPort  = "8080"
	defaultModel = "googleai/gemini-flash-latest"
	webDir       = "./web"
)

func Run() error {
	if err := loadEnv(); err != nil {
		return err
	}

	ctx := context.Background()
	g := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel(defaultModel),
	)

	flow := defineInterpretationFlow(g)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/ws", newWSHandler(flow))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("servindo em http://localhost:%s (webcam) — janelas de %ds enviadas ao Genkit\n", port, windowSeconds)
	return http.ListenAndServe(":"+port, mux)
}

func loadEnv() error {
	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("erro ao ler .env: %w", err)
	}

	if os.Getenv("GEMINI_API_KEY") == "" {
		return errors.New("defina a variável de ambiente GEMINI_API_KEY antes de rodar (ver README.md)")
	}
	return nil
}
