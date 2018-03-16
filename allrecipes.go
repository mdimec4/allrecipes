package allrecipes

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Recipe struct {
	RecipeID    string   `json:"recipe_id"`
	Author      string   `json:"author"`
	SourceURL   string   `json:"source_url"`
	Name        string   `json:"name"`
	ImageURL    string   `json:"image_url"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
	Directions  []string `json:"directions"`
	Footnotes   []string `json:"footnotes"`
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

func GetRecipe(recipeID string) (Recipe, error) {
	// get recipe id from url
	u, err := url.Parse("https://www.allrecipes.com")
	if err != nil {
		return Recipe{}, fmt.Errorf("parse error %s", err)
	}
	u.Path = path.Join("recipe", recipeID)

	// parse html
	resp, err := http.Get(u.String())
	if err != nil {
		return Recipe{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK /*200*/ {
		return Recipe{}, fmt.Errorf("allrecipes.com responded with: %s", resp.Status)
	}
	ret := Recipe{RecipeID: recipeID, SourceURL: resp.Request.URL.String()}

	z := html.NewTokenizer(resp.Body)
endloop:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				break endloop
			}
			return Recipe{}, fmt.Errorf("main parser loop error: %s", z.Err())
		case html.StartTagToken:
			token := z.Token()
			if token.DataAtom == atom.H1 &&
				checkAttr(token.Attr, "itemprop", "name") {
				// <h1 class="recipe-summary__h1" itemprop="name">Spaghetti Pie I</h1>
				tt := z.Next()
				switch tt {
				case html.TextToken:
					token = z.Token()
					ret.Name = delNewLine(html.UnescapeString(token.Data))
					// fmt.Println("Name>", ret.Name)
				case html.ErrorToken:
					return Recipe{}, fmt.Errorf("name text err: %s", z.Err())
				default:
					return Recipe{}, errors.New("allrecipes parser: recipe name text was expected here")
				}
			} else if token.DataAtom == atom.Span &&
				checkAttr(token.Attr, "itemprop", "author") {
				// <span class="submitter__name" itemprop="author">Kimberley</span>
				tt := z.Next()
				switch tt {
				case html.TextToken:
					token = z.Token()
					ret.Author = delNewLine(html.UnescapeString(token.Data))
					// fmt.Println("Author>", ret.Author)
				case html.ErrorToken:
					return Recipe{}, fmt.Errorf("author text err: %s", z.Err())
				default:
					return Recipe{}, errors.New("allrecipes parser: author name text was expected here")
				}
			} else if token.DataAtom == atom.Div &&
				checkAttr(token.Attr, "itemprop", "description") {
				// <div class="submitter__description" itemprop="description"> "Family favorite. Serve with lemon wedges."</div>
				tt := z.Next()
				switch tt {
				case html.TextToken:
					token = z.Token()
					ret.Description = delNewLine(html.UnescapeString(token.Data))
					// fmt.Println("Description>", ret.Description)
				case html.ErrorToken:
					return Recipe{}, fmt.Errorf("description text err: %s", z.Err())
				default:
					return Recipe{}, errors.New("allrecipes parser: description text was expected here")
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
					// fmt.Println("Ingredient>", ret.Ingredients[len(ret.Ingredients)-1])
				case html.ErrorToken:
					return Recipe{}, fmt.Errorf("ingredient text err: %s", z.Err())
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
					ret.Directions = append(ret.Directions,
						delNewLine(html.UnescapeString(token.Data)))
					// fmt.Println("Instructioni>", ret.Directions[len(ret.Directions)-1])
				case html.ErrorToken:
					return Recipe{}, fmt.Errorf("direction text err: %s", z.Err())
				default:
					return Recipe{}, errors.New("allrecipes parser: direction text was expected here")
				}
			} else if token.DataAtom == atom.Span &&
				checkAttr(token.Attr, "class", "recipe-footnotes__header") {
				// did we hit footnotes
				// <span class="recipe-footnotes__header">Nutrition:</span>
				tt := z.Next()
				// next token should be text of the footnotes title
				about := ""
				switch tt {
				case html.TextToken:
					token = z.Token()
					about = html.UnescapeString(token.Data)
					// fmt.Println("footnotes about:", about)
				endloopFootnotes:
					for {
						// fast forward to next "Li" token
						tt := z.Next()
						switch tt {
						case html.StartTagToken:
							token = z.Token()
							if token.DataAtom == atom.Li {
								tt := z.Next()
								// next token should be text of the instruction <li>
								switch tt {
								case html.TextToken:
									token = z.Token()
									ret.Footnotes = append(ret.Footnotes,
										about+" "+delNewLine(html.UnescapeString(token.Data)))
									// fmt.Println("Footnotes>", ret.Footnotes[len(ret.Footnotes)-1])
									break endloopFootnotes
								case html.ErrorToken:
									return Recipe{}, fmt.Errorf("footnotes text err: %s", z.Err())
								default:
									return Recipe{}, errors.New("allrecipes parser: footnotes text was expected here")
								}
							}
						case html.ErrorToken:
							return Recipe{}, fmt.Errorf("footnotest <li> err: %s", z.Err())
						}
					}
				case html.ErrorToken:
					return Recipe{}, fmt.Errorf("footnotest title text err: %s", z.Err())
				default:
					return Recipe{}, errors.New("allrecipes parser: footnotes title text was expected here")
				}
			}
		case html.SelfClosingTagToken:
			token := z.Token()
			if token.DataAtom == atom.Meta &&
				checkAttr(token.Attr, "property", "og:image") {
				// <meta property="og:image" content="https://images.media-allrecipes.com/userphotos/560x315/726090.jpg" />
				imgURL := getAttrVal(token.Attr, "content")
				// fmt.Println("Image>", imgURL)
				ret.ImageURL = imgURL
			}

		}

	}
	return ret, nil
}

/*
func main() {
	//url := "http://allrecipes.com/recipe/231495/texas-boiled-beer-shrimp/"
	//url := "http://allrecipes.com/recipe/11772/spaghetti-pie-i/?clickId=right%20rail0&internalSource=rr_feed_recipe_sb&referringId=231495%20referringContentType%3Drecipe"
	recipe, err := GetRecipe("231495")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err) // TODO stderr
		return
	}
	fmt.Printf("\nrecipe: %+v\n", recipe)

}
*/
