package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_copyTreeAllTargetFilesPresentInSource_AllCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	srcCh := make(chan *string, 10)
	tgtCh := make(chan *string, 10)

	go func(sch chan<- *string, tch chan<- *string) {
		src := []string{"/s/p1/f1.txt", "/s/p1/f2.txt", "/s/p1/p2/f1.txt", "/s/p1/p2/f2.txt"}
		tgt := []string{"/t/p1/f1.txt", "/t/p1/f2.txt", "/t/p1/p2/f1.txt", "/t/p1/p2/f2.txt"}

		for _, s := range src {
			tmp := s
			sch <- &tmp
		}
		close(sch)
		for _, t := range tgt {
			tmp := t
			tch <- &tmp
		}
		close(tch)
	}(srcCh, tgtCh)

	result := make(map[string]string)

	// Act
	r := copyTree(srcCh, tgtCh, "/s", "/t", false, func(src, dst string) error {
		result[dst] = src
		return nil
	})

	// Assert
	ass.Equal(4, len(result))
	ass.Equal("/s/p1/f1.txt", result["/t/p1/f1.txt"])
	ass.Equal("/s/p1/f2.txt", result["/t/p1/f2.txt"])
	ass.Equal("/s/p1/p2/f1.txt", result["/t/p1/p2/f1.txt"])
	ass.Equal("/s/p1/p2/f2.txt", result["/t/p1/p2/f2.txt"])
	ass.Equal(int64(4), r.TotalCopied)
	ass.Equal(int64(0), r.NotFoundInSource)
}

func Test_copyTreeSourcesMoreThenTargets_OnlyMathesCopiedFromSources(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	srcCh := make(chan *string, 10)
	tgtCh := make(chan *string, 10)

	go func(sch chan<- *string, tch chan<- *string) {
		src := []string{"/s/p1/f1.txt", "/s/p1/f2.txt", "/s/p1/p2/f1.txt", "/s/p1/p2/f2.txt"}
		tgt := []string{"/t/p1/f1.txt", "/t/p1/f2.txt"}

		for _, s := range src {
			tmp := s
			sch <- &tmp
		}
		close(sch)
		for _, t := range tgt {
			tmp := t
			tch <- &tmp
		}
		close(tch)
	}(srcCh, tgtCh)

	result := make(map[string]string)

	// Act
	r := copyTree(srcCh, tgtCh, "/s", "/t", false, func(src, dst string) error {
		result[dst] = src
		return nil
	})

	// Assert
	ass.Equal(2, len(result))
	ass.Equal("/s/p1/f1.txt", result["/t/p1/f1.txt"])
	ass.Equal("/s/p1/f2.txt", result["/t/p1/f2.txt"])
	ass.Equal(int64(2), r.TotalCopied)
	ass.Equal(int64(0), r.NotFoundInSource)
}

func Test_copyTreeTargetsContainMissingSourcesElements_OnlyFoundCopiedFromSources(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	srcCh := make(chan *string, 10)
	tgtCh := make(chan *string, 10)

	go func(sch chan<- *string, tch chan<- *string) {
		src := []string{"/s/p1/f1.txt"}
		tgt := []string{"/t/p1/f1.txt", "/t/p1/f2.txt"}

		for _, s := range src {
			tmp := s
			sch <- &tmp
		}
		close(sch)
		for _, t := range tgt {
			tmp := t
			tch <- &tmp
		}
		close(tch)
	}(srcCh, tgtCh)

	result := make(map[string]string)

	// Act
	r := copyTree(srcCh, tgtCh, "/s", "/t", false, func(src, dst string) error {
		result[dst] = src
		return nil
	})

	// Assert
	ass.Equal(1, len(result))
	ass.Equal("/s/p1/f1.txt", result["/t/p1/f1.txt"])
	ass.Equal(int64(1), r.TotalCopied)
	ass.Equal(int64(1), r.NotFoundInSource)
}

func Test_copyTreeSourcesContainsSameNameFilesButInSubfolders_OnlyExactMatchedCopiedFromSources(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	srcCh := make(chan *string, 10)
	tgtCh := make(chan *string, 10)

	go func(sch chan<- *string, tch chan<- *string) {
		src := []string{"/s/p1/f1.txt", "/s/p1/p2/f2.txt"}
		tgt := []string{"/t/p1/f1.txt", "/t/p1/f2.txt"}

		for _, s := range src {
			tmp := s
			sch <- &tmp
		}
		close(sch)
		for _, t := range tgt {
			tmp := t
			tch <- &tmp
		}
		close(tch)
	}(srcCh, tgtCh)

	result := make(map[string]string)

	// Act
	copyTree(srcCh, tgtCh, "/s", "/t", false, func(src, dst string) error {
		result[dst] = src
		return nil
	})

	// Assert
	ass.Equal(1, len(result))
	ass.Equal("/s/p1/f1.txt", result["/t/p1/f1.txt"])
}

func Test_copyTreeSourcesContainsNoMatchingFiles_NothingCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	srcCh := make(chan *string, 10)
	tgtCh := make(chan *string, 10)

	go func(sch chan<- *string, tch chan<- *string) {
		src := []string{"/s/p1/f3.txt"}
		tgt := []string{"/t/p1/f1.txt", "/t/p1/f2.txt"}

		for _, s := range src {
			tmp := s
			sch <- &tmp
		}
		close(sch)
		for _, t := range tgt {
			tmp := t
			tch <- &tmp
		}
		close(tch)
	}(srcCh, tgtCh)

	result := make(map[string]string)

	// Act
	r := copyTree(srcCh, tgtCh, "/s", "/t", false, func(src, dst string) error {
		result[dst] = src
		return nil
	})

	// Assert
	ass.Equal(0, len(result))
	ass.Equal(int64(0), r.TotalCopied)
	ass.Equal(int64(2), r.NotFoundInSource)
}

func Test_copyTreeEmptySources_NothingCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	srcCh := make(chan *string, 10)
	tgtCh := make(chan *string, 10)

	go func(sch chan<- *string, tch chan<- *string) {
		src := []string{}
		tgt := []string{"/t/p1/f1.txt", "/t/p1/f2.txt"}

		for _, s := range src {
			tmp := s
			sch <- &tmp
		}
		close(sch)
		for _, t := range tgt {
			tmp := t
			tch <- &tmp
		}
		close(tch)
	}(srcCh, tgtCh)

	result := make(map[string]string)

	// Act
	copyTree(srcCh, tgtCh, "/s", "/t", false, func(src, dst string) error {
		result[dst] = src
		return nil
	})

	// Assert
	ass.Equal(0, len(result))
}

func Test_copyTreeEmptyTargets_NothingCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	srcCh := make(chan *string, 10)
	tgtCh := make(chan *string, 10)

	go func(sch chan<- *string, tch chan<- *string) {
		src := []string{"/s/p1/f3.txt"}
		tgt := []string{}

		for _, s := range src {
			tmp := s
			sch <- &tmp
		}
		close(sch)
		for _, t := range tgt {
			tmp := t
			tch <- &tmp
		}
		close(tch)
	}(srcCh, tgtCh)

	result := make(map[string]string)

	// Act
	copyTree(srcCh, tgtCh, "/s", "/t", false, func(src, dst string) error {
		result[dst] = src
		return nil
	})

	// Assert
	ass.Equal(0, len(result))
}

func Test_copyTreeDifferentCase_DifferentCaseFilesCopied(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	srcCh := make(chan *string, 10)
	tgtCh := make(chan *string, 10)

	go func(sch chan<- *string, tch chan<- *string) {
		src := []string{"/s/p1/f1.txt"}
		tgt := []string{"/t/p1/F1.txt"}

		for _, s := range src {
			tmp := s
			sch <- &tmp
		}
		close(sch)
		for _, t := range tgt {
			tmp := t
			tch <- &tmp
		}
		close(tch)
	}(srcCh, tgtCh)

	result := make(map[string]string)

	// Act
	copyTree(srcCh, tgtCh, "/s", "/t", false, func(src, dst string) error {
		result[dst] = src
		return nil
	})

	// Assert
	ass.Equal(1, len(result))
	ass.Equal("/s/p1/f1.txt", result["/t/p1/F1.txt"])
}
