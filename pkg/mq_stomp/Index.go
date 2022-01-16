package mq_stomp

import (
	"strconv"
	"sync"

	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"

	"github.com/yockii/qscore/pkg/logger"
)

var defaultStomp = NewStomp()

type mqStomp struct {
	address  string
	username string
	password string

	inited  bool
	started bool
	lock    sync.Mutex

	sendConn *stomp.Conn

	errorChan chan error

	handlers map[string]func([]byte) error
}

func NewStomp() *mqStomp {
	return &mqStomp{
		handlers: make(map[string]func([]byte) error),
	}
}

func (mq *mqStomp) Init() error {
	if mq.inited {
		return nil
	}
	var err error
	{
		options := []func(*stomp.Conn) error{
			stomp.ConnOpt.Host("/"),
		}
		if mq.username != "" && mq.password != "" {
			options = append(options, stomp.ConnOpt.Login(mq.username, mq.password))
		}
		mq.sendConn, err = stomp.Dial("tcp", mq.address, options...)
		if err != nil {
			return err
		}
	}

	mq.inited = true
	return nil
}
func (mq *mqStomp) InitWithUsernamePassword(address, username, password string) error {
	mq.address = address
	mq.username = username
	mq.password = password
	return mq.Init()
}

func (mq *mqStomp) InitWithAddress(address string) error {
	mq.address = address
	return mq.Init()
}

func (mq *mqStomp) RegisterHandlers(queue string, handler func([]byte) error) {
	mq.handlers[queue] = handler
}

func (mq *mqStomp) Send(queue string, data []byte, delay int64) error {
	var delayOpt func(*frame.Frame) error
	if delay > 0 {
		delayOpt = stomp.SendOpt.Header("AMQ_SCHEDULED_DELAY", strconv.FormatInt(delay, 10))
	}
	err := mq.sendConn.Send(
		queue,
		"text/plain",
		data,
		stomp.SendOpt.Receipt,
		delayOpt,
	)
	if err != nil {
		if err == stomp.ErrAlreadyClosed {
			mq.inited = false
			mq.reinit()
			err = mq.Send(queue, data, delay)
		}
		return err
	}
	return nil
}

func (mq *mqStomp) StartRead() {
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

func (mq *mqStomp) reinit() {
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
}

func (mq *mqStomp) Close() error {
	mq.sendConn.Disconnect()
	return nil
}

func (mq *mqStomp) read(queue string) {
	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(mq.username, mq.password),
		stomp.ConnOpt.Host("/"),
	}
	conn, err := stomp.Dial("tcp", mq.address, options...)
	if err != nil {
		logger.Fatal("cannot connect to ", mq.address)
	}
	defer conn.Disconnect()

	sub, err := conn.Subscribe(queue, stomp.AckClient)
	if err != nil {
		logger.Fatal("cannot subscribe to ", queue, err.Error())
		return
	}
	ec := mq.errorChan
	for {
		msg := <-sub.C
		if msg.Err != nil {
			mq.safeSendError(msg.Err, ec)
			return
		}

		e := mq.handlers[queue](msg.Body)
		if e != nil {
			_ = conn.Nack(msg)
		} else {
			_ = conn.Ack(msg)
		}
	}
}

func (mq *mqStomp) safeSendError(err error, ec chan error) {
	defer func() {
		if recover() != nil {
		}
	}()
	ec <- err
	close(ec)
}

////////////////////////////////////////////////////////////////////////////

func RegisterHandlers(queue string, handler func([]byte) error) {
	defaultStomp.RegisterHandlers(queue, handler)
}

func Send(queue string, data []byte, delay int64) error {
	return defaultStomp.Send(queue, data, delay)
}
func InitWithUsernamePassword(address, username, password string) error {
	return defaultStomp.InitWithUsernamePassword(address, username, password)
}
func InitWithAddress(address string) error {
	return defaultStomp.InitWithAddress(address)
}
func Init() error {
	return defaultStomp.Init()
}
func StartRead() {
	defaultStomp.StartRead()
}
func Close() error {
	return defaultStomp.Close()
}
