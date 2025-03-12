package wal

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/anthony-dong/golang/pkg/compress"
	"github.com/anthony-dong/golang/pkg/utils"
)

type MetaInfo struct {
	Key          string            `json:"k"`
	Offset       int64             `json:"o"`
	Size         int               `json:"s"`
	CompressType compress.Type     `json:"c,omitempty"`
	Tags         map[string]string `json:"t,omitempty"`
}

func (m *MetaInfo) Clone() MetaInfo {
	if m == nil {
		return MetaInfo{}
	}
	clone := *m
	clone.Tags = make(map[string]string, len(clone.Tags))
	for k, v := range m.Tags {
		clone.Tags[k] = v
	}
	return clone
}

func marshalMetaInfo(m *MetaInfo) ([]byte, error) {
	marshal, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return append(marshal, '\n'), nil
}

type Wal struct {
	metaInfos  []*MetaInfo
	indexes    map[string][]int
	tagIndexes map[string]map[string][]int
	buffer     Buffer

	dataOffset  int64 // data 偏移量
	indexOffset int64 // index 偏移量

	indexMaxSize int // 索引最大存储的大小
}

func NewWal(buffer Buffer, indexMaxSize int) (*Wal, error) {
	if indexMaxSize <= 0 {
		return nil, errors.New("indexMaxSize must be greater than zero")
	}
	w := &Wal{indexes: make(map[string][]int), buffer: buffer, dataOffset: 0, indexOffset: 0, indexMaxSize: indexMaxSize, metaInfos: make([]*MetaInfo, 0), tagIndexes: make(map[string]map[string][]int)}
	if err := w.initIndex(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Wal) initIndex() error {
	if _, err := w.buffer.Seek(0, io.SeekStart); err != nil {
		return err
	}

	decoder := json.NewDecoder(w.buffer)
	for {
		metaInfo := &MetaInfo{}
		if err := decoder.Decode(metaInfo); err != nil {
			break
		}
		if metaInfo.Key == "" || metaInfo.Offset == 0 {
			break
		}
		w.appendMetaInfo(metaInfo)
	}

	if len(w.metaInfos) == 0 {
		w.dataOffset = int64(w.indexMaxSize)
		w.indexOffset = 0
		if _, err := w.buffer.Seek(0, io.SeekStart); err != nil {
			return err
		}
		return nil
	}

	// todo 优化效率
	indexOffset := 0
	for _, elem := range w.metaInfos {
		info, err := marshalMetaInfo(elem)
		if err != nil {
			return err
		}
		indexOffset = indexOffset + len(info)
	}
	if _, err := w.buffer.Seek(int64(indexOffset), io.SeekStart); err != nil {
		return err
	}

	lastData := w.metaInfos[len(w.metaInfos)-1]
	w.dataOffset = lastData.Offset + int64(lastData.Size)
	w.indexOffset = int64(indexOffset)
	return nil
}

func (w *Wal) SetString(key string, value string, c compress.Type) error {
	return w.Set(key, utils.String2Bytes(value), c, nil)
}

func (w *Wal) writeIndex(metaInfo *MetaInfo) error {
	data, err := marshalMetaInfo(metaInfo)
	if err != nil {
		return err
	}
	if int(w.indexOffset)+len(data) > w.indexMaxSize {
		return indexOutOfRangeErr
	}
	wn, err := w.buffer.WriteAt(data, w.indexOffset)
	if err != nil {
		return err
	}
	if wn != len(data) {
		return errors.New("wal write fail")
	}

	w.appendMetaInfo(metaInfo)
	w.indexOffset = w.indexOffset + int64(wn)
	return nil
}

func (w *Wal) appendMetaInfo(metaInfo *MetaInfo) {
	index := len(w.metaInfos)
	w.metaInfos = append(w.metaInfos, metaInfo)
	w.indexes[metaInfo.Key] = append(w.indexes[metaInfo.Key], index)
	for k, v := range metaInfo.Tags {
		if w.tagIndexes[k] == nil {
			w.tagIndexes[k] = map[string][]int{}
		}
		w.tagIndexes[k][v] = append(w.tagIndexes[k][v], index)
	}
}

func (w *Wal) List() []string {
	keys := make([]string, 0, len(w.metaInfos))
	for key := range w.indexes {
		keys = append(keys, key)
	}
	return keys
}

func (w *Wal) Set(key string, value []byte, c compress.Type, tags map[string]string) error {
	if key == "" {
		return errors.New("key can not be empty")
	}
	data, err := compress.NewCompressor(c).Compress(value)
	if err != nil {
		return err
	}
	metaInfo := &MetaInfo{Key: key, Offset: w.dataOffset, Size: len(data), CompressType: c, Tags: tags}

	// write index
	if err := w.writeIndex(metaInfo); err != nil {
		return err
	}

	// write value
	wn, err := w.buffer.WriteAt(data, metaInfo.Offset)
	if err != nil {
		return err
	}
	if wn != len(data) {
		return errors.New("wal write fail")
	}
	w.dataOffset = w.dataOffset + int64(wn)
	return nil
}

var (
	notFoundErr        = errors.New("not found")
	indexOutOfRangeErr = errors.New("index out of range")
)

func IsNotFoundErr(err error) bool {
	return err == notFoundErr
}

func IsIndexOutOfRangeErr(err error) bool {
	return err == indexOutOfRangeErr
}

func (w *Wal) GetString(key string) (string, error) {
	get, err := w.Get(key)
	if err != nil {
		return "", err
	}
	return string(get), nil
}

func (w *Wal) Get(key string) ([]byte, error) {
	indexes := w.indexes[key]
	if len(indexes) == 0 {
		return nil, notFoundErr
	}
	info := w.metaInfos[indexes[len(indexes)-1]]
	data := make([]byte, info.Size)
	n, err := w.buffer.ReadAt(data, info.Offset)
	if err != nil {
		return nil, err
	}
	data = data[:n]
	return compress.NewCompressor(info.CompressType).Decompress(data)
}

func (w *Wal) byIndex(i int) (MetaInfo, bool) {
	if i >= len(w.metaInfos) {
		return MetaInfo{}, false
	}
	return w.metaInfos[i].Clone(), true
}

func (w *Wal) byMultiIndex(indexes []int) []MetaInfo {
	result := make([]MetaInfo, 0)
	for _, index := range indexes {
		info, isOk := w.byIndex(index)
		if isOk {
			result = append(result, info)
		}
	}
	return result
}

func (w *Wal) SearchByTags(tag map[string]string) []MetaInfo {
	index := make([]int, 0)
	count := 0
	for k, v := range tag {
		count = count + 1
		curIndex := w.searchByTag(k, v)
		if count == 1 {
			index = curIndex
			continue
		}
		index = intersection(index, curIndex)
	}
	return w.byMultiIndex(index)
}

func (w *Wal) searchByTag(k, v string) []int {
	return w.tagIndexes[k][v]
}

func intersection(arr1 []int, arr2 []int) []int {
	set1 := make(map[int]struct{})
	for _, num := range arr1 {
		set1[num] = struct{}{}
	}
	var result []int
	set2 := make(map[int]struct{})
	for _, num := range arr2 {
		if _, found := set1[num]; found {
			if _, alreadyAdded := set2[num]; !alreadyAdded {
				result = append(result, num)
				set2[num] = struct{}{}
			}
		}
	}
	return result
}
