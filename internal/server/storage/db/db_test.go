package db

import (
	_ "errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Sofja96/GophKeeper.git/internal/server/settings"
)

func TestNewAdapter(t *testing.T) {

	t.Run("database connection error", func(t *testing.T) {
		// Создаем реальный объект settings.Settings с неверным DSN
		settings := &settings.Settings{
			DbDsn: "invalid_dsn",
		}

		// Вызываем NewAdapter
		adapter, err := NewAdapter(settings)

		// Проверяем, что возвращается ошибка
		assert.Error(t, err)
		assert.Nil(t, adapter)
	})

}
