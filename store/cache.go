package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	updateInterval  = time.Hour * 1
	firebaseTimeout = time.Minute * 5
)

// Cache provides cahcing layer to limit Firestore usage
// and improve response time
type Cache struct {
	mux  sync.RWMutex
	ctx  context.Context
	data *cacheData
}

// cacheData internal struct that holds actual cache data
type cacheData struct {
	categories       []*Category
	categoriesByName map[string]*Category
	items            map[int]*Item
	itemsByName      map[string]*Item
}

// NewCache returns a pointer to a new Cache instance already populated
// with fresh data
func NewCache(ctx context.Context) (*Cache, error) {
	c := &Cache{
		ctx:  ctx,
		data: new(cacheData),
	}

	if err := c.data.populate(c.ctx); err != nil {
		return nil, fmt.Errorf("failed to populate cache: %w", err)
	}

	go c.updateLoop()

	return c, nil
}

// populate connects to firebase, fetches categories and items
// collections and populates *ByName index maps
func (c *cacheData) populate(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, firebaseTimeout)
	defer cancel()

	db, err := New(timeoutCtx)
	if err != nil {
		return err
	}

	c.categories, err = db.GetCategories(timeoutCtx)
	if err != nil {
		return err
	}

	c.categoriesByName = make(map[string]*Category, len(c.categories))
	for i := range c.categories {
		c.categoriesByName[c.categories[i].Name] = c.categories[i]
	}

	c.items, err = db.GetItems(timeoutCtx)
	if err != nil {
		return err
	}

	c.itemsByName = make(map[string]*Item, len(c.items))
	for i := range c.items {
		c.itemsByName[c.items[i].Name] = c.items[i]
	}

	return nil
}

// updateLoop periodically fetches new data from firebase and
// changes data pointer to point at a fresh data
func (c *Cache) updateLoop() {
	ticker := time.NewTicker(updateInterval)
	for {
		select {
		case <-c.ctx.Done():
			log.Print("INFO: cache updater stopped")
			return
		case <-ticker.C:
		}

		data := new(cacheData)
		if err := data.populate(c.ctx); err != nil {
			log.Printf("ERROR: failed to populate cache: %v", err)
			continue
		}

		c.mux.Lock()
		c.data = data
		c.mux.Unlock()
	}
}

// GetCategoriesPage returns one page of categories from cache
func (c *Cache) GetCategoriesPage(pageNum, pageSize int) []*Category {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if len(c.data.categories) < pageNum*pageSize {
		return nil
	}

	from := pageNum * pageSize
	to := (pageNum + 1) * pageSize
	if to > len(c.data.categories) {
		to = len(c.data.categories)
	}

	return c.data.categories[from:to]
}

// GetItemsPage returns one page of category's items from cache
func (c *Cache) GetItemsPage(categoryName string, pageNum, pageSize int) ([]*Item, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	cat, ok := c.data.categoriesByName[categoryName]
	if !ok {
		return nil, sql.ErrNoRows
	}

	if len(cat.Products) < pageNum*pageSize {
		return nil, nil
	}

	from := pageNum * pageSize
	to := (pageNum + 1) * pageSize
	if to > len(cat.Products) {
		to = len(cat.Products)
	}

	// TODO: can preallocate len
	// and/or cache pages
	products := []*Item{}
	for i := from; i < to; i++ {
		item, ok := c.data.items[cat.Products[i]]
		if !ok {
			log.Printf(
				"ERROR: product_id %d from category_id %d is not found",
				cat.Products[i],
				cat.ID,
			)
			continue
		}
		products = append(products, item)
	}

	return products, nil
}

// GetItem returns menu item by its name
func (c *Cache) GetItem(itemName string) (*Item, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	item, ok := c.data.itemsByName[itemName]
	if !ok {
		return nil, sql.ErrNoRows
	}

	return item, nil
}
