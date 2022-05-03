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
				"3441163888648125": false,
				"3067776772306568": false,
				"1804275416778862": false,
				"1371370161448606": true,
				"1004134757644779": false,
				"4671600734854570": false,
				"3012567225501053": false,
				"7123825227302814": false,
				"7144508437244852": true,
				"3716784018533255": false,
				"0250148785553268": false,
				"6268836324723122": false,
				"0585638804102119": true,
				"2283854773334275": false,
				"5577677268753446": false,
				"8118645888857103": false,
				"4207840610016737": false,
				"4474217331260538": false,
				"1267747462531452": false,
				"4832244585810039": false,
				"1560302845565449": false,
				"8367581205307712": false,
				"4764208667301532": false,
				"1050671374026263": false,
				"7363672347225700": false,
				"2354084557523752": false,
				"1607325485528125": true,
				"0017166005332606": false,
				"4784887705282733": false,
				"6668455704647276": true,
				"1417034260724004": false,
				"4406474523252467": false,
				"7234803454266073": false,
				"0735370377831853": false,
				"7811677408518537": false,
				"3170207271867508": true,
				"0113444231702345": false,
				"1228478465850040": false,
				"6020123787536559": false,
				"5310627162416766": false,
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
