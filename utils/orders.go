package utils

import "orderbot/models"

func ArrangeOrderByPrice (orders []*models.Order) ([]*models.Order) {
	return mergeSortReverse(orders)
}

func mergeSort(items []*models.Order) []*models.Order {
	if items == nil || len(items) == 0 {
		return items
	}
	var num = len(items)

	if num == 1 {
		return items
	}

	middle := num / 2
	var (
		left = make([]*models.Order, middle)
		right = make([]*models.Order, num-middle)
	)
	for i := 0; i < num; i++ {
		if i < middle {
			left[i] = items[i]
		} else {
			right[i-middle] = items[i]
		}
	}

	return merge(mergeSort(left), mergeSort(right))
}

func mergeSortReverse(items []*models.Order) []*models.Order {
	if items == nil || len(items) == 0 {
		return items
	}
	var num = len(items)

	if num == 1 {
		return items
	}

	middle := num / 2
	var (
		left = make([]*models.Order, middle)
		right = make([]*models.Order, num-middle)
	)
	for i := 0; i < num; i++ {
		if i < middle {
			left[i] = items[i]
		} else {
			right[i-middle] = items[i]
		}
	}

	return mergeReverse(mergeSortReverse(left), mergeSortReverse(right))
}

func mergeReverse(left, right []*models.Order) (result []*models.Order) {
	result = make([]*models.Order, len(left) + len(right))

	i := 0
	for len(left) > 0 && len(right) > 0 {
		if left[0].Price > right[0].Price {
			result[i] = left[0]
			left = left[1:]
		} else {
			result[i] = right[0]
			right = right[1:]
		}
		i++
	}

	for j := 0; j < len(left); j++ {
		result[i] = left[j]
		i++
	}
	for j := 0; j < len(right); j++ {
		result[i] = right[j]
		i++
	}

	return
}

func merge(left, right []*models.Order) (result []*models.Order) {
	result = make([]*models.Order, len(left) + len(right))

	i := 0
	for len(left) > 0 && len(right) > 0 {
		if left[0].Price < right[0].Price {
			result[i] = left[0]
			left = left[1:]
		} else {
			result[i] = right[0]
			right = right[1:]
		}
		i++
	}

	for j := 0; j < len(left); j++ {
		result[i] = left[j]
		i++
	}
	for j := 0; j < len(right); j++ {
		result[i] = right[j]
		i++
	}

	return
}
