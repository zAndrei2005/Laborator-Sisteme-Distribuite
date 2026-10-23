// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"strconv"
// )

// func main() {
// 	scanner := bufio.NewScanner(os.Stdin)

// 	scanner.Scan()

// 	v := scanner.Text()
// 	x, err := strconv.Atoi(v)

// 	if err != nil {
// 		fmt.Print("eroare")
// 	}

// 	var ok int = 0

// 	if x < 2 {
// 		fmt.Print("Nu e prim")
// 	} else {
// 		for i := 2; i <=x/2; i++ {
// 			if x % 2 == 0 {
// 				ok = 1
// 				break
// 			}
// 		}

// 		if ok == 1 {
// 			fmt.Print("Nu e prim")
// 		} else {
// 			fmt.Print("Prim")
// 		}
// 	}
// }

// EX 2    package main

// import "fmt"

// func main() {
// 	arr := [5]int{5, 2, 8, 1, 9}
// 	n := len(arr)

// 	for i := 0; i < n-1; i++ {
// 		for j := 0; j < n-i-1; j++ {
// 			if arr[j] > arr[j+1] {
// 				arr[j], arr[j+1] = arr[j+1], arr[j]
// 			}
// 		}
// 	}

// 	fmt.Println(arr)
// }

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main(){
	ex3()
	ex4()
	ex5()
	ex6()
	ex7()
	ex8()
	ex9()
}

func ex3(){
	intSlice := []int {1,2,3,4}
	intSlice2 := []int{}

	for i := 0; i < len(intSlice); i++ {
		var aux int = intSlice[i]
		if aux % 2 == 0 {
			intSlice2 = append(intSlice2, intSlice[i])
		}
	}

	fmt.Print(intSlice2)
}

func ex4(){
	slice := []float32{1.2, 2.3, 3.4, 4.5}
	scanner := bufio.NewScanner(os.Stdin)
	
	scanner.Scan()
	x := scanner.Text()
	index, err := strconv.Atoi(x)

	if err != nil {
		fmt.Println("eroare")
	}

	if index >= 0 && index < len(slice) {
		slice = append(slice[:index], slice[index+1:]...)
	}

	fmt.Print(slice)
}

func ex5(){
	slice := []string{"unu", "doi", "trei", "patru"}
	
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	
	fmt.Println(slice)
}

func ex6(){
	persoane := map[string]int{
		"Andrei": 190, 
		"Maria": 165,
		"Ion": 182,
		"Elena": 172,
		"George": 178, 
		"Ana": 165,
	}

	var nume []string
	for n := range persoane {
		nume = append(nume, n)
	}

	for i := 0; i < len(nume)-1; i++ {
		for j := i + 1; j < len(nume); j++ {
			if persoane[nume[i]] < persoane[nume[j]] {
				nume[i], nume[j] = nume[j], nume[i]
			}
		}
	}

	for _, n := range nume {
		fmt.Printf("%s: %d\n", n, persoane[n])
	}
}

func ex7(){	
	file, err := os.Open("date.txt")
	if err != nil {
		fmt.Println("Eroare la deschidere:", err)
		return
	}
	defer file.Close()

	frecventa := make(map[string]int)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		linie := scanner.Text()
		cuvinte := strings.Fields(linie)
		
		for _, cuv := range cuvinte {
			frecventa[cuv]++
		}
	}

	err = scanner.Err(); 

	if err != nil {
		fmt.Println("eroare la citire", err)
	}

	for cuvant, nr := range frecventa {
		fmt.Printf("%s: %d\n", cuvant, nr)
	}
}

func ex8(){
	angajati := map[string]int{
		"Ana Popescu": 1,
		"Mihai Ionescu": 2,
		"Bogdan Radu": 1, 
		"Elena Dumitrescu": 3,
		"Cristian Vasile": 2, 
		"Ioana Enache": 1,
	}

	departamente := make(map[int][]string)

	for nume, id := range angajati {
		departamente[id] = append(departamente[id], nume)
	}

	fmt.Print(departamente)
}

func ex9(){
	arr := []int{2, 7, 11, 8, 4, 9}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	a := scanner.Text()
	s, err := strconv.Atoi(a)

	if err != nil {
		fmt.Print("eroare")
	}

	vazute := make(map[int]int)

	for i, val := range arr {
		complement := s - val
		
		if idx, gasit := vazute[complement]; gasit {
			fmt.Printf("%d %d\n", idx, i)
		}
		vazute[val] = i
	}
}