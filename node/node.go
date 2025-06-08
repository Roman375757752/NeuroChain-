package node

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strconv"
    "strings"
    "sync"
)

type Transaction struct {
    ID       int
    From     string
    To       string
    Amount   int
    Nonce    int64
    Hash     string
    Mined    bool
}

type Node struct {
    Port          int
    Balances      map[string]int
    BalancesMutex sync.Mutex
    TxIDCounter   int
    TxMutex       sync.Mutex
    ProcessedTx   map[int]*Transaction
}

func NewNode(port int) *Node {
    return &Node{
        Port:        port,
        Balances:    make(map[string]int),
        ProcessedTx: make(map[int]*Transaction),
    }
}

func (n *Node) Start() {
    go n.listen()
}

func (n *Node) listen() {
    // заглушка - слушаем порт, но не используем в этом примере
    fmt.Println("Нейро-узел слушает на порту", n.Port)
}

func (n *Node) MineTransaction(tx *Transaction) bool {
    // Простейший Proof of Work: ищем хеш с 1 нулём в начале (для скорости)
    prefix := "0"
    nonce := 0
    for nonce < 10000 {
        data := tx.From + tx.To + strconv.Itoa(tx.Amount) + strconv.Itoa(nonce)
        hashBytes := sha256.Sum256([]byte(data))
        hash := hex.EncodeToString(hashBytes[:])
        if strings.HasPrefix(hash, prefix) {
            tx.Hash = hash
            tx.Nonce = int64(nonce)
            tx.Mined = true
            return true
        }
        nonce++
    }
    return false
}

func (n *Node) ProcessTransaction(from, to string, amount int) bool {
    // Проверка баланса и обновление балансов
    n.BalancesMutex.Lock()
    defer n.BalancesMutex.Unlock()

    fromBalance, ok := n.Balances[from]
    if !ok || fromBalance < amount {
        return false
    }
    if from == to {
        return false
    }

    n.TxMutex.Lock()
    n.TxIDCounter++
    txID := n.TxIDCounter
    n.TxMutex.Unlock()

    tx := &Transaction{
        ID:     txID,
        From:   from,
        To:     to,
        Amount: amount,
    }

    if !n.MineTransaction(tx) {
        fmt.Println("Майнинг не удался для транзакции ID:", tx.ID)
        return false
    }

    // Обновляем балансы после успешного майнинга
    n.Balances[from] -= amount
    n.Balances[to] += amount

    n.ProcessedTx[txID] = tx
    return true
}

func (n *Node) SetInitialBalance(wallet string, amount int) {
    n.BalancesMutex.Lock()
    defer n.BalancesMutex.Unlock()
    n.Balances[wallet] = amount
}

func (n *Node) GetBalance(wallet string) int {
    n.BalancesMutex.Lock()
    defer n.BalancesMutex.Unlock()
    return n.Balances[wallet]
}
