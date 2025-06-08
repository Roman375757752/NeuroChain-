package main

import (
    "fmt"
    "github.com/Roman375757752/neurochain/node"
    "github.com/Roman375757752/neurochain/simulator"
    "time"
)

func main() {
    n := node.NewNode(9000)
    n.Start()
    fmt.Println("Запущен нейро-узел")

    simulator.RunSimulation(n)

    // Чтобы программа не завершалась сразу, спим, или ждем ввода пользователя
    time.Sleep(10 * time.Second)
}
