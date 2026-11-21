package erasmus

import "strconv"

func entero(texto string) int {
        numero, _ := strconv.Atoi(texto)
        return numero
}