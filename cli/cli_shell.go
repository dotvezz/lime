package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// shell launches the cli as an interactive shell
func (cli cli) shell() {
	fmt.Print("entering shell mode")
	if len(*cli.name) > 1 {
		fmt.Println(" for", *cli.name)
	} else {
		fmt.Println()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s ", *cli.prompt)
		scanner.Scan()
		input := scanner.Text()
		if string(input) == *cli.exitWord {
			break
		}
		args := strings.Split(string(input), " ")
		c, depth, err := match(*cli.commands, args, 0)
		if err != nil {
			if err != errNoMatch || len(input) > 0 {
				fmt.Println(err)
			}
			continue
		}
		err = exec(c, depth, args)
		if err != nil {
			fmt.Println(err)
		}
	}
}
