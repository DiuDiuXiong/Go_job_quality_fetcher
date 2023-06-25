package duplicateremover

import (
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func JaccardSimilarity(a, b map[string]float64) float64 {
	intersect := 0.0
	sa := 0.0
	sb := 0.0
	for _, v := range a {
		sa += v
	}

	for _, v := range b {
		sb += v
	}

	for key, vala := range a {
		if valb, exist := b[key]; exist {
			intersect += math.Min(vala, valb)
		}
	}

	return intersect / (sa + sb - intersect)
}

func GenerateCountForData(s *string) map[string]float64 {
	words := strings.FieldsFunc(*s, func(r rune) bool {
		return !strings.ContainsRune(" -'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", r)
	})

	counts := make(map[string]float64)
	for _, word := range words {
		word = strings.ToLower(word) // Convert to lowercase to handle case-insensitivity
		counts[word]++
	}

	return counts
}

func GenerateCountForFile(fileName string) (map[string]float64, error) {
	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if len(d) == 0 {
		return nil, nil
	}
	ds := string(d)
	return GenerateCountForData(&ds), nil
}

func GenerateCountForFolder(folder string) (map[string]map[string]float64, error) {
	res := make(map[string]map[string]float64)
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			d, err := GenerateCountForFile(path)
			if err != nil {
				return err
			}
			if d != nil {
				res[path] = d
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil

}

func CombineTwoFolderCounts(mA, mB map[string]map[string]float64, threshold float64) {
	match := make(map[string]string)
	for filePathA, countA := range mA {
		for filePathB, countB := range mB {
			if JaccardSimilarity(countA, countB) > threshold {
				match[filePathA] = filePathB
				break
			}
		}
	}
	for filePathA := range mA {
		if _, exist := match[filePathA]; !exist {
			mB[filePathA] = mA[filePathA]
		}
	}
}

func WriteFilesToNewFolder(m map[string]map[string]float64, destFolder string) error {
	idx := 0
	for src, _ := range m {
		if err := CopyFile(src, path.Join(destFolder, strconv.Itoa(idx)+".txt")); err != nil {
			return err
		}
		idx++
	}
	return nil
}

func RemoveDuplicate(m map[string]map[string]float64, threshold float64) {
	isDuplicate := make(map[string]bool)
	fileName := make([]string, 0)
	for k := range m {
		isDuplicate[k] = false
		fileName = append(fileName, k)
	}

	for i := 0; i < len(m); i++ {
		for j := i + 1; j < len(m); j++ {
			if isDuplicate[fileName[i]] || isDuplicate[fileName[j]] {
				continue
			}
			if JaccardSimilarity(m[fileName[i]], m[fileName[j]]) > threshold {
				isDuplicate[fileName[j]] = true
			}
		}
	}

	for fileN, dup := range isDuplicate {
		if dup {
			delete(m, fileN)
		}
	}
}

func RemoveDuplicates(folders []string, destFolder string, threshold float64) error {
	if folders == nil || len(folders) == 0 {
		return nil
	}
	res, err := GenerateCountForFolder(folders[0])
	RemoveDuplicate(res, threshold)
	if err != nil {
		return err
	}
	for idx := 1; idx < len(folders); idx++ {
		newCount, err := GenerateCountForFolder(folders[idx])
		RemoveDuplicate(newCount, threshold)
		if err != nil {
			return err
		}
		CombineTwoFolderCounts(newCount, res, threshold)
	}

	return WriteFilesToNewFolder(res, destFolder)
}

func CopyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
