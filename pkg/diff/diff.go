package diff

type DiffData interface {
	Key() string
}

type wrapperDiffData[T any] struct {
	Data   T
	GetKey func(T) string
}

func (w *wrapperDiffData[T]) Key() string {
	return w.GetKey(w.Data)
}

func unWrapper[T any](data DiffData) T {
	return data.(*wrapperDiffData[T]).Data
}

func wrapper[T any](input T, getKey func(T) string) DiffData {
	return &wrapperDiffData[T]{
		Data:   input,
		GetKey: getKey,
	}
}

func ToDiffData[T any](input []T, getKey func(T) string) []DiffData {
	result := make([]DiffData, len(input))
	for index := range input {
		result[index] = wrapper(input[index], getKey)
	}
	return result
}

func FromDiffData[T any](input []DiffData) []T {
	result := make([]T, len(input))
	for index := range input {
		result[index] = unWrapper[T](input[index])
	}
	return result
}

func FromDiffUpdateData[O, P any](updates []Update[DiffData, DiffData]) []Update[O, P] {
	result := make([]Update[O, P], 0, len(updates))
	for _, elem := range updates {
		result = append(result, Update[O, P]{
			Origin: unWrapper[O](elem.Origin),
			Patch:  unWrapper[P](elem.Patch),
		})
	}
	return result
}

type Update[O, P any] struct {
	Origin O
	Patch  P
}

func DiffCurdData[O, P any](origin []O, patch []P, getOKey func(O) string, getPKey func(P) string) (creates []P, updates []Update[O, P], deletes []O) {
	c, u, d := DiffCurd(ToDiffData(origin, getOKey), ToDiffData(patch, getPKey))
	creates = FromDiffData[P](c)
	updates = make([]Update[O, P], 0)
	updates = FromDiffUpdateData[O, P](u)
	deletes = FromDiffData[O](d)
	return
}

func Nop[T any](s T) T { return s }

var NopString = Nop[string]

func DiffCurd(origin []DiffData, patch []DiffData) (creates []DiffData, updates []Update[DiffData, DiffData], deletes []DiffData) {
	toMap := func(datas []DiffData) map[string]DiffData {
		result := make(map[string]DiffData, len(datas))
		for _, data := range datas {
			result[data.Key()] = data
		}
		return result
	}
	originMap := toMap(origin)
	patchMap := toMap(patch)

	for _, value := range patch {
		key := value.Key()
		originValue, ok := originMap[key]
		if !ok {
			creates = append(creates, value)
			continue
		}
		updates = append(updates, Update[DiffData, DiffData]{
			Origin: originValue,
			Patch:  value,
		})
	}
	for _, value := range origin {
		key := value.Key()
		if _, ok := patchMap[key]; !ok {
			deletes = append(deletes, value)
			continue
		}
	}
	return
}
