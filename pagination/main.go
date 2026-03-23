package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

// Sample data structure
type User struct {
	ID    int
	Name  string
	Email string
	Role  string
}

// Paginator struct
type Paginator struct {
	TotalRows   int
	PerPage     int
	CurrentPage int
	TotalPages  int
	HasNext     bool
	HasPrev     bool
	NextPage    int
	PrevPage    int
	Offset      int // Starting position (0-based)
	Limit       int // Number of items to fetch
	Links       []int
}

// NewPaginator creates a new Paginator using LIMIT and OFFSET concept
func NewPaginator(totalRows, perPage, currentPage int) *Paginator {
	totalPages := int(math.Ceil(float64(totalRows) / float64(perPage)))
	hasNext := currentPage < totalPages
	hasPrev := currentPage > 1
	nextPage := currentPage + 1
	prevPage := currentPage - 1

	// LIMIT and OFFSET calculation
	// OFFSET = (currentPage - 1) * perPage
	// LIMIT = perPage
	offset := (currentPage - 1) * perPage

	// Generate links (you can customize the number of links)
	numLinks := 5
	start := currentPage - numLinks/2
	if start < 1 {
		start = 1
	}
	end := start + numLinks - 1
	if end > totalPages {
		end = totalPages
	}
	links := []int{}
	for i := start; i <= end; i++ {
		links = append(links, i)
	}

	return &Paginator{
		TotalRows:   totalRows,
		PerPage:     perPage,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
		NextPage:    nextPage,
		PrevPage:    prevPage,
		Offset:      offset,  // This is the OFFSET for database queries
		Limit:       perPage, // This is the LIMIT for database queries
		Links:       links,
	}
}

func getIntFromQuery(r *http.Request, key string, defaultValue int) int {
	values := r.URL.Query()
	if value, ok := values[key]; ok {
		if intValue, err := strconv.Atoi(value[0]); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Generate sample data
func generateSampleData() []User {
	users := []User{}
	for i := 1; i <= 100; i++ {
		users = append(users, User{
			ID:    i,
			Name:  fmt.Sprintf("User %d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Role:  []string{"Admin", "User", "Moderator", "Guest"}[i%4],
		})
	}
	return users
}

// Get paginated data using LIMIT and OFFSET concept
// This simulates a database query: SELECT * FROM users LIMIT limit OFFSET offset
func getPaginatedData(users []User, offset, limit int) []User {
	// Validate offset
	if offset >= len(users) {
		return []User{}
	}

	// Calculate end position
	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	// Return the slice from offset to end (simulating LIMIT/OFFSET)
	return users[offset:end]
}

// Example of how this would work with a real database
func getDatabaseQueryExample(offset, limit int) string {
	return fmt.Sprintf(`
-- Example SQL query using LIMIT and OFFSET:
SELECT id, name, email, role 
FROM users 
ORDER BY id 
LIMIT %d OFFSET %d;

-- This query would return %d rows starting from position %d
-- Page 1: LIMIT 10 OFFSET 0   (rows 1-10)
-- Page 2: LIMIT 10 OFFSET 10  (rows 11-20)
-- Page 3: LIMIT 10 OFFSET 20  (rows 21-30)
-- etc.
`, limit, offset, limit, offset)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Generate sample data
	allUsers := generateSampleData()

	// Get pagination parameters
	perPage := getIntFromQuery(r, "per_page", 10)
	currentPage := getIntFromQuery(r, "page", 1)

	// Create paginator (calculates LIMIT and OFFSET)
	paginator := NewPaginator(len(allUsers), perPage, currentPage)

	// Get paginated data using LIMIT and OFFSET
	paginatedUsers := getPaginatedData(allUsers, paginator.Offset, paginator.Limit)

	// Parse template from file
	tmpl, err := template.ParseFiles("templates/pagination.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Users     []User
		Paginator *Paginator
	}{
		Users:     paginatedUsers,
		Paginator: paginator,
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Add a new endpoint to show the LIMIT/OFFSET concept
func debugHandler(w http.ResponseWriter, r *http.Request) {
	perPage := getIntFromQuery(r, "per_page", 10)
	currentPage := getIntFromQuery(r, "page", 1)

	paginator := NewPaginator(100, perPage, currentPage)

	queryExample := getDatabaseQueryExample(paginator.Offset, paginator.Limit)

	// Parse template from file
	tmpl, err := template.ParseFiles("templates/debug.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Paginator  *Paginator
		SQLExample string
	}{
		Paginator:  paginator,
		SQLExample: queryExample,
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/debug", debugHandler)
	fmt.Println("Server is running on port 8080...")
	fmt.Println("Visit http://localhost:8080 to see the pagination example")
	fmt.Println("Visit http://localhost:8080/debug to see LIMIT/OFFSET explanation")
	http.ListenAndServe(":8080", nil)
}
