package generatedold__test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/microavia/go-messgen/internal/generator/golang/testdata/generatedold/v/p/message"
)

func TestGenerate(t *testing.T) {
	out := message.Message6{
		Field2: [7]int32{1298498081, 2019727887, 1427131847, 939984059, 911902081, 1474941318, 140954425},
		Field3: []int16{32584, 18340, 31942, 23215, -94, 24561, 8536, -27366, 30347, 23701, 16165, -29726},
		Field4: []message.Message10{
			{
				FieldCharSlice:    []uint8{0x68, 0x92, 0x7f, 0x2b, 0x2f, 0xf8, 0x36, 0xf7, 0x35, 0x78},
				FieldFloat32Slice: []float32{0.282081, 0.7886049, 0.36180547, 0.8805431, 0.29711226, 0.89436173, 0.097454615, 0.97691685, 0.074291, 0.22228941, 0.6810783, 0.24151509, 0.31152245, 0.9328464},
				FieldFloat64Slice: []float64{0.8010550141334534, 0.7302314639091492, 0.18292491137981415, 0.4283570945262909, 0.8969919681549072, 0.6826534867286682, 0.978929340839386, 0.9222122430801392, 0.0908372774720192, 0.4931420087814331, 0.926986813545227, 0.95494544506073, 0.3479539752006531, 0.6908388137817383},
				FieldInt16Slice:   []int16{-4601, 30844, 15540, -11489, -11399, 11521, 6045, 21385, 3115},
				FieldInt32Slice:   []int32{1384138643, 183653891, 1437902002, 1337298878, 793909336, 508572546, 1149509107, 402107940, 512906503},
				FieldInt64Slice:   []int64{1169089424364679180, 2594813965004488500, 3784560248718450071, 4011359550169803385, 5765484004404056823, 5074209722772702441, 5751776211841778805, 6725505124774569258},
				FieldUint16Slice:  []uint16{0xd632, 0x7ef5, 0xaafe, 0x2470, 0x278a, 0x5f29, 0x44bf, 0x6306, 0x2928, 0xa201},
				FieldUint32Slice:  []uint32{0x479a2bf9, 0x685f3257, 0x7062b076},
				FieldUint64Slice:  []uint64{0x4cd239ea0c8dc214, 0x35ca80d72521a90, 0x6c443faf8eb3e4a1, 0x1ff5f26283efc6c6, 0x5225fcd6090ec04f, 0x1facfc5dc1540864, 0x163a5aceec2c8aaa, 0x4bdb185b70ab53ba, 0x683e14a538d3b494, 0x58cfb024878d4063, 0x3e19bf7a317ae3f, 0x4504d6353cb62f07, 0x7ce2e98ef360412c, 0x601900fb4ffbf3a9},
				FieldUint8Slice:   []uint8{0xb8, 0x32, 0x20, 0xcf, 0x58},
			},
			{
				FieldCharSlice:    []uint8{0x6c, 0xbc, 0x40},
				FieldFloat32Slice: []float32{},
				FieldFloat64Slice: []float64{},
				FieldInt16Slice:   []int16{},
				FieldInt32Slice:   []int32{},
				FieldInt64Slice:   []int64{},
				FieldUint16Slice:  []uint16{},
				FieldUint32Slice:  []uint32{},
				FieldUint64Slice:  []uint64{},
				FieldUint8Slice:   []uint8{0xb3, 0x5b, 0x74, 0xf2, 0x4b, 0x76, 0x9c, 0x8b, 0xf0},
			},
		},
		Field5: "dOCu/.SK>.$QW\\7B5TEv$nG}B6%6\".|O6*y1hSahIc.uFg$P\\&'urx%T<j7Q$*h1<}ZUAqYHkozoB8PZ1cYyFsJVBnvJJlBrjT%LQ`EP\\52#EIGHoR_tVGwK7j^hq*NVTYd[dM5[ia\"9@\\s/*}GjQU$?|3R0D=iO1Sf*LCh8f#Nf-I&{0^t\\duv%tc`^Hu-*Ra+0xix*5@9]IW2z2|;|=R#vrP#|\\qKuz_kUQ!wF8fE]:C",
		Field6: 0.7224104120134726,
	}

	b := make([]byte, 65536)
	n, err := out.Pack(b)
	require.NoError(t, err, "packing old: %v", err)

	err = os.WriteFile("../message6.old.bin", b[:n], 0644)
	require.NoError(t, err, "dumping old: %v", err)
}
