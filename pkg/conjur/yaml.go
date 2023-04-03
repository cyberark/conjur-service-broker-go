package conjur

import (
	"bytes"
	"io"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type Tag interface {
	Value() interface{}
}

type tag[T any] struct {
	v T
}

func (t *tag[T]) Value() interface{} {
	return t.v
}

func NewTag[T any](v T) Tag {
	return &tag[T]{v}
}

func (t *tag[T]) MarshalYAML() (interface{}, error) {
	type aliasType *T
	data := (aliasType)(&t.v)

	node := &yaml.Node{Kind: yaml.MappingNode}
	if err := node.Encode(data); err != nil {
		return nil, err
	}
	node.Tag = typeName(t.v)
	node.Style = yaml.TaggedStyle

	return node, nil
}

type Ref interface {
	Kind() string
	Value() string
}

type ref struct {
	kind  string
	value string
}

func (r *ref) Kind() string {
	return r.kind
}

func (r *ref) Value() string {
	return r.value
}

func NewRef[T any](v string) Ref {
	return &ref{typeName(new(T)), v}
}

func (r *ref) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: r.Value(),
		Tag:   r.Kind(),
		Style: yaml.TaggedStyle,
	}, nil
}

func policyReader(policy PolicyDocument) (io.Reader, error) {
	res := new(bytes.Buffer)
	encoder := yaml.NewEncoder(res)
	err := encoder.Encode(policy)
	if err != nil {
		return nil, err
	}
	return res, err
}

func typeName(v interface{}) string {
	name := reflect.TypeOf(v).String()
	dot := strings.LastIndex(name, ".")
	if dot == -1 {
		return ""
	}
	if dot >= len(name) {
		return ""
	}
	return "!" + strings.ToLower(name[dot+1:])
}
