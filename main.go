package main

import (
	"log"
	"strings"

	"github.com/ksiezykm/FerretMate/db"
	"github.com/ksiezykm/FerretMate/list"
	"github.com/ksiezykm/FerretMate/model"
	"github.com/ksiezykm/FerretMate/notepad"
	"github.com/ksiezykm/FerretMate/popup"

	"github.com/awesome-gocui/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = false

	connections, err := model.LoadConnections()
	if err != nil {
		log.Panicln(err)
	}

	var connNames []string
	for _, c := range connections {
		connNames = append(connNames, c.Name)
	}

	m := &model.Model{
		SelectedListView:   "connections",
		LoadedConnections:  connections,
		Connections:        connNames,
		SelectedConnection: "",

		DBs:        []string{"MongoDB", "FerretDB", "Postgres", "MySQL", "SQLite1"},
		SelectedDB: "",

		Collections:        []string{"users", "orders", "products"},
		SelectedCollection: "",

		Documents:        []string{"user_001", "user_002", "order_12345", "product_abc"},
		SelectedDocument: "",

		// Mockup MongoDB documents as JSON
		DocumentContent: map[string]string{
			"user_001": `{
  "_id": "507f1f77bcf86cd799439011",
  "username": "john_doe",
  "email": "john.doe@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "age": 28,
  "address": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "zipCode": "10001",
    "country": "USA"
  },
  "phoneNumbers": [
    {
      "type": "home",
      "number": "+1-555-123-4567"
    },
    {
      "type": "mobile",
      "number": "+1-555-987-6543"
    }
  ],
  "isActive": true,
  "roles": ["user", "customer"],
  "createdAt": "2024-01-15T10:30:00Z",
  "lastLogin": "2025-10-12T08:45:22Z"
}`,
			"user_002": `{
  "_id": "507f1f77bcf86cd799439012",
  "username": "jane_smith",
  "email": "jane.smith@example.com",
  "firstName": "Jane",
  "lastName": "Smith",
  "age": 34,
  "address": {
    "street": "456 Oak Ave",
    "city": "Los Angeles",
    "state": "CA",
    "zipCode": "90001",
    "country": "USA"
  },
  "phoneNumbers": [
    {
      "type": "mobile",
      "number": "+1-555-444-8899"
    }
  ],
  "isActive": true,
  "roles": ["user", "admin", "moderator"],
  "preferences": {
    "notifications": true,
    "theme": "dark",
    "language": "en"
  },
  "createdAt": "2023-06-22T14:20:00Z",
  "lastLogin": "2025-10-11T16:30:15Z"
}`,
			"order_12345": `{
  "_id": "65f1a2b3c4d5e6f7g8h9i0j1",
  "orderId": "ORD-2025-12345",
  "customerId": "507f1f77bcf86cd799439011",
  "orderDate": "2025-10-10T14:30:00Z",
  "status": "shipped",
  "items": [
    {
      "productId": "prod_abc123",
      "name": "Wireless Mouse",
      "quantity": 2,
      "price": 29.99,
      "total": 59.98
    },
    {
      "productId": "prod_def456",
      "name": "USB-C Cable",
      "quantity": 3,
      "price": 12.50,
      "total": 37.50
    }
  ],
  "subtotal": 97.48,
  "tax": 8.78,
  "shipping": 5.99,
  "total": 112.25,
  "shippingAddress": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "zipCode": "10001",
    "country": "USA"
  },
  "paymentMethod": "credit_card",
  "trackingNumber": "1Z999AA10123456784"
}`,
			"product_abc": `{
  "_id": "prod_abc123def456",
  "sku": "WMOUSE-BLK-001",
  "name": "Wireless Ergonomic Mouse",
  "category": "Electronics",
  "subcategory": "Computer Accessories",
  "description": "High-precision wireless mouse with ergonomic design",
  "price": 29.99,
  "currency": "USD",
  "inStock": true,
  "quantity": 150,
  "specifications": {
    "color": "Black",
    "connectivity": "Bluetooth 5.0",
    "batteryLife": "18 months",
    "dpi": [800, 1200, 1600, 2400],
    "weight": "95g",
    "dimensions": {
      "length": 120,
      "width": 65,
      "height": 40,
      "unit": "mm"
    }
  },
  "tags": ["wireless", "ergonomic", "bluetooth", "gaming"],
  "ratings": {
    "average": 4.5,
    "count": 328,
    "distribution": {
      "5": 210,
      "4": 85,
      "3": 20,
      "2": 8,
      "1": 5
    }
  },
  "vendor": {
    "id": "vendor_xyz789",
    "name": "TechGear Inc.",
    "country": "Taiwan"
  },
  "createdAt": "2024-03-10T09:00:00Z",
  "updatedAt": "2025-10-05T11:20:00Z"
}`,
		},
	}

	// Create notepad
	note := &notepad.Notepad{
		Name:     "editor",
		Title:    "Editor",
		Editable: false,
		Content:  "Pick something from the list...",
	}

	// Create popup for line editing
	var editPopup *popup.Popup
	var currentEditLine int

	// Set up notepad's edit line callback
	note.OnEditLine = func(lineNum int, oldLine string) {
		currentEditLine = lineNum

		editPopup = &popup.Popup{
			Name:    "editPopup",
			Title:   "Edit Line (Ctrl+S to save, ESC to cancel)",
			Content: oldLine,
			OnSave: func(newContent string) {
				// Update the line in notepad
				note.Lines[currentEditLine] = newContent

				// Rebuild the full content
				newFullContent := strings.Join(note.Lines, "\n")
				note.Update(g, newFullContent)

				// Update the document in model
				if m.SelectedDocument != "" {
					m.DocumentContent[m.SelectedDocument] = newFullContent
				}

				// Restore notepad border color
				note.SetActive(g, true)

				// Return focus to notepad
				if _, err := g.SetCurrentView(note.Name); err != nil {
					log.Panicln(err)
				}
			},
			OnCancel: func() {
				// Restore notepad border color
				note.SetActive(g, true)

				// Return focus to notepad
				if _, err := g.SetCurrentView(note.Name); err != nil {
					log.Panicln(err)
				}
			},
		}

		// Show the popup
		if err := editPopup.Show(g); err != nil {
			log.Panicln(err)
		}

		// Bind popup keys
		editPopup.BindKeys(g)
	}

	var listView *list.List

	// Set up notepad's back callback
	note.OnBack = func() {
		// Go back to document list
		listView.Title = "Documents"
		listView.Items = m.Documents
		listView.Update(g)

		// Update border colors
		note.SetActive(g, false)
		listView.SetActive(g, true)

		// Switch focus to list
		if _, err := g.SetCurrentView(listView.Name); err != nil {
			log.Panicln(err)
		}
	}

	// Create list with callback
	listView = &list.List{
		Name:     "listView",
		Title:    "Connections",
		Items:    m.Connections,
		Selected: 0,
		OnSelect: func(item string) {

			if m.SelectedListView == "connections" {
				m.SelectedConnection = item

				var selectedConn model.Connection
				for _, c := range m.LoadedConnections {
					if c.Name == item {
						selectedConn = c
						break
					}
				}

				if err := db.Connect(selectedConn); err != nil {
					log.Printf("Connection failed: %v", err)
					return
				}

				dbs, err := db.ListDatabases(db.Client)
				if err != nil {
					log.Printf("Failed to list databases: %v", err)
					return
				}
				m.DBs = dbs

				m.SelectedListView = "dbs"

				listView.Title = "DBs"
				listView.Items = m.DBs

				listView.Update(g)
			} else if m.SelectedListView == "dbs" {
				m.SelectedDB = item

				colls, err := db.ListCollections(db.Client, item)
				if err != nil {
					log.Printf("Failed to list collections: %v", err)
					return
				}
				m.Collections = colls

				m.SelectedListView = "collections"

				listView.Title = "Collections"
				listView.Items = m.Collections

				listView.Update(g)
			} else if m.SelectedListView == "collections" {
				m.SelectedCollection = item

				docs, err := db.ListDocuments(db.Client, m.SelectedDB, item)
				if err != nil {
					log.Printf("Failed to list documents: %v", err)
					return
				}

				m.DocumentContent = make(map[string]string)
				m.Documents = []string{}
				for i, doc := range docs {
					name := item + "_" + string(rune('0'+i))
					m.Documents = append(m.Documents, name)
					m.DocumentContent[name] = doc
				}

				m.SelectedListView = "documents"

				listView.Title = "Documents"
				listView.Items = m.Documents

				listView.Update(g)
			} else if m.SelectedListView == "documents" {
				// Display the selected document in the notepad
				m.SelectedDocument = item

				// Get the document content from mockup data
				var content string
				if docContent, exists := m.DocumentContent[item]; exists {
					content = docContent
				} else {
					content = "{\n  \"_id\": \"" + item + "\",\n  \"error\": \"Document not found in mockup data\"\n}"
				}

				// Update notepad title and content
				v, err := g.View(note.Name)
				if err == nil {
					v.Title = "Document: " + item
				}
				note.Update(g, content)

				// Update border colors
				listView.SetActive(g, false)
				note.SetActive(g, true)

				// Switch focus to editor to view the document
				if _, err := g.SetCurrentView(note.Name); err != nil {
					log.Panicln(err)
				}
			}

			// switch focus to editor
			// if _, err := g.SetCurrentView(note.Name); err != nil {
			// 	log.Panicln(err)
			// }
		},
		OnBack: func() {
			if m.SelectedListView == "documents" {
				// If we're viewing a document, first check if we need to go back to the list
				currentView := g.CurrentView()
				if currentView != nil && currentView.Name() == note.Name {
					// We're in the editor, go back to the document list
					note.SetActive(g, false)
					listView.SetActive(g, true)
					if _, err := g.SetCurrentView(listView.Name); err != nil {
						log.Panicln(err)
					}
					return
				}

				// Otherwise, go back to collections
				m.SelectedListView = "collections"
				m.SelectedDocument = ""

				listView.Title = "Collections"
				listView.Items = m.Collections
				listView.Update(g)
			} else if m.SelectedListView == "collections" {
				// Go back to DBs
				m.SelectedListView = "dbs"
				m.SelectedCollection = ""

				listView.Title = "DBs"
				listView.Items = m.DBs
				listView.Update(g)
			} else if m.SelectedListView == "dbs" {
				// Go back to connections
				m.SelectedListView = "connections"
				m.SelectedDB = ""

				listView.Title = "Connections"
				listView.Items = m.Connections
				listView.Update(g)
			}
			// If already at connections, do nothing (or could quit)
		},
	}

	// Layout manager
	g.SetManagerFunc(func(g *gocui.Gui) error {
		if err := listView.Layout(g); err != nil {
			return err
		}
		if err := note.Layout(g); err != nil {
			return err
		}
		return nil
	})

	// Bind keys
	listView.BindKeys(g)
	note.BindKeys(g)

	// Set initial border colors (list is active by default)
	listView.SetActive(g, true)
	note.SetActive(g, false)

	// global quit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
