package main

import (
	"text/template"
	"log"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"

)

type Block struct {
	Name string `yaml:"name"`
	Vars map[string]string `yaml:"vars,omitempty"`
}

type TemplateDataSchema struct {
	Blocks []Block
}

const pipelineTemplateName = "pipelineTemplateName"

func main() {
	dataFile, err := ioutil.ReadFile("resources/data-remove-tile.yml")
	check(err)

	var dataSchema TemplateDataSchema
	err = yaml.Unmarshal(dataFile, &dataSchema)
	check(err)

	blockTemplateData, err := ioutil.ReadFile("resources/block-remove-tile.yml")
	check(err)

	pipelineTemplateData, err := ioutil.ReadFile("resources/template-pipeline-laundromat.yml")
	check(err)

	pipelineTemplate := template.New(pipelineTemplateName)
	flattenedData := make(map[string]map[string]string)
	for _, block := range dataSchema.Blocks  {
		template.Must(pipelineTemplate.New(block.Name).Parse( string(blockTemplateData)))
		flattenedData[block.Name] = block.Vars
	}

	template.Must(pipelineTemplate.New(pipelineTemplateName).Parse(string(pipelineTemplateData)))
	err = pipelineTemplate.ExecuteTemplate(os.Stdout, "pipelineTemplate", flattenedData)
	check(err)

}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}


