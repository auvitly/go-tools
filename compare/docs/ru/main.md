<h4> 
    <a href="./../../../README.md" align="center"> github.com/auvitly/go-tools </a> 
    > 
    <a href="./../../README.md" align="center"> compare </a>
    >
    ru
</h4>


## Оглавление
1. [Постановка проблемы](#problem)
2. [Описание](#desc)

---

<a name="problem"></a>
### 1. Постановка проблемы

Для сравнения нескольких объектов, которые не являются `comparable` используются пакеты из стандартного набора:
для срезов `slices`, для мап - `maps`. Однако для сравнения структур существует только пакет `reflect`, 
в котором имеется функция `reflect.DeepEqual`. Зачастую задача решается по разному, что уменьшает 
однородность кода внутри проекта.

<a name="desc"></a>
### 2. Описание

Для обеспечения механизма сравнения пакет `compare` предлагает 3 интерфейса:
```go
type ComparableEqual interface {
	Equal(c ComparableEqual) bool
}

type ComparableLess interface {
	LessThan(c ComparableLess) bool
}

type ComparableGreater interface {
	GreaterThan(c ComparableGreater) bool
}
```

На основе данных интерфейсов можно реализовать **полиморфное** поведение, благодаря которому можно сравнивать
объекты разных типов между собой. 

Доступные методы сравнения:
```go
func Equal[T ComparableEqual](a, b T) bool

func NotEqual[T ComparableEqual](a, b T) bool 

func Greater[T ComparableGreater](a, b T) bool

func NotGreater[T ComparableGreater](a, b T) bool

func Less[T ComparableLess](a, b T) bool

func NotLess[T ComparableLess](a, b T) bool

func GreaterOrEqual[T interface {ComparableGreater; ComparableEqual}](a, b T) bool

func LessOrEqual[T interface {ComparableLess; ComparableEqual}](a, b T) bool
```