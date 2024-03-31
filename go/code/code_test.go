package code

import "testing"

func Test_Make(t *testing.T) {
	testData := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
	}

	for _, td := range testData {
		instruction := Make(td.op, td.operands...)
		if len(instruction) != len(td.expected) {
			t.Errorf("instruction has wrong length.\n\twant=%d\n\tgot=%d", len(td.expected), len(instruction))
		}

		for i, b := range td.expected {
			if instruction[i] != td.expected[i] {
				t.Errorf("wrong byte at pos %d.\n\twant=%d\n\tgot=%d", i, b, instruction[i])
			}
		}
	}
}

func Test_InstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}

	expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`
	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\n\twant=%q\n\tgot=%q", expected, concatted.String())
	}
}

func Test_ReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandRead, n := ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong\n\twant=%d\n\tgot=%d", tt.bytesRead, n)
		}

		for i, want := range tt.operands {
			if operandRead[i] != want {
				t.Errorf("operand wrong\n\twant=%d\n\tgot=%d", want, operandRead[i])
			}
		}
	}
}
