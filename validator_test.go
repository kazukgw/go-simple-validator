package validator

import (
	"testing"
	"time"
)

func TestEmpty(t *testing.T) {
	trueargs := []interface{}{
		false,
		0,
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		0.0,
		float32(0),
		float64(0),
		"",
		time.Time{},
		&time.Time{},
		struct{}{},
		struct {
			Str string
			Num int
		}{"", 0},
		&struct {
			Str string
			Num int
		}{"", 0},
		[]string{},
		map[string]string{},
	}

	for _, arg := range trueargs {
		if !Empty(arg) {
			t.Errorf("%v is not empty!", arg)
		}
	}

	falseargs := []interface{}{
		true,
		1,
		int64(1),
		uint64(1),
		1.1,
		float64(1),
		time.Now(),
		struct {
			Str string
			Num int
		}{"hoge", 0},
		[]int{1},
		map[string]string{"foo": "bar"},
	}

	for _, arg := range falseargs {
		if Empty(arg) {
			t.Errorf("%v is not NotEmpty!", arg)
		}
	}
}

func TestRange(t *testing.T) {
	if !Range(1, 0, 1) {
		t.Errorf("Range(1,0,1) must true")
	}

	if !Range(50, 0, 100) {
		t.Errorf("Range(50,0,100) must true")
	}

	if Range(30, 50, 100) {
		t.Errorf("Range(30,50,100) must false")
	}
}

func TestStringSize(t *testing.T) {
	if !StringSize("hoge", 4, 5) {
		t.Errorf("StringSize(\"hoge\", 4, 5) must true")
	}

	if !StringSize("ほげ", 1, 3) {
		t.Errorf("StringSize(\"ほげ\", 1, 3) must true")
	}

	if StringSize("ほげふが", 1, 3) {
		t.Errorf("StringSize(\"ほげふが\", 1, 3) must false")
	}

	if !StringSize("ほげふがabc", 1, 7) {
		t.Errorf("StringSize(\"ほげふがabc\", 1, 7) must true")
	}
}

func TestRegexp(t *testing.T) {
	if !Regexp("090-1234-1234", "[0-9]{3}-[0-9]{4}-[0-9]{4}") {
		t.Errorf("Regexp(\"090-1234-1234\", \"[0-9]{3}-[0-9]{4}-[0-9]{4}\") must true")
	}

	if Regexp("090-1234-1234", "[0-9]{3}-[0-9]{2}-[0-9]{3}") {
		t.Errorf("Regexp(\"090-1234-1234\", \"[0-9]{3}-[0-9]{2}-[0-9]{3}\") must false")
	}

	if !Regexp("foo bar hoge fuga", "^foo .+ fuga$") {
		t.Errorf("Regexp(\"foo bar hoge fuga\", \"^foo .+ fuga$\") must true")
	}

	if Regexp("foo bar hoge fuga", "^foo .+ huga$") {
		t.Errorf("Regexp(\"foo bar hoge fuga\", \"^foo .+ huga$\") must false")
	}
}

func TestEqual(t *testing.T) {
	if !Equal(nil, nil) {
		t.Errorf("nil should equal to nil")
	}
	if Equal(nil, false) {
		t.Errorf("nil should not equal to 0")
	}
}

func TestContain(t *testing.T) {
	if !Contain("hoge", []string{"foo", "bar", "hoge", "fuga"}) {
		t.Errorf("Contain(\"hoge\", []string{\"foo\",\"bar\",\"hoge\",\"fuga\"}) must true")
	}

	if !Contain(140, []int{10, 543, 140, 12}) {
		t.Errorf("Contain(140, []string{10, 543, 140, 12}) must true")
	}

	type Person struct {
		Name string
		Age  int
	}
	taro := Person{"Taro", 33}
	hanako := Person{"Hanako", 28}
	john := Person{"John", 47}
	mary := Person{"Mary", 52}
	if !Contain(taro, []Person{taro, hanako, john, mary}) {
		t.Errorf("Contain func must use with struct")
	}
	if !Contain(&taro, []*Person{&taro, &hanako, &john, &mary}) {
		t.Errorf("Contain func must use with struct")
	}
	if Contain(taro, []Person{hanako, hanako, john, mary}) {
		t.Errorf("Contain func must use with struct")
	}
}

func TestTimeRange(t *testing.T) {
	t1 := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	t2 := time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)
	t3 := time.Date(2011, time.November, 10, 23, 0, 0, 0, time.UTC)

	if !TimeRange(t2, t1, t3) {
		t.Errorf("TiemRange func does not work")
	}

	if !TimeRange(t1, t1, t3) {
		t.Errorf("TiemRange func does not work")
	}

	if TimeRange(t1, t2, t3) {
		t.Errorf("TiemRange func does not work")
	}
}
