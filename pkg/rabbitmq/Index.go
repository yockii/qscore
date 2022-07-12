package rabbitmq

import (
	"errors"
	"sync"

	"github.com/rabbitmq/amqp091-go"

	"github.com/yockii/qscore/pkg/logger"
)

var defaultRabbitMq = NewRabbitMq()

type rabbitMq struct {
	urlWithUsernameAndPassword string

	conn     *amqp091.Connection
	channel  *amqp091.Channel
	queues   map[string]amqp091.Queue
	handlers map[string]func([]byte) error

	lock   sync.Mutex
	inited bool

	errorChan chan error
}

func NewRabbitMq() *rabbitMq {
	return &rabbitMq{
		queues:   make(map[string]amqp091.Queue),
		handlers: make(map[string]func([]byte) error),
	}
}

func (mq *rabbitMq) SetAddress(addressWithUsernamePassword string) {
	mq.urlWithUsernameAndPassword = addressWithUsernamePassword
}

func (mq *rabbitMq) RegisterHandler(queue string, handler func([]byte) error) {
	mq.queues[queue] = amqp091.Queue{Name: queue}
	mq.handlers[queue] = handler
}

func (mq *rabbitMq) Send(queue string, data []byte /*, delay int64*/) error {
	if _, ok := mq.queues[queue]; !ok {
		mq.queues[queue], _ = mq.channel.QueueDeclare(
			queue,
			false, false, false, false, nil,
		)
	}
	dc, err := mq.channel.PublishWithDeferredConfirm(
		"",
		queue,
		false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	if err != nil {
		if err == amqp091.ErrClosed {
			// 出现链接关闭，尝试重新链接处理
			mq.inited = false
			mq.init()
			err = mq.Send(queue, data)
		}
		return err
	}
	if dc.Wait() {
		return nil
	}
	// TODO return ?
	return nil
}

func (mq *rabbitMq) init() error {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	if mq.inited {
		return nil
	}
	_ = mq.Close()
	var err error
	if mq.conn == nil || mq.conn.IsClosed() {
		mq.conn, err = amqp091.Dial(mq.urlWithUsernameAndPassword)
		if err != nil {
			return err
		}
	}
	if mq.channel == nil || mq.channel.IsClosed() {
		mq.channel, err = mq.conn.Channel()
		if err != nil {
			return err
		}
	}
	err = mq.channel.Confirm(false)
	if err != nil {
		return err
	}

	err = mq.channel.Qos(10, 0, false)
	if err != nil {
		return err
	}

	for name, _ := range mq.queues {
		mq.queues[name], _ = mq.channel.QueueDeclare(name, false, false, false, false, nil)
	}
	mq.inited = true
	mq.errorChan = make(chan error)
	return nil
}

func (mq *rabbitMq) Init() error {
	if mq.inited {
		return nil
	}
	return mq.init()
}

func (mq *rabbitMq) Close() error {
	if mq.channel != nil {
		err := mq.channel.Close()
		if err != nil {
			return err
		}
	}
	if mq.conn != nil {
		err := mq.conn.Close()
		if err != nil {
			return err
		}
	}
	if mq.errorChan != nil {
		close(mq.errorChan)
	}
	return nil
}

// StartRead 异步执行监听
func (mq *rabbitMq) StartRead() {
	for queue, handler := range mq.handlers {
		if handler != nil {
			msgChan, err := mq.channel.Consume(queue, "", false, false, false, false, nil)
			if err != nil {
				logger.Fatal("Failed to get mq channel:", queue)
			}
			go mq.read(msgChan, handler)
		}
	}
	<-mq.errorChan
	mq.inited = false
	mq.init()
}

func (mq *rabbitMq) read(msgChan <-chan amqp091.Delivery, handler func([]byte) error) {
	for deliveryInfo := range msgChan {
		if err := handler(deliveryInfo.Body); err != nil {
			deliveryInfo.Reject(true)
		}
		deliveryInfo.Ack(false)
	}
	mq.errorChan <- errors.New("closed")
}

func SetAddress(addressWithUsernamePassword string) {
	defaultRabbitMq.SetAddress(addressWithUsernamePassword)
}
func RegisterHandler(queue string, handler func([]byte) error) {
	defaultRabbitMq.RegisterHandler(queue, handler)
}
func Init() error  { return defaultRabbitMq.Init() }
func Close() error { return defaultRabbitMq.Close() }
func StartRead()   { defaultRabbitMq.StartRead() }
func Send(queue string, data []byte /*, delay int64*/) error {
	return defaultRabbitMq.Send(queue, data)
}
