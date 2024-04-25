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

	rows := strings.Split(string(dat), "\n")

	fmt.Println("old : ", string(dat))
	///split the rows by \n
	///splitting each row by " "

	result2 := ""

	for j := 0; j < len(rows); j++ {

		newString := removeSpaces(rows[j])

		// fmt.Println("\nnew : ", string(newString))

		slices := strings.Split(newString, " ")
		//		handleFlags(rows[j])

		result := ""

		checkVowels(slices)

		for i := 0; i < len(slices); i++ {

			fmt.Println(slices[i])
			switch slices[i] {
			case "(cap)":
				handleLowCapUp("low", slices, i)
				slices = handleLowCapUp("cap", slices, i)

			case "(cap,":

				nbrOfWords, _ := strconv.Atoi(string(slices[i+1][:len(slices[i+1])-1]))
				handlePreviousLowCapUp("(low,", slices, i, nbrOfWords)
				slices = handlePreviousLowCapUp("(cap,", slices, i, nbrOfWords)
			case "(up)":
				slices = handleLowCapUp("up", slices, i)

			case "(up,":
				num := strings.Index(slices[i+1], ")")
				numOfWords := slices[i+1][0:num]
				charsAfter := slices[i+1][num+1:]
				nbrOfWords, _ := strconv.Atoi(string(numOfWords))

				slices = handlePreviousLowCapUp("(up,", slices, i, nbrOfWords)
				slices = append(slices, charsAfter)

			case "(low)":
				slices = handleLowCapUp("low", slices, i)

			case "(low,":
				nbrOfWords, _ := strconv.Atoi(string(slices[i+1][:len(slices[i+1])-1]))
				slices = handlePreviousLowCapUp("(low,", slices, i, nbrOfWords)

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

		slicesStr := SliceToString(slices)
		punkSlice := strings.Split(slicesStr, " ")

		for i := 0; i < len(punkSlice); i++ {

			for strings.Index(punkSlice[i], ".") == 0 || strings.Index(punkSlice[i], ",") == 0 ||
				strings.Index(punkSlice[i], "!") == 0 || strings.Index(punkSlice[i], "?") == 0 ||
				strings.Index(punkSlice[i], ":") == 0 || strings.Index(punkSlice[i], ";") == 0 {

				punkSlice = handlePunctuations(punkSlice, i)
			}
		}
		result = SliceToString(punkSlice)
		myslice := strings.Split(string(result), " ")

		handlePunctuationMark(myslice)
		result2 += SliceToString(myslice)

		///if there is a \n
		if j < len(rows)-1 {
			result2 += "\n"
		}

	}

	fmt.Println("\nresult : ", result2)
	data := []byte(result2)
	//create file if not exist and save it
	err = os.WriteFile(myArgs[1], data, 0644)
	if err != nil {
		panic(err)
	}
}

//get all flags
func handleFlags(sentence string) {
	regRule := regexp.MustCompile(`\(cap*,?\s?\d*\)|\(up,?\s?\d?\)`)

	///find all instances

	flags := regRule.FindAllString(sentence, -1)
	fmt.Println("reg : ", flags)
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

func handleLowCapUp(ToApply string, mySlice []string, i int) []string {
	// remove the (rule)
	//mySlice[i] = ""
	if i-1 >= 0 {
		mySlice[i] = ""
		//word := mySlice[i-1]
		//mySlice = append(mySlice[:i], mySlice[i+1:]...)
		// apply the rule
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
	return mySlice
}

func handlePreviousLowCapUp(ToApply string, mySlice []string, i int, nbrOfWords int) []string {
	// remove the (rule)

	// mySlice = append(mySlice[:i], mySlice[i+1:]...)
	// mySlice = append(mySlice[:i], mySlice[i+1:]...)

	mySlice[i] = ""
	mySlice[i+1] = ""
	// apply the rule
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
func handlePunctuations(slices []string, i int) []string {
	if strings.Index(slices[i], ".") == 0 {
		slices[i-1] = slices[i-1] + "."
		slices[i] = slices[i][1:]
	} else if strings.Index(slices[i], "!") == 0 {
		slices[i-1] = slices[i-1] + "!"
		slices[i] = slices[i][1:]
	} else if strings.Index(slices[i], ",") == 0 {
		slices[i-1] = slices[i-1] + ","
		slices[i] = slices[i][1:]
	} else if strings.Index(slices[i], "?") == 0 {
		slices[i-1] = slices[i-1] + "?"
		slices[i] = slices[i][1:]
	} else if strings.Index(slices[i], ":") == 0 {
		slices[i-1] = slices[i-1] + ":"
		slices[i] = slices[i][1:]
	} else if strings.Index(slices[i], ";") == 0 {
		slices[i-1] = slices[i-1] + ";"
		slices[i] = slices[i][1:]
	}
	mm := SliceToString(slices)
	myStr := strings.Split(mm, " ")

	return myStr
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
