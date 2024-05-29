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
* Удобство обработки.

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
* `WithHandlers` - позволяет установить обработчики исключения, которые должны типу `recovery.Handler`.
* `WithAsyncHandlers` - позволяет установить асинхронные обработчики исключения, которые должны типу `recovery.AsyncHandler`.
* `WithoutHandlers` - позволяет сбросить все установленные обработчики. Необходим, когда требуется исключить 
использование глобальные обработчиков.
* `Do`, `DoContext` - методы позволяют выполнить перехват исключения. **По умолчанию** учитывает зарегистрированные 
глобальные обработчики, для отказа от глобальных обработчиков используйте метод `WithoutHandlers`.

Регистрация глобальных обработчиков реализуется методами:
* `RegistryHandlers` - позволяет установить глобальные обработчики исключения, которые должны типу `recovery.Handler`.
* `RegistryAsyncHandlers` - позволяет установить глобальные асинхронные обработчики исключения, которые должны типу `recovery.AsyncHandler`.

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

Для реакции на исключение используется `WithHandlers` или `WithAsyncHandler`. 

Пример с синхронным и асинхронными обработчиками:

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

#### 3.2 Передача в обработчики специфических параметров

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

Как следуюет из сингатуры `recovery.Handler`: синхронные обработчики возвращают ошибку. 
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

func fn(ctx context.Context) (err *stderrs.Error) {
    defer recovery.WithHandlers(syncHandler).On(&err).Do()
    
    panic("I'm the exception")
}

func main() {
    err := fn(context.TODO())
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

Если вам требуется обогатить собственную имлементацию, то можно воспользоваться механзимом обработчиков и 
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

```text
2024/05/29 22:04:13 ERROR fn error: enrichMyError error: I'm the error: I'm the exception
```

#### 3.6 Прочие примеры

* Демонстрация механизма перехвата паники c пользовательскими обработчиками [[ссылка](../../../examples/relax/main.go)]