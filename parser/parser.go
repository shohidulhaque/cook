package parser

import (
	"errors"
	"strconv"
	"strings"

	lg "github.com/hellozee/cook/logger"
)

//compiler  Data Structure to hold the details for the compiler
type compiler struct {
	Binary     string
	Name       string
	Start      string
	LdFlags    string
	Includes   string
	OtherFlags string
}

//params  Data Structure to hold the details for the file entity
type params struct {
	File string
	Deps []string
}

//Parser  Data Structure for holding parser details and parsed object
type Parser struct {
	input           []item
	pos             int
	currentItem     item
	prevItem        item
	nextItem        item
	CompilerDetails compiler
	FileDetails     map[string]params
	Imports         []string
	Logger          *lg.Logger
}

//next  Function to shift to the next item in the list
func (par *Parser) next() item {
	par.prevItem = par.currentItem
	par.currentItem = par.nextItem
	par.pos++
	par.nextItem = par.input[par.pos]
	return par.currentItem
}

//Parse  Function to parse the the Recipe File
func (par *Parser) Parse() error {

	for par.nextItem.typ != itemEOF {
		par.next()
		//Parsing the imports
		if par.currentItem.typ == itemImport {
			par.next()
			if par.currentItem.typ != itemDoubleQuotes {
				par.reportError("\"")
			}
			par.next()
			if par.currentItem.typ != itemString {
				par.reportError("folder name")
			}
			par.Imports = append(par.Imports, par.currentItem.val)
			par.next()
			par.next()
			if par.currentItem.typ != itemSemicolon {
				par.reportError("; or new line")
			}
		}
		//Parsing the entity
		if par.currentItem.typ == itemEntity {
			par.parseEntity()
		}
	}
	par.Logger.ReportSuccess("Successfully parsed Recipe file")
	return nil
}

//reportError  Function for reporting syntax errors
func (par *Parser) reportError(expected string) error {
	syntaxError := errors.New("Syntax error on line " + strconv.Itoa(par.currentItem.line) +
		": Expected " + expected + " , found " + par.nextItem.val)
	par.Logger.ReportError(syntaxError.Error())
	return syntaxError
}

//ParseEntity  Function for parsing an entity
func (par *Parser) parseEntity() {
	par.next()
	if par.currentItem.typ != itemString {
		par.reportError("entity name")
	}
	name := par.currentItem.val
	isCompiler := false
	identifier := itemNULL
	params := ""
	if name == "#" {
		isCompiler = true
	}
	par.next()
	if par.currentItem.typ != itemLeftBrace {
		par.reportError("{")
	}

	for par.next().typ != itemRightBrace {
		if par.currentItem.typ > itemKeyWord {
			identifier = par.currentItem.typ
			params = ""
		}

		if par.currentItem.typ == itemString {
			params = par.currentItem.val
		}

		if par.currentItem.typ == itemSemicolon {
			if isCompiler == true {
				par.fillCompilerDetails(identifier, params)
			} else {
				par.fillFileDetails(name, identifier, params)
			}
		}
	}
}

//fillCompilerDetails  Function to store the compiler details
func (par *Parser) fillCompilerDetails(identifier itemType, param string) {
	if identifier == itemBinary {
		par.CompilerDetails.Binary = param
	}
	if identifier == itemName {
		par.CompilerDetails.Name = param
	}
	if identifier == itemStart {
		par.CompilerDetails.Start = param
	}
	if identifier == itemLdFlags {
		par.CompilerDetails.LdFlags = param
	}
	if identifier == itemIncludes {
		par.CompilerDetails.Includes = param
	}
	if identifier == itemOthers {
		par.CompilerDetails.OtherFlags = param
	}
}

//fillFileDetails  Function to fill the file details
func (par *Parser) fillFileDetails(name string, identifier itemType, param string) {
	var temp params

	if identifier == itemFile {
		temp.File = param
	} else if param != "" {
		temp = par.FileDetails[name]
	}

	if param == "" {
		return
	}

	if identifier == itemDeps {
		paramArray := strings.Split(param, " ")
		temp.Deps = paramArray
	}

	par.FileDetails[name] = temp
}

//NewParser  Function to help create a parser
func NewParser(file string, log *lg.Logger) Parser {
	lex := newLexer(file)
	lex.analyze()
	par := Parser{
		input:       lex.items,
		pos:         0,
		nextItem:    lex.items[0],
		FileDetails: make(map[string]params),
		Logger:      log,
	}
	return par
}
