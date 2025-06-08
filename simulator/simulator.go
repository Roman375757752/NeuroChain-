package simulator

import (
    "fmt"
    "math/rand"
    "github.com/Roman375757752/neurochain/node"
    "os"
    "strconv"
    "sync"
    "time"
)

func RunSimulation(n *node.Node) {
    walletCount := 5000
    initialBalance := 1000
    txCount := 10000

    wallets := make([]string, walletCount)
    for i := 0; i < walletCount; i++ {
        wallets[i] = "wallet_" + strconv.Itoa(i+1)
        n.SetInitialBalance(wallets[i], initialBalance)
    }

    fmt.Println("Создано кошельков:", walletCount)

    var wg sync.WaitGroup
    rand.Seed(time.Now().UnixNano())

    successes := 0
    failures := 0

    start := time.Now()

    txChan := make(chan bool, txCount)

    for i := 0; i < txCount; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            fromIdx := rand.Intn(walletCount)
            toIdx := rand.Intn(walletCount)
            for toIdx == fromIdx {
                toIdx = rand.Intn(walletCount)
            }
            amount := rand.Intn(initialBalance) + 1

            ok := n.ProcessTransaction(wallets[fromIdx], wallets[toIdx], amount)
            txChan <- ok
        }()
    }

    wg.Wait()
    close(txChan)

    for ok := range txChan {
        if ok {
            successes++
        } else {
            failures++
        }
    }

    duration := time.Since(start).Seconds()
    tps := float64(successes) / duration

    fmt.Printf("Транзакций успешно: %d\n", successes)
    fmt.Printf("Транзакций неудачно: %d\n", failures)
    fmt.Printf("Общее время (сек): %.2f\n", duration)
    fmt.Printf("TPS (транзакций в секунду): %.2f\n", tps)

    exportTransactionsCSV(n)
}

func exportTransactionsCSV(n *node.Node) {
    f, err := os.Create("transactions.csv")
    if err != nil {
        fmt.Println("Ошибка при создании CSV-файла:", err)
        return
    }
    defer f.Close()

    fmt.Fprintln(f, "from,to,amount")
    for _, tx := range n.ProcessedTx {
        fmt.Fprintf(f, "%s,%s,%d\n", tx.From, tx.To, tx.Amount)
    }

    fmt.Println("CSV файл transactions.csv создан.")
}
