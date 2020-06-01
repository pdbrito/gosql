package main

import (
	"bufio"
	"fmt"
	"github.com/pdbrito/gosql"
	"log"
	"os"
	"strings"
)

func main() {
	mb := gosql.NewMemoryBackend()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to gosql.")
	for {
		fmt.Print("# ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		text = strings.Replace(text, "\n", "", -1)
		ast, err := gosql.Parse(text)
		if err != nil {
			log.Panic(err)
		}

		for _, stmt := range ast.Statements {
			switch stmt.Kind {
			case gosql.CreateTableKind:
				err := mb.CreateTable(stmt.CreateTableStatement)
				if err != nil {
					log.Panic(err)
				}
				fmt.Println("ok")
			case gosql.InsertKind:
				err = mb.Insert(stmt.InsertStatement)
				if err != nil {
					log.Panic(err)
				}
				fmt.Println("ok")
			case gosql.SelectKind:
				results, err := mb.Select(stmt.SelectStatement)
				if err != nil {
					log.Panic(err)
				}
				for _, column := range results.Columns {
					fmt.Printf("| %s ", column.Name)
				}
				fmt.Println("|")

				for i := 0; i < 20; i++ {
					fmt.Printf("=")
				}
				fmt.Println()

				for _, result := range results.Rows {
					fmt.Printf("|")
					for i, cell := range result {
						typ := results.Columns[i].Type
						s := ""
						switch typ {
						case gosql.IntType:
							s = fmt.Sprintf("%d", cell.AsInt())
						case gosql.TextType:
							s = cell.AsText()
						}

						fmt.Printf(" %s | ", s)
					}

					fmt.Println()
				}

				fmt.Println("ok")
			}
		}
	}
}
