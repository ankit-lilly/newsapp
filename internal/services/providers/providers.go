package providers

import (
	"fmt"
	"sort"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/davecheney"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/fiercepharma"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/hindu"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/hindustan"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/martinfowler"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/natgeo"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/techcrunch"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/verge"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/wired"
)

/*
Providers interface is used to define the methods that each provider should implement.
*/
type Providers interface {
	Categories() map[string]string
	ID() string
	Name() string
	FeedURL(category string) string
	Fetch(category string) ([]models.Article, error)
	Parse(url string) (models.Article, error)
	HasCategories() bool
	IsCategoryValid(category string) bool
}

type ProviderCategory struct {
	ID          string
	Name        string
	HasChildren bool
	Categories  map[string]string
}

// Registry is a map of provider ID to the provider instance.
var Registry = map[string]Providers{}

func Register(p Providers) {
	Registry[p.ID()] = p
}

func Get(id string) (Providers, error) {
	p, ok := Registry[id]
	if !ok {
		return nil, fmt.Errorf("portal with id %q not found", id)
	}
	return p, nil
}

/*
GetProviderCategories returns the list of provider categories.
This is used to display the list of providers in the UI(in the Navbar)
*/

func GetProviderCategories() []ProviderCategory {

	var providersCategory []ProviderCategory

	for _, portal := range Registry {

		var portalcats ProviderCategory

		portalcats.ID = portal.ID()
		portalcats.Name = portal.Name()
		portalcats.Categories = portal.Categories()
		portalcats.HasChildren = portal.HasCategories()
		providersCategory = append(providersCategory, portalcats)
	}

	sort.Slice(providersCategory, func(i, j int) bool {
		return providersCategory[i].Name < providersCategory[j].Name
	})

	return providersCategory
}

/*
Init registers all the providers.
*/

func Init() {
	thehindu := hindu.NewTheHinduCom()
	hindustanTimes := hindustan.NewHindusTanTimes()
	wired := wired.NewWired()
	techcrunch := techcrunch.NewTechcrunch()
	verge := verge.NewVerge()
	fiercepharma := fiercepharma.NewFiercePharma()
	martinfowler := martinfowler.NewMartinFowler()
	davecheney := davecheney.NewDaveCheney()
  natgeo := natgeo.NewNatGeo()

	Register(thehindu)
	Register(hindustanTimes)
	Register(techcrunch)
	Register(wired)
	Register(verge)
	Register(fiercepharma)
	Register(davecheney)
	Register(martinfowler)
  Register(natgeo)
}
