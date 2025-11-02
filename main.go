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
		DocumentObjects:  make(map[string]interface{}),
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

				// Update the document in model
				if m.SelectedDocument != "" {
					m.DocumentContent[m.SelectedDocument] = newFullContent

					// Save to database
					if err := db.UpdateDocument(db.Client, m.SelectedDB, m.SelectedCollection, newFullContent); err != nil {
						log.Printf("Failed to save document: %v", err)
						popup.ShowInfo(g, "Failed to save document")
					} else {
						// Re-fetch document from database
						if docID, ok := m.DocumentObjects[m.SelectedDocument]; ok {
							if freshDoc, err := db.GetDocument(db.Client, m.SelectedDB, m.SelectedCollection, docID); err == nil {
								m.DocumentContent[m.SelectedDocument] = freshDoc
								note.Update(g, freshDoc)
								popup.ShowInfo(g, "Document saved successfully")
							} else {
								note.Update(g, newFullContent)
								popup.ShowInfo(g, "Saved but failed to refresh")
							}
						} else {
							note.Update(g, newFullContent)
							popup.ShowInfo(g, "Saved successfully")
						}
					}
				} else {
					note.Update(g, newFullContent)
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
				m.SelectedConnectionIndex = listView.Selected

				var selectedConn model.Connection
				for _, c := range m.LoadedConnections {
					if c.Name == item {
						selectedConn = c
						break
					}
				}

				popup.ShowConnect(g, selectedConn, func() error {
					dbs, err := db.ListDatabases(db.Client)
					if err != nil {
						return err
					}
					m.DBs = dbs
					m.SelectedListView = "dbs"

					g.Update(func(g *gocui.Gui) error {
						listView.Title = "DBs"
						listView.Items = m.DBs
						listView.Selected = m.SelectedDBIndex
						return listView.Update(g)
					})
					return nil
				})
				return
			} else if m.SelectedListView == "dbs" {
				m.SelectedDB = item
				m.SelectedDBIndex = listView.Selected

				colls, err := db.ListCollections(db.Client, item)
				if err != nil {
					log.Printf("Failed to list collections: %v", err)
					return
				}
				m.Collections = colls

				m.SelectedListView = "collections"

				listView.Title = "Collections"
				listView.Items = m.Collections
				listView.Selected = m.SelectedCollectionIndex

				listView.Update(g)
			} else if m.SelectedListView == "collections" {
				m.SelectedCollection = item
				m.SelectedCollectionIndex = listView.Selected

				docs, err := db.ListDocuments(db.Client, m.SelectedDB, item)
				if err != nil {
					log.Printf("Failed to list documents: %v", err)
					return
				}

				m.DocumentContent = make(map[string]string)
				m.DocumentObjects = make(map[string]interface{})
				m.Documents = []string{}
				for i, doc := range docs {
					name := doc.Summary
					if name == "" {
						name = item + "_" + string(rune('0'+i))
					}
					m.Documents = append(m.Documents, name)
					m.DocumentContent[name] = doc.JSON
					m.DocumentObjects[name] = doc.ID
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
				listView.Selected = m.SelectedCollectionIndex
				listView.Update(g)

				note.Update(g, "Pick something from the list...")
				v, err := g.View(note.Name)
				if err == nil {
					v.Title = "Editor"
				}
			} else if m.SelectedListView == "collections" {
				// Go back to DBs
				m.SelectedListView = "dbs"
				m.SelectedCollection = ""

				listView.Title = "DBs"
				listView.Items = m.DBs
				listView.Selected = m.SelectedDBIndex
				listView.Update(g)
			} else if m.SelectedListView == "dbs" {
				// Go back to connections
				m.SelectedListView = "connections"
				m.SelectedDB = ""

				listView.Title = "Connections"
				listView.Items = m.Connections
				listView.Selected = m.SelectedConnectionIndex
				listView.Update(g)
			}
			// If already at connections, do nothing (or could quit)
		},
	}

	// Layout manager
	g.SetManagerFunc(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()

		if v, err := g.SetView("header", 0, 0, maxX-1, 2, 0); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Frame = true
			v.Title = ""
			v.Clear()
			title := "FerretMate - MongoDB/FerretDB TUI Client"
			padding := (maxX - len(title) - 2) / 2
			if padding < 0 {
				padding = 0
			}
			v.Write([]byte(strings.Repeat(" ", padding) + title))
		}

		if v, err := g.SetView("footer", 0, maxY-3, maxX-1, maxY-1, 0); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Frame = false
			v.Title = ""
			v.Clear()
			v.Write([]byte(" ↑↓: Navigate | Enter: Select | N: New | Del: Delete | ESC: Back | Ctrl+C: Quit"))
		}

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

	// Key binding for creating new items
	if err := g.SetKeybinding("", 'n', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		switch m.SelectedListView {
		case "dbs":
			// Show popup for new database name
			editPopup := &popup.Popup{
				Name:       "newDatabasePopup",
				Title:      "Create Database - Step 1/2 (Enter or Ctrl+S to continue, ESC to cancel)",
				Content:    "",
				SingleLine: true,
				OnSave: func(dbName string) {
					if dbName == "" {
						return
					}

					if db.Client == nil {
						popup.ShowInfo(g, "Not connected to any server")
						return
					}

					// Store the database name temporarily
					tempDBName := dbName

					// Now ask for the collection name
					collPopup := &popup.Popup{
						Name:       "newCollectionPopup",
						Title:      "Create Database - Step 2/2: Collection Name (Enter or Ctrl+S to create, ESC to cancel)",
						Content:    "",
						SingleLine: true,
						OnSave: func(collName string) {
							if collName == "" {
								return
							}

							// Create the database with the first collection
							if err := db.CreateDatabase(db.Client, tempDBName, collName); err != nil {
								popup.ShowInfo(g, "Failed to create database")
								log.Printf("Failed to create database: %v", err)
								return
							}

							popup.ShowInfo(g, "Database created successfully")

							// Refresh database list
							dbs, err := db.ListDatabases(db.Client)
							if err == nil {
								m.DBs = dbs
								m.SelectedDB = tempDBName
								// Find and select the newly created database
								for i, dbItem := range m.DBs {
									if dbItem == tempDBName {
										m.SelectedDBIndex = i
										listView.Selected = i
										break
									}
								}
								listView.Items = m.DBs
								listView.Update(g)

								// Set focus back to list view
								g.SetCurrentView(listView.Name)
								g.Cursor = false
							}
						},
						OnCancel: func() {
							// Set focus back to list view on cancel
							g.SetCurrentView(listView.Name)
							g.Cursor = false
						},
					}
					collPopup.Show(g)
					collPopup.BindKeys(g)
				},
				OnCancel: func() {
					// Set focus back to list view on cancel
					g.SetCurrentView(listView.Name)
					g.Cursor = false
				},
			}
			editPopup.Show(g)
			editPopup.BindKeys(g)

		case "collections":
			// Show popup for new collection name
			editPopup := &popup.Popup{
				Name:       "newCollectionPopup",
				Title:      "Create Collection (Enter or Ctrl+S to create, ESC to cancel)",
				Content:    "",
				SingleLine: true,
				OnSave: func(collName string) {
					if collName == "" {
						return
					}

					if db.Client == nil {
						popup.ShowInfo(g, "Not connected to any server")
						return
					}

					dbName := m.DBs[m.SelectedDBIndex]
					if err := db.CreateCollection(db.Client, dbName, collName); err != nil {
						popup.ShowInfo(g, "Failed to create collection")
						log.Printf("Failed to create collection: %v", err)
						return
					}

					popup.ShowInfo(g, "Collection created successfully")

					// Refresh collection list
					colls, err := db.ListCollections(db.Client, dbName)
					if err == nil {
						m.Collections = colls
						m.SelectedCollection = collName
						m.SelectedCollectionIndex = len(m.Collections) - 1
						listView.Items = m.Collections
						listView.Selected = len(m.Collections) - 1 // Select the newly created collection
						listView.Update(g)

						// Set focus back to list view
						g.SetCurrentView(listView.Name)
						g.Cursor = false
					}
				},
				OnCancel: func() {
					// Set focus back to list view on cancel
					g.SetCurrentView(listView.Name)
					g.Cursor = false
				},
			}
			editPopup.Show(g)
			editPopup.BindKeys(g)

		case "documents":
			// Show popup with template document JSON
			templateDoc := `{
  "_id": {
    "$oid": "000000000000000000000000"
  },
  "new": "document",
  "status": "pending"
}`
			editPopup := &popup.Popup{
				Name:    "newDocumentPopup",
				Title:   "Create Document (Ctrl+S to create, ESC to cancel)",
				Content: templateDoc,
				OnSave: func(docJSON string) {
					if docJSON == "" {
						return
					}

					if db.Client == nil {
						popup.ShowInfo(g, "Not connected to any server")
						return
					}

					dbName := m.DBs[m.SelectedDBIndex]
					collName := m.Collections[m.SelectedCollectionIndex]

					if err := db.CreateDocument(db.Client, dbName, collName, docJSON); err != nil {
						popup.ShowInfo(g, "Failed to create document")
						log.Printf("Failed to create document: %v", err)
						return
					}

					popup.ShowInfo(g, "Document created successfully")

					// Refresh document list
					docs, err := db.ListDocuments(db.Client, dbName, collName)
					if err == nil {
						m.Documents = []string{}
						for _, doc := range docs {
							name := doc.Summary
							m.DocumentObjects[name] = doc.ID
							m.DocumentContent[name] = doc.JSON
							m.Documents = append(m.Documents, name)
						}
						listView.Items = m.Documents
						listView.Selected = len(m.Documents) - 1 // Select the newly created document
						listView.Update(g)

						// Set focus back to list view
						g.SetCurrentView(listView.Name)
						g.Cursor = false
					}
				},
				OnCancel: func() {
					// Set focus back to list view on cancel
					g.SetCurrentView(listView.Name)
					g.Cursor = false
				},
			}
			editPopup.Show(g)
			editPopup.BindKeys(g)
		}

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	// Key binding for deleting items
	if err := g.SetKeybinding("", gocui.KeyDelete, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		switch m.SelectedListView {
		case "dbs":
			// Delete database - use current cursor position
			if len(m.DBs) == 0 || listView.Selected >= len(m.DBs) {
				return nil
			}
			dbName := m.DBs[listView.Selected]

			popup.ShowConfirmation(g, "Delete database '"+dbName+"'?", func() {
				if err := db.DeleteDatabase(db.Client, dbName); err != nil {
					popup.ShowInfo(g, "Failed to delete database")
					log.Printf("Failed to delete database: %v", err)
					return
				}

				popup.ShowInfo(g, "Database deleted successfully")

				// Refresh database list
				dbs, err := db.ListDatabases(db.Client)
				if err == nil {
					m.DBs = dbs
					// Adjust cursor position after deletion
					if listView.Selected >= len(m.DBs) {
						listView.Selected = len(m.DBs) - 1
					}
					if listView.Selected < 0 {
						listView.Selected = 0
					}
					m.SelectedDBIndex = listView.Selected
					listView.Items = m.DBs
					listView.Update(g)
				}
			}, func() {
				// Cancelled - do nothing
			})

		case "collections":
			// Delete collection - use current cursor position
			if len(m.Collections) == 0 || listView.Selected >= len(m.Collections) {
				return nil
			}
			collName := m.Collections[listView.Selected]
			dbName := m.DBs[m.SelectedDBIndex]

			popup.ShowConfirmation(g, "Delete collection '"+collName+"'?", func() {
				if err := db.DeleteCollection(db.Client, dbName, collName); err != nil {
					popup.ShowInfo(g, "Failed to delete collection")
					log.Printf("Failed to delete collection: %v", err)
					return
				}

				popup.ShowInfo(g, "Collection deleted successfully")

				// Refresh collection list
				colls, err := db.ListCollections(db.Client, dbName)
				if err == nil {
					m.Collections = colls
					// Adjust cursor position after deletion
					if listView.Selected >= len(m.Collections) {
						listView.Selected = len(m.Collections) - 1
					}
					if listView.Selected < 0 {
						listView.Selected = 0
					}
					m.SelectedCollectionIndex = listView.Selected
					listView.Items = m.Collections
					listView.Update(g)
				}
			}, func() {
				// Cancelled - do nothing
			})

		case "documents":
			// Delete document - use current cursor position
			if len(m.Documents) == 0 || listView.Selected >= len(m.Documents) {
				return nil
			}
			docName := m.Documents[listView.Selected]
			docID := m.DocumentObjects[docName]
			dbName := m.DBs[m.SelectedDBIndex]
			collName := m.Collections[m.SelectedCollectionIndex]

			popup.ShowConfirmation(g, "Delete document '"+docName+"'?", func() {
				if err := db.DeleteDocument(db.Client, dbName, collName, docID); err != nil {
					popup.ShowInfo(g, "Failed to delete document")
					log.Printf("Failed to delete document: %v", err)
					return
				}

				popup.ShowInfo(g, "Document deleted successfully")

				// Refresh document list
				docs, err := db.ListDocuments(db.Client, dbName, collName)
				if err == nil {
					m.Documents = []string{}
					m.DocumentObjects = make(map[string]interface{})
					m.DocumentContent = make(map[string]string)
					for _, doc := range docs {
						name := doc.Summary
						m.DocumentObjects[name] = doc.ID
						m.DocumentContent[name] = doc.JSON
						m.Documents = append(m.Documents, name)
					}
					// Adjust cursor position after deletion
					if listView.Selected >= len(m.Documents) {
						listView.Selected = len(m.Documents) - 1
					}
					if listView.Selected < 0 {
						listView.Selected = 0
					}
					m.SelectedDocumentIndex = listView.Selected
					listView.Items = m.Documents
					listView.Update(g)

					// Clear the notepad if the deleted document was being viewed
					note.Update(g, "Pick something from the list...")
				}
			}, func() {
				// Cancelled - do nothing
			})
		}

		return nil
	}); err != nil {
		log.Panicln(err)
	}

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
