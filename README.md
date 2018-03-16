# allrecipes
Library parses allrecipes.com recipes. 

This library has only one function which accepts allrecipes.com recipe ID and will as result returns paresd recipe as an object.

## Include parser in your project

```go
import "github.com/mdimec4/allrecipes.com_parser"

var recipe allrecipes.Recipe
var err error

recipe, error := alrecipes.GetRecipe("11772")
fmt.Println(recipe.Ingredients)
```

## via provided web api which parses recipe into JSON format
Example:
To parse recipe: ```https://www.allrecipes.com/recipe/11772/spaghetti-pie-i/``` you would query this service with

```
$ curl http://localhost:4007/api/recipe/11772
```

to get the recipe in JSON format. 11772 is recipe ID.


Use ALLRECIPES_PARSER_LISTEN_ADDR environment variable to modify service listen address and/or port.
Default ALLRECIPES_PARSER_LISTEN_ADDR value is :4007.
