package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCfg struct {
	URL      string `yaml:"url"`
	Password string `yaml:"password"`
}

// New инициализирует Redis-клиент и проверяет его работоспособность
func New(c RedisCfg) (*redis.Client, error) {
	// Создаем новый Redis клиент
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.URL,
		Password: c.Password, // Если пароль отсутствует, то использовать пустую строку
		DB:       0,          // Используем базу по умолчанию
	})

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Пингуем Redis сервер для проверки подключения
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Если пинг успешен, возвращаем клиент
	return rdb, nil
}
