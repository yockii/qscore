package activemq

import (
	"context"
	"sync"
	"time"

	"github.com/Azure/go-amqp"
)

var defaultActiveMq = NewActiveMq()

type activeMq struct {
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

func NewActiveMq() *activeMq {
	return &activeMq{
		handlers: make(map[string]func([]byte) error),
		senders:  make(map[string]*amqp.Sender),
	}
}

func (mq *activeMq) RegisterHandlers(queue string, handler func([]byte) error) {
	mq.handlers[queue] = handler
}

func (mq *activeMq) Send(queue string, data []byte, delay int64) error {
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

func (mq *activeMq) reinit() {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	if mq.inited {
		return
	}
	_ = mq.Close()
	err := mq.Init()
	for err != nil {
		_ = mq.Close()
		err = mq.Init()
	}
	for q := range mq.senders {
		var opts []amqp.LinkOption
		opts = append(opts, amqp.LinkTargetAddress(q))
		mq.senders[q], _ = mq.session.NewSender(opts...)
	}
}

func (mq *activeMq) InitWithUsernamePassword(address, username, password string) error {
	mq.address = address
	mq.username = username
	mq.password = password
	return mq.Init()
}

func (mq *activeMq) InitWithAddress(address string) error {
	mq.address = address
	return mq.Init()
}

func (mq *activeMq) Init() error {
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

func (mq *activeMq) StartRead() {
	// 处理注册的handlers
	for queue := range mq.handlers {
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

func (mq *activeMq) Close() error {
	mq.started = false

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, sender := range mq.senders {
		_ = sender.Close(ctx)
	}
	_ = mq.session.Close(ctx)
	return mq.client.Close()
}

func (mq *activeMq) createNewSession() error {
	session, err := mq.client.NewSession()
	if err != nil {
		return err
	}
	mq.session = session
	return nil
}

func (mq *activeMq) read(queue string) {
	ctx := context.Background()

	receiver, err := mq.session.NewReceiver(
		amqp.LinkSourceAddress(queue),
		amqp.LinkCredit(1),
	)
	if err != nil {
		//mq.safeSendError(err)
		return
	}
	defer func() {
		c, cancel := context.WithTimeout(ctx, 1*time.Second)
		if receiver != nil {
			_ = receiver.Close(c)
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
			_ = receiver.ReleaseMessage(ctx, msg)
		} else {
			_ = receiver.AcceptMessage(ctx, msg)
		}
	}
}

func (mq *activeMq) safeSendError(err error, ec chan error) {
	defer func() {
		if recover() != nil {
		}
	}()
	ec <- err
	close(ec)
}

////////////////////////////////////////////////////////////////////////////

func RegisterHandlers(queue string, handler func([]byte) error) {
	defaultActiveMq.RegisterHandlers(queue, handler)
}
func Send(queue string, data []byte, delay int64) error {
	return defaultActiveMq.Send(queue, data, delay)
}
func InitWithUsernamePassword(address, username, password string) error {
	return defaultActiveMq.InitWithUsernamePassword(address, username, password)
}
func InitWithAddress(address string) error {
	return defaultActiveMq.InitWithAddress(address)
}
func Init() error {
	return defaultActiveMq.Init()
}
func StartRead() {
	defaultActiveMq.StartRead()
}
func Close() error {
	return defaultActiveMq.Close()
}
