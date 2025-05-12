package calc_test

import (
	"project_yandex_lms/lms-calc-with-gorutine/calc"
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    float64
		wantErr bool
	}{

		{"simple addition", "1+2", 3, false},
		{"simple subtraction", "5-2", 3, false},
		{"simple multiplication", "3*4", 12, false},
		{"simple division", "8/2", 4, false},
		{"division with fraction", "5/2", 2.5, false},

		{"multiple operations", "2+3*4", 14, false},
		{"operations with same precedence", "6-3+2", 5, false},
		{"complex expression", "(2+3)*4", 20, false},
		{"nested parentheses", "((2+3)*4)/5", 4, false},
		{"floating point numbers", "3.5*2", 7, false},

		{"with spaces", " 2 + 3 * 4 ", 14, false},

		{"negative numbers", "-2+3", 1, false},
		{"negative in parentheses", "(-2+3)*4", 4, false},

		{"empty expression", "", 0, true},
		{"invalid characters", "2+a", 0, true},
		{"unmatched parentheses", "(2+3", 0, true},
		{"division by zero", "3/0", 0, true},
		{"invalid decimal", "2.3.4+5", 0, true},
		{"only operator", "+", 0, true},
		{"trailing operator", "2+", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calc.Calc(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Calc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	edgeTests := []struct {
		name    string
		expr    string
		want    float64
		wantErr bool
	}{
		{"single number", "42", 42, false},
		{"zero division in expression", "1/(5-5)", 0, true},
		{"multiple parentheses", "((((2))))", 2, false},
		{"large numbers", "1000000*1000000", 1000000000000, false},
		{"very small numbers", "0.000001*0.000001", 0.000000000001, false},
		{"expression with all operators", "1+2*3-4/2", 5, false},
	}

	for _, tt := range edgeTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calc.Calc(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Calc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		expr string
		want float64
	}{
		{"2+3*4", 14},
		{"2*3+4", 10},
		{"2+3+4*5", 25},
		{"2*3+4*5", 26},
		{"2+3*4+5", 19},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			got, err := calc.Calc(tt.expr)
			if err != nil {
				t.Errorf("Calc() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("Calc() = %v, want %v", got, tt.want)
			}
		})
	}
}
