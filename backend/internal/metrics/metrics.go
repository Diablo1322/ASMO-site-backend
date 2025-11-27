package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Метрики объявляем как глобальные переменные
	DatabaseQueryDuration *prometheus.HistogramVec

	// Защита от двойной регистрации
	metricsOnce sync.Once
)

// InitMetrics инициализирует метрики только один раз
func InitMetrics() {
	metricsOnce.Do(func() {
		DatabaseQueryDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "database_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "table"},
		)

		// Регистрируем метрики
		prometheus.MustRegister(DatabaseQueryDuration)
	})
}

// RecordDatabaseQuery записывает метрику для запроса к БД
func RecordDatabaseQuery(operation, table string, duration time.Duration) {
	// Гарантируем инициализацию метрик
	InitMetrics()

	if DatabaseQueryDuration != nil {
		DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
	}
}