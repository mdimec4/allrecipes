# allrecipes
Service and library that parses allrecipes.com recipes. 

This library has only one function which accepts allrecipes.com recipe ID and will as result returns paresd recipe as an object.

## Include parser in your project

```go
import "github.com/mdimec4/allrecipes"

var recipe allrecipes.Recipe
var err error

recipe, error := alrecipes.GetRecipe("11772")
fmt.Println(recipe.Ingredients)
```

## Usage via provided web api which parses recipe into JSON format
```
$ go build github.com/mdimec4/allrecipes/cmd/webapi
$ ./webapi &
```
Example:
To parse recipe: ```https://www.allrecipes.com/recipe/11772/spaghetti-pie-i/``` you would query this service with

```
$ curl http://localhost:4007/api/recipe/11772
```

to get the recipe in JSON format. 11772 is recipe ID.


Use ```ALLRECIPES_PARSER_LISTEN_ADDR``` environment variable to modify service listen address and/or port.
Default ```ALLRECIPES_PARSER_LISTEN_ADDR``` value is ```:4007```.

## Docker image

Use provided ``Dockerfile`` to buld docker image.
To build for ARM, you need to first generate ```Dockerfile_arm``` with ```generate_Dockerfile_arm.sh```.

Dockerhub build (not neceserily up to date) is also avaliable: https://hub.docker.com/r/mihad/allrecipes_parser/
