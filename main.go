package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	Reloaded()
}
n
func check(e error) {
	if e != nil {
		return
	}
}

//purge the input
func removeSpaces(oldString string) string {

	regRule := regexp.MustCompile(`\s+`)
	return regRule.ReplaceAllString(oldString, " ")

}

func Reloaded() {
	myArgs := os.Args[1:]

	dat, err := os.ReadFile(myArgs[0])

	///if no result file is provided
	if len(myArgs) < 2 {
		fmt.Println("please enter result file!!!")
		return
	}
	if len(myArgs) > 2 {
		fmt.Println("you entered more than 2 args!!!")
	}

	check(err)

	///split the rows by \n
	rows := strings.Split(string(dat), "\n")

	fmt.Println("old : ", string(dat))

	finalResult := ""

	for j := 0; j < len(rows); j++ {

		newString := removeSpaces(rows[j])

		///splitting each row by " "
		slices := strings.Split(newString, " ")
		handleFlags(rows[j])

		result := ""

		checkVowels(slices)

		slices = generalCapUpLow(slices)
		fmt.Println("slices : ", slices)
		///punctuations
		slicesStr := SliceToString(slices)
		punkSlice := strings.Split(slicesStr, " ")
		punks := ".,!?:;"

		for i := 0; i < len(punkSlice); i++ {
			for j := 0; j < len(punks); j++ {
				if strings.Index(punkSlice[i], string(punks[j])) == 0 {

					if i == len(punkSlice)-1 {

						punkSlice[i-1] += punkSlice[i]
						punkSlice = punkSlice[:i]

					} else {

						punkSlice[i-1] = punkSlice[i-1] + string(punks[j])
						punkSlice[i] = punkSlice[i][1:]
					}
					i--
				}
			}
		}
		result = SliceToString(punkSlice)
		myslice := strings.Split(string(result), " ")

		handlePunctuationMark(myslice)
		finalResult += SliceToString(myslice)
		handleFlags("(cap)")

		///if there is a \n
		if j < len(rows)-1 {
			finalResult += "\n"
		}
	}

	fmt.Println(finalResult)
	data := []byte(finalResult)

	//create file if not exist and save it
	err = os.WriteFile(myArgs[1], data, 0644)
	if err != nil {
		panic(err)
	}
}

func generalCapUpLow(slices []string) []string {
	for i := 0; i < len(slices); i++ {

		switch slices[i] {
		case "(cap)":
			handleLowCapUp("low", slices, i, false)
			slices = handleLowCapUp("cap", slices, i, true)
			i--
		case "(cap,":
			nbrOfWords, _ := strconv.Atoi(string(slices[i+1][:len(slices[i+1])-1]))
			//lower all then capitalize
			handlePreviousLowCapUp("(low,", slices, i, nbrOfWords, false)
			slices = handlePreviousLowCapUp("(cap,", slices, i, nbrOfWords, true)
			i -= 2

		case "(up)":
			slices = handleLowCapUp("up", slices, i, true)
			i--
		case "(up,":
			// num := strings.Index(slices[i+1], ")")
			// numOfWords := slices[i+1][0:num]
			// charsAfter := slices[i+1][num+1:]
			// nbrOfWords, _ := strconv.Atoi(string(numOfWords))
			nbrOfWords, _ := strconv.Atoi(string(slices[i+1][:len(slices[i+1])-1]))

			slices = handlePreviousLowCapUp("(up,", slices, i, nbrOfWords, true)
			//slices = append(slices, charsAfter)
			i -= 2

		case "(low)":
			slices = handleLowCapUp("low", slices, i, true)
			i--
		case "(low,":
			nbrOfWords, _ := strconv.Atoi(string(slices[i+1][:len(slices[i+1])-1]))
			slices = handlePreviousLowCapUp("(low,", slices, i, nbrOfWords, true)
			i -= 2

		case "(hex)":
			if toDecimal(slices[i-1], 16) == 0 {
				slices[i] = ""
			} else {
				slices[i-1] = fmt.Sprint(toDecimal(slices[i-1], 16))
				slices[i] = ""
			}
		case "(bin)":
			if toDecimal(slices[i-1], 2) == 0 {

				slices[i] = ""
			} else {
				slices[i-1] = fmt.Sprint(toDecimal(slices[i-1], 2))
				slices[i] = ""
			}
		}
	}
	return slices
}

//get all flags
func handleFlags(slice string) {
	regRule := regexp.MustCompile(`\(cap*,?\s?\d*\)|\(up,?\s?\d?\)`)

	///find all instances
	// for i := 0; i < len(slice); i++ {

	flags := regRule.ReplaceAllString(slice, ".....")

	fmt.Println("flags : ", flags)
	// }
	// flags := regRule.FindAllString(sentence, -1)
	// flagsIndex := regRule.FindAllStringIndex(sentence, -1)
	/*
	   take all flags ---> put them in a slice
	   split the string by space
	   check every slice[i] if it equal the flag
	*/
	fmt.Println("reg : ")
	// fmt.Println("reg : ", 1, flags[1])
	// fmt.Println("reg : ", flagsIndex)

	///regRule.ReplaceAllString(oldString, " ")
}

func toDecimal(value string, base int) int64 {
	hexadecimal_num := value

	// use the parseInt() function to convert
	decimal_num, err := strconv.ParseInt(hexadecimal_num, base, 64)
	// in case of any error

	if err != nil {
		fmt.Println("- - - - \nError!!\nThe", value, "doesn't match the provided base (", base, ").\n- - - - ")
		return 0
		//	os.Exit(1)
	}
	return decimal_num
}

func handleLowCapUp(ToApply string, mySlice []string, i int, shrinkSlice bool) []string {

	if i-1 >= 0 {

		if ToApply == "cap" {
			mySlice[i-1] = strings.Title(mySlice[i-1])
			//don't , aren't
			if strings.Contains(mySlice[i-1], "'") {
				wordSlice := strings.Split(mySlice[i-1], "'")
				wordSlice[1] = strings.ToLower(wordSlice[1])
				mySlice[i-1] = wordSlice[0] + "'" + wordSlice[1]
			}
		} else if ToApply == "up" {
			mySlice[i-1] = strings.ToUpper(mySlice[i-1])
		} else if ToApply == "low" {
			mySlice[i-1] = strings.ToLower(mySlice[i-1])
		}
	}
	if shrinkSlice {
		mySlice = append(mySlice[:i], mySlice[i+1:]...)
	}

	return mySlice
}

func handlePreviousLowCapUp(ToApply string, mySlice []string, i int, nbrOfWords int, shrinkSlice bool) []string {
	// apply the rule
	fmt.Println("rrrrrrrrrrr")
	if i == 0 { 
	}
	if i-1 >= 0 {

		for j := 0; j < nbrOfWords; j++ {
			if i-1-j >= 0 {
				if ToApply == "(cap," {
					if strings.Contains(mySlice[i-1], "'") {
						wordSlice := strings.Split(mySlice[i-1], "'")
						wordSlice[1] = strings.ToLower(wordSlice[1])
						mySlice[i-1] = wordSlice[0] + "'" + wordSlice[1]
					}
					mySlice[i-1-j] = strings.Title(mySlice[i-1-j])

				} else if ToApply == "(up," {
					mySlice[i-1-j] = strings.ToUpper(mySlice[i-1-j])

				} else if ToApply == "(low," {
					mySlice[i-1-j] = strings.ToLower(mySlice[i-1-j])
				}
			}
		}
	}
	if shrinkSlice {
		// remove the (rule)
		mySlice = append(mySlice[:i], mySlice[i+2:]...)
	}
	return mySlice
}

func handlePunctuationMark(value []string) []string {
	first := 0
	for i := 0; i < len(value); i++ {

		if value[i] == "'" && first == 0 {
			value[i+1] = "'" + value[i+1]
			value[i] = ""
			first = 1
			continue
		}
		if value[i] == "'" && first == 1 {
			value[i-1] += "'"
			value[i] = ""
			first = 0
		}
	}
	return value

}

func checkVowels(mySlice []string) {
	vowels := "aeoiuhAEOIUH"
	for i := 0; i < len(mySlice); i++ {
		for j := 0; j < len(vowels); j++ {

			if string(mySlice[i]) != "" && i > 0 && mySlice[i][0] == vowels[j] {
				if string(mySlice[i-1]) == "A" || string(mySlice[i-1]) == "a" {

					mySlice[i-1] += "n"
				}
			}
		}
	}
}

///convert slice to string
func SliceToString(mySlice []string) string {
	myString := ""
	for i := 0; i < len(mySlice); i++ {

		if mySlice[i] != "" {
			myString += mySlice[i]
		}
		if i == len(mySlice)-1 {
			break
		}
		//if i != len -1 && mySlice[i] != ""
		if mySlice[i] != "" && mySlice[i] != " " {
			myString += " "
		}
	}
	return myString
}
