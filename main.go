package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"sync"
	"time"
)

const (
	numUEs       = 100 // мобилки
	numPreambles = 64  // в 1 prach сетке
	maxRetries   = 5   // попыток при коллизиях
	//Preamble sequence
	Nzc = 839
	u   = 1
	// ----
)

var (
	preambles    [numPreambles]Preamble
	connectedUEs []int
	mux          sync.Mutex
	randSource   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Preamble struct {
	ID        int
	UsedByUEs []int
	sequence  []complex128
}

type RandomAccessResponse struct {
	TimingAdvance int
	TempCRNTI     int
	Success       bool
}

type UE struct {
	ID          int
	Attempts    int
	PreambleID  int
	TempCRNTI   int
	IsConnected bool
}

func generateZadoffChuSequence(length, u int) []complex128 {
	sequence := make([]complex128, length)
	for n := 0; n < length; n++ {
		exponent := -math.Pi * float64(u*n*(n+1)) / float64(length)
		sequence[n] = cmplx.Exp(complex(0, exponent))
	}
	return sequence
}

// InitializePreambles инициализирует доступные преамбулы
func InitializePreambles() {
	for i := 0; i < numPreambles; i++ {
		preambles[i] = Preamble{ID: i, UsedByUEs: []int{}}
	}
}

// ResetPreambles освобождает преамбулы для нового subframe
func ResetPreambles() {
	for i := 0; i < numPreambles; i++ {
		preambles[i].UsedByUEs = []int{}
	}
}

// ChoosePreamble выбирает случайную преамбулу для UE
func (ue *UE) ChoosePreamble() {
	ue.PreambleID = randSource.Intn(numPreambles)
	mux.Lock()
	preambles[ue.PreambleID].UsedByUEs = append(preambles[ue.PreambleID].UsedByUEs, ue.ID)
	mux.Unlock()
}

// TransmitPreamble имитирует передачу преамбулы и получение ответа от eNodeB
func (ue *UE) TransmitPreamble() RandomAccessResponse {
	mux.Lock()
	users := preambles[ue.PreambleID].UsedByUEs
	mux.Unlock()

	if len(users) > 1 {
		// Коллизия
		return RandomAccessResponse{Success: false}
	}

	// Успешный ответ
	return RandomAccessResponse{
		TimingAdvance: randSource.Intn(100),
		TempCRNTI:     ue.ID,
		Success:       true,
	}
}

// PerformRandomAccessForSubframe реализует процесс случайного доступа для одного subframe
func (ue *UE) PerformRandomAccessForSubframe() bool {
	if ue.IsConnected {
		return true
	}

	ue.ChoosePreamble()
	response := ue.TransmitPreamble()

	if response.Success {
		ue.TempCRNTI = response.TempCRNTI
		ue.IsConnected = true
		mux.Lock()
		connectedUEs = append(connectedUEs, ue.ID)
		mux.Unlock()
		fmt.Printf("UE %d успешно подключено с временным C-RNTI %d\n", ue.ID, ue.TempCRNTI)
		return true
	}

	ue.Attempts++
	return false
}

func main() {
	var ues []UE
	for i := 0; i < numUEs; i++ {
		ues = append(ues, UE{ID: i})
	}

	subframe := 0
	for len(connectedUEs) < numUEs {
		subframe++
		fmt.Printf("--- Subframe %d ---\n", subframe)
		InitializePreambles()

		for i := range ues {
			if !ues[i].IsConnected {
				ues[i].PerformRandomAccessForSubframe()
			}
		}
		ResetPreambles()
	}

	fmt.Printf("Итог: %d устройств подключены из %d за %d subframe\n", len(connectedUEs), numUEs, subframe)
}
