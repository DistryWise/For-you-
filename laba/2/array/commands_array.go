package array

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Array []int

func LoadArray(arrayName, fileName string, arrays map[string]Array) error {
	file, err := os.Open("laba2_array/" + fileName)
	if err != nil {
		fmt.Printf("Error open file: %v\n", err)
		return err
	}
	defer file.Close()

	array := arrays[arrayName]
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		number, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Error ATOI: %v", err)
			return err
		}

		array = append(array, number)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file: %v", err)
		return err
	}

	arrays[arrayName] = array
	fmt.Printf("Загрузка в массив %s успешно завершена!\n", arrayName)
	return nil
}

func SaveArray(arrayName, fileName string, arrays map[string]Array) error {
	file, err := os.Create("laba2_array/" + fileName)
	if err != nil {
		fmt.Printf("Error create file: %v\n", err)
		return err
	}
	defer file.Close()

	array := arrays[arrayName]
	for _, element := range array {
		file.WriteString(strconv.Itoa(element) + "\n")
	}

	fmt.Printf("Сохранение в файл %s успешно завершена!\n", fileName)
	return nil
}

func RandArray(arrayName string, count, lb, rb int, arrays map[string]Array) error {
	rand.Seed(time.Now().UnixNano())
	array := arrays[arrayName]

	for i := 0; i < count; i++ {
		array = append(array, lb+rand.Intn(rb-lb+1))
	}

	arrays[arrayName] = array
	return nil
}

func ConcatArray(arrayName1, arrayName2 string, arrays map[string]Array) error {
	array1 := arrays[arrayName1]
	array2 := arrays[arrayName2]

	array1 = append(array1, array2...)

	arrays[arrayName1] = array1
	return nil
}

func FreeArray(arrayName string, arrays map[string]Array) error {
	delete(arrays, arrayName)
	return nil
}

func RemoveArray(arrayName string, index, count int, arrays map[string]Array) error {
	array := arrays[arrayName]

	if index < 0 || index >= len(array) {
		fmt.Println("Error invalid index")
		return errors.New("invalid index")
	}

	if index+count > len(array) {
		array = array[:index]
	} else {
		array = append(array[:index], array[index+count:]...)
	}

	arrays[arrayName] = array
	return nil
}

func CopyArray(arrayName1, arrayName2 string, lb, rb int, arrays map[string]Array) error {
	array1 := arrays[arrayName1]

	if lb < 0 || rb < 0 || lb >= len(array1) || rb >= len(array1) || lb > rb {
		fmt.Println("Error invalid index")
		return errors.New("invalid index")
	}

	array2 := append(arrays[arrayName2], array1[lb:rb+1]...)
	arrays[arrayName2] = array2
	return nil
}

func QuickSort(array []int, left, right int, option byte) {
	if left >= right {
		return
	}
	i, j := left, right
	pivot := array[(left+right)/2]

	for i <= j {
		if option == '+' {
			for array[i] < pivot {
				i++
			}
			for array[j] > pivot {
				j--
			}
		} else {
			for array[i] > pivot {
				i++
			}
			for array[j] < pivot {
				j--
			}
		}

		if i <= j {
			array[i], array[j] = array[j], array[i]
			i++
			j--
		}
	}

	QuickSort(array, left, j, option)
	QuickSort(array, i, right, option)
}

func SortArray(arrayName string, arrays map[string]Array) error {
	array := arrays[string(arrayName[0])]
	option := arrayName[1]

	QuickSort(array, 0, len(array)-1, option)

	arrays[string(arrayName[0])] = array
	return nil
}

func ShuffleArray(arrayName string, arrays map[string]Array) error {
	array := arrays[arrayName]

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(array), func(i, j int) {
		array[i], array[j] = array[j], array[i]
	})

	arrays[arrayName] = array
	return nil
}

func StatsArray(arrayName string, arrays map[string]Array) error {
	array := arrays[arrayName]

	if len(array) == 0 {
		fmt.Printf("Array %s is empty\n", arrayName)
		return nil
	}

	maxEl, minEl := math.MinInt, math.MaxInt
	maxElIndex, minElIndex := 0, 0
	sum := 0
	frequency := make(map[int]int)
	maxCount := 0
	mostCountElement := array[0]

	for index, element := range array {
		if element > maxEl {
			maxEl = element
			maxElIndex = index
		}
		if element < minEl {
			minEl = element
			minElIndex = index
		}

		sum += element

		frequency[element]++
		if frequency[element] > maxCount {
			maxCount = frequency[element]
			mostCountElement = element
		} else if frequency[element] == maxCount && element > mostCountElement {
			mostCountElement = element
		}
	}

	mean := float64(sum) / float64(len(array))
	maxDeviation := 0.0
	for _, value := range array {
		deviation := math.Abs(float64(value) - mean)
		if deviation > maxDeviation {
			maxDeviation = deviation
		}
	}

	fmt.Printf("Размер массива %s - %d\n", arrayName, len(array))
	fmt.Printf("Минимальный элемент: %d (индекс %d)\n", minEl, minElIndex)
	fmt.Printf("Максимальный элемент: %d (индекс %d)\n", maxEl, maxElIndex)
	fmt.Printf("Наиболее часто встречающийся элемент: %d\n", mostCountElement)
	fmt.Printf("Среднее значение элементов: %.2f\n", mean)
	fmt.Printf("Максимальное отклонение от среднего значения: %.2f\n", maxDeviation)
	return nil
}

func PrintArray(arrayName string, index string, arrays map[string]Array) error {
	array := arrays[arrayName]

	if index == "all" {
		fmt.Printf("Массив %s - %v\n", arrayName, array)
		return nil
	}

	indexDigit, err := strconv.Atoi(index)
	if err != nil {
		fmt.Println("Error ATOI\n")
		return err
	}

	if indexDigit < 0 || indexDigit >= len(array) {
		fmt.Println("Error invalid index\n")
		return errors.New("invalid index")
	}

	fmt.Printf("Элемент под индексом %s в массиве %s - %d\n", index, arrayName, array[indexDigit])
	return nil
}

func PrintRangeArray(arrayName string, lb, rb int, arrays map[string]Array) error {
	array := arrays[arrayName]

	if lb < 0 || rb < 0 || lb >= len(array) || rb >= len(array) || lb > rb {
		fmt.Println("Error invalid index")
		return errors.New("invalid index")
	}

	fmt.Printf("Элементы массива %s с %d по %d - %v\n", arrayName, lb, rb, array[lb:rb+1])
	return nil
}
