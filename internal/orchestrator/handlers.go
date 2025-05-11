package orchestrator

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/dimakirio/calculatorv1/internal/models"
	"github.com/dimakirio/calculatorv1/pkg/config"
	"github.com/dimakirio/calculatorv1/pkg/logger"
	"github.com/Knetic/govaluate" // Импорт библиотеки для вычисления выражений
	"github.com/google/uuid"
)

var (
	expressions = make(map[string]models.Expression)
	mu          sync.Mutex
)

type Orchestrator struct {
	log *logger.Logger
	cfg *config.Config
}

func NewOrchestrator(log *logger.Logger, cfg *config.Config) *Orchestrator {
	return &Orchestrator{log: log, cfg: cfg}
}

func (o *Orchestrator) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	// Проверяем корректность выражения
	if !isValidExpression(req.Expression) {
		http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
		return
	}

	// Вычисляем выражение
	result, err := evaluateExpression(req.Expression)
	if err != nil {
		http.Error(w, "Failed to evaluate expression", http.StatusUnprocessableEntity)
		return
	}

	id := uuid.New().String()
	mu.Lock()
	expressions[id] = models.Expression{
		ID:     id,
		Status: "completed",
		Result: result,
	}
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (o *Orchestrator) HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var exprs []models.Expression
	for _, expr := range expressions {
		exprs = append(exprs, expr)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": exprs})
}

func (o *Orchestrator) HandleGetExpressionByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/expressions/"):]
	mu.Lock()
	expr, exists := expressions[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
}

// evaluateExpression вычисляет значение выражения
func evaluateExpression(expression string) (float64, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return 0, err
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		return 0, err
	}

	return result.(float64), nil
}

// isValidExpression проверяет корректность выражения
func isValidExpression(expression string) bool {
	// Простая проверка на наличие некорректных символов
	for _, char := range expression {
		if !isValidCharacter(char) {
			return false
		}
	}
	return true
}

// isValidCharacter проверяет, является ли символ допустимым
func isValidCharacter(char rune) bool {
	// Разрешенные символы: цифры, операторы (+,-,*,/), пробелы, скобки
	return (char >= '0' && char <= '9') ||
		char == '+' || char == '-' || char == '*' || char == '/' ||
		char == ' ' || char == '(' || char == ')'
}
