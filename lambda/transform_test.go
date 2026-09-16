package lambda_test

import (
	"reflect"
	"testing"

	"github.com/gfx-labs/utilgo/lambda"
	"github.com/stretchr/testify/assert"
)

func TestFlattenIntSix(t *testing.T) {
	arr := [][]int{{1}, {2}, {3}, {4}, {6}, {5}}
	ans := lambda.Flatten(arr)
	exp := []int{1, 2, 3, 4, 6, 5}
	if !reflect.DeepEqual(exp, ans) {
		t.Fatalf("expect %v to be %v", exp, ans)
	}
}

func TestSplitIntSix0(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	ans := lambda.Split(arr, 0)
	if ans != nil {
		t.Fatalf("expect %v to be %v", ans, nil)
	}
}

func TestSplitIntSix1(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	ans := lambda.Split(arr, 1)
	exp := [][]int{{1, 2, 3, 4, 6, 5}}
	assert.EqualValues(t, exp, ans)
}

func TestSplitIntSix2(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	ans := lambda.Split(arr, 2)
	exp := [][]int{{1, 2, 3}, {4, 6, 5}}
	assert.EqualValues(t, exp, ans)
}
func TestSplitIntSix6(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	ans := lambda.Split(arr, 6)
	exp := [][]int{{1}, {2}, {3}, {4}, {6}, {5}}
	assert.EqualValues(t, exp, ans)
}

func TestSplitIntSix5(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	ans := lambda.Split(arr, 5)
	exp := [][]int{{1}, {2}, {3}, {4}, {6, 5}}
	assert.EqualValues(t, exp, ans)
}

func TestSplitIntSix10(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	ans := lambda.Split(arr, 10)
	exp := [][]int{{1, 2, 3, 4, 6, 5}}
	assert.EqualValues(t, exp, ans)
}
func TestSegmentIntTen3(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ans := lambda.Segment(arr, 3)
	exp := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}
	assert.EqualValues(t, exp, ans)
}

func TestSegmentIntTen4(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ans := lambda.Segment(arr, 4)
	exp := [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10}}
	assert.EqualValues(t, exp, ans)
}
func TestSegmentIntTen11(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ans := lambda.Segment(arr, 11)
	exp := [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	assert.EqualValues(t, exp, ans)
}

func TestSegmentIntTen2(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ans := lambda.Segment(arr, 2)
	exp := [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}}
	assert.EqualValues(t, exp, ans)
}
