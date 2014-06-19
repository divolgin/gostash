package gostash

import (
	"net"
)

type Config struct {
	LogstashHost string
	LogstashPort string
	InputType    string
	ThisHostName string
	CodeVersion  string
}

type LogstashClient struct {
	logstashHost string
	logstashPort string
	serverAddr   *net.UDPAddr
	conn         *net.UDPConn
	err          error

	formatter *Formatter
}

func NewLogstashClient(config *Config) *LogstashClient {
	c := LogstashClient{}
	c.formatter = NewFormatter(config)
	c.logstashHost = config.LogstashHost
	c.logstashPort = config.LogstashPort
	c.serverAddr, c.err = net.ResolveUDPAddr("udp", c.logstashHost+":"+c.logstashPort)
	if c.err != nil {
		// This is an invalid connection with error info
		return &c
	}
	c.conn, c.err = net.DialUDP("udp", nil, c.serverAddr)
	return &c
}

func (c *LogstashClient) Formatter() *Formatter {
	return c.formatter
}

func (c *LogstashClient) LastError() error {
	return c.err
}

func (c *LogstashClient) SendMessage(msg string, metadata map[string]string) {
	if c.err != nil {
		return
	}

	var encoded []byte
	encoded, c.err = c.formatter.encode(msg, metadata)
	if c.err != nil {
		return
	}
	c.Write(encoded)
}

func (c *LogstashClient) Write(msg []byte) (n int, err error) {
	if msg == nil {
		return 0, nil
	}

	if c.conn != nil {
		_, c.err = c.conn.Write(msg)
	}

	// Allow logging to fail silently when used with loggo
	return len(msg), nil
}

func (c *LogstashClient) Close() {
	if c.conn != nil {
		c.err = c.conn.Close()
	}
}
