
Simple PoC of how YAML can be templated using reusable `blocks` of content with associated variable data.
This project is a simple wrapper of the Go template engine so all the features 
in the documentation should work: [https://golang.org/pkg/text/template/]() 

## Usage
Checkout into $GOPATH/src
* `go install`
* `$GOPATH/bin/template-runner -d data-vars.yml -t template-layout.yml -r resources` 
* or `go run template-runner.go` -d "" -r "" -t ""
    * -d is the variable definition file
    * -t is the top-level layout YAML file
    * -r is is the directory where the YAML files are

### Block definition
```
blocks:
- name: removeTile112
  template: block-remove-tile.yml
  vars:
    pcfVersion: 1-12
  subBlocks:
  - name: moveToLaundromat
    template: block-move-to-laundromat.yml
    vars:
      pool: dirty-pool-no-tile

```
Each block is a partial YAML file that will be injected into the parent YAML layout file
along with associated variables. The `name:` key must be referenced in the layout YAML
file as a `{{template}}`. 

So for example, a block with the name `myBlock` and template `myBlock.yml` will need a corresponding template inclusion: `{{- template "myBlock" .myBlock -}}`
This syntax injects the contents of `myBlock.yml` and also passes a `map[string]string` of variables to the injected block

The block template file can reference any variables from the parent, and also include a `SubBlock`. 
To reference variables in a `block` template:
```
{{- $myVar := index . "myvar" -}}
Hello, this is the value of {{$myVar}}
```
### SubBlock definition

A SubBlock points to a `template` file that will be injected into a `block` (in exactly the same way as a block is injected into a layout template)
The parameters for defining a SubBlock are the same for a `block`:
```
name: *Name used to reference the (sub)block in the {{template}} inclusion*
template: *name of the YAML to be included*
subBlocks: *List of nested blocks to be included in the block (SubBlocks cannot contain SubBlocks)
vars: *A map of vars to be passed to the injected template*
    myVar: foo
    otherVar: bar
```
### Layout file definition

And finally the layout file which wires up the blocks with the data, eg;

```
{{- define "pipelineTemplate"}}
jobs:
{{- "\n"}}
{{- template "removeTile112" .removeTile112 -}}
{{- template "removeTile20" .removeTile20 -}}
{{end -}}
```

The final result of this example would be;

```
jobs:
- name: remove-tile-1-12
  max_in_flight: 6
  on_success:
  - do:
    - put: dirty-pool-no-tile-1-12
      params:
        add: dirty-pool-no-tile-1-12
    - put: dirty-pool-no-tile-1-12
      params:
        remove: dirty-pool-no-tile-1-12
- name: remove-tile-2-0
  max_in_flight: 6
  on_success:
  - do:
    - put: dirty-pool-no-tile-2-0
      params:
        add: dirty-pool-no-tile-2-0
    - put: dirty-pool-no-tile-2-0
      params:
        remove: dirty-pool-no-tile-2-0
```

### Issues
YAML indentation is sometimes not correct. Ensure that whitespace is removed by using hyphens
 before / after template markers eg. `{{-` and `-}}`.
 
Indentation can be forced if needed using `{{- "" | indent 2 -}}` before the text that needs indenting.
Maybe this could be fixed by parsing the template as YAML..