package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

//Test if solution gives right answers for questions.
func Test_saveToCache(t *testing.T) {
	cache := "cacheTest"
	assert.NoError(t, saveToCache("hello", cache)) // No error in saveToCache
	defer os.Remove(cache)

	_, err := os.Stat(cache)
	assert.False(t, os.IsNotExist(err)) // file exists
	assert.NoError(t, err)              // No other error retrieving file information

	f, err := os.Open(cache)
	defer f.Close()
	assert.NoError(t, err)

	r, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(r))
}

func Test_retrieveCachedCode(t *testing.T) {
	cache := "cacheTest"
	//Shouldn't throw error or retrieve any code, file does not exist.
	code, err := retrieveCachedCode(cache)
	assert.NoError(t, err)
	assert.Equal(t, "", code)

	// Mock cache file
	f, err := os.Create(cache)
	assert.NoError(t, err)
	f.Write([]byte("hello"))
	defer f.Close()
	defer os.Remove(cache)

	//Should retrieve correct code.
	code, err = retrieveCachedCode(cache)
	assert.NoError(t, err)
	assert.Equal(t, "hello", code)
}

func Test_findQuestion(t *testing.T) {
	passHTML := "<span>question</span>"
	passDoc, err := html.Parse(strings.NewReader(passHTML))
	if err != nil {
		t.Fatal(err)
	}

	failHTML := ""
	failDoc, err := html.Parse(strings.NewReader(failHTML))
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		doc   *html.Node
		xpath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Find question - pass", args{passDoc, "//span"}, "question", false},
		{"Find question - fail", args{failDoc, "//span"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findQuestion(tt.args.doc, tt.args.xpath)
			if (err != nil) != tt.wantErr {
				t.Errorf("findQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("findQuestion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_retrieveAcessCode(t *testing.T) {
	page := `<input type="hidden" name="chaveDeAcesso" value="chaveDeAcesso1526349532185444432"` // AccessCode should have 32 characters
	type args struct {
		page string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"valid access code", args{page}, "chaveDeAcesso1526349532185444432"},
		{"access code not found", args{``}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := retrieveAcessCode(tt.args.page); got != tt.want {
				t.Errorf("retrieveAcessCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
