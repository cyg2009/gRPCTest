package main
import "fmt"

// Greet returns a simple greeting.
func Handler(data string) string {
    return fmt.Sprintf("trigger data: %s\n", data)
}