package logger

import (
	"context"
	"fmt"
	"sync"
)

// MockLogger - mock для тестирования с возможностью проверки вызовов.
type MockLogger struct {
	mu    sync.RWMutex
	calls []LogCall
}

type LogCall struct {
	Level   string
	Message string
	Fields  []Field
	Context context.Context
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		calls: make([]LogCall, 0),
	}
}

func (m *MockLogger) Info(msg string, fields ...Field) {
	m.addCall("INFO", msg, fields, nil)
}

func (m *MockLogger) Error(msg string, fields ...Field) {
	m.addCall("ERROR", msg, fields, nil)
}

func (m *MockLogger) Debug(msg string, fields ...Field) {
	m.addCall("DEBUG", msg, fields, nil)
}

func (m *MockLogger) Warn(msg string, fields ...Field) {
	m.addCall("WARN", msg, fields, nil)
}

func (m *MockLogger) Fatal(msg string, fields ...Field) {
	m.addCall("FATAL", msg, fields, nil)
}

func (m *MockLogger) Panic(msg string, fields ...Field) {
	m.addCall("PANIC", msg, fields, nil)
}

func (m *MockLogger) With(fields ...Field) ILogger {
	// Для простоты возвращаем тот же логгер
	// В реальном mock можно создавать новый экземпляр с дополнительными полями
	return m
}

func (m *MockLogger) WithContext(ctx context.Context) ILogger {
	return m
}

func (m *MockLogger) InfoContext(ctx context.Context, msg string, fields ...Field) {
	m.addCall("INFO", msg, fields, ctx)
}

func (m *MockLogger) ErrorContext(ctx context.Context, msg string, fields ...Field) {
	m.addCall("ERROR", msg, fields, ctx)
}

func (m *MockLogger) DebugContext(ctx context.Context, msg string, fields ...Field) {
	m.addCall("DEBUG", msg, fields, ctx)
}

func (m *MockLogger) WarnContext(ctx context.Context, msg string, fields ...Field) {
	m.addCall("WARN", msg, fields, ctx)
}

// Вспомогательные методы для тестирования.
func (m *MockLogger) addCall(level, msg string, fields []Field, ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.calls = append(m.calls, LogCall{
		Level:   level,
		Message: msg,
		Fields:  fields,
		Context: ctx,
	})
}

func (m *MockLogger) GetCalls() []LogCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	calls := make([]LogCall, len(m.calls))
	copy(calls, m.calls)

	return calls
}

func (m *MockLogger) GetLastCall() *LogCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.calls) == 0 {
		return nil
	}

	return &m.calls[len(m.calls)-1]
}

func (m *MockLogger) GetCallsCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.calls)
}

func (m *MockLogger) GetCallsByLevel(level string) []LogCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []LogCall
	for _, call := range m.calls {
		if call.Level == level {
			result = append(result, call)
		}
	}

	return result
}

func (m *MockLogger) HasCallWithMessage(msg string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, call := range m.calls {
		if call.Message == msg {
			return true
		}
	}

	return false
}

func (m *MockLogger) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = m.calls[:0]
}

func (m *MockLogger) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := fmt.Sprintf("MockLogger with %d calls:\n", len(m.calls))
	for i, call := range m.calls {
		result += fmt.Sprintf("  %d. [%s] %s (fields: %d)\n",
			i+1, call.Level, call.Message, len(call.Fields))
	}

	return result
}

// Убеждаемся, что MockLogger реализует все интерфейсы.
var (
	_ ILogger       = (*MockLogger)(nil)
	_ FieldLogger   = (*MockLogger)(nil)
	_ ContextLogger = (*MockLogger)(nil)
)
