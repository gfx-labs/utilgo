package lambda_test

import (
	"reflect"
	"testing"

	"git.tuxpa.in/a/lambda"
)

func TestFlattenIntSix(t *testing.T) {
	var arr [][]int
	arr = [][]int{{1}, {2}, {3}, {4}, {6}, {5}}
	ans := lambda.Flatten(arr)
	exp := []int{1, 2, 3, 4, 6, 5}
	if !reflect.DeepEqual(ans, exp) {
		t.Fatalf("expect %v to be %v", ans, exp)
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
	if !reflect.DeepEqual(ans, exp) {
		t.Fatalf("expect %v to be %v", ans, exp)
	}
}

func TestSplitIntSix2(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	ans := lambda.Split(arr, 2)
	exp := [][]int{{1, 2, 3}, {4, 6, 5}}
	if !reflect.DeepEqual(ans, exp) {
		t.Fatalf("expect %v to be %v", ans, exp)
	}
}
func TestSplitIntSix6(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	ans := lambda.Split(arr, 6)
	exp := [][]int{{1}, {2}, {3}, {4}, {6}, {5}}
	if !reflect.DeepEqual(ans, exp) {
		t.Fatalf("expect %v to be %v", ans, exp)
	}
}

func TestSplitIntSix5(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	ans := lambda.Split(arr, 5)
	exp := [][]int{{1}, {2}, {3}, {4}, {6, 5}}
	if !reflect.DeepEqual(ans, exp) {
		t.Fatalf("expect %v to be %v", ans, exp)
	}
}

func TestSplitIntSix10(t *testing.T) {
	arr := []int{1, 2, 3, 4, 6, 5}
	ans := lambda.Split(arr, 10)
	exp := [][]int{{1, 2, 3, 4, 6, 5}}
	if !reflect.DeepEqual(ans, exp) {
		t.Fatalf("expect %v to be %v", ans, exp)
	}
}
