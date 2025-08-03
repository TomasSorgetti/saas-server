package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// en produccion deberia usar un servicio externo
func NewLogger() *zerolog.Logger {
    file, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        panic(err) // Change in production
    }
    
    // Configurar zerolog para escribir en archivo y consola
    multi := zerolog.MultiLevelWriter(
        zerolog.ConsoleWriter{Out: os.Stdout}, // Para desarrollo
        file,                                  // Para persistencia
    )
    
    logger := zerolog.New(multi).
        With().
        Timestamp().
        Logger()
    return &logger
}