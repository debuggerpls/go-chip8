package chip8

import (
	"fmt"
	"testing"
)

type MockDisplay struct {
	clear, drawbytes int
	x, y             byte
}

func (d *MockDisplay) Create() error {
	return nil
}

func (d *MockDisplay) Destroy() {
}

func (d *MockDisplay) Clear() {
	d.clear += 1
}

func (d *MockDisplay) Update() {
}

func (d *MockDisplay) Draw(x, y byte, sprite []byte) (collision byte) {
	d.x = x
	d.y = y
	d.drawbytes = len(sprite)
	collision = 0
	return
}

func (d *MockDisplay) String() string {
	return fmt.Sprintf("x=%d y=%d drawbytes=%d", d.x, d.y, d.drawbytes)
}

// Test opcode getter functions
func TestOpFuncs(t *testing.T) {
	testData := []struct {
		op, expected uint16
		fn           func(uint16) uint16
		fnname       string
	}{
		{0x7abc, 7, OpNr, "OpNr"},
		{0x7abc, 0xabc, OpNNN, "OpNNN"},
		{0x7abc, 0xc, OpN, "OpN"},
		{0x7abc, 0xa, OpX, "OpX"},
		{0x7abc, 0xb, OpY, "OpY"},
		{0x7abc, 0xbc, OpKK, "OpKK"},
	}

	for _, data := range testData {
		actual := data.fn(data.op)
		if actual != data.expected {
			t.Errorf("%s: op=%04x expected=%04x != actual=%04x", data.fnname, data.op, data.expected, actual)
		}
	}
}

func TestOpNr0(t *testing.T) {
	r := Registers{}
	m := Memory{}
	d := &MockDisplay{}
	var opcode uint16 = 0x0000

	if err := OpNr0(opcode, &r, &m, d); err != nil {
		t.Error(err)
	}

	opcode = 0x00e0
	if err := OpNr0(opcode, &r, &m, d); err != nil {
		t.Error(err)
	}
	if d.clear != 1 {
		t.Errorf("Display was not cleared, clear=%d", d.clear)
	}

	opcode = 0x00ee
	r.SP = 1
	r.Stack[0] = 0xabcd
	if err := OpNr0(opcode, &r, &m, d); err != nil {
		t.Error(err)
	}
	if r.SP != 0 {
		t.Errorf("Wrong SP, expected=%02x\n%s", 0, r.String())
	}
	if r.PC != 0xabcd {
		t.Errorf("Wrong PC, expected=%04x\n%s", 0xabcd, r.String())
	}
}

func TestOpNr1(t *testing.T) {
	r := Registers{}
	m := Memory{}
	var opcode uint16 = 0x1abc

	if err := OpNr1(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.PC != 0xabc {
		t.Errorf("Wrong PC, expected=%04x\n%s", 0xabc, r.String())
	}
}

func TestOpNr2(t *testing.T) {
	r := Registers{}
	m := Memory{}
	var opcode uint16 = 0x2abc

	if err := OpNr2(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.SP != 1 {
		t.Errorf("Wrong SP, expected=%02x\n%s", 1, r.String())
	}
	if r.Stack[0] != 0 {
		t.Errorf("Wrong Stack[0], expected=%02x\n%s", 0, r.String())
	}
	if r.PC != 0xabc {
		t.Errorf("Wrong PC, expected=%04x\n%s", 0xabc, r.String())
	}
}

func TestOpNr3(t *testing.T) {
	r := Registers{}
	m := Memory{}
	var opcode uint16 = 0x3abc

	if err := OpNr3(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.PC != 0 {
		t.Errorf("Wrong PC, expected=%04x\n%s", 0, r.String())
	}

	r.V[0xa] = 0xbc
	if err := OpNr3(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.PC != 2 {
		t.Errorf("Wrong PC, expected=%04x\n%s", 2, r.String())
	}
}

func TestOpNr4(t *testing.T) {
	r := Registers{}
	m := Memory{}
	var opcode uint16 = 0x4abc

	r.V[0xa] = 0xbc
	if err := OpNr4(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.PC != 0 {
		t.Errorf("Wrong PC, expected=%04x\n%s", 0, r.String())
	}

	r.V[0xa] = 0
	if err := OpNr4(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.PC != 2 {
		t.Errorf("Wrong PC, expected=%04x\n%s", 2, r.String())
	}
}

func TestOpNr5(t *testing.T) {
	r := Registers{}
	m := Memory{}
	var opcode uint16 = 0x5010

	r.V[0] = 1
	if err := OpNr5(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.PC != 0 {
		t.Errorf("Wrong PC, expected=%04x\n%s", 0, r.String())
	}

	r.V[0] = 0
	if err := OpNr5(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.PC != 2 {
		t.Errorf("Wrong PC, expected=%04x\n%s", 2, r.String())
	}
}

func TestOpNr6(t *testing.T) {
	r := Registers{}
	m := Memory{}
	var opcode uint16 = 0x6010

	if err := OpNr6(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != 0x10 {
		t.Errorf("Wrong V0, expected=%04x\n%s", 0x10, r.String())
	}
}

func TestOpNr7(t *testing.T) {
	r := Registers{}
	m := Memory{}
	var opcode uint16 = 0x7010

	if err := OpNr7(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != 0x10 {
		t.Errorf("Wrong V0, expected=%04x\n%s", 0x10, r.String())
	}
	if err := OpNr7(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != 0x20 {
		t.Errorf("Wrong V0, expected=%04x\n%s", 0x20, r.String())
	}
}

func TestOpNr8(t *testing.T) {
	r := Registers{}
	m := Memory{}

	var opcode uint16 = 0x8010
	r.V[1] = 0x1
	if err := OpNr8(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != 1 {
		t.Errorf("Wrong V0, expected=%04x\n%s", 1, r.String())
	}

	opcode = 0x8011
	r.V[0] = 0
	r.V[1] = 0xf0
	expected := r.V[0] | r.V[1]
	if err := OpNr8(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != expected {
		t.Errorf("Wrong V0, expected=%04x\n%s", expected, r.String())
	}

	opcode = 0x8012
	r.V[0] = 0
	r.V[1] = 0xf0
	expected = r.V[0] & r.V[1]
	if err := OpNr8(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != expected {
		t.Errorf("Wrong V0, expected=%04x\n%s", expected, r.String())
	}

	opcode = 0x8013
	r.V[0] = 0
	r.V[1] = 0xf0
	expected = r.V[0] ^ r.V[1]
	if err := OpNr8(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != expected {
		t.Errorf("Wrong V0, expected=%04x\n%s", expected, r.String())
	}

	opcode = 0x8014
	r.V[0] = 0
	r.V[1] = 0xf0
	expected = r.V[0] + r.V[1]
	if err := OpNr8(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != expected {
		t.Errorf("Wrong V0, expected=%04x\n%s", expected, r.String())
	}
	if r.V[0xf] != 0 {
		t.Errorf("Wrong Vf, expected=%04x\n%s", 0, r.String())
	}

	r.V[0] = 0xf0
	r.V[1] = 0xf0
	expected = r.V[0] + r.V[1]
	if err := OpNr8(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != expected {
		t.Errorf("Wrong V0, expected=%04x\n%s", expected, r.String())
	}
	if r.V[0xf] != 1 {
		t.Errorf("Wrong Vf, expected=%04x\n%s", 1, r.String())
	}

	opcode = 0x8015
	r.V[0] = 0
	r.V[1] = 0xf0
	expected = r.V[0] - r.V[1]
	if err := OpNr8(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != expected {
		t.Errorf("Wrong V0, expected=%04x\n%s", expected, r.String())
	}
	if r.V[0xf] != 0 {
		t.Errorf("Wrong Vf, expected=%04x\n%s", 0, r.String())
	}

	opcode = 0x8016
	r.V[0] = 11
	r.V[1] = 0x0
	expected = r.V[0] / 2
	if err := OpNr8(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != expected {
		t.Errorf("Wrong V0, expected=%04x\n%s", expected, r.String())
	}
	if r.V[0xf] != 1 {
		t.Errorf("Wrong Vf, expected=%04x\n%s", 1, r.String())
	}

	opcode = 0x8017
	r.V[0] = 11
	r.V[1] = 0x0
	expected = r.V[1] - r.V[0]
	if err := OpNr8(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != expected {
		t.Errorf("Wrong V0, expected=%04x\n%s", expected, r.String())
	}
	if r.V[0xf] != 0 {
		t.Errorf("Wrong Vf, expected=%04x\n%s", 0, r.String())
	}

	opcode = 0x801e
	r.V[0] = 0xf0
	r.V[1] = 0x0
	expected = r.V[0] * 2
	if err := OpNr8(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != expected {
		t.Errorf("Wrong V0, expected=%04x\n%s", expected, r.String())
	}
	if r.V[0xf] != 1 {
		t.Errorf("Wrong Vf, expected=%04x\n%s", 1, r.String())
	}
}

func TestOpNr9(t *testing.T) {
	r := Registers{}
	m := Memory{}
	var opcode uint16 = 0x9010

	r.V[1] = 1
	if err := OpNr9(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.PC != 2 {
		t.Errorf("Wrong PC, expected=%04x\n%s", 2, r.String())
	}
}

func TestOpNrA(t *testing.T) {
	r := Registers{}
	m := Memory{}
	var opcode uint16 = 0xa010

	expected := uint16(0x10)
	if err := OpNrA(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.I != expected {
		t.Errorf("Wrong I, expected=%04x\n%s", expected, r.String())
	}
}

func TestOpNrB(t *testing.T) {
	r := Registers{}
	m := Memory{}

	var opcode uint16 = 0xb010
	r.V[0] = 0xab
	expected := uint16(0x10) + uint16(r.V[0])
	if err := OpNrB(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.PC != expected {
		t.Errorf("Wrong PC, expected=%04x\n%s", expected, r.String())
	}
}

func TestOpNrD(t *testing.T) {
	r := Registers{}
	m := Memory{}
	d := &MockDisplay{}

	var opcode uint16 = 0xd015
	r.V[0] = 10
	r.V[1] = 15
	if err := OpNrD(opcode, &r, &m, d); err != nil {
		t.Error(err)
	}
	if d.x != 10 || d.y != 15 || d.drawbytes != 5 {
		t.Errorf("Wrong Display state: %s", d.String())
	}
}

func TestOpNrF(t *testing.T) {
	r := Registers{}
	m := Memory{}

	var opcode uint16 = 0xf007
	r.DT = 0xa
	expected := r.DT
	if err := OpNrF(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != expected {
		t.Errorf("Wrong V0, expected=%04x\n%s", expected, r.String())
	}

	opcode = 0xf015
	r.V[0] = 0xa
	r.DT = 0
	expected = r.V[0]
	if err := OpNrF(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.DT != expected {
		t.Errorf("Wrong DT, expected=%04x\n%s", expected, r.String())
	}

	opcode = 0xf018
	r.V[0] = 0xa
	r.ST = 0
	expected = r.V[0]
	if err := OpNrF(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.ST != expected {
		t.Errorf("Wrong S, expected=%04x\n%s", expected, r.String())
	}

	opcode = 0xf01e
	r.V[0] = 0xa
	r.I = 0
	expected = r.V[0]
	if err := OpNrF(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.I != uint16(expected) {
		t.Errorf("Wrong I, expected=%04x\n%s", expected, r.String())
	}

	opcode = 0xf033
	r.V[0] = 234
	r.I = 1
	if err := OpNrF(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if m[1] != 2 && m[2] != 3 && m[3] != 4 {
		t.Errorf("Wrong Memory\n")
	}

	opcode = 0xf055
	r.V[0] = 0xab
	r.V[1] = 0xcd
	r.V[2] = 0xef
	r.I = 1
	if err := OpNrF(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if m[1] != 0xab && m[2] != 0xcd && m[3] != 0xef {
		t.Errorf("Wrong Memory\n")
	}

	opcode = 0xf065
	r.V[0] = 0
	r.V[1] = 0
	r.V[2] = 0
	r.I = 1
	if err := OpNrF(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.V[0] != 0xab && r.V[1] != 0xcd && r.V[2] != 0xef {
		t.Errorf("Wrong V[]\n%s", r.String())
	}

	opcode = 0xf029
	r.V[0] = 5
	r.I = 0
	if err := OpNrF(opcode, &r, &m); err != nil {
		t.Error(err)
	}
	if r.I != 25 {
		t.Errorf("Wrong I, expected=%04x\n%s", 25, r.String())
	}

}
