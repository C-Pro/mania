package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
)

// DB core object for firebase storage interface
type DB struct {
	app *firebase.App
	cl  *firestore.Client
}

// New creates new DB object
func New(ctx context.Context) (*DB, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return nil, err
	}

	cl, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error connecting to db: %v\n", err)
		return nil, err
	}

	d := DB{
		app: app,
		cl:  cl,
	}

	return &d, nil
}

// Item holds menu item data
type Item struct {
	ID          int `json:"product_id"`
	Name        string
	Image       string
	Price       float64
	Composition string
	Description string
}

// Category is a menu category description
type Category struct {
	ID       int `json:"category_id"`
	Icon     string
	Name     string
	ParentID int `json:"parent_id"`
	Products []int
}

// GetCategories returns categories from Firestore
func (db *DB) GetCategories(ctx context.Context) ([]*Category, error) {
	cats := []*Category{}
	iter := db.cl.Collection("categories").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return cats, err
		}
		cat, err := mapToCategory(doc.Data())
		if err != nil {
			return cats, err
		}

		cats = append(cats, &cat)
	}
	return cats, nil
}

// mapToCategory deals with mess of an interfaces
// firebase admin api returns
// doc.DataAs does not work well out of the box, I checked :)
func mapToCategory(m map[string]interface{}) (Category, error) {
	cat := Category{}

	s, ok := m["category_id"].(string)
	if !ok {
		return cat, errors.New("bad category id")
	}
	id, err := strconv.Atoi(s)
	if err != nil {
		return cat, fmt.Errorf("failed to convert category_id: %w", err)
	}

	cat.ID = id

	s, ok = m["parent_id"].(string)
	if !ok {
		return cat, errors.New("bad category parent_id")
	}

	id, err = strconv.Atoi(s)
	if err != nil {
		return cat, err
	}

	cat.ParentID = id

	cat.Name, ok = m["name"].(string)
	if !ok {
		return cat, errors.New("bad category name")
	}
	cat.Icon = m["icon"].(string)
	if !ok {
		return cat, errors.New("bad category icon")
	}

	products, ok := m["products"].([]interface{})
	if !ok {
		return cat, errors.New("bad products list")
	}
	for _, pi := range products {
		p, ok := pi.(map[string]interface{})
		if !ok {
			return cat, errors.New("bad product value")
		}
		s, ok := p["product_id"].(string)
		if !ok {
			return cat, errors.New("bad product id")
		}
		id, err := strconv.Atoi(s)
		if err != nil {
			return cat, fmt.Errorf("failed to convert product_id: %w", err)
		}

		cat.Products = append(cat.Products, id)

	}

	return cat, nil
}

// GetItems returns menu items from Firestore
func (db *DB) GetItems(ctx context.Context) (map[int]*Item, error) {
	items := make(map[int]*Item)
	iter := db.cl.Collection("products").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return items, err
		}
		item, err := mapToItem(doc.Data())
		if err != nil {
			return items, err
		}

		items[item.ID] = &item
	}
	return items, nil
}

// mapToItem deals with mess of an interfaces
// firebase admin api returns
// doc.DataAs does not work well out of the box, I checked :)
func mapToItem(m map[string]interface{}) (Item, error) {
	item := Item{}

	s, ok := m["product_id"].(string)
	if !ok {
		return item, errors.New("bad product id")
	}
	id, err := strconv.Atoi(s)
	if err != nil {
		return item, fmt.Errorf("failed to convert product_id: %w", err)
	}

	item.ID = id

	item.Name, ok = m["name"].(string)
	if !ok {
		return item, errors.New("bad product name")
	}
	item.Name = cleanupString(item.Name)

	item.Composition, ok = m["composition"].(string)
	if !ok {
		return item, errors.New("bad product composition")
	}
	item.Composition = fixNutrients(cleanupString(item.Composition))

	item.Description, ok = m["description"].(string)
	if !ok {
		return item, errors.New("bad product description")
	}
	item.Description = cleanupString(item.Description)

	s, ok = m["price"].(string)
	if !ok {
		return item, errors.New("bad product price")
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return item, fmt.Errorf("failed to convert product price: %w", err)
	}

	item.Price = f

	return item, nil
}

var (
	// matches html entities and unicode BOM mark
	removeRe = regexp.MustCompile(`&[^;]+;|\x{feff}`)
	spacesRe = regexp.MustCompile(`\s+`)
	// matches nutrients composition strings lacking spaces like
	// "соусБ-20,5Ж-27,4У-6,3 Ккал-257"
	nutrientsRe = regexp.MustCompile(`(Б-[0-9,.]+)(Ж-[0-9,.]+)(У-[0-9,.]+)`)
)

// cleanupString removes newlines, html entities, normalizes spaces
func cleanupString(s string) string {
	s = removeRe.ReplaceAllString(s, "")
	return spacesRe.ReplaceAllString(s, " ")
}

// fixNutrients fixes composition nutrients info formatting
// composition strings for some reason do not have spaces between
// nutrients, like: "фирменный соусБ-20,5Ж-27,4У-6,3 Ккал-257"
func fixNutrients(s string) string {
	return nutrientsRe.ReplaceAllString(s, " $1%, $2%, $3%, пищевая ценность: $4")
}
