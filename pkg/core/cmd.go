package core

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// ::name key=value,key=value::message

type ICmd interface {
	Name() string
	GetMsg() string
	SetMsg(msg string)
	GetProp(k string) interface{}
	SetProp(k string, v interface{})
	DelProp(k string)
}

type Cmd struct {
	name string
	msg  string
	m    map[string]interface{}
}

var cmdPattern = regexp.MustCompile(`^::(?P<name>\S+)\s?(?P<props>.*)::(?P<msg>.*)`)

func CmdFromStr(s string) (ICmd, error) {
	match := cmdPattern.FindStringSubmatch(s)
	if len(match) == 4 {
		c := NewCmd(match[1])
		c.SetMsg(match[3])
		props := strings.Split(match[2], ",")
		for _, prop := range props {
			if prop == "" {
				continue
			}
			kv := strings.Split(prop, "=")
			if len(kv) != 2 {
				continue
			}
			k, v := kv[0], kv[1]
			if k == "" {
				return nil, fmt.Errorf("Empty key")
			}
			c.SetProp(k, unEscapeProp(v))
		}
		return c, nil
	}
	return nil, fmt.Errorf("Invalid string: \"%s\"", s)
}

func NewCmd(name string) ICmd {
	return &Cmd{name: name, m: make(map[string]interface{})}
}

func (c *Cmd) Name() string {
	return c.name
}

func (c *Cmd) GetMsg() string {
	return c.msg
}

func (c *Cmd) SetMsg(msg string) {
	c.msg = msg
}

func (c *Cmd) GetProp(k string) interface{} {
	v, _ := c.m[k]
	return v
}

func (c *Cmd) SetProp(k string, v interface{}) {
	c.m[k] = v
}

func (c *Cmd) DelProp(k string) {
	delete(c.m, k)
}

func (c *Cmd) String() string {
	kvs := []string{}
	for k, v := range c.m {
		kvs = append(kvs, fmt.Sprintf("%s=%s", k, escapeProp(prop2json(v))))
	}
	sort.Strings(kvs)
	props := ""
	if len(kvs) > 0 {
		props = " " + strings.Join(kvs, ",")
	}
	return fmt.Sprintf("::%s%s::%s",
		c.Name(),
		props,
		escapeMsg(c.GetMsg()),
	)
}

type ICmdHandler func(c ICmd) error
type ICmdTable interface {
	Register(name string, handler ICmdHandler)
	UnRegister(name string)
	Handle(c ICmd) error
	Size() int
}

type CmdTable struct {
	m map[string]ICmdHandler
}

func NewCmdTable() ICmdTable {
	return &CmdTable{make(map[string]ICmdHandler)}
}

func (t *CmdTable) Size() int {
	return len(t.m)
}

func (t *CmdTable) Register(name string, handler ICmdHandler) {
	t.m[name] = handler
}

func (t *CmdTable) UnRegister(name string) {
	delete(t.m, name)
}

func (t *CmdTable) Handle(c ICmd) error {
	hdlr, ok := t.m[c.Name()]
	if ok {
		return hdlr(c)
	}
	return nil
}

func prop2json(v interface{}) string {
	if nil == v {
		return ""
	}
	switch _v := v.(type) {
	case string:
		return _v
	default:
		b, err := json.Marshal(_v)
		if err != nil {
			panic(err)
		}
		return string(b)
	}
}

func unEscapeProp(s string) string {
	s = strings.ReplaceAll(s, "%3A", ":")
	s = strings.ReplaceAll(s, "%2C", ",")
	s = strings.ReplaceAll(s, "%25", "%")
	s = strings.ReplaceAll(s, "%0D", "\r")
	s = strings.ReplaceAll(s, "%0A", "\n")
	return s
}

func escapeProp(s string) string {
	s = escapeMsg(s)
	s = strings.ReplaceAll(s, ":", "%3A")
	s = strings.ReplaceAll(s, ",", "%2C")
	return s
}

func escapeMsg(s string) string {
	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "\r", "%0D")
	s = strings.ReplaceAll(s, "\n", "%0A")
	return s
}
