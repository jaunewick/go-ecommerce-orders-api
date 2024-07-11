package application

import (
    "os"
    "strconv"
)

type Config struct {
    RedisAddr  string
    ServerPort uint16
}

func LoadConfig() Config {
    cfg := Config{
        RedisAddr:  "localhost:6379",
        ServerPort: 3000,
    }

    if redisAddr, exist := os.LookupEnv("REDIS_ADDR"); exist {
        cfg.RedisAddr = redisAddr
    }

    if serverPort, exist := os.LookupEnv("SERVER_PORT"); exist {
        if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
            cfg.ServerPort = uint16(port)
        }
    }

    return cfg
}