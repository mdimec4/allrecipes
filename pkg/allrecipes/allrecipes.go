package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Recipe struct {
	RecipeID     string   `json:"recipe_id"`
	Publisher    string   `json:"publisher"`
	SourceURL    string   `json:"source_url"`
	Title        string   `json:"title"`
	ImageURL     string   `json:"image_url"`
	Description  string   `json:"description"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Tips         []string `json:"tips"`
}

func delNewLine(s string) string {
	return strings.Replace(
		strings.Replace(s, "\n", "", -1),
		"\r", "", -1)
}

func checkAttr(attr []html.Attribute, key, val string) bool {
	for _, a := range attr {
		if a.Key == key && a.Val == val {
			return true
		}
	}
	return false
}

func getAttrVal(attr []html.Attribute, key string) string {
	for _, a := range attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func GetRecipe(recipeUrl string) (Recipe, error) {
	// get recipe id from url
	u, err := url.Parse(recipeUrl)
	if err != nil {
		return Recipe{}, err
	}
	if u.Host != "allrecipes.com" {
		return Recipe{}, errors.New("expected allrecipes.com host name")
	}
	//remove Qury part from URL
	u.RawQuery = ""

	// parse html
	resp, err := http.Get(u.String())
	if err != nil {
		return Recipe{}, err
	}
	defer resp.Body.Close()
	ret := Recipe{RecipeID: u.String(), SourceURL: u.String()}

	z := html.NewTokenizer(resp.Body)
endloop:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				break endloop
			}
			return Recipe{}, z.Err()
		case html.StartTagToken:
			token := z.Token()
			if token.DataAtom == atom.H1 &&
				checkAttr(token.Attr, "itemprop", "name") {
				// <h1 class="recipe-summary__h1" itemprop="name">Spaghetti Pie I</h1>
				tt := z.Next()
				switch tt {
				case html.TextToken:
					token = z.Token()
					ret.Title = delNewLine(html.UnescapeString(token.Data))
					fmt.Println("Title>", ret.Title)
				case html.ErrorToken:
					return Recipe{}, z.Err()
				default:
					return Recipe{}, errors.New("allrecipes parser: author name was expected here")
				}
			} else if token.DataAtom == atom.Span &&
				checkAttr(token.Attr, "itemprop", "author") {
				// <span class="submitter__name" itemprop="author">Kimberley</span>
				tt := z.Next()
				switch tt {
				case html.TextToken:
					token = z.Token()
					ret.Publisher = delNewLine(html.UnescapeString(token.Data))
					fmt.Println("Publisher>", ret.Publisher)
				case html.ErrorToken:
					return Recipe{}, z.Err()
				default:
					return Recipe{}, errors.New("allrecipes parser: author name was expected here")
				}
			} else if token.DataAtom == atom.Div &&
				checkAttr(token.Attr, "itemprop", "description") {
				// <div class="submitter__description" itemprop="description"> "Family favorite. Serve with lemon wedges."</div>
				tt := z.Next()
				switch tt {
				case html.TextToken:
					token = z.Token()
					ret.Description = delNewLine(html.UnescapeString(token.Data))
					fmt.Println("Description>", ret.Description)
				case html.ErrorToken:
					return Recipe{}, z.Err()
				default:
					return Recipe{}, errors.New("allrecipes parser: author name was expected here")
				}
			} else if token.DataAtom == atom.Span &&
				checkAttr(token.Attr, "itemprop", "ingredients") {
				// did we hit one of the ingredients
				// <span class="recipe-ingred_txt added" ... itemprop="ingredients">
				tt := z.Next()
				// next token should be text of the ingredient span
				switch tt {
				case html.TextToken:
					token = z.Token()
					ret.Ingredients = append(ret.Ingredients,
						delNewLine(html.UnescapeString(token.Data)))
					fmt.Println("Ingredient>", ret.Ingredients[len(ret.Ingredients)-1])
				case html.ErrorToken:
					return Recipe{}, z.Err()
				default:
					return Recipe{}, errors.New("allrecipes parser: ingredient text was expected here")
				}
			} else if token.DataAtom == atom.Span &&
				checkAttr(token.Attr, "class", "recipe-directions__list--item") &&
				!checkAttr(token.Attr, "ng-bind", "model.itemNote") {
				// did we hit one of the instructions
				// <span class="recipe-directions__list--item" ...>
				tt := z.Next()
				// next token should be text of the instruction span
				switch tt {
				case html.TextToken:
					token = z.Token()
					ret.Instructions = append(ret.Instructions,
						delNewLine(html.UnescapeString(token.Data)))
					fmt.Println("Instruction", ret.Instructions[len(ret.Instructions)-1])
				case html.ErrorToken:
					return Recipe{}, z.Err()
				default:
					return Recipe{}, errors.New("allrecipes parser: instruction text was expected here")
				}
			}
		case html.SelfClosingTagToken:
			token := z.Token()
			if token.DataAtom == atom.Meta &&
				checkAttr(token.Attr, "property", "og:image") {
				// <meta property="og:image" content="https://images.media-allrecipes.com/userphotos/560x315/726090.jpg" />
				imgURL := getAttrVal(token.Attr, "content")
				fmt.Println("Image>", imgURL)
				ret.ImageURL = imgURL
			}

		}

	}
	return ret, nil
}

func main() {
	//url := "http://allrecipes.com/recipe/231495/texas-boiled-beer-shrimp/"
	url := "http://allrecipes.com/recipe/11772/spaghetti-pie-i/?clickId=right%20rail0&internalSource=rr_feed_recipe_sb&referringId=231495%20referringContentType%3Drecipe"
	recipe, err := GetRecipe(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err) // TODO stderr
		return
	}
	fmt.Printf("\nrecipe: %+v\n", recipe)

}
