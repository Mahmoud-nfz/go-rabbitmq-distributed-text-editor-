package rabbitmq 

import (
    "fmt"
    "math/rand"
    "time"
)

func RandomID() string {
	length := 10
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, length+2)
    rand.Read(b)
    return fmt.Sprintf("%x", b)[2 : length+2]
}