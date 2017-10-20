package cards_service

import (
	fmt "fmt"
	"strings"
)

func (c *Card) Display() string {
	w := len(c.GetTitle()) + 2

	if w < 20 {
		w = 20
	}

	sep := strings.Repeat("-", w)
	p := fmt.Sprintf("\n+%s+\n", sep)
	p += fmt.Sprintf("| Card: %[2]*[1]d |\n", c.Id, w-8)
	p += fmt.Sprintf("+%s+\n", sep)
	p += fmt.Sprintf("| %-[2]*[1]s |\n", c.Title, w-2)
	p += fmt.Sprintf("+%s+\n", sep)

	if len(c.Cards) > 0 {

		p += fmt.Sprintf("| includes: %[2]*[1]s |\n", "", w-12)

		for _, i := range c.Cards {
			p += fmt.Sprintf("| - %-[2]*[1]d |\n", i.Id, w-4)
		}

		p += fmt.Sprintf("+%s+\n", sep)
	}

	return p
}
