package dupremover

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func Test_ProcessWithDuplicatesDeletion(t *testing.T) {

	oldDir := ".\\testbackup\\"
	newDir := ".\\test\\"

	cmd := exec.Command("xcopy", oldDir, newDir, "/s/e")
	cmd.Run()

	defer os.RemoveAll(newDir)

	dir := "test"

	remover := New(dir)

	dups, err := remover.Process()
	if err != nil {
		t.Fatalf("Process error: %v\n", err)
	}

	if len(dups) != 2 {
		t.Fatalf("Wrong number of duplicate files.\n")
	}

	file1 := "b"
	file1_found := false

	file2 := "c.txt"
	file2_found := false

	for _, v := range dups {
		temp := strings.Split(v, "\\")
		name := temp[len(temp)-1]
		switch name {
		case file1:
			file1_found = true
		case file2:
			file2_found = true
		}
	}
	if !file1_found || !file2_found {
		t.Fatalf("Not all files found")
	}

	errs := remover.Remove()
	if len(errs) > 0 {
		t.Fatalf("Remove errors: %v", errs)
	}

	var file1_int, file2_int int = 0, 0

	err = filepath.Walk(newDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			switch info.Name() {
			case file1:
				file1_int++
			case file2:
				file2_int++
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to read a processed directory")
	}

	if file1_int+file2_int > 2 {
		t.Fatalf("Duplicate files found after the deletion")
	}

}

func Test_ProcessWithoutDuplicatesDeletion(t *testing.T) {

	oldDir := ".\\testbackup\\"
	newDir := ".\\test\\"

	cmd := exec.Command("xcopy", oldDir, newDir, "/s/e")
	cmd.Run()

	defer os.RemoveAll(newDir)

	dir := "test"

	remover := New(dir)

	dups, err := remover.Process()
	if err != nil {
		t.Fatalf("Process error: %v\n", err)
	}

	if len(dups) != 2 {
		t.Fatalf("Wrong number of duplicate files.\n")
	}

	file1 := "b"
	file1_found := false

	file2 := "c.txt"
	file2_found := false

	for _, v := range dups {
		temp := strings.Split(v, "\\")
		name := temp[len(temp)-1]
		switch name {
		case file1:
			file1_found = true
		case file2:
			file2_found = true
		}
	}

	if !file1_found || !file2_found {
		t.Fatalf("Not all files found")
	}

	var file1_int, file2_int int = 0, 0

	err = filepath.Walk(newDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			switch info.Name() {
			case file1:
				file1_int++
			case file2:
				file2_int++
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to read a processed directory")
	}

	if file1_int+file2_int != 4 {
		t.Fatalf("Duplicate files are not found")
	}
}
