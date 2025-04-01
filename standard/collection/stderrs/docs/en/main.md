<h4> 
    <a href="../../../../README.md" align="center"> github.com/auvitly/go-tools </a> 
    > 
    <a href="../../README.md" align="center"> stderrs </a>
    >
    en
</h4>

## Chapters
1. [Formulation of the problem](#problem)
2. [Description](#desc)

---

<a name="problem"></a>
### 1. Formulation of the problem

During application development, you need to handle errors.
The usual format of interaction is based on the use of standard packages `fmt` and `errors`.

Typically, after any external call you need to wrap an error:
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

The following problems can be identified:
* There are many conventions for wrapping errors, however the user must
  get a simple message that makes it easy to understand what went wrong. If you give
  user the error that the service collected during the process of multiple wrapping, then
  The information content of the message is greatly reduced.
* Determining the error code status when a large number of internal calls generates 
  the need to route error models, which becomes a separate task.
* When exceptions (panic) occur in services, the user should receive an understandable
  error about service unavailability, developer comprehensive description of where it occurred
  problems, and the service is an action to roll back an unsuccessful transaction.
> Exception handling is handled by the `recovery` package from the `go-tools` set.
> You can view `recovery` package at [link](../../../recovery/README.md).

<a name="desc"></a>
### 2. Description
The package represents a unified model (hereinafter referred to as standard error):

```go
package stderrs

import "google.golang.org/grpc/codes"

type Error struct {
	// Code - error code.
	Code    string         `json:"code"`
	// Message - the message that the user will receive.
	Message string         `json:"message"`
	// Embed - built-in errors of the standard type.
	// Note: allows you to add a standard error interface to the error.
	Embed   error          `json:"embed"`
	// Fields - allows you to pass specific parameters that determine the initial error.
	// For example, when accessing remote resources, you can pass a request model.
	// Fields is processed by the json package, so the values for the keys must
	// be regular implementation, or structures, or respond to the json.Marshaler interface.
	// Otherwise Fields will be skipped.
	Fields  map[string]any `json:"fields"`
	// Codes - status codes for various transport levels.
	Codes   struct {
		GRPC codes.Code `json:"grpc"`
		HTTP int        `json:"http"`
	} `json:"codes"`
}
```

The standard error can be initialized explicitly, or using a constructor based on the `New` function, for example:
```go
var Canceled = New("canceled").
    SetGRPCCode(codes.Canceled).
    SetHTTPCode(499)
```

Standard error has the following set of methods:
* `Is` - method for providing a standard interface. Can be used explicitly instead of `errors.Is`;
* `Unwrap` - method to provide a standard interface. Recovers all built-in errors;
* `SetCode` - set the error code, HTTP and GRPC status based on the standard error;
* `SetMessage` - set the client message;
* `SetHTTPCode` - set the HTTP status code;
* `SetGRPCCode` - set the GRPC status code;
* `EmbedErrors` - embed an error;
* `WithField` - add a value by key;
* `WithFields` - add values;
* `WithFieldIf` - add a value by key if the condition is met.

A set of standard errors is constructed based on the model:
| Standard Error       | Standard Code       | GRPC Status         | HTTP Status               | Default Message                                                 |
|----------------------|---------------------|---------------------|---------------------------|-----------------------------------------------------------------|
| `Canceled`           | canceled            | Canceled            | StatusClientClosedRequest | canceled                                                        |   
| `Unknown`            | unknown             | Unknown             | StatusInternalServerError | internal server error                                           |
| `InvalidArgument`    | invalid_argument    | InvalidArgument     | StatusBadRequest          | bad request                                                     |
| `DeadlineExceeded`   | deadline_exceeded   | DeadlineExceeded    | StatusBadGateway          | deadline exceeded                                               |
| `NotFound`           | not_found           | NotFound            | StatusNotFound            | not found                                                       |
| `AlreadyExists`      | already_exists      | AlreadyExists       | StatusConflict            | already exists                                                  |
| `PermissionDenied`   | permission_denied   | PermissionDenied    | StatusForbidden           | permission denied                                               |
| `ResourceExhausted`  | resource_exhausted  | ResourceExhausted   | StatusTooManyRequests     | resource has been exhausted                                     |
| `FailedPrecondition` | failed_precondition | FailedPrecondition  | StatusBadRequest          | system is not in a state required for the operation's execution |
| `Aborted`            | aborted             | Aborted             | StatusConflict            | aborted                                                         |
| `OutOfRange`         | out_of_range        | OutOfRange          | StatusBadRequest          | attempted past the valid range                                  |
| `Unimplemented`      | unimplemented       | Unimplemented       | StatusNotImplemented      | not implemented or not supported/enabled                        |
| `Internal`           | internal            | Internal            | StatusInternalServerError | internal server error                                           |
| `Unavailable`        | unavailable         | Unavailable         | StatusServiceUnavailable  | service unavailable                                             |
| `DataLoss`           | data_loss           | DataLoss            | StatusInternalServerError | unrecoverable data loss or corruption                           |
| `Unauthenticated`    | unauthenticated     | Unauthenticated     | StatusUnauthorized        | request does not have valid authentication credentials          |
| `Undefined`          | -                   | Internal            | StatusInternalServerError | internal server error                                           |
| `Panic`              | panic               | Internal            | StatusInternalServerError | internal server error                                           |
To restore the standard error from the `error` interface, it is proposed to use the `From` method:
```go
// From - error recovery function from the standard interface.
func From(err error) (*Error, bool) 
```

The method allows you to recover the error from the GRPC response and check it against the standard model:
```go
func TestFrom(t *testing.T) {
    var err = status.Error(codes.Internal, "message")

    std, ok := stderrs.From(err)
    require.True(t, ok)
    require.True(t, std.Is(stderrs.Internal))
}
```

If a custom error recovery method is needed, then there are two possible embedding methods:
registering a handler and implementing an interface.

Example of registering a handler:

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

Example of implementation of the handler interface:

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

