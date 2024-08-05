package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

var memCounter = 16
var memSymbolTable map[string]string = map[string]string{
		"R0":     "0000000000000000",
		"R1":     "0000000000000001",
		"R2":     "0000000000000010",
		"R3":     "0000000000000011",
		"R4":     "0000000000000100",
		"R5":     "0000000000000101",
		"R6":     "0000000000000110",
		"R7":     "0000000000000111",
		"R8":     "0000000000001000",
		"R9":     "0000000000001001",
		"R10":    "0000000000001010",
		"R11":    "0000000000001011",
		"R12":    "0000000000001100",
		"R13":    "0000000000001101",
		"R14":    "0000000000001110",
		"R15":    "0000000000001111",
		"SP":     "0000000000000000",
		"LCL":    "0000000000000001",
		"ARG":    "0000000000000010",
		"THIS":   "0000000000000011",
		"THAT":   "0000000000000100",
		"SCREEN": "0100000000000000",
		"KBD":    "0110000000000000",
	}


var compSymbolTable map[string]string = map[string]string{
		"0":   "101010",
		"1":   "111111",
		"-1":  "111010",
		"D":   "001100",
		"A":   "110000",
		"!D":  "001101",
		"!A":  "110001",
		"-D":  "001111",
		"-A":  "110011",
		"D+1": "011111",
		"A+1": "110111",
		"D-1": "001110",
		"A-1": "110010",
		"D+A": "000010",
		"D-A": "010011",
		"A-D": "000111",
		"D&A": "000000",
		"D|A": "010101",
		"M":   "110000",
		"!M":  "110001",
		"-M":  "110011",
		"M+1": "110111",
		"M-1": "110010",
		"D+M": "000010",
		"D-M": "010011",
		"M-D": "000111",
		"D&M": "000000",
		"D|M": "010101",
}


var compTableABit map[string]string = map[string]string{
		"0":   "101010",
		"1":   "111111",
		"-1":  "111010",
		"D":   "001100",
		"A":   "110000",
		"!D":  "001101",
		"!A":  "110001",
		"-D":  "001111",
		"-A":  "110011",
		"D+1": "011111",
		"A+1": "110111",
		"D-1": "001110",
		"A-1": "110010",
		"D+A": "000010",
		"D-A": "010011",
		"A-D": "000111",
		"D&A": "000000",
		"D|A": "010101",
	}

var	destSymbolTable map[string]string = map[string]string{
		"null": "000",
		"M":    "001",
		"D":    "010",
		"MD":   "011",
		"A":    "100",
		"AM":   "101",
		"AD":   "110",
		"AMD":  "111",
	}

var jumpSymbolTable map[string]string = map[string]string{
		"null": "000",
		"JGT":  "001",
		"JEQ":  "010",
		"JGE":  "011",
		"JLT":  "100",
		"JNE":  "101",
		"JLE":  "110",
		"JMP":  "111",
	}

//muh structs, muh oop
func main(){
    fmt.Println("Starting...");
    fileName := loadArgs()
    fContents := readFile(fileName)
    fContents = passOne(fContents)
    fContents = passTwo(fContents)
    saveToFile(fContents, fileName)
}

func loadArgs() string {
    argsLen := len(os.Args)
    switch {
    case argsLen==0:
       log.Fatalln("Error, no args given. \nPlease provide the asm filename.")
    case argsLen>=3:
        log.Fatalf("Too many args given, %d:", len(os.Args))
    }
    return os.Args[1]
}

//add (symbol_text) to memSymbolTable map
func passOne(fContents []string) []string {
    newFContents := make([]string, 0)
    counter := 0
    for _,line := range fContents {
        if strings.Contains(line, "(") && strings.Contains(line, ")") {
            usrDefSymbol := line[1:len(line)-1]
            binary := intToBinary(counter)
            memSymbolTable[usrDefSymbol] = binary
            continue
        }
        counter += 1
        newFContents = append(newFContents, line)
    }
    return newFContents
}

func passTwo(fContents []string) []string {
    newFContents := make([]string, 0)
    for _,line := range fContents {
        if strings.Index(line, "@") == 0 {
            line = aCommand(line)
        } else {
            line = cCommand(line)
        }
        newFContents = append(newFContents, line)
    }
    return newFContents
}

func readFile(fileName string) []string {
    log.Printf("Reading file: %q \n", fileName)
    file, err := os.Open(fileName)
    if err != nil {
        log.Printf("Error reading file: %q\nErr: %v\n", fileName, err)
        if fNotExist := errors.Is(err, fs.ErrNotExist); fNotExist {
            log.Println("File does not exist!")
        }
        os.Exit(1)
    }

    defer file.Close()
    scanner := bufio.NewScanner(file)
    fileContents := make([]string, 0)

    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println(line)
        line = trimWhiteSpace(line)
        if isComment(line) || isWhitespace(line) {
            continue
        }
        fileContents = append(fileContents, line)
    }
    if scanner.Err() != nil {
        fmt.Printf("Error encountered while scanning file: %v", scanner.Err())
    }

    fmt.Println("Printing file contents without whitespace/comments:")
    for _,line := range fileContents {
        fmt.Println(line)
    }
    return fileContents
}

func trimWhiteSpace(line string) string {
    trimmedLine := strings.ReplaceAll(line, " ", "")
    return trimmedLine
}

func isComment(line string) bool {
    if index := strings.IndexAny(line, "/"); index == 0 {
        return true
    }
    return false
}

func isWhitespace(line string) bool {
    trimmedLine := strings.TrimSpace(line)
    return len(trimmedLine) == 0
}

//1 - predefined symbol, return map[key]
//2 - number, return number
//3 - user defined symbol, allocate mem from 16 onwards
func aCommand(line string) string {
    aCmd := line[1:]
    fmt.Println(aCmd)
    //predefined
    if value, exists := memSymbolTable[aCmd]; exists {
        fmt.Println(value)
        return rightPad(value, 16)
    } else if isNumeric(aCmd) {
        //number
        aCmdInt := atoiWrap(aCmd)
        binary := intToBinary(aCmdInt)
        memSymbolTable[aCmd] = binary
        return rightPad(binary, 16)
    } else { //user defined symbol
        binary := intToBinary(memCounter)
        memSymbolTable[aCmd] = binary
        memCounter++
        return rightPad(binary, 16)
    }
}

//dest = comp ; jump
//comp is mandatory rest are not
// = and ; can be ommited
func cCommand(line string) string  {
    padding := "111"
    nullDest := "000"
    nullJump := "000"

    compDest := strings.Split(line, "=")
    compDest_Jump := strings.Split(line, ";")

	// Only comp, no '=' or ';' use entire line, rest are null
	if len(compDest) == 1 && len(compDest_Jump) == 1 {
        comp := line
        aBit := detABit(comp)
		return padding + aBit + compSymbolTable[comp] + nullDest + nullJump
	} else if len(compDest) > 1 && len(compDest_Jump) == 1 {
        // Dest + comp, '=' exists no ';' [0] is dest [1] is comp ex:(A=M+D)
		dest := compDest[0]
        comp := compDest[1]
        aBit := detABit(comp)
		return padding + aBit + compSymbolTable[comp] + destSymbolTable[dest] + nullJump
	} else if len(compDest) == 1 && len(compDest_Jump) > 1 {
		//comp + jump, no dest
        comp := compDest_Jump[0]
		jump := compDest_Jump[1]
        aBit := detABit(comp)
		return padding + aBit + compSymbolTable[comp] + nullDest + jumpSymbolTable[jump]
    }
    return "ERROR"
}

func isNumeric(s string) bool {
	for _, char := range s {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

func detABit(compCmd string) string {
    for k := range compTableABit {
        if strings.EqualFold(compCmd, k) {
            return "0"
        }
    }
    return "1"
}

func saveToFile(fileContents []string, fileName string){
    base := filepath.Base(fileName)
    newFName := strings.TrimSuffix(base, filepath.Ext(base))
    newFPath := filepath.Join(filepath.Dir(fileName), newFName+".hack")

 	file, err := os.Create(newFPath)
	if err != nil {
        log.Fatalf("Error creating file: %q, err: %v", fileName, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range fileContents {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
            log.Fatalf("Error writing to file, fName: %q, line: %q, err: %v", fileName, line, err)
		}
	}

	// Flush any buffered data to the file
	if err := writer.Flush(); err != nil {
        log.Printf("Error flushing file, err: %v", err)
	}
}

func intToBinary(num int) string {
    fmt.Println(num)
    num64 := int64(num)
    fmt.Println(num64)
    binary := strconv.FormatInt(num64, 2)
    fmt.Println(binary)
    return binary
}

func atoiWrap(text string) int {
    convertedInt, err := strconv.Atoi(text)
    if err != nil {
        log.Fatalf("Error atoi: string %q, err %v", text, err)
    }
    return convertedInt
}

func rightPad(text string, padding int) string {
    for {
        if len(text) >= padding {
            return text
        }
        text = "0" + text
    }
}
