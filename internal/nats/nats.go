package nats

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/aurelius15/product-reviews/internal/config"
	"github.com/aurelius15/product-reviews/internal/web/rest/apimodel"
)

type Publisher interface {
	Publish(data *apimodel.Review, isDeleted bool) error
	Close()
}

type Nats struct {
	conn    *nats.Conn
	js      nats.JetStreamContext
	subject string
}

var instance *Nats

func NewNats(cnf *config.NatsCnf) (Publisher, error) {
	if instance != nil {
		return instance, nil
	}

	conn, err := nats.Connect(cnf.Host())
	if err != nil {
		return nil, err
	}

	if conn.Status() != nats.CONNECTED {
		return nil, err
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     fmt.Sprintf("%s-stream", cnf.Subject()),
		Subjects: []string{cnf.Subject()},
		Storage:  nats.FileStorage,
	})

	if err != nil {
		return nil, err
	}

	instance = &Nats{
		conn:    conn,
		js:      js,
		subject: cnf.Subject(),
	}

	return instance, nil
}

func (n *Nats) Publish(obj *apimodel.Review, isDeleted bool) error {
	data, _ := json.Marshal(obj)
	msg := &nats.Msg{
		Subject: n.subject,
		Data:    data,
		Header:  make(nats.Header),
	}

	if isDeleted {
		msg.Header.Add("isDeleted", "true")
	}

	_, err := n.js.PublishMsg(msg)

	return err
}

func (n *Nats) Close() {
	if n.conn == nil {
		return
	}

	n.conn.Close()
}
