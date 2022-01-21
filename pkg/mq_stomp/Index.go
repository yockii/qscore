package mq_stomp

import (
	"crypto/tls"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"

	"github.com/yockii/qscore/pkg/logger"
)

var defaultStomp = NewStomp()

type mqStomp struct {
	username           string
	password           string
	addressList        []string
	insecureSkipVerify bool
	addressIndex       int

	inited  bool
	started bool
	lock    sync.Mutex

	sendConn *stomp.Conn

	errorChan chan error

	handlers    map[string]func([]byte) error
	receiveConn *stomp.Conn
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
			//stomp.ConnOpt.Host("/"),
		}
		if mq.username != "" && mq.password != "" {
			options = append(options, stomp.ConnOpt.Login(mq.username, mq.password))
		}

		for addressIndex := 0; addressIndex < len(mq.addressList); addressIndex++ {
			address := mq.addressList[addressIndex]
			if isSTOMPTLS(address) {
				var netConn *tls.Conn
				netConn, err = tls.Dial("tcp", address, &tls.Config{})
				if err != nil {
					logger.Error(err)
					continue
				}
				mq.sendConn, err = stomp.Connect(netConn, options...)
			} else {
				mq.sendConn, err = stomp.Dial("tcp", address, options...)
			}
			if err != nil {
				logger.Error(err)
				continue
			}
			mq.addressIndex = addressIndex
			break
		}
		if err != nil {
			return err
		}
	}

	mq.inited = true
	return nil
}
func (mq *mqStomp) InitWithUsernamePassword(address, username, password string) error {
	mq.SetAddresses(address)
	mq.username = username
	mq.password = password
	return mq.Init()
}

func (mq *mqStomp) InitWithAddress(address string) error {
	mq.SetAddresses(address)
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
		if strings.Contains(err.Error(), "closed") || err == stomp.ErrAlreadyClosed {
			mq.inited = false
			mq.reinit()
			err = mq.Send(queue, data, delay)
		}
		return err
	}
	return nil
}

func (mq *mqStomp) StartRead() {
	// 建立receiveConn
	options := []func(*stomp.Conn) error{
		//stomp.ConnOpt.Host("/"),
	}
	if mq.username != "" && mq.password != "" {
		options = append(options, stomp.ConnOpt.Login(mq.username, mq.password))
	}
	var err error
	address := mq.addressList[mq.addressIndex]
	if isSTOMPTLS(address) {
		var netConn *tls.Conn
		netConn, err = tls.Dial("tcp", address, &tls.Config{})
		if err == nil {
			mq.receiveConn, err = stomp.Connect(netConn, options...)
		}
	} else {
		mq.receiveConn, err = stomp.Dial("tcp", address, options...)
	}
	if err != nil {
		logger.Error(err)
		for addressIndex := 0; addressIndex < len(mq.addressList); addressIndex++ {
			address = mq.addressList[addressIndex]
			if isSTOMPTLS(address) {
				var netConn *tls.Conn
				netConn, err = tls.Dial("tcp", address, &tls.Config{})
				if err != nil {
					logger.Error(err)
					continue
				}
				mq.receiveConn, err = stomp.Connect(netConn, options...)
			} else {
				mq.receiveConn, err = stomp.Dial("tcp", address, options...)
			}
			//mq.receiveConn, err = stomp.Dial("tcp", address, options...)
			if err != nil {
				logger.Error(err)
				continue
			}
			mq.addressIndex = addressIndex
			break
		}
		if err != nil {
			logger.Error(err)
			return
		}
	}

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
	mq.receiveConn.Disconnect()
	return nil
}

func (mq *mqStomp) read(queue string) {
	//options := []func(*stomp.Conn) error{
	//	stomp.ConnOpt.Login(mq.username, mq.password),
	//	stomp.ConnOpt.Host("/"),
	//}
	//conn, err := stomp.Dial("tcp", mq.addressList[mq.addressIndex], options...)
	//if err != nil {
	//	logger.Fatal("cannot connect to ", mq.addressList[mq.addressIndex])
	//}
	//defer conn.Disconnect()

	sub, err := mq.receiveConn.Subscribe(queue, stomp.AckClient)
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
			_ = mq.receiveConn.Nack(msg)
		} else {
			_ = mq.receiveConn.Ack(msg)
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

func (mq *mqStomp) SetUsername(username string) {
	mq.username = username
}

func (mq *mqStomp) SetPassword(password string) {
	mq.password = password
}

func (mq *mqStomp) SetAddresses(addresses ...string) {
	mq.addressList = addresses
}
func (mq *mqStomp) SetInsecureSkipVerify(insecureSkipVerify bool) {
	mq.insecureSkipVerify = insecureSkipVerify
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
func SetAddresses(addresses ...string) {
	defaultStomp.SetAddresses(addresses...)
}
func SetUsername(username string) {
	defaultStomp.SetUsername(username)
}
func SetPassword(password string) {
	defaultStomp.SetPassword(password)
}
func SetInsecureSkipVerify(insecureSkipVerify bool) {
	defaultStomp.SetInsecureSkipVerify(insecureSkipVerify)
}
func StartRead() {
	defaultStomp.StartRead()
}
func Close() error {
	return defaultStomp.Close()
}

/////////////////////

func isSTOMPTLS(address string) bool {
	_, port, err := net.SplitHostPort(address)
	if err != nil {
		return false
	}

	switch port {
	case "61613":
		return false
	case "61614":
		return true
	default:
		return false
	}
}
