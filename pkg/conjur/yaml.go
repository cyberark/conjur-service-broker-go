package conjur

import (
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func typeName(v interface{}) string {
	return "!" + strings.ToLower(reflect.TypeOf(v).Name())
}

type Tag[T any] struct {
	v T
}

func NewTag[T any](v T) *Tag[T] {
	return &Tag[T]{v}
}

func (t *Tag[T]) MarshalYAML() (interface{}, error) {
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

type Ref[T any] string

func (r *Ref[T]) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: string(*r),
		Tag:   typeName(*new(T)),
		Style: yaml.TaggedStyle,
	}, nil
}
