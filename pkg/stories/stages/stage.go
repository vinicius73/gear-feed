package stages

import "os"

type Stage struct {
	Width      int
	Height     int
	Full       string
	Background string
	Foreground string
}

func (s Stage) Files() []string {
	return []string{s.Full, s.Background, s.Foreground}
}

func (s Stage) RemoveAll() error {
	for _, file := range s.Files() {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}
