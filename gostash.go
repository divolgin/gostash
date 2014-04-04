package gostash

import (
	"encoding/json"
	"net"
	"time"
)

type logstashMessage struct {
	Timestamp string            `json:"@timestamp"`
	Message   string            `json:"message"`
	Type      string            `json:"type"`
	Fields    map[string]string `json:"fields"`
}

type LogstashClient struct {
	logstashHost string
	logstashPort string
	serverAddr   *net.UDPAddr
	conn         *net.UDPConn
	err          error
	msg          logstashMessage
}

func NewLogstashClient(host string, port string, inputType string) *LogstashClient {
	c := LogstashClient{}
	c.logstashHost = host
	c.logstashPort = port
	c.serverAddr, c.err = net.ResolveUDPAddr("udp", host+":"+port)
	if c.err != nil {
		// This is an invalid connection with error info
		return &c
	}
	c.conn, c.err = net.DialUDP("udp", nil, c.serverAddr)

	c.msg.Type = inputType

	return &c
}

func (c *LogstashClient) Error() error {
  return c.err
}

func (c *LogstashClient) SendMessage(msg string, metadata map[string]string) {
	if c.err != nil {
		return
	}
	encoded, err := c.encode(msg, metadata)
	if err != nil {
		return
	}
	go c.conn.Write(encoded)
}

func (c *LogstashClient) Close() {
	if c.err == nil {
		return
	}
	c.conn.Close()
}

func (c *LogstashClient) encode(msg string, metadata map[string]string) ([]byte, error) {
	c.msg.Timestamp = getTimeStamp()
	c.msg.Message = msg
	c.msg.Fields = metadata
	return json.Marshal(c.msg)
}

func getTimeStamp() string {
	now := time.Now().UTC()
	return now.Format("2006-01-02T15:04:05.000-07:00")
}
