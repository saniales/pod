package util_test

import (
	"math"
	"testing"

	. "github.com/stalker-loki/pod/pkg/util"
)

func TestAmountCreation(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		valid    bool
		expected Amount
	}{
		// Positive tests.
		{
			name:     "zero",
			amount:   0,
			valid:    true,
			expected: 0,
		},
		{
			name:     "max producible",
			amount:   21e6,
			valid:    true,
			expected: MaxSatoshi,
		},
		{
			name:     "min producible",
			amount:   -21e6,
			valid:    true,
			expected: -MaxSatoshi,
		},
		{
			name:     "exceeds max producible",
			amount:   21e6 + 1e-8,
			valid:    true,
			expected: MaxSatoshi + 1,
		},
		{
			name:     "exceeds min producible",
			amount:   -21e6 - 1e-8,
			valid:    true,
			expected: -MaxSatoshi - 1,
		},
		{
			name:     "one hundred",
			amount:   100,
			valid:    true,
			expected: 100 * SatoshiPerBitcoin,
		},
		{
			name:     "fraction",
			amount:   0.01234567,
			valid:    true,
			expected: 1234567,
		},
		{
			name:     "rounding up",
			amount:   54.999999999999943157,
			valid:    true,
			expected: 55 * SatoshiPerBitcoin,
		},
		{
			name:     "rounding down",
			amount:   55.000000000000056843,
			valid:    true,
			expected: 55 * SatoshiPerBitcoin,
		},
		// Negative tests.
		{
			name:   "not-a-number",
			amount: math.NaN(),
			valid:  false,
		},
		{
			name:   "-infinity",
			amount: math.Inf(-1),
			valid:  false,
		},
		{
			name:   "+infinity",
			amount: math.Inf(1),
			valid:  false,
		},
	}
	for _, test := range tests {
		a, err := NewAmount(test.amount)
		switch {
		case test.valid && err != nil:
			t.Errorf("%v: Positive test Amount creation failed with: %v", test.name, err)
			continue
		case !test.valid && err == nil:
			t.Errorf("%v: Negative test Amount creation succeeded (value %v) when should fail", test.name, a)
			continue
		}
		if a != test.expected {
			t.Errorf("%v: Created amount %v does not match expected %v", test.name, a, test.expected)
			continue
		}
	}
}
func TestAmountUnitConversions(t *testing.T) {
	tests := []struct {
		name      string
		amount    Amount
		unit      AmountUnit
		converted float64
		s         string
	}{
		{
			name:      "MDUO",
			amount:    MaxSatoshi,
			unit:      AmountMegaDUO,
			converted: 21,
			s:         "21 MDUO",
		},
		{
			name:      "kDUO",
			amount:    44433322211100,
			unit:      AmountKiloDUO,
			converted: 444.33322211100,
			s:         "444.333222111 kDUO",
		},
		{
			name:      "DUO",
			amount:    44433322211100,
			unit:      AmountDUO,
			converted: 444333.22211100,
			s:         "444333.222111 DUO",
		},
		{
			name:      "mDUO",
			amount:    44433322211100,
			unit:      AmountMilliDUO,
			converted: 444333222.11100,
			s:         "444333222.111 mDUO",
		},
		{
			name:      "μDUO",
			amount:    44433322211100,
			unit:      AmountMicroDUO,
			converted: 444333222111.00,
			s:         "444333222111 μDUO",
		},
		{
			name:      "satoshi",
			amount:    44433322211100,
			unit:      AmountSatoshi,
			converted: 44433322211100,
			s:         "44433322211100 Satoshi",
		},
		{
			name:      "non-standard unit",
			amount:    44433322211100,
			unit:      AmountUnit(-1),
			converted: 4443332.2211100,
			s:         "4443332.22111 1e-1 DUO",
		},
	}
	for _, test := range tests {
		f := test.amount.ToUnit(test.unit)
		if f != test.converted {
			t.Errorf("%v: converted value %v does not match expected %v", test.name, f, test.converted)
			continue
		}
		s := test.amount.Format(test.unit)
		if s != test.s {
			t.Errorf("%v: format '%v' does not match expected '%v'", test.name, s, test.s)
			continue
		}
		// Verify that Amount.ToDUO works as advertised.
		f1 := test.amount.ToUnit(AmountDUO)
		f2 := test.amount.ToDUO()
		if f1 != f2 {
			t.Errorf("%v: ToDUO does not match ToUnit(AmountDUO): %v != %v", test.name, f1, f2)
		}
		// Verify that Amount.String works as advertised.
		s1 := test.amount.Format(AmountDUO)
		s2 := test.amount.String()
		if s1 != s2 {
			t.Errorf("%v: String does not match Format(AmountBitcoin): %v != %v", test.name, s1, s2)
		}
	}
}
func TestAmountMulF64(t *testing.T) {
	tests := []struct {
		name string
		amt  Amount
		mul  float64
		res  Amount
	}{
		{
			name: "Multiply 0.1 DUO by 2",
			amt:  100e5, // 0.1 DUO
			mul:  2,
			res:  200e5, // 0.2 DUO
		},
		{
			name: "Multiply 0.2 DUO by 0.02",
			amt:  200e5, // 0.2 DUO
			mul:  1.02,
			res:  204e5, // 0.204 DUO
		},
		{
			name: "Multiply 0.1 DUO by -2",
			amt:  100e5, // 0.1 DUO
			mul:  -2,
			res:  -200e5, // -0.2 DUO
		},
		{
			name: "Multiply 0.2 DUO by -0.02",
			amt:  200e5, // 0.2 DUO
			mul:  -1.02,
			res:  -204e5, // -0.204 DUO
		},
		{
			name: "Multiply -0.1 DUO by 2",
			amt:  -100e5, // -0.1 DUO
			mul:  2,
			res:  -200e5, // -0.2 DUO
		},
		{
			name: "Multiply -0.2 DUO by 0.02",
			amt:  -200e5, // -0.2 DUO
			mul:  1.02,
			res:  -204e5, // -0.204 DUO
		},
		{
			name: "Multiply -0.1 DUO by -2",
			amt:  -100e5, // -0.1 DUO
			mul:  -2,
			res:  200e5, // 0.2 DUO
		},
		{
			name: "Multiply -0.2 DUO by -0.02",
			amt:  -200e5, // -0.2 DUO
			mul:  -1.02,
			res:  204e5, // 0.204 DUO
		},
		{
			name: "Round down",
			amt:  49, // 49 Satoshis
			mul:  0.01,
			res:  0,
		},
		{
			name: "Round up",
			amt:  50, // 50 Satoshis
			mul:  0.01,
			res:  1, // 1 Satoshi
		},
		{
			name: "Multiply by 0.",
			amt:  1e8, // 1 DUO
			mul:  0,
			res:  0, // 0 DUO
		},
		{
			name: "Multiply 1 by 0.5.",
			amt:  1, // 1 Satoshi
			mul:  0.5,
			res:  1, // 1 Satoshi
		},
		{
			name: "Multiply 100 by 66%.",
			amt:  100, // 100 Satoshis
			mul:  0.66,
			res:  66, // 66 Satoshis
		},
		{
			name: "Multiply 100 by 66.6%.",
			amt:  100, // 100 Satoshis
			mul:  0.666,
			res:  67, // 67 Satoshis
		},
		{
			name: "Multiply 100 by 2/3.",
			amt:  100, // 100 Satoshis
			mul:  2.0 / 3,
			res:  67, // 67 Satoshis
		},
	}
	for _, test := range tests {
		a := test.amt.MulF64(test.mul)
		if a != test.res {
			t.Errorf("%v: expected %v got %v", test.name, test.res, a)
		}
	}
}
