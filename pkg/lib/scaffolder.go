package lib

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/logrusorgru/aurora"
)

//go:generate mockgen -package=lib -destination ./scaffolder_mock.go -source=scaffolder.go
type Scaffolder interface {
	Scaffold(ctx context.Context, applicationType ApplicationType, language Language, applicationName string) (err error)
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

func (s *scaffolder) Scaffold(ctx context.Context, applicationType ApplicationType, language Language, applicationName string) (err error) {
	if !applicationType.IsSupported() {
		return fmt.Errorf("application type is unknown; supported values are %v", strings.Join(SupportedApplicationTypes.ToStringArray(), ", "))
	}
	if !language.IsSupported() {
		return fmt.Errorf("language type is unknown; supported values are %v", strings.Join(SupportedLanguages.ToStringArray(), ", "))
	}
	if applicationName == "" {
		return fmt.Errorf("application name is empty")
	}

	templateURL := fmt.Sprintf("%v%v-%v.gotmpl", s.templateBaseURL, applicationType, language)

	// fetch template
	resp, err := http.Get(templateURL)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("template url %v returned invalid status code %v", templateURL, resp.StatusCode)
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
