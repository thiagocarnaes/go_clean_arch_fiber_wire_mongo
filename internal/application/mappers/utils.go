package mappers

import "math"

// calculateTotalPages calcula o número total de páginas
func calculateTotalPages(total int64, perPage int64) int64 {
	return int64(math.Ceil(float64(total) / float64(perPage)))
}
