package yamlembed

import "strings"

type Foo struct {
	A string `yaml:"aa"`
	p int64  `yaml:"p"`
}

type Bar struct {
	I      int64    `yaml:",omitempty"`
	B      string   `yaml:"b"`
	UpperB string   `yaml:"-"`
	OI     []string `yaml:",omitempty"`
	F      []any    `yaml:"f,flow"`
}

func (b *Bar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error

	originMap := make(map[string]any)

	err = unmarshal(originMap)

	b.B = originMap["b"].(string)
	b.UpperB = strings.ToUpper(b.B)

	oi := originMap["oi"].([]any)
	for _, v := range oi {
		b.OI = append(b.OI, v.(string))
	}

	b.F = originMap["f"].([]any)

	return err
}

type Baz struct {
	Foo `yaml:",inline"`
	Bar `yaml:",inline"`
}

func (b *Baz) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error

	originMap := make(map[string]any)

	err = unmarshal(&originMap)

	if err == nil {
		b.Foo.A = originMap["aa"].(string)

		err = b.Bar.UnmarshalYAML(unmarshal)
	}

	return err
}
