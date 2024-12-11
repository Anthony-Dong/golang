package diff

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffCurdData(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		creates, updates, deletes := DiffCurdData([]string{"1", "2", "3", "4"}, []string{"3", "4", "5", "6"}, NopString, NopString)
		assert.Equal(t, creates, []string{"5", "6"})
		assert.Equal(t, updates, []Update[string, string]{
			{
				Origin: "3", Patch: "3",
			},
			{
				Origin: "4", Patch: "4",
			},
		})
		assert.Equal(t, deletes, []string{"1", "2"})
	})

	t.Run("test2", func(t *testing.T) {
		creates, updates, deletes := DiffCurdData([]string{"1", "2", "3", "4"}, []int{3, 4, 5, 6}, NopString, func(i int) string {
			return strconv.Itoa(i)
		})
		assert.Equal(t, creates, []int{5, 6})
		assert.Equal(t, updates, []Update[string, int]{
			{
				Origin: "3", Patch: 3,
			},
			{
				Origin: "4", Patch: 4,
			},
		})
		assert.Equal(t, deletes, []string{"1", "2"})
	})

	t.Run("1", func(t *testing.T) {
		t.Logf("%.2f\n", float64(13213123123213)/float64(10011101))
	})
}
