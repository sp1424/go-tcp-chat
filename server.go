package main

import (
	"bufio"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net"
	"sync"
)

func main() {
	var loggerConfig = zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)

	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	listener, err := net.Listen("tcp", "localhost:4000")
	if err != nil {
		return
	}

	defer listener.Close()
	var connMap = &sync.Map{}

	for  {
		conn, err := listener.Accept()

		if err != nil {
			logger.Error("error accepting connection", zap.Error(err))
			return
		}
		id := uuid.New().String()
		connMap.Store(id, conn)
		go handleConnection(conn, id, connMap, logger)
	}
}

func handleConnection(conn net.Conn, id string, connMap *sync.Map, logger *zap.Logger)  {
	defer conn.Close() //any way function is finished connection is closed


	for  {
		input, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			logger.Error("error reading from client", zap.Error(err))
			return
		}
		connMap.Range(func(key, value interface{}) bool {
			if conn, ok := value.(net.Conn); ok {
				if _, err := conn.Write([]byte(input)); err != nil {
					logger.Error("error on writing to connection", zap.Error(err))
				}
			}
			return true
		})
	}
}

