package quicksort

type QuicksortBenchmark struct {
	seed uint
}

func NewQuicksortBenchmark(seed uint) *QuicksortBenchmark {
	return &QuicksortBenchmark{seed: seed}
}

func (qs *QuicksortBenchmark) random() uint {
	qs.seed = (1103515245*qs.seed + 12345) % (1 << 31)
	return qs.seed
}

func (qs *QuicksortBenchmark) randomizeArray(arr []uint) {
	for i := range arr {
		arr[i] = qs.random() % 1000
	}
}

func (qs *QuicksortBenchmark) quicksort(arr []uint, left, right int) {
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
		qs.quicksort(arr, left, j)
	}
	if i < right {
		qs.quicksort(arr, i, right)
	}
}

func (qs *QuicksortBenchmark) Run(arrLen int, iter int) uint {
	var checksum uint
	arr := make([]uint, arrLen)
	for i := 0; i < iter; i++ {
		qs.randomizeArray(arr)
		qs.quicksort(arr, 0, arrLen-1)
		checksum += arr[arrLen/2]
	}
	return checksum
}
