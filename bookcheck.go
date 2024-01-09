package bookcheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type Book struct {
	Name        string
	Done        bool
	Added       time.Time
	CompletedAt time.Time
}

type Books []Book

func (b *Books) Add(name string) {
	book := Book{
		Name:        name,
		Done:        false,
		Added:       time.Now(),
		CompletedAt: time.Time{},
	}

	*b = append(*b, book)
}

func (b *Books) Complete(index int) error {
	ls := *b
	if index <= 0 || index > len(ls) {
		return errors.New("invalid index")
	}

	ls[index-1].CompletedAt = time.Now()
	ls[index-1].Done = true

	return nil
}

func (b *Books) Load(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return nil
	}

	if len(file) == 0 {
		return nil
	}

	err = json.Unmarshal(file, b)

	if err != nil {
		return nil
	}

	return nil
}

func (b *Books) Store(filename string) error {
	data, err := json.Marshal(b)
	if err != nil {
		return nil
	}

	return os.WriteFile(filename, data, 0644)
}

func (b *Books) Delete(index int) error {
	ls := *b
	if index <= 0 || index > len(ls) {
		return errors.New("invalid index")
	}
	*b = append(ls[:index-1], ls[index:]...)

	return nil
}

func (b *Books) Print() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Book"},
			{Align: simpletable.AlignCenter, Text: "Done"},
			{Align: simpletable.AlignRight, Text: "CreatedAt"},
			{Align: simpletable.AlignRight, Text: "CompletedAt"},
		},
	}

	var cells [][]*simpletable.Cell

	for idx, item := range *b {
		idx++
		name := blue(item.Name)
		done := blue("no")
		if item.Done {
			name = green(fmt.Sprintf("\u2705 %s", item.Name))
			done = green("yes")
		}
		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", idx)},
			{Text: name},
			{Text: done},
			{Text: item.Added.Format(time.RFC822)},
			{Text: item.CompletedAt.Format(time.RFC822)},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: red(fmt.Sprintf("you have %d pending books", b.CountPending()))},
	}}

	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}

func (b *Books) CountPending() int {
	total := 0

	for _, item := range *b {
		if item.Done {
			total++
		}
	}
	return total
}
