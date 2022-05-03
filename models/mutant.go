package models

import (
	"fmt"
	"strings"
)

type Mutante struct {
	adnMutant    [][]string
	adnCadena    []string
	isMutantFlat bool
	iterLetter   string
	function     string
}

/**
 * @param adn []string array de string con cadenas de adn
 */
func (mutant *Mutante) insertDna(adn []string) {
	for _, value := range adn {
		mutant.adnMutant = append(mutant.adnMutant, strings.Split(value, ""))
	}
}

/**
  Metodo que recorre el adnMutant de manera horizontal y verifica si es mutante
*/
func (mutant *Mutante) validateHorizontal() {
	countIter := 0
	for i := 0; i < len(mutant.adnMutant); i++ {
		if mutant.iterLetter != "" && countIter < 4 {
			mutant.iterLetter = ""
			countIter = 0
		}
		for j := 0; j <= len(mutant.adnMutant[i])-1; j++ {
			if j < len(mutant.adnMutant[j])-1 {
				if mutant.iterLetter == "" {
					if mutant.adnMutant[i][j] == mutant.adnMutant[i][j+1] {
						mutant.iterLetter = mutant.adnMutant[i][j]
						countIter++
					}
				} else {
					if mutant.iterLetter == mutant.adnMutant[i][j] {
						countIter++
					}
				}
				if countIter == 4 {
					mutant.isMutantFlat = true
					mutant.function = "horizontal"
					break
				}
			}
		}
		if mutant.isMutantFlat {
			break
		}
	}
}

/**
  Metodo que recorre el adnMutant de manera vertical y verifica si es mutante
*/
func (mutant *Mutante) validateVertical() {
	if !mutant.isMutantFlat {
		countIter := 0
		for i := 0; i < len(mutant.adnMutant); i++ {
			for j := 0; j < len(mutant.adnMutant[i]); j++ {
				if j < len(mutant.adnMutant[i])-1 {
					if mutant.iterLetter == "" {
						if mutant.adnMutant[j][i] == mutant.adnMutant[j+1][i] {
							mutant.iterLetter = mutant.adnMutant[j][i]
							countIter++
						}
					} else {
						if mutant.iterLetter == mutant.adnMutant[j][i] {
							countIter++
						}
					}
					if countIter == 4 {
						mutant.isMutantFlat = true
						mutant.function = "vertical"
						break
					}
				}
			}
			if mutant.isMutantFlat {
				break
			}
		}
	}
}

/**
  Metodo que recorre el adnMutant de manera horizontal y verifica si es mutante
*/
func (mutant *Mutante) validateDiagonal() {
	if !mutant.isMutantFlat {
		if len(mutant.adnMutant)%2 == 0 {
			countIter := 0
			for i := 0; i < len(mutant.adnMutant); i++ {
				if i < len(mutant.adnMutant)-1 && i < len(mutant.adnMutant[i])-1 {
					mutant.iterLetter = mutant.adnMutant[i][i]
					if mutant.iterLetter == mutant.adnMutant[i+1][i+1] {
						mutant.iterLetter = mutant.adnMutant[i][i]
						countIter++
					}
					if countIter == 3 {
						mutant.isMutantFlat = true
						mutant.function = "diag-izq/der"
						break
					}
				}
			}
			countIter = 0
			if !mutant.isMutantFlat {
				aux := 0
				for i := len(mutant.adnMutant) - 1; i > 0; i-- {
					mutant.iterLetter = mutant.adnMutant[aux][i]
					if mutant.iterLetter == mutant.adnMutant[aux+1][i-1] {
						mutant.iterLetter = mutant.adnMutant[aux][i]
						countIter++
					}
					if countIter == 4 {
						mutant.isMutantFlat = true
						mutant.function = "diag-der/izq"
						break
					}
					aux++
				}
			}
		}
	}
}

/**
  Metodo que imprime información del mutante
*/
func (mutant *Mutante) writeData() {
	fmt.Println("ADN: ", mutant.adnMutant)
	fmt.Println("isMutant: ", mutant.isMutantFlat)
	fmt.Println("iterLetter: ", mutant.iterLetter, " function: ", mutant.function)
}

/**
 *
 * @param []dna array de string con cadenas de adn
 * @returns bool si es mutante, true de lo contrario false
 */
func IsMutant(dna []string) bool {
	mutant := Mutante{}
	mutant.adnCadena = dna
	mutant.insertDna(dna)
	mutant.validateHorizontal()
	mutant.validateVertical()
	mutant.validateDiagonal()
	mutant.writeData()
	SaveDna(mutant.adnCadena, mutant.isMutantFlat)
	return mutant.isMutantFlat
}
