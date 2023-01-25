package definition

import (
	"encoding/json"
	"fmt"
	"sort"
)

var ErrDuplicateKey = fmt.Errorf("duplicate key: %w", ErrBadData)

func MapAppend[K comparable, V any](m map[K]V, k K, v V) (map[K]V, error) {
	if m == nil {
		m = make(map[K]V)
	}

	if _, ok := m[k]; ok {
		return nil, fmt.Errorf("%+v: %w", k, ErrDuplicateKey)
	}

	m[k] = v

	return m, nil
}

func SliceToMap[K comparable, V any, R any](
	in []V,
	convertValue func(V) (R, error),
	extractKey func(R) K,
) (map[K]R, error) {
	out := make(map[K]R, len(in))

	for i, v := range in {
		r, err := convertValue(v)
		if err != nil {
			return out, fmt.Errorf("converting %d of %d from %T to %T: %w", i+1, len(in), v, r, err)
		}

		_, err = MapAppend(out, extractKey(r), r)
		if err != nil {
			return out, fmt.Errorf("adding %d of %d: %w", i, len(in), err)
		}
	}

	return out, nil
}

func MapToSlice[K comparable, V any](in map[K]V, less func(i, j K) bool) []V {
	out := make([]V, 0, len(in))

	for _, k := range sortKeys(in, less) {
		out = append(out, in[k])
	}

	return out
}

func sortKeys[K comparable, V any](in map[K]V, less func(i, j K) bool) []K {
	keys := make([]K, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return less(keys[i], keys[j]) })

	return keys
}

func MapUnmarshalJSON[K comparable, V any, R any](
	b []byte,
	v *map[K]R,
	convertValue func(V) (R, error),
	extractKey func(R) K,
) error {
	var values []V
	if err := json.Unmarshal(b, &values); err != nil {
		return err
	}

	res, err := SliceToMap(values, convertValue, extractKey)
	if err != nil {
		return err
	}

	*v = res

	return nil
}
