package main

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func isEverythingNil(list []string) bool {
	nils := 0
	for _, v := range list {
		if v == "" {
			nils++
		}
	}

	if len(list) == nils {
		return true
	}

	return false
}

func del(list []string, index int) []string {
	n := []string{}
	for i, v := range list {
		if i == index {
			continue
		}
		n = append(n, v)
	}

	return n
}

func parse(list []string) []string {
	pq := false
	args := []string{}
	q := []string{}
	for _, v := range list {
		if pq == false && !strings.HasPrefix(v, "\"") {
			args = append(args, v)
		} else if pq == false && strings.HasPrefix(v, "\"") {
			pq = true
			q = append(q, strings.TrimPrefix(v, "\""))
		} else if pq == true && !strings.HasSuffix(v, "\"") {
			q = append(q, v)
		} else if pq == true && strings.HasSuffix(v, "\"") {
			q = append(q, strings.TrimSuffix(v, "\""))
			args = append(args, strings.Join(q, " "))
			q = []string{}
			pq = false
		}
	}

	return args
}

func org(args []string, items []string) map[string]string {
	m := map[string]string{}
	for i, v := range args {
		m[items[i]] = v
	}

	return m
}

func mapToSlice(m map[string]string) []string {
	a := []string{}
	for _, v := range m {
		a = append(a, v)
	}

	return a
}

func fileExists(f string) bool {
	_, err := os.Stat(f)
	if err == nil {
		return true
	}
	return false
}

func isHex(s string) bool {
	r := regexp.MustCompile("\\A\\b[0-9a-fA-F]+\\b")
	return r.MatchString(s)
}

func chkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func copy(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		return err
	}

	return nil
}

func download(url, output string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	in := bytes.NewReader(body)

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func unzipFile(in, file, out string) error {
	files, err := zip.OpenReader(in)
	if err != nil {
		return err
	}
	defer files.Close()

	for _, v := range files.File {
		if v.Name == file {
			source, err := v.Open()
			if err != nil {
				return err
			}
			defer source.Close()

			dest, err := os.Create(out)
			if err != nil {
				return err
			}
			defer dest.Close()

			_, err = io.Copy(dest, source)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
