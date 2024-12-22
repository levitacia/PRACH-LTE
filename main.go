package main

import (
	"fmt"
	"math/rand"
)

const (
	initialUEs        = 100 // Стартовое количество абонентов
	numPreambles      = 64  // Число преамбул в PRACH
	maxRetries        = 5   // Максимальное количество попыток
	maxSubframes      = 100 // Общее количество сабфреймов для моделирования
	maxNewUEsPerFrame = 70  // Максимальное число новых абонентов за сабфрейм
)

var (
	preambles    [numPreambles]Preamble
	connectedUEs []int
	excludedUEs  []int // Список устройств, исключённых из системы
	totalUEs     = 0   // Общее количество устройств в системе
)

type Preamble struct {
	ID        int
	UsedByUEs []int
}

type UE struct {
	ID          int
	Attempts    int
	PreambleID  int
	TempCRNTI   int
	IsConnected bool
	Excluded    bool
}

func InitializePreambles() {
	for i := 0; i < numPreambles; i++ {
		preambles[i] = Preamble{ID: i, UsedByUEs: []int{}}
	}
}

func ResetPreambles() {
	for i := 0; i < numPreambles; i++ {
		preambles[i].UsedByUEs = []int{}
	}
}

func (ue *UE) ChoosePreamble() {
	ue.PreambleID = rand.Intn(numPreambles)
	preambles[ue.PreambleID].UsedByUEs = append(preambles[ue.PreambleID].UsedByUEs, ue.ID)
}

func ProcessSubframe(ues *[]UE) (int, int, int) {
	collisions := 0
	successfulConnections := 0
	exclusions := 0

	for i := range *ues {
		// Обрабатываем только устройства, которые еще не подключены и не исключены
		if !(*ues)[i].IsConnected && !(*ues)[i].Excluded {
			(*ues)[i].ChoosePreamble()
		}
	}

	// Проверка на коллизии
	for i := 0; i < numPreambles; i++ {
		if len(preambles[i].UsedByUEs) == 1 {
			// Успешное подключение
			ueID := preambles[i].UsedByUEs[0]
			for j := range *ues {
				if (*ues)[j].ID == ueID {
					(*ues)[j].IsConnected = true
					(*ues)[j].TempCRNTI = ueID // Временный C-RNTI
					connectedUEs = append(connectedUEs, ueID)
					successfulConnections++
					break
				}
			}
		} else if len(preambles[i].UsedByUEs) > 1 {
			// Коллизия
			collisions++
			for _, ueID := range preambles[i].UsedByUEs {
				for j := range *ues {
					if (*ues)[j].ID == ueID {
						(*ues)[j].Attempts++
						if (*ues)[j].Attempts > maxRetries {
							(*ues)[j].Excluded = true
							excludedUEs = append(excludedUEs, ueID)
							exclusions++
						}
					}
				}
			}
		}
	}

	return successfulConnections, collisions, exclusions
}

func AddNewUEs(ues *[]UE, numNewUEs int) {
	startID := totalUEs
	for i := 0; i < numNewUEs; i++ {
		newUE := UE{ID: startID + i}
		*ues = append(*ues, newUE)
	}
	totalUEs += numNewUEs
}

func main() {
	var ues []UE
	totalUEs = initialUEs

	// Добавляем начальные устройства
	for i := 0; i < initialUEs; i++ {
		ues = append(ues, UE{ID: i})
	}

	subframe := 0
	for subframe < maxSubframes {
		subframe++
		fmt.Printf("\n--- Subframe %d ---\n", subframe)
		InitializePreambles()

		// Добавляем случайное количество новых абонентов
		newUEs := rand.Intn(maxNewUEsPerFrame + 1)
		if newUEs > 0 {
			AddNewUEs(&ues, newUEs)
			fmt.Printf("Добавлено %d новых абонентов, всего устройств: %d\n", newUEs, totalUEs)
		}

		// Подсчитываем количество устройств, которые пытаются подключиться
		activeUEs := 0
		for _, ue := range ues {
			if !ue.IsConnected && !ue.Excluded {
				activeUEs++
			}
		}

		fmt.Printf("Устройств, пытающихся подключиться: %d\n", activeUEs)

		if activeUEs > 0 {
			// Обрабатываем текущий сабфрейм
			successfulConnections, collisions, exclusions := ProcessSubframe(&ues)

			// Выводим подробную информацию о сабфрейме
			fmt.Printf("Успешные подключения: %d\n", successfulConnections)
			fmt.Printf("Коллизии: %d\n", collisions)
			fmt.Printf("Исключенные устройства за этот сабфрейм: %d\n", exclusions)
		} else {
			fmt.Println("Нет активных устройств для обработки, ожидание новых устройств...")
		}

		ResetPreambles()
	}

	// Итоги моделирования
	fmt.Printf("\nИтог: %d устройств подключено, %d исключено, всего %d сабфреймов\n",
		len(connectedUEs), len(excludedUEs), subframe)

	if len(excludedUEs) > 0 {
		fmt.Println("Устройства, которые не удалось подключить совсем:")
		for _, ueID := range excludedUEs {
			fmt.Printf("UE %d\n", ueID)
		}
	}
}
