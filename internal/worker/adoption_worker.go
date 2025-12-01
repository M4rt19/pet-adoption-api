package worker

import (
	"context"
	"log"
	"time"
)

// event type pushed by handlers
type AdoptionEvent struct {
	RequestID uint
	UserID    uint
	PetID     uint
	Status    string
	Message   string
}

type AdoptionWorker struct {
	Events chan AdoptionEvent
}

func NewAdoptionWorker(buffer int) *AdoptionWorker {
	return &AdoptionWorker{
		Events: make(chan AdoptionEvent, buffer),
	}
}

// Run the worker in background
func (w *AdoptionWorker) Start(ctx context.Context) {
	log.Println("[WORKER] Adoption worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("[WORKER] Adoption worker shutting down...")
			return

		case evt := <-w.Events:
			// simulate sending email (or any async task)
			log.Printf(
				"[WORKER] Processing adoption event → requestID=%d userID=%d petID=%d status=%s message=%s\n",
				evt.RequestID, evt.UserID, evt.PetID, evt.Status, evt.Message,
			)

			// simulate slow notification
			time.Sleep(1 * time.Second)
		}
	}
}
