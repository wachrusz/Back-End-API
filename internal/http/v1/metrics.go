package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

// RequestCountResponse представляет структуру ответа для GetRequestCountHandler
type RequestCountResponse struct {
	TotalRPS  int               `json:"total_rps"` // Общее количество запросов за последнюю секунду
	Details   map[string]string `json:"details"`   // Детали запросов по IP
	Timestamp string            `json:"timestamp"` // Временная метка ответа
	LocalID   string            `json:"local_id"`  // Идентификатор локального сервера
}

// GetRequestCountHandler returns the number of requests in the last second in JSON format
//
//	@Summary		Returns the number of requests in the last second as a stream
//	@Description	This streaming handler returns the number of requests in the last second and details by IP for rate limiting and usage monitoring.
//	@Tags			Metrics
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	RequestCountResponse	"Successful response with request count"
//	@Failure		500	{object}	jsonresponse.ErrorResponse	"Internal server error"
//	@Failure		429	{object}	jsonresponse.ErrorResponse	"Rate limit exceeded"
//	@Router			/metrics/rps [get]
func (h *MyHandler) GetRequestCountHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовки для поточного ответа
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	// Создаем флушер
	flusher, ok := w.(http.Flusher)
	if !ok {
		h.errResp(w, fmt.Errorf("streaming not supported"), http.StatusInternalServerError)
		return
	}

	// Бесконечный цикл для периодической отправки данных
	for {
		// Создаем контекст с таймаутом для работы с Redis (каждый запрос имеет свой контекст)
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel() // Отменяем контекст после обработки

		// Получаем все данные из Redis Hash
		details, err := h.rdb.HGetAll(ctx, "rate_limit:details").Result()
		if err != nil {
			h.errResp(w, fmt.Errorf("internal server error: %w", err), http.StatusInternalServerError)
			return
		}

		// Получаем общее количество запросов
		totalRPS, err := h.rdb.Get(ctx, "rate_limit:total").Int()
		if err != nil && !errors.Is(err, redis.Nil) {
			h.errResp(w, fmt.Errorf("internal server error: %w", err), http.StatusInternalServerError)
			return
		}
		if errors.Is(err, redis.Nil) {
			totalRPS = 0 // Если ключ отсутствует, присваиваем 0
		}

		// Создаем ответ
		response := RequestCountResponse{
			TotalRPS:  totalRPS,                        // Преобразуем к числу
			Details:   details,                         // Здесь уже данные о запросах по каждому IP
			Timestamp: time.Now().Format(time.RFC3339), // Текущее время в формате RFC3339
			LocalID:   "server_1",                      // Или используйте реальный идентификатор сервера
		}

		// Отправляем данные клиенту
		if err := json.NewEncoder(w).Encode(response); err != nil {
			h.errResp(w, fmt.Errorf("failed to encode response: %w", err), http.StatusInternalServerError)
			return
		}

		// Принудительно сбрасываем данные на клиент
		flusher.Flush()

		// Ждем 1 секунду перед следующим обновлением
		time.Sleep(time.Duration(h.rateLimitCfg) * time.Second)
	}
}
