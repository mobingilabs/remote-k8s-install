package template

import (
	stdtemplate "text/template"
)

func Parse(text string, data interface{}) ([]byte, error) {
	t, err := stdtemplate.New("").Parse(text)
	if err != nil {
		return nil, err
	}

	bw := &bytesWriter{}
	if err := t.Execute(bw, data); err != nil {
		return nil, err
	}
	return bw.data, nil
}

type bytesWriter struct {
	data []byte
}

func (bw *bytesWriter) Write(p []byte) (int, error) {
	bw.data = append(bw.data, p...)
	return len(p), nil
}
