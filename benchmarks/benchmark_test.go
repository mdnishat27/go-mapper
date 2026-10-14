// Package benchmarks compares go-mapper against other ways of converting between structs.
package benchmarks

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/jinzhu/copier"
	"github.com/mdnishat27/go-mapper"
	"github.com/shopspring/decimal"
)

// Field names are identical in both structs so that reflection- and
// tag-based approaches can map them automatically.
type Source struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Amount    decimal.Decimal `json:"amount"`
	Active    bool            `json:"active"`
	CreatedAt time.Time       `json:"created_at"`
}

type Dest struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Amount    decimal.Decimal `json:"amount"`
	Active    bool            `json:"active"`
	CreatedAt time.Time       `json:"created_at"`
}

func convertManual(s Source) Dest {
	return Dest{
		ID:        s.ID,
		Name:      s.Name,
		Email:     s.Email,
		Amount:    s.Amount,
		Active:    s.Active,
		CreatedAt: s.CreatedAt,
	}
}

func convertFunc(s Source) (Dest, error) {
	return Dest{
		ID:        s.ID,
		Name:      s.Name,
		Email:     s.Email,
		Amount:    s.Amount,
		Active:    s.Active,
		CreatedAt: s.CreatedAt,
	}, nil
}

// convertReflect is a hand-rolled reflection copy matching fields by name.
func convertReflect(s Source) Dest {
	var d Dest
	sv := reflect.ValueOf(s)
	st := sv.Type()
	dv := reflect.ValueOf(&d).Elem()

	for i := 0; i < sv.NumField(); i++ {
		df := dv.FieldByName(st.Field(i).Name)
		if df.IsValid() && df.CanSet() && df.Type() == sv.Field(i).Type() {
			df.Set(sv.Field(i))
		}
	}
	return d
}

// convertCopier uses the popular github.com/jinzhu/copier library.
func convertCopier(s Source) (Dest, error) {
	var d Dest
	err := copier.Copy(&d, &s)
	return d, err
}

// convertJSON is the naive marshal/unmarshal round trip.
func convertJSON(s Source) (Dest, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return Dest{}, err
	}
	var d Dest
	err = json.Unmarshal(b, &d)
	return d, err
}

func convertMapper(s Source) (Dest, error) {
	return mapper.Map[Source, Dest](s)
}

func fixture() Source {
	return Source{
		ID:        42,
		Name:      "Alex Russo",
		Email:     "alex@example.com",
		Amount:    decimal.RequireFromString("12345.6789"),
		Active:    true,
		CreatedAt: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
	}
}

// equalDest compares field by field; decimal and time need Equal(), not ==.
func equalDest(a, b Dest) bool {
	return a.ID == b.ID &&
		a.Name == b.Name &&
		a.Email == b.Email &&
		a.Amount.Equal(b.Amount) &&
		a.Active == b.Active &&
		a.CreatedAt.Equal(b.CreatedAt)
}

// A benchmark is only fair if every approach produces the same output.
func TestConvertersAgree(t *testing.T) {
	mapper.RegisterMapper[Source, Dest](convertFunc)
	defer mapper.UnregisterMapper[Source, Dest]()

	src := fixture()
	want := convertManual(src)

	cases := map[string]func(Source) (Dest, error){
		"reflect":   func(s Source) (Dest, error) { return convertReflect(s), nil },
		"copier":    convertCopier,
		"json":      convertJSON,
		"go-mapper": convertMapper,
	}

	for name, fn := range cases {
		name, fn := name, fn
		t.Run(name, func(t *testing.T) {
			got, err := fn(src)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !equalDest(got, want) {
				t.Fatalf("mismatch\n got: %+v\nwant: %+v", got, want)
			}
		})
	}
}

var sink Dest

func BenchmarkConvert(b *testing.B) {
	src := fixture()

	b.Run("manual", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sink = convertManual(src)
		}
	})

	b.Run("go-mapper", func(b *testing.B) {
		mapper.RegisterMapper[Source, Dest](convertFunc)
		defer mapper.UnregisterMapper[Source, Dest]()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sink, _ = mapper.Map[Source, Dest](src)
		}
	})

	b.Run("reflect", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sink = convertReflect(src)
		}
	})

	b.Run("copier", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sink, _ = convertCopier(src)
		}
	})

	b.Run("json", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sink, _ = convertJSON(src)
		}
	})
}
