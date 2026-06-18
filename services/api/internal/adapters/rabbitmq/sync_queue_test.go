package rabbitmq

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	app "github.com/sutad-p/oh-my-trading/services/api/internal/application/marketdata"
)

func TestSyncPublisherPublishesRequest(t *testing.T) {
	url := os.Getenv("OMT_TEST_RABBITMQ_URL")
	if url == "" {
		url = "amqp://omt:omt_local_password@localhost:5672/"
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		t.Skipf("rabbitmq integration broker is unavailable: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("open channel: %v", err)
	}
	defer ch.Close()

	queueName := "test.market-data-sync"
	if _, err := ch.QueueDelete(queueName, false, false, false); err != nil {
		t.Fatalf("delete queue: %v", err)
	}

	publisher := NewSyncPublisher(ch, queueName)
	request := app.SyncRequest{
		SymbolID:   "018f4f8a-0000-7000-9000-000000000401",
		SymbolCode: "XAUUSD",
		Timeframes: []string{"1h"},
		From:       time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		To:         time.Date(2026, 1, 1, 1, 0, 0, 0, time.UTC),
	}
	if err := publisher.PublishSyncRequest(context.Background(), request); err != nil {
		t.Fatalf("PublishSyncRequest returned error: %v", err)
	}

	msg, ok, err := ch.Get(queueName, true)
	if err != nil {
		t.Fatalf("get message: %v", err)
	}
	if !ok {
		t.Fatal("expected a published message")
	}

	var got app.SyncRequest
	if err := json.Unmarshal(msg.Body, &got); err != nil {
		t.Fatalf("decode message: %v", err)
	}
	if got.SymbolCode != "XAUUSD" {
		t.Fatalf("SymbolCode = %q, want XAUUSD", got.SymbolCode)
	}
}
