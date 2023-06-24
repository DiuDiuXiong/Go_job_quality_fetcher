package extractor

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

func TestExtractor(t *testing.T) {
	tests := []struct {
		start string
		end   string
		src   string
	}{
		{
			start: "What\nWhere\nFind Jobs\n",
			end:   "Report job\nHiring Lab\nCareer Advice\n",
			src:   "Indeed.txt",
		},
		{
			start: "Company reviews\nSign in\nEmployer site",
			end:   "Job seekers\nJob search\n",
			src:   "Seek.txt",
		},
		{
			start: "Dismiss\nJoin now\nSign in",
			end:   "Referrals increase your chances of interviewing at ",
			src:   "Linkedin.txt",
		},
	}

	for _, tt := range tests {
		data, err := ioutil.ReadFile(path.Join("tmp", tt.src))
		if err != nil {
			t.Errorf("Failed to load job description data, please run text crawler test first to generate the date")
		}
		dstr := string(data)
		if c := strings.Count(dstr, tt.start); c != 1 {
			t.Errorf("The regrex for extract job description for %s is broke since found %d amounts of starting string.", tt.src, c)
		}

		if c := strings.Count(dstr, tt.end); c != 1 {
			t.Errorf("The regrex for extract job description for %s is broke since found %d amounts of ending string.", tt.src, c)
		}
	}

}
