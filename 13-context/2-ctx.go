package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	i, err := Slow(ctx, "ac")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(i)
}
func Slow(ctx context.Context, input string) (int, error) {
	//sql.DB{}.ExecContext(ctx, "SELECT 1")
	i, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("error converting %s to int: %w", input, err)
	}
	time.Sleep(2 * time.Millisecond)
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return i, nil
	}
}
