package file_template

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"os"
	"text/template"
)

var (
	NO_FUNCS = template.FuncMap{}

	FUNCTIONS = template.FuncMap{
		"base64Encode": func(value string) string {
			return Base64EncodeString(value)
		},
		"basicCredentials": func(username, password string) string {
			creds := fmt.Sprintf("%s:%s", username, password)
			return Base64EncodeString(creds)
		},
	}
)

type FileTemplate struct {
	Name    string
	Content string
	tmpl    *template.Template
}

func CreateFilePathTemplate(name, filePath string, fmap template.FuncMap) *FileTemplate {
	appContent, err := ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	tmpl, err := NewFileTemplate(name, string(appContent), fmap)
	if err != nil {
		panic(err)
	}

	return tmpl
}

func NewFileTemplate(name, content string, fmap template.FuncMap) (*FileTemplate, error) {
	//name := strings.Split(filepath.Base(path), ".")[0]
	//data, err := os.ReadFile(path)
	//if err != nil {
	//	return nil, err
	//}
	//content := string(data)
	return NewTextTemplate(name, content, fmap)
}

func NewTextTemplate(name, content string, fmap template.FuncMap) (*FileTemplate, error) {
	var tmpl *template.Template
	var err error

	if fmap != nil && len(fmap) > 0 {
		tmpl, err = template.New(name).Funcs(fmap).Parse(content)
	} else {
		tmpl, err = template.New(name).Parse(content)
	}
	if err != nil {
		return nil, err
	}

	return &FileTemplate{Name: name, Content: content, tmpl: tmpl}, nil
}

func (f *FileTemplate) ExecuteTemplate(name string, data interface{}) (string, error) {

	var buffer = new(bytes.Buffer)
	var err error

	err = f.tmpl.ExecuteTemplate(buffer, name, data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (f *FileTemplate) ExecuteTemplates(names []string, data interface{}) ([]string, error) {

	results := make([]string, len(names))

	for index, name := range names {
		var buffer = new(bytes.Buffer)
		var err error

		err = f.tmpl.ExecuteTemplate(buffer, name, data)
		if err != nil {
			return nil, fmt.Errorf("failed to generate template for %s - %s", name, err)
		}
		results[index] = buffer.String()
	}

	return results, nil
}

func (f *FileTemplate) Execute(data interface{}, fmap template.FuncMap) (string, error) {

	var buffer = new(bytes.Buffer)
	var err error

	if len(fmap) > 0 {
		err = template.Must(f.tmpl.Clone()).Funcs(fmap).Execute(buffer, data)
	} else {
		err = f.tmpl.Execute(buffer, data)
	}

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func Base64EncodeString(value string) string {
	return string([]byte(b64.StdEncoding.EncodeToString([]byte(value))))
}

func ReadFile(path string) ([]byte, error) {
	//name := strings.Split(filepath.Base(path), ".")[0]
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}
