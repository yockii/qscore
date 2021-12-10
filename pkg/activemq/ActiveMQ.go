package activemq

import (
	"context"
	"sync"
	"time"

	"github.com/Azure/go-amqp"
)

var ActiveMQ = &activeMqNew{
	handlers: make(map[string]func([]byte) error),
	senders:  make(map[string]*amqp.Sender),
}

type activeMqNew struct {
	address   string
	username  string
	password  string
	anonymous bool

	inited    bool
	started   bool
	client    *amqp.Client
	session   *amqp.Session
	lock      sync.Mutex
	handlers  map[string]func([]byte) error
	errorChan chan error
	senders   map[string]*amqp.Sender
}

func (mq *activeMqNew) RegisterHandlers(queue string, handler func([]byte) error) {
	mq.handlers[queue] = handler
}

func (mq *activeMqNew) Send(queue string, data []byte, delay int64) error {
	sender, ok := mq.senders[queue]
	var err error
	if !ok {
		var opts []amqp.LinkOption
		opts = append(opts, amqp.LinkTargetAddress(queue))
		mq.senders[queue], err = mq.session.NewSender(opts...)
		if err != nil {
			return err
		}
		sender = mq.senders[queue]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := amqp.NewMessage(data)
	msg.Header = &amqp.MessageHeader{Durable: true}
	if delay > 0 {
		msg.Annotations = amqp.Annotations{
			"x-opt-delivery-delay": delay,
		}
	}
	err = sender.Send(ctx, msg)
	if err != nil {
		if err == amqp.ErrConnClosed {
			mq.inited = false
			mq.reinit()
			err = mq.Send(queue, data, delay)
		}
		return err
	}
	return nil
}

func (mq *activeMqNew) reinit() {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	if mq.inited {
		return
	}
	mq.Close()
	err := mq.Init()
	for err != nil {
		mq.Close()
		err = mq.Init()
	}
	for q, _ := range mq.senders {
		var opts []amqp.LinkOption
		opts = append(opts, amqp.LinkTargetAddress(q))
		mq.senders[q], _ = mq.session.NewSender(opts...)
	}
}

func (mq *activeMqNew) InitWithUsernamePassword(username, password string) error {
	mq.username = username
	mq.password = password
	return mq.Init()
}

func (mq *activeMqNew) Init() error {
	if mq.inited {
		return nil
	}
	var client *amqp.Client
	var err error
	if mq.anonymous {
		client, err = amqp.Dial(
			mq.address,
			amqp.ConnSASLAnonymous(),
		)
	} else if mq.username != "" && mq.password != "" {
		client, err = amqp.Dial(
			mq.address,
			amqp.ConnSASLPlain(mq.username, mq.password),
		)
	} else {
		mq.anonymous = true
		client, err = amqp.Dial(
			mq.address,
			amqp.ConnSASLAnonymous(),
		)
	}
	if err != nil {
		return err
	}
	mq.client = client
	err = mq.createNewSession()
	if err != nil {
		return err
	}
	mq.errorChan = make(chan error)
	mq.inited = true
	mq.started = false
	return nil
}

func (mq *activeMqNew) StartRead() {
	// 处理注册的handlers
	for queue, _ := range mq.handlers {
		go mq.read(queue)
	}
	mq.started = true
	<-mq.errorChan
	mq.inited = false

	mq.reinit()

	if !mq.started {
		mq.StartRead()
	}
}

func (mq *activeMqNew) Close() error {
	mq.started = false

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, sender := range mq.senders {
		sender.Close(ctx)
	}
	mq.session.Close(ctx)
	return mq.client.Close()
}

func (mq *activeMqNew) createNewSession() error {
	session, err := mq.client.NewSession()
	if err != nil {
		return err
	}
	mq.session = session
	return nil
}

func (mq *activeMqNew) read(queue string) {
	ctx := context.Background()

	receiver, err := mq.session.NewReceiver(
		amqp.LinkSourceAddress(queue),
		amqp.LinkCredit(10),
	)
	if err != nil {
		//mq.safeSendError(err)
		return
	}
	defer func() {
		c, cancel := context.WithTimeout(ctx, time.Second)
		if receiver != nil {
			receiver.Close(c)
		}
		cancel()
	}()

	ec := mq.errorChan
	for {
		msg, err2 := receiver.Receive(ctx)
		if err2 != nil {
			mq.safeSendError(err2, ec)
			return
		}
		e := mq.handlers[queue](msg.GetData())
		if e != nil {
			receiver.ReleaseMessage(ctx, msg)
		} else {
			receiver.AcceptMessage(ctx, msg)
		}
	}
}

func (mq *activeMqNew) safeSendError(err error, ec chan error) {
	defer func() {
		if recover() != nil {
		}
	}()
	ec <- err
	close(ec)
}
