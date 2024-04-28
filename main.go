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
		///punctuations
		slicesStr := SliceToString(slices)
		punkSlice := strings.Split(slicesStr, " ")

		punks := ".,!?:;"

		for i := 0; i < len(punkSlice); i++ {
			for j := 0; j < len(punks); j++ {
				////check if punkSlice[i] start with punk index = 0
				if strings.Index(punkSlice[i], string(punks[j])) == 0 {
					if i-1 >= 0 {

						if len(punkSlice[i]) == 1 {
							// if there is a char only
							punkSlice[i-1] += string(punks[j])
							punkSlice = append(punkSlice[:i], punkSlice[i+1:]...)
						} else {
							punkSlice[i-1] += string(punks[j])
							punkSlice[i] = punkSlice[i][1:]
						}

						i--
					}
				}
			}
		}

		afterSingle := handleSingleQuote(SliceToString(punkSlice))

		result := ""
		found := false
		a := 0

		////handle concatenated punks
		for i := 0; i < len(afterSingle); i++ {
			for j := 0; j < len(punks); j++ {
				if afterSingle[i] == punks[j] {
					//////handle out of range i+1
					if i < len(afterSingle)-1 {
						if string(afterSingle[i+1]) != " " && checkForPunk(string(afterSingle[i+1])) {
							result += string(afterSingle[i]) + " "
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
				result += string(afterSingle[i])
			}
		}

		finalResult += result

		///if there is a \n
		if j < len(rows)-1 {
			finalResult += "\n"
		}
	}

	fmt.Println("\n")
	fmt.Println("new :", finalResult)
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
		if i+1 < len(slices)-1 {
			switch slices[i] {
			case "(cap)":
				handleLowCapUp("low", slices, i, false)
				slices = handleLowCapUp("cap", slices, i, true)
				i--
			case "(cap,":

				nbrOfWords, _ := strconv.Atoi(string(slices[i+1][:len(slices[i+1])-1]))
				// lower all then capitalize
				handlePreviousLowCapUp("(low,", slices, i, nbrOfWords, false)
				slices = handlePreviousLowCapUp("(cap,", slices, i, nbrOfWords, true)
				i -= 1

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
				// slices = append(slices, charsAfter)
				i -= 1

			case "(low)":
				slices = handleLowCapUp("low", slices, i, true)
				i--
			case "(low,":
				nbrOfWords, _ := strconv.Atoi(string(slices[i+1][:len(slices[i+1])-1]))
				slices = handlePreviousLowCapUp("(low,", slices, i, nbrOfWords, true)
				i -= 1

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

func handleBinHex(base int, mySlice []string, i int) []string {
	if i-1 >= 0 {
		mySlice[i-1] = fmt.Sprint(toDecimal(mySlice[i-1], base))
	}
	mySlice = append(mySlice[:i], mySlice[i+1:]...)
	return mySlice
}

func toDecimal(value string, base int) int64 {
	hexadecimal_num := value

	// use the parseInt() function to convert
	decimal_num, err := strconv.ParseInt(hexadecimal_num, base, 64)
	// in case of any error
	if err != nil {
		fmt.Println("- - - - \nError!!\nThe '", value, "' doesn't match the provided base (", base, ").\n- - - - ")
		return 0
	}
	return decimal_num
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

func handlePreviousLowCapUp(ToApply string, mySlice []string, i int, nbrOfWords int, shrinkSlice bool) []string {
	// apply the rule

	if i-1 >= 0 {
		if nbrOfWords < 0 {
			fmt.Println("- - - - \nError!!\nYou entered a negative number (", nbrOfWords, ").\n- - - - ")
		} else {
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
	}
	if shrinkSlice {
		// remove the (rule)
		mySlice = append(mySlice[:i], mySlice[i+2:]...)
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

	openQuote := false
	openIndex := -1
	closeQuote := false

	for i := 0; i < len(mySlice); i++ {
		if i < len(mySlice)-1 && mySlice[i] == "'" && !openQuote && mySlice[i+1] != "'" {
			openQuote = true
			openIndex = i
		} else if mySlice[i] == "'" && openQuote && !closeQuote {
			closeQuote = true
			mySlice[openIndex+1] = mySlice[openIndex] + mySlice[openIndex+1]
			mySlice[i-1] += mySlice[i]
			mySlice[openIndex] = ""
			mySlice[i] = ""
			fmt.Println("len++ : ", len(mySlice))

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
				if string(mySlice[i-1]) == "A" || string(mySlice[i-1]) == "a" {
					mySlice[i-1] += "n"
				}
				if len(mySlice[i-1]) == 2 {
					if string(mySlice[i-1][1]) == "a" {
						if string(mySlice[i-1][0]) == "(" || string(mySlice[i-1][0]) == "\"" || string(mySlice[i-1][0]) == "'" {
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
	return strings.Trim(strings.Join(mySlice, " "), " ")
}
