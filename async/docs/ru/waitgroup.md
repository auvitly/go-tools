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

`WaitContext` может быть использован, когда требуется дождаться окончания ожидания с ограничением времени жизни контекста.

```go
func main() {
    ctx, cancel := context.WithTimeout()	
}
```
