package gamescene

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Note struct {
	Direction string `json:"direction"`
	Onset     int    `json:"onset"`
}

type Song struct {
	Title    string  `json:"title"`
	Filename string  `json:"filename"`
	Bpm      int     `json:"bpm"`
	Notes    []*Note `json:"notes"`

	PlayingNotes map[int]*Note
	noteIndex    int
	offsetBpm    float64
	count        int
}

func NewSong(path string) *Song {
	// TODO: Read song data from JSON file.
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	var song Song
	if err := json.Unmarshal(byteValue, &song); err != nil {
		log.Fatal(err)
	}
	song.PlayingNotes = make(map[int]*Note)

	return &song
}

func (s *Song) Update() error {
	s.count++

	// Remove old notes from PlayingNotes
	for i, n := range s.PlayingNotes {
		if s.GetPositionInBPM() > float64(n.Onset)+1 { // 1 beat buffer
			delete(s.PlayingNotes, i)
		}
	}

	// Get playing notes
	for s.noteIndex < len(s.Notes) {
		n := s.Notes[s.noteIndex]
		if s.GetPositionInBPM()+s.offsetBpm >= float64(n.Onset) {
			s.PlayingNotes[s.noteIndex] = n
			s.noteIndex++
		} else {
			// The upcoming notes are not ready yet.
			break
		}
	}

	return nil
}

func (s *Song) NextNote() *Note {
	s.noteIndex++
	// Check the current note
	// If it should be drawed on track, returns and increase noteIndex
	// else return nil
	if s.noteIndex >= len(s.Notes) {
		return nil
	}

	note := s.Notes[s.noteIndex]
	// Check if note should be drawed

	return note
}

// bpm = beats per minute
// counts runs 60 bpms per default
// In 60 bpm, 1 bpm happens one time per second
// then, 60 counts is 1 beat
func (s *Song) GetPositionInBPM() float64 {
	ticksPerBeat := s.GetTicksPerBeat()
	return float64(s.count) / ticksPerBeat
}

func (s *Song) GetTicksPerBeat() float64 {
	return (60 * 60) / float64(s.Bpm)
}
