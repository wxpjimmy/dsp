package service

import (
	"net"
	"time"
	"net/http"
	"sync"
	"dsp_demo/runner"
	"github.com/astaxie/beego/logs"
)

type Service struct {
	http.Server
	// Other things
	Quit           chan bool
	WaitGroup      *sync.WaitGroup
	mu             *sync.Mutex
	conns          map[string]net.Conn
	ln             net.Listener
	RouterRegister func()
	Dispatcher     *runner.Dispatcher
}

func NewService() *Service {
	//disp := runner.NewDispatcher(10)
	s := &Service{
		// Init Other things
		Quit:      make(chan bool),
		WaitGroup: &sync.WaitGroup{},
		mu:        &sync.Mutex{},
		conns:     make(map[string]net.Conn),
		//Dispatcher: &disp,
	}

	return s
}

func (s *Service) Stop() error {
	//ln.Close()
	logs.Info("Start stopping...")
	close(s.Quit)
	s.SetKeepAlivesEnabled(false)
	s.mu.Lock()
	// close listenser
	if err := s.ln.Close(); err != nil {
		return err
	}
	//将当前idle的connections设置read timeout，便于后续关闭。
	t := time.Now().Add(150 * time.Millisecond)
	for _, c := range s.conns {
		c.SetReadDeadline(t)
	}
	s.conns = make(map[string]net.Conn)
	s.mu.Unlock()
	logs.Debug("wait for group end!")
	s.WaitGroup.Wait()
	logs.Info("Stop!")
	return nil
}

func (s *Service) Start() {
	s.WaitGroup.Add(1)
	defer s.WaitGroup.Done()


	s.RouterRegister()

	var err error
	s.ln, err = net.Listen("tcp", ":8080")
	if err != nil {
		logs.Critical("Start service failed!")
		return
	}
	s.Server = http.Server{Handler: nil}
	s.ConnState = func(conn net.Conn, state http.ConnState) {
		switch state {
		case http.StateNew:
			// 新的连接，计数加1
			logs.Debug("New connection create!")
			s.WaitGroup.Add(1)
		case http.StateActive:
			// 有新的请求，从idle conn pool中移除
			logs.Debug("New Request Accepted!")
			s.mu.Lock()
			delete(s.conns, conn.LocalAddr().String())
			s.mu.Unlock()
		case http.StateIdle:
			logs.Debug("Connection Idle!")
			select {
			case <-s.Quit:
				// 如果要关闭了，直接Close，否则加入idle conn pool中。
				conn.Close()
			default:
				s.mu.Lock()
				s.conns[conn.LocalAddr().String()] = conn
				s.mu.Unlock()
			}
		case http.StateHijacked, http.StateClosed:
			// conn已经closed了，计数减一
			logs.Debug("Connection closed!")
			s.WaitGroup.Done()
		}
	}
	s.Serve(s.ln);
	logs.Info("Server stopped!")
}