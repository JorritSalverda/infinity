package lib

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/logrusorgru/aurora"
)

type Scaffolder interface {
	Scaffold(ctx context.Context, templateName, applicationName string) (err error)
}

type scaffolder struct {
	verbose               bool
	buildManifestFilename string
	templateBaseURL       string
}

func NewScaffolder(verbose bool, buildManifestFilename, templateBaseURL string) Scaffolder {
	return &scaffolder{
		verbose:               verbose,
		buildManifestFilename: buildManifestFilename,
		templateBaseURL:       templateBaseURL,
	}
}

func (s *scaffolder) Scaffold(ctx context.Context, templateName, applicationName string) (err error) {

	templateURL := fmt.Sprintf("%v%v.gotmpl", s.templateBaseURL, templateName)

	// fetch template
	resp, err := http.Get(templateURL)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Template url %v returned invalid status code %v", templateURL, resp.StatusCode)
	}
	defer resp.Body.Close()

	// read body into string
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	templateString := buf.String()

	if s.verbose {
		log.Println(aurora.BrightBlue(templateURL))
		log.Println(templateString)
	}

	// render template as gotemplate
	data := struct {
		ApplicationName string
	}{applicationName}
	tmpl, err := template.New("renderTemplate").Parse(templateString)
	if err != nil {
		return
	}
	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, data)
	if err != nil {
		return
	}

	// write to .infinity.yaml file
	err = ioutil.WriteFile(s.buildManifestFilename, renderedTemplate.Bytes(), 0600)
	if err != nil {
		return
	}

	return nil
}
