<h4> 
    <a href="./../../../README.md" align="center"> github.com/auvitly/go-tools </a> 
    > 
    <a href="./../../README.md" align="center"> async </a>
    >
    <a href="./main.md" align="center"> ru </a>
    >
    WaitGroup
</h4>

<a name="desc"></a>
### 1. Описание

Является адаптером над `sync.WaitGroup` и предоставляет несколько дополнительных методов для ожидания:
* `WaitContext` - ожидание окончания работы группы, по времени жизни контекста. Возвращает ошибку, если контекста 
завершился раньше, чем произошло ожидание, то будет возвращена одна из ошибок пакета `context`: 
`context.DeadlineExceeded` или `context.Canceled`.
* `WaitDone` - возвращает канал, который закрывается, когда ожидание будет выполнено

### 2. Примеры использования

#### 2.1 Завершение ожидания WaitGroup по сроки жизни контекста

`WaitContext` может быть использован, когда требуется дождаться окончания ожидания с ограничением времени жизни контекста.

```go
func goroutine(i int, wg *async.WaitGroup) {
    defer wg.Done()
    
    var d = rand.Int63() * int64(time.Second)
    
    time.Sleep(time.Duration(d))
    
    slog.Info(fmt.Sprintf("I'm goroutine #%d", i))
}

func main() {
    var wg async.WaitGroup
    
    ctx, _ := context.WithTimeout(context.Background(), time.Second)
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        
        go goroutine(i, &wg)
    }
    
    if err := wg.WaitContext(ctx); err != nil {
        panic(err)
    }
}
```

Один из вариантов вывода:
```
2024/06/01 19:00:04 INFO I'm goroutine #1
2024/06/01 19:00:04 INFO I'm goroutine #3
2024/06/01 19:00:04 INFO I'm goroutine #4
2024/06/01 19:00:04 INFO I'm goroutine #7
2024/06/01 19:00:04 INFO I'm goroutine #9
panic: context deadline exceeded

goroutine 1 [running]:
main.main()
        F:/Work/projects/git/auvitly/go-tools/examples/wg/main.go:34 +0xeb

Process finished with the exit code 2
```

#### 2.2 Обработка события завершения ожидания через `select`

```go
func goroutine(wg *async.WaitGroup) {
    defer wg.Done()
    
    time.Sleep(5 * time.Second)
}

func main() {
    var (
        wg     async.WaitGroup
        ticker = time.NewTicker(time.Second)
    )
    
    wg.Add(1)
    
    go goroutine(&wg)
    
    for {
        select {
        case <-wg.WaitDone():
            slog.Info("all goroutines done")
    
            return
        case <-ticker.C:
            slog.Info("1 second has passed")
        }
    }
}
```

Результат:
```
2024/06/01 19:58:14 INFO 1 second has passed
2024/06/01 19:58:15 INFO 1 second has passed
2024/06/01 19:58:16 INFO 1 second has passed
2024/06/01 19:58:17 INFO 1 second has passed
2024/06/01 19:58:18 INFO all goroutines done
```