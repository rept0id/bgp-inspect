package main

type TRecord struct {
	Name string
	Loops int
	Repeats int
	ConsecutiveRepeats int
	NonConsecutiveRepeats int
}
type TRecords map[uint32]*TRecord
