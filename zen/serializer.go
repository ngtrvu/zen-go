package zen

type ModelSerializerInterface interface {
	serialize(item interface{}) interface{}
}

type ModelSerializer struct {
	ModelItem interface{}
}

func NewModelSerializer(item interface{}) *ModelSerializer {
	return &ModelSerializer{
		ModelItem: item,
	}
}

func (m ModelSerializer) serialize(item interface{}) interface{} {
	return item
}
