package events

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/nats-io/nats.go"
	"github.com/solrac97gr/cqrs/models"
)

const (
	BufferSize = 64
)

type NatsEventStore struct {
	conn            *nats.Conn
	feedCreatedSub  *nats.Subscription
	feedCreatedChan chan CreatedFeedMessage
}

var _ EventStore = &NatsEventStore{}

func NewNatsEventStore(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{conn: conn}, nil
}

func (es *NatsEventStore) Close() {
	if es.conn != nil {
		es.conn.Close()
	}
	if es.feedCreatedSub != nil {
		es.feedCreatedSub.Unsubscribe()
	}

	close(es.feedCreatedChan)
}

func (es *NatsEventStore) encodeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), err
}

func (es *NatsEventStore) PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	msg := CreatedFeedMessage{
		ID:          feed.ID,
		Title:       feed.Title,
		Description: feed.Description,
		CreatedAt:   feed.CreatedAt,
	}
	data, err := es.encodeMessage(msg)
	if err != nil {
		return err
	}
	return es.conn.Publish(msg.Type(), data)
}

func (es *NatsEventStore) decodeMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

func (es *NatsEventStore) OnCreateFeed(f func(CreatedFeedMessage)) (err error) {
	msg := CreatedFeedMessage{}
	es.feedCreatedSub, err = es.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		err := es.decodeMessage(m.Data, &msg)
		if err != nil {
			return
		}
		f(msg)
	})
	return
}

func (es *NatsEventStore) SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	m := CreatedFeedMessage{}
	es.feedCreatedChan = make(chan CreatedFeedMessage, BufferSize)
	ch := make(chan *nats.Msg, BufferSize)
	var err error
	es.feedCreatedSub, err = es.conn.ChanSubscribe(m.Type(), ch)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case msg := <-ch:
				err := es.decodeMessage(msg.Data, &m)
				if err != nil {
					continue
				}
				es.feedCreatedChan <- m
			}
		}
	}()
	return (<-chan CreatedFeedMessage)(es.feedCreatedChan), nil
}
