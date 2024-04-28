package connection

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"
	"velocityApi/config"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func ConnectClickhouse() (driver.Conn, error) {
	clickhouseHost := config.ViperConfig.GetString("clickHouse.host")
	clickhousePort := config.ViperConfig.GetString("clickHouse.port")
	clickhouseUser := config.ViperConfig.GetString("clickHouse.user")
	clickhousePass := config.ViperConfig.GetString("clickHouse.pass")
	clickhouseDatabase := config.ViperConfig.GetString("clickHouse.Database")

	dialTimeout := config.ViperConfig.GetInt32("clickHouse.dial_timeout")
	maxConnectionDuration := config.ViperConfig.GetInt32("clickhouse.max_connection_lifetime_duration")
	maxExecutionTime := config.ViperConfig.GetInt32("clickhouse.max_execution_time")
	maxOpenConnections := config.ViperConfig.GetInt("clickhouse.max_open_connections")
	maxIdleConnections := config.ViperConfig.GetInt("clickhouse.max_idle_connections")
	blockBufferSize := config.ViperConfig.GetInt32("clickhouse.block_buffer_size")
	MaxCompressionBufferSize := config.ViperConfig.GetInt("clickhouse.block_buffer_size")
	InsecureSkipVerify := config.ViperConfig.GetBool("clickhouse.insecure_skip_verify")

	var dialCount int32
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{clickhouseHost + ":" + clickhousePort},
			Auth: clickhouse.Auth{
				Database: clickhouseDatabase,
				Username: clickhouseUser,
				Password: clickhousePass,
			},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "inbox-intelligence-clickhouse-client", Version: "0.1"},
				},
			},
			DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
				dialCount++
				var d net.Dialer
				return d.DialContext(ctx, "tcp", addr)
			},
			Debug: true,
			Settings: clickhouse.Settings{
				"max_execution_time": maxExecutionTime,
			},
			Compression: &clickhouse.Compression{
				Method: clickhouse.CompressionLZ4,
			},
			DialTimeout:          time.Duration(dialTimeout) * time.Second,
			MaxOpenConns:         maxOpenConnections,
			MaxIdleConns:         maxIdleConnections,
			ConnMaxLifetime:      time.Duration(maxConnectionDuration) * time.Minute,
			ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
			BlockBufferSize:      uint8(blockBufferSize),
			MaxCompressionBuffer: MaxCompressionBufferSize,
			TLS: &tls.Config{
				InsecureSkipVerify: InsecureSkipVerify,
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}

func DisconnectClickhouse(conn driver.Conn) error {
	err := conn.Close()
	if err != nil {
		return err
	}
	return nil
}
