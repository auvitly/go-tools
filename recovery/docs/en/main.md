<h4> 
    <a href="./../../../README.md" align="center"> github.com/auvitly/go-tools </a> 
    > 
    <a href="./../../README.md" align="center"> recovery </a>
    >
    en
</h4>

## Chapters
1. [Formulation of the problem](#problem)
2. [Description](#desc)
3. [Examples](#example)

---

<a name="problem"></a>
### 1. Formulation of the problem

Exceptions can be a big problem in both development and operations. To provide
The stability of the products being developed requires exception handling. The technique usually comes down to
to the following construction:
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

Developers often neglect the practice of exception handling because adding such handlers 
into each function is labor-intensive and represents a routine task. Usually panics are caught
at the stage of writing tests, however, even 100% code coverage by tests cannot always guarantee 
exception protection (for example, when reflection is used). Also sometimes you need to do
actions to eliminate the consequences of panic.

The following problems can be identified:
* Routine of the process;
* Lack of a unified approach;
* Inconvenient exception handling.

<a name="desc"></a>
### 2. Description

When working with exceptions, the package provides the following capabilities:
* prevent exception
* treat the exception as an error
* perform actions based on exception

To generate an exception handler, use the type `recovery.Builder`, 
which is saturated with parameters using constructor functions:
* `On`, `OnError` - when an exception occurs, generates a *target error* of type [`stderrs.Error`](./../../../stderrs/README.md) 
and `stderrs.Panic` from the standard error set, then injects the passed error, 
if not equal`nil`. The difference between `On` and `OnError` is the type of error accepted: `On` - `**stderrs.Error`, 
`OnError` - `*error`. 
* `SetMessage` - sets the Message field for the generated [standard error](./../../../stderrs/README.md).
* `WithField`, `WithFieldIf`, `WithFields`, `WithFieldsIf` - methods for adding fields for the target
  [standard error](./../../../stderrs/README.md). The `If` suffix allows you to set the condition for adding.
* `WithHandlers`, `WithHandlersIf` - allows you to set exception handlers that must be of the `recovery.Handler` type.
* `WithAsyncHandlers`, `WithAsyncHandlersIf` - allows you to set asynchronous exception handlers that must be of type 
`recovery.AsyncHandler`.* `WithoutHandlers` - reset all installed handlers. Required when you want to exclude use of 
global handlers.
* `Do`, `DoContext` - мThe methods allow you to catch an exception. **Default** takes into account registered 
global handlers; to disable global handlers, use the `WithoutHandlers` method.

Registration of global handlers is implemented using the following methods:
* `RegistryHandlers` - allows you to set global exception handlers that must be of the `recovery.Handler` type.
* `RegistryAsyncHandlers` - allows you to set global asynchronous exception handlers that must be of type `recovery.AsyncHandler`.

<a name="example"></a>
### 3. Examples

#### 3.1 Preventing an exception

To prevent an exception, simply add `defer recovery.Do()` to the beginning of the function:

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

Out:

```text
2024/05/29 21:14:58 INFO Hello!
```

#### 3.2 Preventing an Exception with Handlers

For custom exception handling, use the `WithHandlers` or `WithAsyncHandler` methods:

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

Out:

```text
2024/05/29 21:23:51 INFO I'm a synchronous processor  
2024/05/29 21:23:51 INFO I'm an asynchronous processor
```

#### 3.3 Passing specific parameters to handlers

Using the higher order function pattern, you can pass the necessary data to the handler function, for example
context:

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

Out:

```text
2024/05/29 21:23:51 INFO I'm a synchronous processor
```

> Note that the context passed to the `DoContext` method allows you to limit the waiting time for asynchronous 
> handlers, unlike the `Do` method, which waits for all asynchronous handlers to finish.

#### 3.4 Обработка ошибок обработчиков

As follows from the `recovery.Handler` signature: synchronous handlers return an error. 
The errors that handlers return enrich the target error. Let us turn to the model [standard error](./../../../stderrs/README.md):
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

A standard error has an `Embed` field that allows you to store all embedded errors. Thus, 
the returned error `stderrs.Panic` will contain inline errors if synchronous handlers return them.

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

Out:

```text
2024/05/29 22:11:59 ERROR {"code": "panic", "message": "internal server error: unhandled exception", "fields": {"panic":"I'm the exception"}, "embed": ["syncHandler error: I'm the error"]}
```

#### 3.5 Handling error handlers for custom error types

If you need to enrich your own implementation, you can use handlers and 
mechanism of higher order functions:

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

Out:

```text
2024/05/29 22:04:13 ERROR fn error: enrichMyError error: I'm the error: I'm the exception
```

#### 3.6 Preventing exceptions in handlers

Custom handlers can also contain exceptions. 
It seems like the idea would be to use the `recovery` package, however 
this is not required for synchronous handlers, it will 
exceptions are automatically caught and
provided stack call to target error:

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

#### 3.7 Preventing exceptions in asynchronous user handlers

If handlers are run via the `Do` method, then all exceptions will be added to the **target error**.

Asynchronous handlers may fail with an exception outside the lifetime of the context passed to the `DoContext`. 
Global exception handlers can lead to `stack overflow`, if you use `recovery.Do`. 
For tasks where it is necessary to handle exceptions without using
global handlers there is a method `WithoutHandlers`. When this method is called, all user 
handlers are cleared to be called locally:

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

Out:

```
2024/05/30 00:33:37 ERROR {"code": "panic", "message": "internal server error: unhandled exception", "fields": {"panic":"I'm the exception"}}
2024/05/30 00:33:38 ERROR exception msg=globalPanicHandler
```

> Try not to throw exceptions in exception handlers!

#### 3.8 Catching exceptions as a data processing protection mechanism

Task:

There is a sample of records, a small part of which (~1%) contains unexpected content.
When trying to convert a model from a selection to another model, an exception occurs due to 
references to the null pointer. If you run the basic processing scenario, then 
the remaining (~99%) records will become unavailable. To ensure the return of all valid 
records, you can use the exception handling mechanism:

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

Out:

```
2024/05/30 01:28:32 ERROR we obtain panic panic="runtime error: invalid memory address or nil pointer dereference" stack="goroutine 6 [running]:\nruntime/debug.Stack()\n\tC:/Program Files/Go/src/runtime/debug/stack.go:24 +0x5e\nmain.convert.log.func1({0xcf42e0?, 0xf7f800?})\n\tF:/Work/projects/git/auvitly/go-tools/examples/test/main.go:23 +0xbc\ngithub.com/auvitly/go-tools/recovery.Builder.useAsync({{0xfe8320, 0x0, 0x0}, {0xc000044058, 0x1, 0x1}, 0x0, 0x0, {0xd462c5, 0x2a}, ...}, ...)\n\tF:/Work/projects/git/auvitly/go-tools/recovery/builder.go:175 +0x8d\ngithub.com/auvitly/go-tools/recovery.Builder.handle.func1(0x0?)\n\tF:/Work/projects/git/auvitly/go-tools/recovery/builder.go:253 +0xc5\ncreated by github.com/auvitly/go-tools/recovery.Builder.handle in goroutine 1\n\tF:/Work/projects/git/auvitly/go-tools/recovery/builder.go:250 +0xa7\n" values="[{String:not valid PtrInt:<nil>}]"
2024/05/30 01:28:32 INFO response results="[{String:valid Int:0}]"

```

#### 3.9 Прочие примеры

* Demonstration of the panic interception mechanism with custom handlers [[link](../../../examples/relax/main.go)]