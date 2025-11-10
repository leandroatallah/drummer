package gamescene

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

const (
	// NoteOffset adjusts the timing of the notes to match the audio.
	// A positive value makes the notes appear later (you hit them earlier).
	// A negative value makes the notes appear earlier (you hit them later).
	NoteOffset = 0.23
)

type Note struct {
	Direction string  `json:"direction"`
	Onset     float64 `json:"onset"`
	skip      bool
}

type Song struct {
	Title    string  `json:"title"`
	Filename string  `json:"filename"`
	Bpm      int     `json:"bpm"`
	Duration float64 `json:"duration"`
	Notes    []*Note `json:"notes"`
	scene    *PlayScene

	PlayingNotes map[int]*Note
	noteIndex    int
	offsetBpm    float64
	count        int
}

func NewSong(path string, scene *PlayScene) *Song {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	return NewSongFromData(byteValue, scene)
}

func NewSongFromData(data []byte, scene *PlayScene) *Song {
	var song Song
	if err := json.Unmarshal(data, &song); err != nil {
		log.Fatal(err)
	}
	song.PlayingNotes = make(map[int]*Note)
	song.scene = scene

	return &song
}

func (s *Song) Update() error {
	s.count++

	// Remove old notes from PlayingNotes
	for i, n := range s.PlayingNotes {
		if s.GetPositionInBPM() > n.Onset+1 { // 1 beat buffer
			if !n.skip {
				s.scene.handleMistake()
			}
			delete(s.PlayingNotes, i)
		}
	}

	// Get playing notes
	for s.noteIndex < len(s.Notes) {
		n := s.Notes[s.noteIndex]
		if s.GetPositionInBPM()+s.offsetBpm >= n.Onset {
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
	if s.scene.songPlayer == nil {
		return 0
	}
	// Get the current position of the audio player
	currentTime := s.scene.songPlayer.Current()

	// Convert the time to seconds
	seconds := currentTime.Seconds()

	// Calculate the position in beats
	beats := seconds * (float64(s.Bpm) / 60.0)

	return beats + NoteOffset
}

func (s *Song) GetTicksPerBeat() float64 {
	return (60 * 60) / float64(s.Bpm)
}

func (s *Song) IsOver() bool {
	panic("implement me")
}

func (s *Song) SetPositionInBPM(beats float64) {
	// 1. Calculate the time in seconds from the beat position.
	seconds := beats / (float64(s.Bpm) / 60.0)
	seekDuration := time.Duration(seconds * float64(time.Second))

	// 2. Seek the audio player.
	if s.scene.songPlayer != nil {
		s.scene.songPlayer.Seek(seekDuration)
	}

	// 3. Update the noteIndex.
	s.noteIndex = 0
	for i, n := range s.Notes {
		if n.Onset < beats {
			s.noteIndex = i + 1
		} else {
			break
		}
	}

	// 4. Clear PlayingNotes.
	s.PlayingNotes = make(map[int]*Note)
}
