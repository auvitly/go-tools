<h4> 
    <a href="./../../../README.md" align="center"> github.com/auvitly/go-tools </a> 
    > 
    <a href="./../../README.md" align="center"> recovery </a>
    >
    ru
</h4>

## Оглавление
1. [Постановка проблемы](#problem)
2. [Описание](#desc)
3. [Пример использования](#example)

---

<a name="problem"></a>
### 1. Постановка проблемы

Исключения могут стать большой проблемой как в разработке, так и эксплуатации. Для обеспечения
стабильности работы разрабатываемых продуктов требуется обработка исключений. Методика обычно сводится
к следующей конструкции:
```go
func main() {
    defer func(){
        if msg := recover; msg != nil {
            // Do something.
        }   	
    }()
	
    panic("hello")
}
```

Часто разработчики принебрегают практикой обработки исключений, так как добавлять подобные обработчики 
в каждую функцию трудозатратно и представляет рутинную задачу. Обычно паники отлавливаются
на этапе написания тестов, однако даже 100% покрытие кода тестами не всегда может гарантировать 
защиту от исключений (например, когда используется рефлексия). Также иногда требуется выполнить
действия по устранению последствий возникновения паники.

Можно выделить следующие проблемы:
* Рутинность процесса;
* Отсутствие единого подхода;
* Неудобство обработки исключений.

<a name="desc"></a>
### 2. Описание

При работе с исключенями пакет предоставляет следующие возможности:
* предотвратить исключение
* обработать исключение как ошибку
* выполнить действия на основе исключения

Для формирования обработчика исключения используется тип `recovery.Builder`, 
который насыщается параметрами при помощи функций-конструкторов:
* `On`, `OnError` - при возникновении исключения формирует *целевую ошибку* типа [`stderrs.Error`](./../../../stderrs/README.md) 
и `stderrs.Panic` из стандартного набора ошибок, затем выполняет встраивание переданной ошибки, 
если не равна `nil`. Различие между `On` и `OnError` заключается в типе принимаемой ошибки: `On` - `**stderrs.Error`, 
`OnError` - `*error`. 
* `SetMessage` - устанавливает поле Message для сформированной [стандартной ошибки](./../../../stderrs/README.md).
* `WithField`, `WithFieldIf`, `WithFields`, `WithFieldsIf` - методы добавления полей для целевой 
[стандартной ошибки](./../../../stderrs/README.md). Суффикс `If` позволяет установить условие добавления.
* `WithHandlers`, `WithHandlersIf` - позволяет установить обработчики исключения, которые должны типу `recovery.Handler`.
* `WithAsyncHandlers`, `WithAsyncHandlersIf` - позволяет установить асинхронные обработчики исключения, которые должны 
типу `recovery.AsyncHandler`.
* `WithoutHandlers` - позволяет сбросить все установленные обработчики. Необходим, когда требуется исключить 
использование глобальные обработчиков.
* `Do` - выполняет перехват исключения. **По умолчанию** учитывает зарегистрированные 
глобальные обработчики, для отказа от глобальных обработчиков используйте метод `WithoutHandlers`.

Регистрация глобальных обработчиков реализуется методами:
* `RegistryHandlers` - позволяет установить глобальные обработчики исключения, которые должны типу `recovery.Handler`.
* `RegistryAsyncHandlers` - позволяет установить глобальные асинхронные обработчики исключения, которые должны типу 
`recovery.AsyncHandler`.

<a name="example"></a>
### 3. Примеры использования

#### 3.1 Предотвращение исключения

Для предотвращения исключения достаточно просто добавить `defer recovery.Do()` в начало функции:

```go
func fn() {
    defer recovery.Do()
	
    panic("I'm the exception")
}

func main() {
    fn()
	
    slog.Info("Hello!")
}
```

Результат:

```text
2024/05/29 21:14:58 INFO Hello!
```

#### 3.2 Предотвращение исключения с обработчиками

Для пользовательской обработки исключений используются методы `WithHandlers` или `WithAsyncHandler`:

```go
func asyncHandler(any) {
    slog.Info("I'm an asynchronous processor")
}

func syncHandler(any) error {
    slog.Info("I'm a synchronous processor")

    return nil
}

func fn() {
    defer recovery.
        WithHandlers(syncHandler).
        WithAsyncHandlers(asyncHandler).
        Do()
    
    panic("I'm the exception")
}

func main() {
    fn()
}
```

Результат:

```text
2024/05/29 21:23:51 INFO I'm a synchronous processor  
2024/05/29 21:23:51 INFO I'm an asynchronous processor
```

#### 3.3 Передача в обработчики специфических параметров

Используя паттерн функций высшего порядка можно передавать в функцию обработчик необходимые данные, например
контекст:

```go
func syncHandler(any) error {
    return errors.New("syncHandler error: I'm the error")
}

func fn() (err *stderrs.Error) {
    defer recovery.WithHandlers(syncHandler).On(&err).Do()
    
    panic("I'm the exception")
}

func main() {
    err := fn()
    if err != nil {
        slog.Error(err.Error())
    }
}
```

Результат:

```text
2024/05/29 21:23:51 INFO I'm a synchronous processor
```

> Отметим, что контекст передаваемый в метод `DoContext` позволяет ограничить время ожидаения асинхронных 
> обработчиков, в отличии от метода `Do`, который ожидает окончания всех асинхронных обработчиков.

#### 3.4 Обработка ошибок обработчиков

Как следует из сигнатуры `recovery.Handler`: синхронные обработчики возвращают ошибку. 
Ошибки, которые возвращают обработчики, обогащают целевую ошибку. Обратимся к модели [стандартной ошибки](./../../../stderrs/README.md):
```go
// Error - unified model.
type Error struct {
    Code    string         `json:"code"`
    Message string         `json:"message"`
    Embed   error          `json:"embed"`
    Wraps   []string       `json:"wraps"`
    Fields  map[string]any `json:"fields"`
    Codes   struct {
        GRPC codes.Code `json:"grpc"`
        HTTP int        `json:"http"`
    } `json:"codes"`
}
```

Стандартная ошибка обладает полем `Embed`, которое позволяет хранить все встраиваемые ошибки. Таким образом, 
возвращаемая ошибка `stderrs.Panic` будет содержать встроенные ошибки, если синхронные обработчики их вернут.

```go
func syncHandler(any) error {
    return errors.New("syncHandler error: I'm the error")
}

func fn() (err *stderrs.Error) {
    defer recovery.WithHandlers(syncHandler).On(&err).Do()
    
    panic("I'm the exception")
}

func main() {
    err := fn()
    if err != nil {
        slog.Error(err.Error())
    }
}
```

Результат:

```text
2024/05/29 22:11:59 ERROR {"code": "panic", "message": "internal server error: unhandled exception", "fields": {"panic":"I'm the exception"}, "embed": ["syncHandler error: I'm the error"]}
```

#### 3.5 Обработка ошибок обработчиков для пользовательских типов ошибок

Если вам требуется обогатить собственную имлементацию, то можно воспользоваться обработчиками и 
механизмом функций высшего порядка:

```go
type MyError string

func (e MyError) Error() string {
    return string(e)
}

func NewMyError(msg string) *MyError {
    return (*MyError)(&msg)
}

func enrichMyError(ptr **MyError) func(msg any) error {
    return func(msg any) (err error) {
        defer func() {
            if err == nil {
                return
            }
            
            if *ptr == nil {
                *ptr = new(MyError)
            }
            
            **ptr = MyError(
                fmt.Sprintf(
                    "%s: %s: %s",
                    **ptr, err.Error(), msg,
                ),
            )
        }()
        
        return errors.New("enrichMyError error: I'm the error")
    }
}

func fn() (err *MyError) {
    defer recovery.WithHandlers(enrichMyError(&err)).Do()
    
    err = NewMyError("fn error")
    
    panic("I'm the exception")
}

func main() {
    err := fn()
    if err != nil {
        slog.Error(err.Error())
    }
}
```

Результат:

```text
2024/05/29 22:04:13 ERROR fn error: enrichMyError error: I'm the error: I'm the exception
```

#### 3.6 Предотвращение исключений в пользовательских обработчиках

Пользовательские обработчики также могут содержать исключения. 
Кажется идеей будет использовать пакет `recovery`, однако 
это не требуется для синхронных обработчиков, будет 
выполнен автоматический перехват исключений и
предоставлен стек вызов в целевую ошибку ошибку:

```go
func globalPanicHandler(any) (err error) {
    panic("globalPanicHandler")
    
    return nil
}

func fn() (err *stderrs.Error) {
    defer recovery.On(&err).Do()
    
    panic("I'm the exception")
}

func main() {
    recovery.RegistryHandlers(globalPanicHandler)
    
    err := fn()
    if err != nil {
        slog.Error(err.Error())
    }
}
```

Результат:

```
2024/05/30 00:19:15 ERROR {"code": "panic", "message": "internal server error: unhandled exception", "fields": {"panic":"I'm the exception"}, "embed": [{"code": "panic", "fields": {"panic":"Hello world!","stack":"goroutine 1 [running]:\nruntime/debug.Stack()\n\tC:/Program Files/Go/src/runtime/debug/stack.go:24 +0x5e\ngithub.com/auvitly/go-tools/recovery.Builder.useSync.func1()\n\tF:/Work/projects/git/auvitly/go-tools/recovery/builder.go:125 +0x38a\npanic({0x113bd00?, 0x120c870?})\n\tC:/Program Files/Go/src/runtime/panic.go:914 +0x21f\nmain.globalPanicHandler({0x17c01f00108, 0x10})\n\tF:/Work/projects/git/auvitly/go-tools/examples/test/main.go:10 +0x25\ngithub.com/auvitly/go-tools/recovery.Builder.useSync({{0xc000044060, 0x1, 0x1}, {0x1447320, 0x0, 0x0}, 0x0, 0xc000044058, {0x11a4fad, 0x2a}, ...}, ...)\n\tF:/Work/projects/git/auvitly/go-tools/recovery/builder.go:145 +0x77\ngithub.com/auvitly/go-tools/recovery.Builder.handle({{0xc000044060, 0x1, 0x1}, {0x1447320, 0x0, 0x0}, 0x0, 0xc000044058, {0x11a4fad, 0x2a}, ...}, ...)\n\tF:/Work/projects/git/auvitly/go-tools/recovery/builder.go:265 +0x39c\ngithub.com/auvitly/go-tools/recovery.Builder.recovery({{0xc000044060, 0x1, 0x1}, {0x1447320, 0x0, 0x0}, 0x0, 0xc000044058, {0x11a4fad, 0x2a}, ...}, ...)\n\tF:/Work/projects/git/auvitly/go-tools/recovery/builder.go:189 +0x110\ngithub.com/auvitly/go-tools/recovery.Builder.Do({{0xc000044060, 0x1, 0x1}, {0x1447320, 0x0, 0x0}, 0x0, 0xc000044058, {0x0, 0x0}, ...})\n\tF:/Work/projects/git/auvitly/go-tools/recovery/builder.go:102 +0x6c\npanic({0x113bd00?, 0x120c880?})\n\tC:/Program Files/Go/src/runtime/panic.go:920 +0x270\nmain.fn()\n\tF:/Work/projects/git/auvitly/go-tools/examples/test/main.go:18 +0x405\nmain.main()\n\tF:/Work/projects/git/auvitly/go-tools/examples/test/main.go:24 +0x3fb\n"}}]}
```

#### 3.7 Предотвращение исключений в асинхронных пользовательских обработчиках

Если обработчики запускаются через метод `Do`, то все исключения будут добавлены в **целевую ошибку**.

Асинхронные обработчики могут завершиться исключением вне срока жизни контекста переданного в `DoContext`. 
Глобальные обработчики содержащие исключения могут привести к `stack overflow` из-за рекурсивности,
если использовать `recovery.Do`. Для задач, когда необходимо обработать исключения без задействования
глобальных обработчиков существует метод `WithoutHandlers`. При вызове этого метода, все пользовательские 
обработчики очищаются для локального вызова:

```go
func log(msg any) (err error) {
    slog.Error("exception", "msg", msg)
    
    return nil
}

func globalAsyncPanicHandler(any) {
    defer recovery.WithoutHandlers().WithHandlers(log).Do()
    
    time.Sleep(2 * time.Second)
    
    panic("globalPanicHandler")
}

func fn(ctx context.Context) (err *stderrs.Error) {
    defer recovery.On(&err).DoContext(ctx)
    
    panic("I'm the exception")
}

func main() {
    recovery.RegistryAsyncHandlers(globalAsyncPanicHandler)
    
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    
    err := fn(ctx)
    if err != nil {
        slog.Error(err.Error())
    }
    
    time.Sleep(2 * time.Second)
}
```

Результат:
```
2024/05/30 00:33:37 ERROR {"code": "panic", "message": "internal server error: unhandled exception", "fields": {"panic":"I'm the exception"}}
2024/05/30 00:33:38 ERROR exception msg=globalPanicHandler
```

> Постарайтесь не допускать исключений в обработчиках исключений!

#### 3.8 Перехват исключений как механизм защиты обработки данных

Задача:

Имеется выборка записей, малая часть которых (~1%) содержит неожиданное наполнение.
При попытке конвертации модели из выборки в другую модель, возникает исключение из-за 
обращения к пустому указателю. Если выполнить базовый сценарий обработки, то 
остальные (~99%) записей станут недоступны. Для обеспечения возврата всех валидных 
записей можно воспользоваться механизмом обработки исключений:

```go
type modelDB struct {
	String string
	PtrInt *int
}

type modelAPI struct {
	String string
	Int    int
}

func log(values ...any) func(any) {
    return func(msg any) {
        slog.Error("we obtain panic",
            slog.Any("panic", msg),
            slog.Any("stack", string(debug.Stack())),
            slog.Any("values", values),
        )
    }
}

func convert(item modelDB) (result *modelAPI) {
    defer recovery.WithAsyncHandlers(log(item)).Do()
    
    return &modelAPI{
        String: item.String,
        Int:    *item.PtrInt,
    }
}

func main() {
    var (
        records = []modelDB{
            {
                String: "valid",
                PtrInt: new(int),
            },
            {
                String: "not valid",
                PtrInt: nil,
            },
        }
        results []modelAPI
    )
    
    for _, record := range records {
        if result := convert(record); result != nil {
            results = append(results, *result)
        }
    }
    
    slog.Info("response", "results", results)
}
```

Результат:

```
2024/05/30 01:28:32 ERROR we obtain panic panic="runtime error: invalid memory address or nil pointer dereference" stack="goroutine 6 [running]:\nruntime/debug.Stack()\n\tC:/Program Files/Go/src/runtime/debug/stack.go:24 +0x5e\nmain.convert.log.func1({0xcf42e0?, 0xf7f800?})\n\tF:/Work/projects/git/auvitly/go-tools/examples/test/main.go:23 +0xbc\ngithub.com/auvitly/go-tools/recovery.Builder.useAsync({{0xfe8320, 0x0, 0x0}, {0xc000044058, 0x1, 0x1}, 0x0, 0x0, {0xd462c5, 0x2a}, ...}, ...)\n\tF:/Work/projects/git/auvitly/go-tools/recovery/builder.go:175 +0x8d\ngithub.com/auvitly/go-tools/recovery.Builder.handle.func1(0x0?)\n\tF:/Work/projects/git/auvitly/go-tools/recovery/builder.go:253 +0xc5\ncreated by github.com/auvitly/go-tools/recovery.Builder.handle in goroutine 1\n\tF:/Work/projects/git/auvitly/go-tools/recovery/builder.go:250 +0xa7\n" values="[{String:not valid PtrInt:<nil>}]"
2024/05/30 01:28:32 INFO response results="[{String:valid Int:0}]"

```

#### 3.9 Прочие примеры

* Демонстрация механизма перехвата паники c пользовательскими обработчиками [[ссылка](../../../examples/relax/main.go)]