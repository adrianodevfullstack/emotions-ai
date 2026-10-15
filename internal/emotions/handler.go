package emotions

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const windowSeconds = 8

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func newWSHandler(flow *interpretationFlow) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("erro no upgrade do websocket:", err)
			return
		}
		defer conn.Close()

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		var writeMu sync.Mutex
		agg := &aggregator{}

		go publishWindows(ctx, cancel, conn, &writeMu, agg, flow)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var sample EmotionSample
			if err := json.Unmarshal(msg, &sample); err != nil {
				log.Println("mensagem inválida:", err)
				continue
			}
			agg.add(sample)
		}
	}
}

func publishWindows(
	ctx context.Context,
	cancel context.CancelFunc,
	conn *websocket.Conn,
	writeMu *sync.Mutex,
	agg *aggregator,
	flow *interpretationFlow,
) {
	ticker := time.NewTicker(windowSeconds * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			window, ok := agg.closeWindow()
			if !ok {
				continue
			}

			interpretation, err := flow.Run(ctx, window)
			if err != nil {
				log.Println("erro ao interpretar janela:", err)
				window.Interpretation = "(não foi possível gerar interpretação agora)"
			} else {
				window.Interpretation = interpretation.Summary + " " + interpretation.Caveat
			}

			payload, _ := json.Marshal(window)
			writeMu.Lock()
			err = conn.WriteMessage(websocket.TextMessage, payload)
			writeMu.Unlock()
			if err != nil {
				log.Println("erro ao enviar resumo:", err)
				cancel()
				return
			}
		}
	}
}
