package store

import (
	"errors"
	"log/slog"
	"maps"
	"sort"
	"strconv"
	"strings"
)

type Contact struct {
	Id    string
	First string
	Last  string
	Phone string
	Email string
}
type Contacts struct {
	contacts map[string]Contact
	nextId   int
}

func NewContacts() *Contacts {
	return &Contacts{contacts: make(map[string]Contact)}
}

func (c *Contacts) Get(id string) Contact {
	return c.contacts[id]
}

func (c *Contacts) All() []Contact {
	result := make([]Contact, 0, len(c.contacts))
	for _, contact := range c.contacts {
		result = append(result, contact)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Last < result[j].Last
	})
	return result
}

func (c *Contacts) Search(query string) []Contact {
	var result []Contact
	query = strings.ToLower(query)
	for _, contact := range c.contacts {
		fullname := contact.First + " " + contact.Last
		if strings.Contains(strings.ToLower(fullname), query) {
			result = append(result, contact)
		}
	}
	return result
}

func (c Contacts) keys() string {
	keys := make([]string, 0, len(c.contacts))
	for key := range maps.Keys(c.contacts) {
		keys = append(keys, key)
	}
	return strings.Join(keys, ",")
}

func (c *Contacts) Add(contact Contact) (string, error) {
	c.nextId = c.nextId + 1
	contact.Id = strconv.Itoa(c.nextId)
	c.contacts[contact.Id] = contact
	slog.Info("store updated","keys", c.keys())
	return contact.Id, nil
}

func (c *Contacts) Update(id string, contact Contact) error {
	if _, ok := c.contacts[id]; !ok {
		return errors.New("contact not found")
	}
	contact.Id = id
	c.contacts[id] = contact
	return nil
}

func (c *Contacts) Delete(id string) error {
	if _, ok := c.contacts[id]; !ok {
		return errors.New("contact not found")
	}
	delete(c.contacts, id)
	slog.Info("store updated","keys", c.keys())
	return nil
}

type Store interface {
	Add(contact Contact) (string, error)
	Search(query string) []Contact
	All() []Contact
	Get(id string) Contact
	Update(id string, contact Contact) error
	Delete(id string) error
}
