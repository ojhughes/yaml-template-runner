package main

import (
	"text/template"
	"log"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"

	"strings"
	"path/filepath"
	"flag"
)

type SubBlock struct {
	Name string `yaml:"name"`
	Template string `yaml:"template"`
	Vars map[string]string `yaml:"vars,omitempty"`
}

type Block struct {
	Name string `yaml:"name"`
	Template string `yaml:"template"`
	Vars map[string]string `yaml:"vars,omitempty"`
	SubBlocks []SubBlock `yaml:"subBlocks,omitempty"`

}

type TemplateDataSchema struct {
	Blocks []Block
}

const (
	pipelineTemplateName   = "pipelineTemplateName"
)

func main() {
	dataFilePath := flag.String("d", "", "Path to YAML templating schema definition")
	templateFilePath := flag.String("t", "", "Path to layout template")
	templatesDirectoryPath := flag.String("r", "resources", "Path to layout template")
	flag.Parse()

	dataFile, err := ioutil.ReadFile(filepath.Join(*templatesDirectoryPath, *dataFilePath))
	check(err)

	var dataSchema TemplateDataSchema
	err = yaml.Unmarshal(dataFile, &dataSchema)
	check(err)

	check(err)

	pipelineTemplateData, err := ioutil.ReadFile(filepath.Join(*templatesDirectoryPath, *templateFilePath))
	check(err)

	pipelineTemplate := template.New(pipelineTemplateName).Funcs(template.FuncMap{"indent": indent})


	flattenedData := make(map[string]map[string]interface{})
	for _, block := range dataSchema.Blocks  {
		blockTemplateData, err := ioutil.ReadFile(filepath.Join(*templatesDirectoryPath, block.Template))
		check(err)
		blockVars := make(map[string]interface{})
		template.Must(pipelineTemplate.New(block.Name).Parse(string(blockTemplateData)))
		for _, subBlock := range block.SubBlocks  {
			template.Must(pipelineTemplate.New(subBlock.Name).ParseFiles(filepath.Join(*templatesDirectoryPath, subBlock.Template)))
			subBlockVars := make(map[string]interface{})
			addAllEntries(block.Vars, subBlockVars)
			addAllEntries(subBlock.Vars, subBlockVars)
			blockVars[subBlock.Name] = subBlockVars
		}
		addAllEntries(block.Vars, blockVars)
		flattenedData[block.Name] = blockVars
	}
	template.Must(pipelineTemplate.New(pipelineTemplateName).Parse(string(pipelineTemplateData)))
	err = pipelineTemplate.ExecuteTemplate(os.Stdout, "pipelineTemplate", flattenedData)
	check(err)

}

func addAllEntries(sourceMap map[string]string, destMap map[string]interface{}) {
	for key, value := range sourceMap {
		destMap[key] = value
	}
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}


