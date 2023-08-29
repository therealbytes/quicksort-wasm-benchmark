package quicksort

type QuickSort struct {
	seed uint64
}

func NewQuicksortBenchmark(seed uint64) *QuickSort {
	return &QuickSort{seed: seed}
}

func (qs *QuickSort) random() uint {
	qs.seed = (1103515245*qs.seed + 12345) % (1 << 32)
	return uint(qs.seed)
}

func (qs *QuickSort) randomizeArray(arr []uint) {
	for i := range arr {
		arr[i] = qs.random()
	}
}

func (qs *QuickSort) quickSort(arr []uint, left, right int) {
	i, j := left, right
	if i == j {
		return
	}
	pivot := arr[left+(right-left)/2]

	for i <= j {
		for arr[i] < pivot {
			i++
		}
		for pivot < arr[j] {
			j--
		}
		if i <= j {
			arr[i], arr[j] = arr[j], arr[i]
			i++
			j--
		}
	}
	if left < j {
		qs.quickSort(arr, left, j)
	}
	if i < right {
		qs.quickSort(arr, i, right)
	}
}

func (qs *QuickSort) Benchmark() uint64 {
	var checksum uint64 = 0
	arr := make([]uint, 1000)
	for i := 0; i < 100; i++ {
		qs.randomizeArray(arr)
		qs.quickSort(arr, 0, 999)
		checksum += uint64(arr[100])
	}
	return checksum
}
