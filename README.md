
Simple PoC of how YAML can be templated with reusable blocks and associated data

3 files need to be provided;
A data YAML file that describes the variables that partial YAML blocks with use eg;

```
blocks:
  - name: removeTile112
    vars:
      pcfVersion: 1-12

  - name: removeTile20
    vars:
      pcfVersion: 2-0
```

The actual blocks that will be injected into the layout template with the data variables, eg
```
- name: remove-tile-{{index . "pcfVersion"}}
  max_in_flight: 6
```

And finally the layout file which wires up the blocks with the data, eg;

```
{{define "pipelineTemplate"}}
jobs:
{{template "removeTile112" .removeTile112}}
{{template "removeTile20" .removeTile20}}
{{end}}
```

The final result of this example would be;

```
jobs:
- name: remove-tile-1-12
  max_in_flight: 6

- name: remove-tile-2-0
  max_in_flight: 6
```

The next feature to implement would be the ability to define nested blocks