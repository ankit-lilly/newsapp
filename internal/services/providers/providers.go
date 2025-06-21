package providers

import (
	"fmt"
	"sort"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/acmqueue"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/davecheney"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/devto"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/fiercepharma"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/hackernoon"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/highscalability"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/hindu"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/martinfowler"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/natgeo"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/newyorker"
	scientificamerican "github.com/ankit-lilly/newsapp/internal/services/providers/sources/scientificAmerican"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/signalsAndthreads"
	softwareengineeringdaily "github.com/ankit-lilly/newsapp/internal/services/providers/sources/softwareEngineeringDaily"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/techcrunch"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources/wired"
)

/*
Providers interface is used to define the methods that each provider should implement.
*/
type Providers interface {
	GetCategories() map[string]string
	GetID() string
	GetName() string
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

func Register(providers []Providers) {
	for _, p := range providers {
		Registry[p.GetID()] = p
	}
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

		portalcats.ID = portal.GetID()
		portalcats.Name = portal.GetName()
		portalcats.Categories = portal.GetCategories()
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
	fiercepharma := fiercepharma.NewFiercePharma()
	davecheney := davecheney.NewDaveCheney()
	martinfowler := martinfowler.NewMartinFowler()
	techcrunch := techcrunch.NewTechcrunch()
	wired := wired.NewWired()
	natgeo := natgeo.NewNatGeo()
	scientificamerican := scientificamerican.NewScientificAmerican()
	highscalability := highscalability.NewHighScalability()
	hackernoon := hackernoon.NewHackerNoon()
	newyorker := newyorker.NewNewyorker()
	acmqueue := acmqueue.NewACMQueue()
	devto := devto.NewDevTo()
	softwareengineeringdaily := softwareengineeringdaily.NewSFD()
	signalsAndthreads := signalsAndthreads.NewSignalsAndThreads()

	Register([]Providers{
		thehindu,
		natgeo,
		fiercepharma,
		davecheney,
		wired,
		martinfowler,
		techcrunch,
		scientificamerican,
		highscalability,
		hackernoon,
		newyorker,
		acmqueue,
		devto,
		softwareengineeringdaily,
		signalsAndthreads,
	})
}
