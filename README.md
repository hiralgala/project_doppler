# project_doppler

A CLI tool which executes a "substitute" command which replaces variable expressions (e.g. ${DATABASE_URL}) in
static files with the respective Doppler secret.

## Requirements

* Go version 1.15

## How to run this code

* Input file path (Required)
* Output Directory (Optional)
* Variable formats called Pattern (Optional)

### From command line:    
```bash
$ go run <file path to your source code> input <input file path> 
```

### Optional Arguments
```
optional arguments:
  output <output dir>
  pattern <variable string>
```

You need to replace:

1. {file path to your file} (required) - file path to your source code
2. {input file path} (required) - path to the file in which you want to replace the secrets
3. {output dir} (optional) - path to the directory where you want to save the file after substitutions
4. {variable string} (optional) - the variable format you want to look for in order to find their respective substitutions


#### Different patterns to choose from:
* dollar
* dollar-curly
* handlebars
* dollar-handlebars


#### Steps to test for your doppler project
* Replace the var token in the source code with your doppler token.
* Include some of your secret keys formed using chosen pattern in your input files. For example, if your secret key is DATABASE and your choice of pattern is "handlebars", then include "{{DATABASE}}" in your input files to see the final substitution once the code is run.
* Run the code to see your substitutions work.

#### Sample command with all the arguments:    
```bash
$ go run <file path to your source code> input <input file path> output <output dir> pattern <variable string>  
```
```bash
$ go run C:\Users\Hiral\go\bin\project.go input "C:\Users\Hiral\go\Input.txt" output "C:\Users\Hiral\go\Output" pattern "handlebars"  
```


