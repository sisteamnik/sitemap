package sitemap

import (
	"io"
)

type Writer struct {
	w             io.Writer
	isIndexWriter bool
	headerWrited  bool
}

func NewWriter(w io.Writer) (wr *Writer, e error) {
	wr = new(Writer)
	wr.w = w
	return
}

func NewIndexWriter(w io.Writer) (wr *Writer, e error) {
	wr, e = NewWriter(w)
	wr.isIndexWriter = true
	return
}

func (w *Writer) writeHeader() error {
	toWrite := []byte(header)
	if w.isIndexWriter {
		toWrite = []byte(indexHeader)
	}
	_, e := w.w.Write(toWrite)
	if e != nil {
		return e
	}
	w.headerWrited = true
	return nil
}

func (w *Writer) Put(i Item) error {
	if !w.headerWrited {
		e := w.writeHeader()
		if e != nil {
			return e
		}
	}

	// check and fix item isIndex
	if w.isIndexWriter != i.isIndex {
		i.isIndex = w.isIndexWriter
	}

	_, e := w.w.Write([]byte(i.String()))
	return e
}

func (w *Writer) Release() error {
	if !w.headerWrited {
		return nil
	}

	toWrite := []byte(footer)
	if w.isIndexWriter {
		toWrite = []byte(indexFooter)
	}
	_, e := w.w.Write(toWrite)
	if e != nil {
		return e
	}
	return nil
}
