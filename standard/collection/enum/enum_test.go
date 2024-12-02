package enum_test

import (
	"github.com/auvitly/go-tools/enum"
	"testing"
)

func TestEnum(t *testing.T) {
	type MyType string

	var (
		MyTypeEnum enum.Enum[MyType]
	)

	var (
		MyTypeEnumValue1 = MyTypeEnum.MustRegistry("value_1")
		MyTypeEnumValue2 = MyTypeEnum.MustRegistry("value_2")
	)

	t.Logf(
		"%v, %v, %v",
		MyTypeEnum.Contains(MyTypeEnumValue1),
		MyTypeEnum.Contains(MyTypeEnumValue2),
		MyTypeEnum.Contains("value_3"),
	)
}
