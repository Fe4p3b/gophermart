package luhn

import (
	"testing"
)

func TestLuhn(t *testing.T) {

	tests := []struct {
		name string
		m    map[string]bool
	}{
		{
			name: "Test case #1",
			m: map[string]bool{
				"1607325485528125": true,
				"0017166005332606": true,
				"4784887705382733": false,
				"6668455704647276": true,
				"1417034260725004": false,
				"4406474523212467": false,
				"7234803454266073": true,
				"0735370377831253": false,
				"7811677408518537": true,
				"3170207271867508": true,
				"0113444231702345": true,
				"1228478465810040": false,
				"6020123787536559": true,
				"5310627162416766": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.m {
				if got := Luhn([]byte(k)); got != v {
					t.Errorf("Luhn(%s) = %v, want %v", k, got, v)
				}
			}
		})
	}
}

func TestOnlyDigits(t *testing.T) {

	tests := []struct {
		name string
		m    map[string]bool
	}{
		{
			name: "Test case #1",
			m: map[string]bool{
				"0704565275704246a": false,
				"4058201207054138":  true,
				"4351207430727113a": false,
				"7557646002372151a": false,
				"qwertyuip":         false,
				"asdfsdf":           false,
				"zxcvbn":            false,
				"1234567890":        true,
				"7414646441113151a": false,
				"5657025220724555a": false,
				"4052105231662780a": false,
				"7814724758642603":  true,
				"6167632328628678a": false,
				"7265586422565765a": false,
				"0536464185687424":  true,
				"7755623860484303a": false,
				"0261227722126359":  true,
				"7436535524361247a": false,
				"4022277650033061a": false,
				"5681845882205860":  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.m {
				if got := OnlyDigits([]byte(k)); got != v {
					t.Errorf("OnlyDigits() = %v, want %v", got, v)
				}
			}
		})
	}
}
