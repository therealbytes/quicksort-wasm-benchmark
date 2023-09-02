package quicksort

const (
	SEED     = 7
	L        = 1000
	N        = 100
	CHECKSUM = 107829970005
)

type Quicksort struct {
	seed uint64
}

func NewQuicksortBenchmark() *Quicksort {
	return &Quicksort{seed: SEED}
}

func (qs *Quicksort) random() uint {
	qs.seed = (1103515245*qs.seed + 12345) % (1 << 31)
	return uint(qs.seed)
}

func (qs *Quicksort) randomizeArray(arr []uint) {
	for i := range arr {
		arr[i] = qs.random()
	}
}

func (qs *Quicksort) quicksort(arr []uint, left, right int) {
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

func (qs *Quicksort) Benchmark() uint64 {
	var checksum uint64 = 0
	arr := make([]uint, L)
	for i := 0; i < N; i++ {
		qs.randomizeArray(arr)
		qs.quicksort(arr, 0, len(arr)-1)
		checksum += uint64(arr[L/2])
	}
	return checksum
}
