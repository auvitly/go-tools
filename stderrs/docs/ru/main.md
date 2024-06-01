<h4> 
    <a href="./../../../README.md" align="center"> github.com/auvitly/go-tools </a> 
    > 
    <a href="./../../README.md" align="center"> stderrs </a>
    >
    ru
</h4>

## Оглавление
1. [Постановка проблемы](#problem)
2. [Описание](#desc)

---

<a name="problem"></a>
### 1. Постановка проблемы

В процессе разработки API необходимо обрабатывать исключения. 
Привычный формат взаимодействия строится на использовании стандартных пакетов `fmt` и `errors`.

Обычно, после любого внешнего вызова необходимо оборачивать ошибку:
```go
package example

import (
	"fmt"
	".../external"
)

func ExternalCall() (any, error) {
    result, err := external.DoSomething()
    if err != nil {
        return nil, fmt.Errorf("external.DoSomething: %w", err)
    }
    
    return result, nil
}
```

Можно выделить следующие проблемы:
* Существует множество соглашений по оборачиванию ошибок, однако пользователь должен 
получить простое сообщение, по которому легко можно понять, что пошло не так. Если отдать
пользователю ошибку, которую собрал сервис в процессе множественного оборачивания, то
информативность сообщения сильно падает. 
* Текст ошибки, отдаваемый пользователю, может быть представлен на нескольких языках, 
однако логирование ошибок всегда должно быть выполнено на каком-то одном языке и иметь
полноту информации для определения источника проблемы.
* Определение статус кода ошибки при большом числе внутренних вызовов порождает 
необходимость выполнять маршрутизацию моделей ошибок, что превращается в отдельную задачу.
* При возникновении исключений (panic) в сервисах, пользователь должен получить понятную 
ошибку о недоступности сервиса, разработчик исчерпывающее описание места возникновения 
проблемы, а сервис действия по откату неудачной транзакции. 
> Обработка исключений решается пакетом `recovery` из набора `go-tools`.
> Ознакомиться с пакетом `recovery` можно по [ссылке](./../../../recovery/README.md).

<a name="desc"></a>
### 2. Описание
Пакет представляет унифицированную модель (далее - стандартная ошибка):

```go
package stderrs

import "google.golang.org/grpc/codes"

type Error struct {
	// Code - код ошибки.
	Code    string         `json:"code"`
	// Message - сообщение, которое получит пользователь.
	Message string         `json:"message"`
	// Embed - встроенные ошибки стандартного типа.
	// Note: позволяет добавить к ошибке стандартный интерфейс error.
	Embed   error          `json:"embed"`
	// Wrap - возможность обернуть ошибку в моменте вызова функции Error.
	// Note: используется для передачи объекта вызова и метода, в котором была получена. 
	Wraps   []string       `json:"wraps"`
	// Fields - позволяет передать спефицичные параметры определяющие исходную ошибку.
	// Например, при обращении к удаленным ресурсам можно передать модель запроса.
	// Fields обрабатывается пакетом json, поэтому должно значения по ключам должны
	// быть обычными типами, либо структурами, либо отвечать интерфейсу json.Marshaler.
	// В противном случае Fields будет пропущен.
	Fields  map[string]any `json:"fields"`
	// Codes - статус коды для различных транспортных уровней.
	Codes   struct {
		GRPC codes.Code `json:"grpc"`
		HTTP int        `json:"http"`
	} `json:"codes"`
}
```

Стандартная ошибка может проинициализирована явно, либо при помощи конструктора на основе функции `New`, например:
```go
var Canceled = New("canceled").
    SetGRPCCode(codes.Canceled).
    SetHTTPCode(499)
```

Стандартная ошибка имеет следующий набор методов:
* `Is` - метод для обеспечения стандартного интерфейса. Может быть использован явно вместо `errors.Is`;
* `Unwrap` - метод для обеспечения стандартного интерфейса. Восстанавливает все встроенные ошибки;
* `SetCode` - установить код ошибки, HTTP и GRPC статус на основе стандартной ошибки;
* `SetMessage` - установить клиентское сообщение;
* `SetHTTPCode` - установить HTTP статус код;
* `SetGRPCCode` - установить GRPC статус код;
* `EmbedErrors` - встроить ошибку;
* `Wrap` - обернуть ошибку сообщением. Обеспечивается следующим форматом `message > %w`;
* `WithField` - добавить значение по ключу;
* `WithFields` - добавить значения;
* `WithFieldIf` - добавить значение по ключу, если выполняется условие.

На основе модели построен набор стандартных ошибок:

| Standard Error       | Standard Code       | GRPC Status         | HTTP Status               |
|----------------------|---------------------|---------------------|---------------------------|
| `Canceled`           | canceled            | Canceled            | StatusClientClosedRequest |
| `Unknown`            | unknown             | Unknown             | StatusInternalServerError |
| `InvalidArgument`    | invalid_argument    | InvalidArgument     | StatusBadRequest          |
| `DeadlineExceeded`   | deadline_exceeded   | DeadlineExceeded    | StatusBadGateway          |
| `NotFound`           | not_found           | NotFound            | StatusNotFound            |
| `AlreadyExists`      | already_exists      | AlreadyExists       | StatusConflict            |
| `PermissionDenied`   | permission_denied   | PermissionDenied    | StatusForbidden           |
| `ResourceExhausted`  | resource_exhausted  | ResourceExhausted   | StatusTooManyRequests     |
| `FailedPrecondition` | failed_precondition | FailedPrecondition  | StatusBadRequest          |
| `Aborted`            | aborted             | Aborted             | StatusConflict            |
| `OutOfRange`         | out_of_range        | OutOfRange          | StatusBadRequest          | 
| `Unimplemented`      | unimplemented       | Unimplemented       | StatusNotImplemented      | 
| `Internal`           | internal            | Internal            | StatusInternalServerError | 
| `Unavailable`        | unavailable         | Unavailable         | StatusServiceUnavailable  | 
| `DataLoss`           | data_loss           | DataLoss            | StatusInternalServerError |
| `Unauthenticated`    | unauthenticated     | Unauthenticated     | StatusUnauthorized        |
| `Undefined`          | -                   | Internal            | StatusInternalServerError |
| `Panic`              | panic               | Internal            | StatusInternalServerError |

Для восстановления стандартной ошибки из интерфейса `error` предлагается использовать метод `From`:
```go
// From - функция восстановления ошибки из стандартного интерфейса.
func From(err error) (*Error, bool) 
```

Метод позволяет восстановить ошибку из GRPC ответа и проверить на стандартную модель:
```go
func TestFrom(t *testing.T) {
    var err = status.Error(codes.Internal, "message")

    std, ok := stderrs.From(err)
    require.True(t, ok)
    require.True(t, std.Is(stderrs.Internal))
}
```

Если необходим кастомный метод восстановления ошибки, то есть два возможных способа встраивания: 
регистрация обработчика и имлементация интерфейса.

Пример регистрации обработчика:

```go
type ForRegistry struct {
    Code    int
    Message string
}

func (e ForRegistry) Error() string {
    return fmt.Sprintf("error with code %d, message %s", e.Code, e.Message)
}

func TestRegistry(t *testing.T) {
    stderrs.RegistryFrom(func(err error) *stderrs.Error {
        var my ForRegistry
        
        if errors.As(err, &my) {
            switch my.Code {
            case 1:
                return stderrs.Internal.SetMessage(my.Message)
            default:
                return stderrs.Unknown.SetMessage(my.Message)
            }
        }
        
        return nil
    })
    
    stderr, ok := stderrs.From(ForRegistry{Code: 1, Message: "message"})
    require.True(t, ok)
    require.True(t, stderr.Is(stderrs.Internal))
    
    stderr, ok = stderrs.From(ForRegistry{Code: 0, Message: "message"})
    require.True(t, ok)
    require.True(t, stderr.Is(stderrs.Unknown))
}
```

Пример имплементации интерфейса обработчика:

```go
type FromImpl struct {
    Code    int
    Message string
}

func (e FromImpl) Error() string {
    return fmt.Sprintf("error with code %d, message %s", e.Code, e.Message)
}

func (e FromImpl) StandardFrom(err error) *stderrs.Error {
    var my FromImpl
    
    if errors.As(err, &my) {
        switch my.Code {
        case 1:
            return stderrs.Internal.SetMessage(my.Message)
        default:
            return stderrs.Unknown.SetMessage(my.Message)
        }
    }
    
    return nil
}

func TestFromImpl(t *testing.T) {
    stderr, ok := stderrs.From(FromImpl{Code: 1, Message: "message"})
    require.True(t, ok)
    require.True(t, stderr.Is(stderrs.Internal))
    
    stderr, ok = stderrs.From(FromImpl{Code: 0, Message: "message"})
    require.True(t, ok)
    require.True(t, stderr.Is(stderrs.Unknown))
}
```

