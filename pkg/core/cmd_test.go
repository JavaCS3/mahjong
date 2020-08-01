package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmd(t *testing.T) {
	c := NewCmd("test-cmd")
	c.SetProp("k1", "v1")
	c.SetProp("k2", "v2")
	c.SetProp("k3", "v3")
	c.SetMsg("test-msg")

	c.DelProp("k3")

	assert.Equal(t, c.GetProp("not-existing"), nil)
	assert.Equal(t, c.GetProp("k1"), "v1")
	assert.Equal(t, c.GetProp("k2"), "v2")
	assert.Equal(t, c.GetMsg(), "test-msg")
	assert.Equal(t, c.Name(), "test-cmd")
	assert.Equal(t, fmt.Sprintf("%s", c), "::test-cmd k1=v1,k2=v2::test-msg")
}

func TestCmdOnly(t *testing.T) {
	c := NewCmd("some-command")

	assert.Equal(t, fmt.Sprintf("%s", c), "::some-command::")
}

func TestCmdEscapesMsg(t *testing.T) {
	c := NewCmd("some-command")
	c.SetMsg("percent % percent % cr \r cr \r lf \n lf \n")

	assert.Equal(t, fmt.Sprintf("%s", c), "::some-command::percent %25 percent %25 cr %0D cr %0D lf %0A lf %0A")

	c.SetMsg("%25 %25 %0D %0D %0A %0A")

	assert.Equal(t, fmt.Sprintf("%s", c), "::some-command::%2525 %2525 %250D %250D %250A %250A")
}

func TestCmdEscapesProps(t *testing.T) {
	c := NewCmd("some-command")
	c.SetProp("name", "percent % percent % cr \r cr \r lf \n lf \n colon : colon : comma , comma ,")

	assert.Equal(t, fmt.Sprintf("%s", c), "::some-command name=percent %25 percent %25 cr %0D cr %0D lf %0A lf %0A colon %3A colon %3A comma %2C comma %2C::")

	c.DelProp("name")
	c.SetMsg("%25 %25 %0D %0D %0A %0A %3A %3A %2C %2C")

	assert.Equal(t, fmt.Sprintf("%s", c), "::some-command::%2525 %2525 %250D %250D %250A %250A %253A %253A %252C %252C")
}

func TestCmdWithMsg(t *testing.T) {
	c := NewCmd("some-command")
	c.SetMsg("some message")

	assert.Equal(t, fmt.Sprintf("%s", c), "::some-command::some message")
}

func TestCmdWithMsgAndProps(t *testing.T) {
	c := NewCmd("some-command")
	c.SetProp("prop1", "value 1")
	c.SetProp("prop2", "value 2")
	c.SetMsg("some message")

	assert.Equal(t, fmt.Sprintf("%s", c), "::some-command prop1=value 1,prop2=value 2::some message")
}

func TestCmdWith2Props(t *testing.T) {
	c := NewCmd("some-command")
	c.SetProp("prop1", "value 1")
	c.SetProp("prop2", "value 2")

	assert.Equal(t, fmt.Sprintf("%s", c), "::some-command prop1=value 1,prop2=value 2::")
}

func TestCmdWithObjectProps(t *testing.T) {
	c := NewCmd("some-command")
	c.SetProp("prop1", struct{ Test string }{"object"})
	c.SetProp("prop2", "123")
	c.SetProp("prop3", true)

	assert.Equal(t, fmt.Sprintf("%s", c), `::some-command prop1={"Test"%3A"object"},prop2=123,prop3=true::`)
}

func TestCmdFromStr(t *testing.T) {
	c, err := CmdFromStr("::some-command prop1=value 1,prop2=value 2::")

	assert.Equal(t, err, nil)
	assert.Equal(t, c.Name(), "some-command")
	assert.Equal(t, c.GetProp("prop1"), "value 1")
	assert.Equal(t, c.GetProp("prop2"), "value 2")
	assert.Equal(t, c.GetMsg(), "")
}

func TestCmdTable(t *testing.T) {
	tbl := NewCmdTable()
	cnt := 0

	assert.Equal(t, tbl.Size(), 0)
	tbl.Register("a", func(c ICmd) error {
		cnt++
		return nil
	})
	assert.Equal(t, tbl.Size(), 1)
	tbl.Handle(NewCmd("a"))
	assert.Equal(t, cnt, 1)
}
