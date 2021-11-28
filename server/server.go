package server

import (
	"fmt"
	"jinx/config"
	"jinx/conn"
	"jinx/router"
	"net"
)

type Server interface {

	// Start 启动服务器
	Start()

	// Stop 停止服务器
	Stop()

	// Serve 开始服务
	Serve()

	// AddRouter 给当前服务添加路由处理业务
	AddRouter(router router.Router)
}

type server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	Router    router.Router
}

func (s *server) AddRouter(router router.Router) {
	s.Router = router
}

func (s *server) Start() {
	fmt.Printf("[Config] ServerName: %s, IP: %s, Port: %d, IPVersion: %s, MaxConn: %d, MaxPackSize: %d byte\n",
		config.ServerConfig.Name, config.ServerConfig.Host, config.ServerConfig.Port,
		config.ServerConfig.IPVersion, config.ServerConfig.MaxConn, config.ServerConfig.MaxPackSize)
	fmt.Printf("[Jinx Start] Server Listener at IP: %s, Port: %d\n", s.IP, s.Port)
	var connID uint = 0
	go func() {
		// 1.获取tcp的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("[Jinx Server] resolve tcp addr error:", err)
		}
		// 2.监听服务器的地址
		tcpListener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("[Jinx Server] listen err ", err)
		}
		// 3.等待客户端链接
		for {
			tcpConn, err := tcpListener.AcceptTCP()
			if err != nil {
				fmt.Println("[Jinx Server] Accept err:", err)
				continue
			}
			connection := conn.NewConnection(tcpConn, connID, s.Router)
			connection.Start()
			connID++
		}
	}()
}

func (s *server) Stop() {

}

func (s *server) Serve() {
	s.Start()

	// 阻塞
	select {}
}

func NewServer(path string) Server {
	if err := config.InitConfig(path); err != nil {
		panic(err)
	}
	s := &server{
		Name:      config.ServerConfig.Name,
		IPVersion: config.ServerConfig.IPVersion,
		IP:        config.ServerConfig.Host,
		Port:      config.ServerConfig.Port,
		Router:    nil,
	}
	return s
}