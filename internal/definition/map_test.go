package definition_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/microavia/go-messgen/internal/definition"
)

type testRowData struct {
	k  string
	f1 string
	f2 int
}

type testRow struct {
	name       string
	dataList   []testRowData
	dataMap    map[string]testRowData
	converter  func(v testRowData) (testRowData, error)
	extractKey func(v testRowData) string
	less       func(s1, s2 string) bool
	err        error
}

func TestAppend(t *testing.T) {
	testRows := []testRow{
		{
			dataList: []testRowData{
				{k: "k2", f1: "f21", f2: 22},
				{k: "k1", f1: "f11", f2: 12},
			},
			dataMap: map[string]testRowData{
				"k2": {k: "k2", f1: "f21", f2: 22},
				"k1": {k: "k1", f1: "f11", f2: 12},
			},
		},
	}

	for i, row := range testRows {
		var m map[string]testRowData
		for _, v := range row.dataList {
			var err error
			m, err = definition.MapAppend(m, v.k, v)
			require.NoError(t, err, "row %d: %+v", i, err)
		}

		require.Equal(t, row.dataMap, m, "row %d", i)
	}
}

func TestAppendError(t *testing.T) {
	testRows := []testRow{
		{
			dataList: []testRowData{
				{k: "k1", f1: "f21", f2: 22},
				{k: "k1", f1: "f11", f2: 12},
			},
			err: definition.ErrDuplicateKey,
		},
	}

	for i, row := range testRows {
		m, err := definition.MapAppend(map[string]testRowData(nil), row.dataList[0].k, row.dataList[0])
		require.NoError(t, err, "row %d: %+v", i, err)
		_, err = definition.MapAppend(m, row.dataList[1].k, row.dataList[1])
		require.ErrorIs(t, err, row.err, "row %d: %+v", i, err)
	}
}

var ErrExpected = errors.New("dataMap error")

func TestSliceToMap(t *testing.T) {
	testRows := []testRow{
		{
			name: "ok",
			dataList: []testRowData{
				{k: "k2", f1: "f21", f2: 22},
				{k: "k1", f1: "f11", f2: 12},
			},
			dataMap: map[string]testRowData{
				"k2": {k: "k2", f1: "f21", f2: 22},
				"k1": {k: "k1", f1: "f11", f2: 12},
			},
			err: nil,
		},
		{
			name: "duplicate key",
			dataList: []testRowData{
				{k: "k1", f1: "f21", f2: 22},
				{k: "k1", f1: "f11", f2: 12},
			},
			err: definition.ErrDuplicateKey,
		},
		{
			name: "convert error",
			dataList: []testRowData{
				{k: "k1", f1: "f21", f2: 22},
				{k: "k1", f1: "f11", f2: 12},
			},
			converter: func(v testRowData) (testRowData, error) { return v, ErrExpected },
			err:       ErrExpected,
		},
	}

	for i, row := range testRows {
		m, err := definition.SliceToMap(
			row.dataList,
			ternary(row.converter != nil, row.converter, func(v testRowData) (testRowData, error) { return v, nil }),
			ternary(row.extractKey != nil, row.extractKey, func(v testRowData) string { return v.k }),
		)
		if row.err == nil {
			require.NoError(t, err, "row %d: %q: %+v", i, row.name, err)
			require.Equal(t, row.dataMap, m, "row %d: %q:", i, row.name)
		} else {
			require.ErrorIs(t, err, row.err, "row %d: %q: %+v", i, row.name, err)
		}
	}
}

func TestMapToSlice(t *testing.T) {
	testRows := []testRow{
		{
			name: "ok",
			dataList: []testRowData{
				{k: "k1", f1: "f11", f2: 12},
				{k: "k2", f1: "f21", f2: 22},
			},
			dataMap: map[string]testRowData{
				"k2": {k: "k2", f1: "f21", f2: 22},
				"k1": {k: "k1", f1: "f11", f2: 12},
			},
		},
	}

	for i, row := range testRows {
		l := definition.MapToSlice(
			row.dataMap,
			ternary(row.less != nil, row.less, func(s1, s2 string) bool { return s1 < s2 }),
		)
		require.Equal(t, row.dataList, l, "row %d: %q:", i, row.name)
	}
}

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}

	return b
}
