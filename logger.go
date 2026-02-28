package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// LogLevel log seviyesi
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// Global log ID counter
var logIdCounter int64

// LogEntry tek bir log kaydı
type LogEntry struct {
	ID        int64                  `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Source    string                 `json:"source"`
	Message   string                 `json:"message"`
	Code      string                 `json:"code,omitempty"`
	Params    map[string]interface{} `json:"params,omitempty"`
}

// LogBuffer log kayıtlarını tutar
type LogBuffer struct {
	entries []LogEntry
	maxSize int
	mu      sync.RWMutex
}

var (
	logBuffer = &LogBuffer{
		entries: make([]LogEntry, 0, 1000),
		maxSize: 1000,
	}
)

// SetLogBufferSize log buffer boyutunu ayarlar
func SetLogBufferSize(size int) {
	logBuffer.mu.Lock()
	defer logBuffer.mu.Unlock()

	logBuffer.maxSize = size

	// Mevcut logları kırp
	if len(logBuffer.entries) > size {
		logBuffer.entries = logBuffer.entries[len(logBuffer.entries)-size:]
	}
}

// AddLog yeni log ekler
func AddLog(level LogLevel, source, message string) {
	logBuffer.mu.Lock()
	defer logBuffer.mu.Unlock()

	entry := LogEntry{
		ID:        atomic.AddInt64(&logIdCounter, 1),
		Timestamp: time.Now(),
		Level:     level,
		Source:    source,
		Message:   message,
	}

	logBuffer.entries = append(logBuffer.entries, entry)

	// Buffer boyutunu kontrol et
	if len(logBuffer.entries) > logBuffer.maxSize {
		logBuffer.entries = logBuffer.entries[1:]
	}

	// Ayrıca standart çıktıya yaz (debug için)
	fmt.Printf("[%s] [%s] %s: %s\n", entry.Timestamp.Format("15:04:05.000"), level, source, message)
}

// GetLogs tüm logları döndürür
func GetLogs() []LogEntry {
	logBuffer.mu.RLock()
	defer logBuffer.mu.RUnlock()

	result := make([]LogEntry, len(logBuffer.entries))
	copy(result, logBuffer.entries)
	return result
}

// GetLogsAfterId belirli bir ID'den sonraki logları döndürür
func GetLogsAfterId(afterId int64) []LogEntry {
	logBuffer.mu.RLock()
	defer logBuffer.mu.RUnlock()

	var result []LogEntry
	for _, entry := range logBuffer.entries {
		if entry.ID > afterId {
			result = append(result, entry)
		}
	}
	return result
}

// getLogsJSON logları JSON string olarak döndürür (WebView2 için)
func getLogsJSON() string {
	logs := GetLogs()
	data, _ := json.Marshal(logs)
	return string(data)
}

// getLogsAfterIdJSON belirli bir ID'den sonraki logları döndürür
func getLogsAfterIdJSON(afterId int64) string {
	logs := GetLogsAfterId(afterId)
	data, _ := json.Marshal(logs)
	return string(data)
}

// clearLogs logları temizler
func clearLogs() string {
	logBuffer.mu.Lock()
	defer logBuffer.mu.Unlock()

	logBuffer.entries = logBuffer.entries[:0]
	return `{"success": true}`
}

// setLogBufferSizeJSON buffer boyutunu ayarlar (WebView2 için)
func setLogBufferSizeJSON(size int) string {
	SetLogBufferSize(size)
	return fmt.Sprintf(`{"success": true, "size": %d}`, size)
}

// getLogBufferSize mevcut buffer boyutunu döndürür
func getLogBufferSize() string {
	logBuffer.mu.RLock()
	defer logBuffer.mu.RUnlock()
	return fmt.Sprintf(`{"size": %d, "count": %d}`, logBuffer.maxSize, len(logBuffer.entries))
}

// Convenience functions
func AddLogWithCode(level LogLevel, source, message, code string, params map[string]interface{}) {
	logBuffer.mu.Lock()
	defer logBuffer.mu.Unlock()

	entry := LogEntry{
		ID:        atomic.AddInt64(&logIdCounter, 1),
		Timestamp: time.Now(),
		Level:     level,
		Source:    source,
		Message:   message,
		Code:      code,
		Params:    params,
	}

	logBuffer.entries = append(logBuffer.entries, entry)

	// Buffer boyutunu kontrol et
	if len(logBuffer.entries) > logBuffer.maxSize {
		logBuffer.entries = logBuffer.entries[1:]
	}

	// Ayrıca standart çıktıya yaz (debug için)
	if params != nil {
		fmt.Printf("[%s] [%s] %s: %s (code=%s, params=%v)\n", entry.Timestamp.Format("15:04:05.000"), level, source, message, code, params)
	} else {
		fmt.Printf("[%s] [%s] %s: %s (code=%s)\n", entry.Timestamp.Format("15:04:05.000"), level, source, message, code)
	}
}

func LogDebug(source, message string) {
	AddLog(LogLevelDebug, source, message)
}

func LogInfo(source, message string) {
	AddLog(LogLevelInfo, source, message)
}

func LogWarn(source, message string) {
	AddLog(LogLevelWarn, source, message)
}

func LogError(source, message string) {
	AddLog(LogLevelError, source, message)
}

func LogErrorWithCode(source, message, code string, params map[string]interface{}) {
	AddLogWithCode(LogLevelError, source, message, code, params)
}

// LogSystem sistem işlemleri için
func LogSystem(message string) {
	AddLog(LogLevelInfo, "System", message)
}

func LogSystemWithCode(message, code string, params map[string]interface{}) {
	AddLogWithCode(LogLevelInfo, "System", message, code, params)
}

// LogHTTP HTTP istekleri için
func LogHTTP(message string) {
	AddLog(LogLevelDebug, "HTTP", message)
}
