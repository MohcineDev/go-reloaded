package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	Reloaded()
}

// purge the input
func removeSpaces(oldString string) string {
	regRule := regexp.MustCompile(`\s+`)
	return regRule.ReplaceAllString(oldString, " ")
}

func Reloaded() {
	myArgs := os.Args[1:]
	if len(myArgs) < 2 {
		fmt.Println("Error!! please enter sample.txt and result file!!!")
		return
	}
	dat, err := os.ReadFile(myArgs[0])
	if err != nil {
		fmt.Println("Error!!", myArgs[0], "file not found !!!")
		return

	}
	///if no result file is provided

	if len(myArgs) > 2 {
		fmt.Println("Error!! you entered more than 2 args!!!")
	}
	if path.Ext(myArgs[1]) != ".txt" {
		fmt.Println("Error!! result file Extension must be (.txt) !!!")
		return
	}

	///split the rows by \n
	rows := strings.Split(string(dat), "\n")

	fmt.Println("old : ", string(dat))

	finalResult := ""

	for j := 0; j < len(rows); j++ {

		newString := removeSpaces(rows[j])

		///splitting each row by " "
		slices := strings.Fields(newString)

		checkVowels(slices)
		slices = generalFlags(slices)
		checkVowels(slices)

		punks := ".,!?:;"

		for i := 0; i < len(slices); i++ {
			for j := 0; j < len(punks); j++ {
				////check if slices[i] start with punk index = 0
				if strings.Index(slices[i], string(punks[j])) == 0 {
					if i-1 >= 0 {

						if len(slices[i]) == 1 {
							// if there is a char only
							slices[i-1] += string(punks[j])
							slices = append(slices[:i], slices[i+1:]...)
						} else {
							slices[i-1] += string(punks[j])
							slices[i] = slices[i][1:]
						}

						i--
					}
				}
			}
		}

		punkString := SliceToString(slices)

		result := ""
		found := false
		a := 0

		////handle concatenated punks
		for i := 0; i < len(punkString); i++ {
			for j := 0; j < len(punks); j++ {
				if punkString[i] == punks[j] {
					//////handle out of range i+1
					if i < len(punkString)-1 {
						if string(punkString[i+1]) != " " && checkForPunk(string(punkString[i+1])) {
							result += string(punkString[i]) + " "
							found = true
							a = 1
						}
					}
				}
			}
			if a > 0 {
				a = 0
				found = false
				continue
			}
			if !found {
				result += string(punkString[i])
			}
		}

		finalResult += handleSingleQuote(result)

		///if there is a \n
		if j < len(rows)-1 {
			finalResult += "\n"
		}
	}

	data := []byte(finalResult)

	// create file if not exist and save it
	err = os.WriteFile(myArgs[1], data, 0o644)
	if err != nil {
		panic(err)
	}
}

func checkForPunk(char string) bool {
	punks := ".,!?:;"
	nextNotPunk := true

	for j := 0; j < len(punks); j++ {
		if char == string(punks[j]) {
			nextNotPunk = false
		}
	}
	return nextNotPunk
}

func generalFlags(slices []string) []string {
	for i := 0; i < len(slices); i++ {
		if i < len(slices) {
			switch slices[i] {
			case "(cap)":
				handleLowCapUp("low", slices, i, false)
				slices = handleLowCapUp("cap", slices, i, true)
				i--
			case "(cap,":
				if i+1 < len(slices) && isRightFlag(slices[i+1]) {

					nbrOfWords, _ := strconv.Atoi(string(slices[i+1][:len(slices[i+1])-1]))
					// lower all then capitalize
					handlePreviousLowCapUp("(low,", slices, i, nbrOfWords)
					slices = handlePreviousLowCapUp("(cap,", slices, i, nbrOfWords)

					slices = append(slices[:i], slices[i+2:]...)
					i -= 1

				}
			case "(up)":
				slices = handleLowCapUp("up", slices, i, true)
				i--
			case "(up,":
				if i+1 < len(slices) && isRightFlag(slices[i+1]) {
					nbrOfWords, _ := strconv.Atoi(string(slices[i+1][:len(slices[i+1])-1]))

					slices = handlePreviousLowCapUp("(up,", slices, i, nbrOfWords)
					slices = append(slices[:i], slices[i+2:]...)
					i -= 1

				}
			case "(low)":
				///fmt.Println(slices)
				slices = handleLowCapUp("low", slices, i, true)
				i--
			case "(low,":
				if i+1 < len(slices) && isRightFlag(slices[i+1]) {
					nbrOfWords, _ := strconv.Atoi(string(slices[i+1][:len(slices[i+1])-1]))
					slices = handlePreviousLowCapUp("(low,", slices, i, nbrOfWords)

					slices = append(slices[:i], slices[i+2:]...)
					i -= 1
				}

			case "(hex)":
				slices = handleBinHex(16, slices, i)
				i--

			case "(bin)":
				slices = handleBinHex(2, slices, i)
				i--
			}
		}
	}
	return slices
}

func isRightFlag(halfFlag string) bool {
	if string(halfFlag[len(halfFlag)-1]) == ")" {
		_, err := strconv.Atoi(string(halfFlag[:len(halfFlag)-1]))
		if err == nil {
			return true
		}
	}

	return false
}

func handleBinHex(base int, mySlice []string, i int) []string {
	if i-1 >= 0 {
		mySlice[i-1] = fmt.Sprint(toDecimal(mySlice[i-1], base))
	}
	mySlice = append(mySlice[:i], mySlice[i+1:]...)
	return mySlice
}

func toDecimal(value string, base int) string {
	// use the parseInt() function to convert
	decimal_num, err := strconv.ParseInt(value, base, 64)
	// in case of any error
	if err != nil {
		fmt.Println("- - - - \nError!!\nThe '", value, "' doesn't match the provided base (", base, ").\n- - - - ")
		return value
	}
	// fmt.Println("0000 : ", decimal_num)
	return strconv.Itoa(int(decimal_num))
}

func handleLowCapUp(ToApply string, mySlice []string, i int, shrinkSlice bool) []string {
	if i-1 >= 0 {
		if ToApply == "cap" {
			mySlice[i-1] = strings.Title(mySlice[i-1])
			// don't , aren't
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

func handlePreviousLowCapUp(ToApply string, mySlice []string, i int, nbrOfWords int) []string {
	// apply the rule

	if i-1 >= 0 {
		if nbrOfWords < 0 {
			fmt.Println("- - - - \nError!!\nYou entered a negative number (", nbrOfWords, ").\n- - - - ")
		} else {
			for j := 0; j < nbrOfWords; j++ {
				if i-1-j >= 0 {
					if ToApply == "(cap," {
						mySlice[i-1-j] = strings.Title(mySlice[i-1-j])

						if strings.Contains(mySlice[i-1-j], "'") {
							wordSlice := strings.Split(mySlice[i-1-j], "'")
							wordSlice[1] = strings.ToLower(wordSlice[1])
							mySlice[i-1-j] = wordSlice[0] + "'" + wordSlice[1]
						}
					} else if ToApply == "(up," {
						mySlice[i-1-j] = strings.ToUpper(mySlice[i-1-j])
					} else if ToApply == "(low," {
						mySlice[i-1-j] = strings.ToLower(mySlice[i-1-j])
					}
				}
			}
		}
	}

	return mySlice
}

// /string to rune
func stringToRune(s string) string {
	myRune := []rune(s)
	new := []rune{}

	for i, v := range myRune {
		///check if v is a single quote
		if v == '\'' && i-1 >= 0 && i+1 < len(myRune) && unicode.IsLetter(myRune[i-1]) && unicode.IsLetter(myRune[i+1]) {
			new = append(new, v)
		} else if v == '\'' {
			new = append(new, ' ')
			new = append(new, v)
			new = append(new, ' ')
		} else {
			new = append(new, v)
		}
	}
	return string(new)
}

func handleSingleQuote(value string) string {
	myStr := stringToRune(value)
	mySlice := strings.Fields(myStr)

	////// check for a ' if it's in the beginning of a word
	////////// YES :  check if there is a second one
	////////////////// YES : check if it's at the end of a word
	//////////////////////// YES : leave it
	/////////////////////// NO : add it to the end
	/////////////////  NO :
	/////////////// YES : add it to the end of the word before it
	/////////////// NO : add space on each side of it if it's not can't don't i-1 !=n and i+1 != t
	////////// NO : add space on each side of it if it's not can't don't i-1 !=n and i+1 != t
	// gf g '  ' fd don't gfsdfg  ' rt
	openQuote := false
	openIndex := -1
	closeQuote := false

	for i := 0; i < len(mySlice); i++ {
		//
		if i < len(mySlice)-1 && mySlice[i] == "'" && !openQuote && mySlice[i+1] != "'" {
			openQuote = true
			openIndex = i
		} else if mySlice[i] == "'" && openQuote && !closeQuote {
			closeQuote = true
			mySlice[openIndex+1] = mySlice[openIndex] + mySlice[openIndex+1]
			mySlice[i-1] += mySlice[i]
			mySlice[openIndex] = ""
			mySlice[i] = ""
			// fmt.Println("len++ : ", len(mySlice))

			openQuote = false
			closeQuote = false
		}
	}

	return SliceToString(mySlice)
}

func checkVowels(mySlice []string) {
	vowels := "aeoiuhAEOIUH"

	for i := 0; i < len(mySlice); i++ {
		for j := 0; j < len(vowels); j++ {
			// check if the first char is a vowel
			if i > 0 && mySlice[i][0] == vowels[j] {
				//fmt.Println(mySlice[i-1])

				if string(mySlice[i-1]) == "A" || string(mySlice[i-1]) == "a" {
					mySlice[i-1] += "n"
				}
				if len(mySlice[i-1]) == 2 {
					if string(mySlice[i-1][1]) == "a" {
						if string(mySlice[i-1][0]) == "(" || string(mySlice[i-1][0]) == "\"" || string(mySlice[i-1][0]) == "'" || string(mySlice[i-1][0]) == "[" || string(mySlice[i-1][0]) == "{" {
							mySlice[i-1] += "n"
						}
					}
				}

			}
		}
	}
}

// /convert slice to string
func SliceToString(mySlice []string) string {
	return removeSpaces(strings.Trim(strings.Join(mySlice, " "), " "))
}
