<h4> 
    <a href="../../../../README.md" align="center"> github.com/auvitly/go-tools </a> 
    > 
    <a href="../../README.md" align="center"> compare </a>
    >
    en
</h4>

## Chapters
1. [Statement of the problem](#problem)
2. [Description](#desc)

---

<a name="problem"></a>
### 1. Statement of the problem

To compare multiple objects that are not `comparable`, packages from the standard set are used:
for slices `slices`, for maps - `maps`. However, for comparing structures there is only the `reflect` package,
which has a function `reflect.DeepEqual`. The existence of the `reflect` package is usually ignored, so developers
they write their own optimized functions for checks, which leads to a decrease in the homogeneity of the code within 
the project.

<a name="desc"></a>
### 2. Description

To provide a comparison mechanism, the `compare` package offers 3 interfaces:
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

Based on these interfaces, you can implement **polymorphic** behavior, thanks to which you can compare
objects of different types among themselves.

Available comparison methods:
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