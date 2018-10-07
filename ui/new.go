package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/TsuyoshiUshio/strikes/template/assets"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type PackageParameter struct {
	TemplateDirPath string
	PackageName     string
	Description     string
	Author          string
	ProjectPage     string
	ProjectRepo     string
	ReleaseNote     string
	ZipFileName     string
}

type ChooseTemplateProcess struct {
	Stdin        *os.File
	ProviderType string
	PackageList  []string
	NextProcess  *Process
	Parameter    PackageParameter
}

func (p *ChooseTemplateProcess) PrintQuestion() error {
	p.PackageList = assets.List(p.ProviderType)
	for i, template := range p.PackageList {
		content, err := assets.ReadTemplateDescription("/" + p.ProviderType + "/" + template)
		if err != nil {
			log.Fatalf("Can not find Template Description for %s, error: %v\n", template, err)
			return nil
		}
		fmt.Printf("%d: %s:%s %s\n", i, template, adjustTabs(template), content)
	}
	fmt.Println("")
	return nil
}
func (p *ChooseTemplateProcess) WaitForInput() (string, error) {
	return readLine((*p).Stdin)
}

func (p *ChooseTemplateProcess) Validate(answer string) bool {
	i, err := strconv.Atoi(answer)
	if err != nil {
		fmt.Printf("Select the proper value. %s is not accepted. \n", answer)
		return false
	}
	if i > (len(p.PackageList) - 1) {
		fmt.Printf("Select the proper value. %s is not accepted. \n", answer)
		return false
	}
	return true
}
func (p *ChooseTemplateProcess) UpdateParameter(answer string, parameter interface{}) (interface{}, error) {
	i, _ := strconv.Atoi(answer) // already validated.
	param := parameter.(PackageParameter)
	param.TemplateDirPath = "/" + p.ProviderType + "/" + p.PackageList[i]
	return param, nil

}
func (p *ChooseTemplateProcess) ShowValidateError(answer string) {
	fmt.Printf("You choose: %s, However, the value should be [%d - %d]\n", answer, 0, len(p.PackageList)-1)
}
func (p *ChooseTemplateProcess) SetNext(process *Process) {
	p.NextProcess = process
}
func (p *ChooseTemplateProcess) Next() *Process {
	return p.NextProcess
}

func (p *ChooseTemplateProcess) IsTargetParameterFilled(parameter interface{}) bool {
	param := parameter.(PackageParameter)
	if param.TemplateDirPath != "" {
		return true
	} else {
		return false
	}
}

func (p *ChooseTemplateProcess) SetParameter(parameter interface{}) {
	p.Parameter = parameter.(PackageParameter)
}

func adjustTabs(name string) string {
	if len(name) > 11 {
		return "\t"
	} else {
		return "\t\t"
	}
}

func readLine(file *os.File) (string, error) {
	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	} else {
		return string(line), nil
	}
}

func NewChooseTemplateProcess(providerType string, file *os.File) *Process {
	var p Process
	p = &ChooseTemplateProcess{
		ProviderType: providerType,
		Stdin:        file,
	}
	return &p
}

type PackageNameProcess struct {
	Stdin       *os.File
	NextProcess *Process
	Parameter   PackageParameter
}

func (p *PackageNameProcess) PrintQuestion() error {
	fmt.Print("PackageName : ")
	return nil
}
func (p *PackageNameProcess) WaitForInput() (string, error) {
	return readLine((*p).Stdin)
}
func (p *PackageNameProcess) Validate(answer string) bool {
	// start with letter
	// letter or number
	// you can include -
	// less than 12 charactor
	err := validation.Validate(answer,
		validation.Required,
		validation.Length(5, 14),
		validation.Match(regexp.MustCompile("^[a-z][a-z0-9-]*")))
	if err != nil {
		return false
	} else {
		return true
	}
}
func (p *PackageNameProcess) IsTargetParameterFilled(parameter interface{}) bool {
	param := parameter.(PackageParameter)
	if param.PackageName != "" {
		return true
	}
	return false
}
func (p *PackageNameProcess) UpdateParameter(answer string, parameter interface{}) (interface{}, error) {
	param := parameter.(PackageParameter)
	param.PackageName = answer
	return param, nil
}
func (p *PackageNameProcess) ShowValidateError(answer string) {
	fmt.Printf("PackageName %s not allowed. \n", answer)
	fmt.Println("PackageName should be start with [a-z], lowercase alphanumeric with '-', lenght should be [4 - 14].")
	fmt.Println("")
}
func (p *PackageNameProcess) SetNext(process *Process) {
	p.NextProcess = process
}
func (p *PackageNameProcess) Next() *Process {
	return p.NextProcess
}

func (p *PackageNameProcess) SetParameter(parameter interface{}) {
	p.Parameter = parameter.(PackageParameter)
}

func NewPackageNameProcess(file *os.File) *Process {
	var p Process
	p = &PackageNameProcess{
		Stdin: file,
	}
	return &p
}

type DescriptionProcess struct {
	Stdin       *os.File
	NextProcess *Process
	Parameter   PackageParameter
}

func (p *DescriptionProcess) PrintQuestion() error {
	fmt.Printf("Description [default: %s package.]: ", p.Parameter.PackageName)
	return nil
}
func (p *DescriptionProcess) WaitForInput() (string, error) {
	return readLine((*p).Stdin)
}
func (p *DescriptionProcess) Validate(answer string) bool {
	// length is [5 - 1024]
	err := validation.Validate(answer,
		validation.Length(5, 1024))
	if err != nil {
		return false
	} else {
		return true
	}
}
func (p *DescriptionProcess) IsTargetParameterFilled(parameter interface{}) bool {
	param := parameter.(PackageParameter)
	if param.Description != "" {
		return true
	}
	return false
}
func (p *DescriptionProcess) UpdateParameter(answer string, parameter interface{}) (interface{}, error) {
	if answer == "" {
		answer = fmt.Sprintf("%s package.", p.Parameter.PackageName)
	}
	param := parameter.(PackageParameter)
	param.Description = answer
	return param, nil
}
func (p *DescriptionProcess) ShowValidateError(answer string) {
	fmt.Printf("Description length should be between 5 - 1024. %s is not allowed. \n", answer)
	fmt.Println("")
}
func (p *DescriptionProcess) SetNext(process *Process) {
	p.NextProcess = process
}
func (p *DescriptionProcess) Next() *Process {
	return p.NextProcess
}

func (p *DescriptionProcess) SetParameter(parameter interface{}) {
	p.Parameter = parameter.(PackageParameter)
}

func NewDescriptionProcess(file *os.File) *Process {
	var p Process
	p = &DescriptionProcess{
		Stdin: file,
	}
	return &p
}

type AuthorProcess struct {
	Stdin       *os.File
	NextProcess *Process
	Parameter   PackageParameter
}

func (p *AuthorProcess) PrintQuestion() error {
	fmt.Printf("Author name : ")
	return nil
}
func (p *AuthorProcess) WaitForInput() (string, error) {
	return readLine((*p).Stdin)
}
func (p *AuthorProcess) Validate(answer string) bool {
	// length is [3 - 64]
	err := validation.Validate(answer,
		validation.Required,
		validation.Length(3, 64))
	if err != nil {
		return false
	} else {
		return true
	}
}
func (p *AuthorProcess) IsTargetParameterFilled(parameter interface{}) bool {
	param := parameter.(PackageParameter)
	if param.Author != "" {
		return true
	}
	return false
}
func (p *AuthorProcess) UpdateParameter(answer string, parameter interface{}) (interface{}, error) {
	param := parameter.(PackageParameter)
	param.Author = answer
	return param, nil
}
func (p *AuthorProcess) ShowValidateError(answer string) {
	fmt.Printf("Author length should be between 3 - 64. Empty is not allowed.: %s \n", answer)
	fmt.Println("")
}
func (p *AuthorProcess) SetNext(process *Process) {
	p.NextProcess = process
}
func (p *AuthorProcess) Next() *Process {
	return p.NextProcess
}

func (p *AuthorProcess) SetParameter(parameter interface{}) {
	p.Parameter = parameter.(PackageParameter)
}

func NewAuthorProcess(file *os.File) *Process {
	var p Process
	p = &AuthorProcess{
		Stdin: file,
	}
	return &p
}

type ProjectPageProcess struct {
	Stdin       *os.File
	NextProcess *Process
	Parameter   PackageParameter
}

func (p *ProjectPageProcess) PrintQuestion() error {
	fmt.Print("Project Page : [ default: https://github.com ] ")
	return nil
}
func (p *ProjectPageProcess) WaitForInput() (string, error) {
	return readLine((*p).Stdin)
}
func (p *ProjectPageProcess) Validate(answer string) bool {
	// should be valid URL
	err := validation.Validate(answer,
		is.URL)
	if err != nil {
		return false
	} else {
		return true
	}
}
func (p *ProjectPageProcess) IsTargetParameterFilled(parameter interface{}) bool {
	param := parameter.(PackageParameter)
	if param.ProjectPage != "" {
		return true
	}
	return false
}
func (p *ProjectPageProcess) UpdateParameter(answer string, parameter interface{}) (interface{}, error) {

	param := parameter.(PackageParameter)
	if answer == "" && param.ProjectPage == "" {
		answer = "https://github.com"
	}
	if param.ProjectPage == "" {
		param.ProjectPage = answer
	}
	return param, nil
}
func (p *ProjectPageProcess) ShowValidateError(answer string) {
	fmt.Printf("Project Page should be valid URL. : %s \n", answer)
	fmt.Println("")
}
func (p *ProjectPageProcess) SetNext(process *Process) {
	p.NextProcess = process
}
func (p *ProjectPageProcess) Next() *Process {
	return p.NextProcess
}

func (p *ProjectPageProcess) SetParameter(parameter interface{}) {
	p.Parameter = parameter.(PackageParameter)
}

func NewProjectPageProcess(file *os.File) *Process {
	var p Process
	p = &ProjectPageProcess{
		Stdin: file,
	}
	return &p
}

type ProjectRepoProcess struct {
	Stdin       *os.File
	NextProcess *Process
	Parameter   PackageParameter
}

func (p *ProjectRepoProcess) PrintQuestion() error {
	fmt.Print("Project Repository : [ default: https://github.com ] ")
	return nil
}
func (p *ProjectRepoProcess) WaitForInput() (string, error) {
	return readLine((*p).Stdin)
}
func (p *ProjectRepoProcess) Validate(answer string) bool {
	// should be valid URL
	err := validation.Validate(answer,
		is.URL)
	if err != nil {
		return false
	} else {
		return true
	}
}
func (p *ProjectRepoProcess) IsTargetParameterFilled(parameter interface{}) bool {
	param := parameter.(PackageParameter)
	if param.ProjectRepo != "" {
		return true
	}
	return false
}
func (p *ProjectRepoProcess) UpdateParameter(answer string, parameter interface{}) (interface{}, error) {

	param := parameter.(PackageParameter)
	if answer == "" && param.ProjectRepo == "" {
		answer = "https://github.com"
	}
	if param.ProjectRepo == "" {
		param.ProjectRepo = answer
	}
	return param, nil
}
func (p *ProjectRepoProcess) ShowValidateError(answer string) {
	fmt.Printf("Project Repository should be valid URL. : %s \n", answer)
	fmt.Println("")
}
func (p *ProjectRepoProcess) SetNext(process *Process) {
	p.NextProcess = process
}
func (p *ProjectRepoProcess) Next() *Process {
	return p.NextProcess
}

func (p *ProjectRepoProcess) SetParameter(parameter interface{}) {
	p.Parameter = parameter.(PackageParameter)
}

func NewProjectRepoProcess(file *os.File) *Process {
	var p Process
	p = &ProjectRepoProcess{
		Stdin: file,
	}
	return &p
}

type ReleaseNoteProcess struct {
	Stdin       *os.File
	NextProcess *Process
	Parameter   PackageParameter
}

func (p *ReleaseNoteProcess) PrintQuestion() error {
	fmt.Print("Release Note : [ default: Initial Release. ] ")
	return nil
}
func (p *ReleaseNoteProcess) WaitForInput() (string, error) {
	return readLine((*p).Stdin)
}
func (p *ReleaseNoteProcess) Validate(answer string) bool {
	return true
}
func (p *ReleaseNoteProcess) IsTargetParameterFilled(parameter interface{}) bool {
	param := parameter.(PackageParameter)
	if param.ReleaseNote != "" {
		return true
	}
	return false
}
func (p *ReleaseNoteProcess) UpdateParameter(answer string, parameter interface{}) (interface{}, error) {

	param := parameter.(PackageParameter)
	if answer == "" && param.ReleaseNote == "" {
		answer = "Initial Release."
	}
	if param.ReleaseNote == "" {
		param.ReleaseNote = answer
	}
	return param, nil
}
func (p *ReleaseNoteProcess) ShowValidateError(answer string) {
	// Not called. no validation.
	fmt.Println("")
}
func (p *ReleaseNoteProcess) SetNext(process *Process) {
	p.NextProcess = process
}
func (p *ReleaseNoteProcess) Next() *Process {
	return p.NextProcess
}

func (p *ReleaseNoteProcess) SetParameter(parameter interface{}) {
	p.Parameter = parameter.(PackageParameter)
}

func NewReleaseNoteProcess(file *os.File) *Process {
	var p Process
	p = &ReleaseNoteProcess{
		Stdin: file,
	}
	return &p
}

type ZipFileNameProcess struct {
	Stdin       *os.File
	NextProcess *Process
	Parameter   PackageParameter
}

func (p *ZipFileNameProcess) PrintQuestion() error {
	fmt.Printf("ZipFileName : [ default: %s.zip ] ", p.Parameter.PackageName)
	return nil
}
func (p *ZipFileNameProcess) WaitForInput() (string, error) {
	return readLine((*p).Stdin)
}
func (p *ZipFileNameProcess) Validate(answer string) bool {
	// should be not
	err := validation.Validate(answer,
		validation.Match(regexp.MustCompile("^[a-z][a-z0-9-]*\\.zip")),
		validation.Length(6, 64))
	if err != nil {
		return false
	} else {
		return true
	}
}
func (p *ZipFileNameProcess) IsTargetParameterFilled(parameter interface{}) bool {
	param := parameter.(PackageParameter)
	if param.ZipFileName != "" {
		return true
	}
	return false
}
func (p *ZipFileNameProcess) UpdateParameter(answer string, parameter interface{}) (interface{}, error) {

	param := parameter.(PackageParameter)
	if answer == "" && param.ZipFileName == "" {
		answer = fmt.Sprintf("%s.zip", p.Parameter.PackageName)
	}
	if param.ZipFileName == "" {
		param.ZipFileName = answer
	}
	return param, nil
}
func (p *ZipFileNameProcess) ShowValidateError(answer string) {
	fmt.Printf("ZipFileName %s not allowed. \n", answer)
	fmt.Println("ZipFileName should be *.zip file. Start with [a-z], lowercase alphanumeric with '-', lenght should be [6 - 64].")
	fmt.Println("")
}
func (p *ZipFileNameProcess) SetNext(process *Process) {
	p.NextProcess = process
}
func (p *ZipFileNameProcess) Next() *Process {
	return p.NextProcess
}

func (p *ZipFileNameProcess) SetParameter(parameter interface{}) {
	p.Parameter = parameter.(PackageParameter)
}

func NewZipFileNameProcess(file *os.File) *Process {
	var p Process
	p = &ZipFileNameProcess{
		Stdin: file,
	}
	return &p
}
