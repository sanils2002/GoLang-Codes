# Go Pagination with LIMIT/OFFSET

A comprehensive pagination implementation in Go using the LIMIT/OFFSET concept with a beautiful HTML interface and dynamic controls.

## 🚀 Features

- **Dynamic Pagination**: Real-time page navigation with Previous/Next buttons
- **Items Per Page Dropdown**: Choose between 5, 10, 20, or 50 items per page
- **Dual Pagination Controls**: Navigation at both top and bottom of the table
- **LIMIT/OFFSET Concept**: Database-ready pagination using offset and limit
- **Responsive Design**: Modern, clean UI that works on all devices
- **Debug Interface**: Visual explanation of how LIMIT/OFFSET works
- **Template Separation**: Clean separation of HTML and Go code

## 📋 Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [LIMIT/OFFSET Concept](#limitoffset-concept)
- [Project Structure](#project-structure)
- [Features Explained](#features-explained)
- [Database Integration](#database-integration)
- [Contributing](#contributing)

## 🛠️ Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd pagination
   ```

2. **Run the application**
   ```bash
   go run main.go
   ```

3. **Access the application**
   - Main pagination: http://localhost:8080
   - Debug interface: http://localhost:8080/debug

## 🎯 Usage

### Main Pagination Interface

Visit `http://localhost:8080` to see the main pagination interface:

- **Table Display**: Shows paginated user data with ID, Name, Email, and Role
- **Dropdown Control**: Change items per page (5, 10, 20, 50)
- **Page Navigation**: Click page numbers or Previous/Next buttons
- **Page Information**: Shows current page and total pages at top and bottom

### Debug Interface

Visit `http://localhost:8080/debug` to understand the LIMIT/OFFSET concept:

- **Current Values**: See real-time LIMIT and OFFSET calculations
- **SQL Examples**: View actual SQL queries for database integration
- **Visual Explanation**: Understand how pagination works under the hood

## 🔗 API Endpoints

| Endpoint | Description |
|----------|-------------|
| `/` | Main pagination interface |
| `/debug` | LIMIT/OFFSET explanation and debug info |

### Query Parameters

- `page`: Current page number (default: 1)
- `per_page`: Items per page (default: 10, options: 5, 10, 20, 50)

### Example URLs

```
http://localhost:8080?page=2&per_page=20
http://localhost:8080/debug?page=3&per_page=10
```

## 🧮 LIMIT/OFFSET Concept

### How It Works

The pagination uses the standard LIMIT/OFFSET approach:

```go
// Formula
OFFSET = (Current Page - 1) × Items Per Page
LIMIT = Items Per Page
```

### Examples

| Page | OFFSET | LIMIT | Rows Returned |
|------|--------|-------|---------------|
| 1    | 0      | 10    | 1-10          |
| 2    | 10     | 10    | 11-20         |
| 3    | 20     | 10    | 21-30         |
| 4    | 30     | 10    | 31-40         |

### Database Query Example

```sql
SELECT id, name, email, role 
FROM users 
ORDER BY id 
LIMIT 10 OFFSET 20;
```

## 📁 Project Structure

```
pagination/
├── main.go                    # Main application logic
├── templates/
│   ├── pagination.html        # Main pagination interface
│   └── debug.html            # Debug and explanation page
├── go.mod                     # Go module file
└── README.md                 # This file
```

## 🎨 Features Explained

### 1. Dynamic Pagination Controls

- **Previous/Next Buttons**: Navigate between pages
- **Page Number Links**: Direct navigation to specific pages
- **Current Page Highlighting**: Visual indication of current page
- **Disabled States**: Buttons are disabled when not applicable

### 2. Items Per Page Dropdown

- **JavaScript Integration**: Smooth dropdown functionality
- **Auto-reset**: Returns to page 1 when changing items per page
- **URL Parameters**: Maintains state through URL parameters

### 3. Dual Pagination

- **Top Controls**: Navigation above the table
- **Bottom Controls**: Same controls below the table
- **Page Information**: Shows "Page X of Y" at both locations

### 4. Responsive Design

- **Modern CSS**: Clean, professional styling
- **Hover Effects**: Interactive elements with smooth transitions
- **Mobile Friendly**: Works well on all screen sizes

## 🗄️ Database Integration

### Go with Database Example

```go
func getUsersFromDB(offset, limit int) []User {
    query := "SELECT id, name, email, role FROM users ORDER BY id LIMIT ? OFFSET ?"
    rows, err := db.Query(query, limit, offset)
    if err != nil {
        return nil
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role)
        if err != nil {
            continue
        }
        users = append(users, user)
    }
    return users
}
```

### Supported Databases

- **PostgreSQL**: `LIMIT offset, limit`
- **MySQL**: `LIMIT offset, limit`
- **SQLite**: `LIMIT limit OFFSET offset`
- **SQL Server**: `OFFSET offset ROWS FETCH NEXT limit ROWS ONLY`

## 🔧 Customization

### Adding More Items Per Page Options

Edit `templates/pagination.html`:

```html
<select id="perPage" onchange="changePerPage(this.value)">
    <option value="5" {{if eq .Paginator.PerPage 5}}selected{{end}}>5</option>
    <option value="10" {{if eq .Paginator.PerPage 10}}selected{{end}}>10</option>
    <option value="20" {{if eq .Paginator.PerPage 20}}selected{{end}}>20</option>
    <option value="50" {{if eq .Paginator.PerPage 50}}selected{{end}}>50</option>
    <option value="100" {{if eq .Paginator.PerPage 100}}selected{{end}}>100</option>
</select>
```

### Changing Page Link Count

Edit `main.go` in the `NewPaginator` function:

```go
// Generate links (you can customize the number of links)
numLinks := 7  // Change from 5 to 7 for more page links
```

## 🚀 Performance Considerations

### For Large Datasets

1. **Database Indexing**: Ensure proper indexes on ORDER BY columns
2. **Count Optimization**: Use `COUNT(*)` efficiently for total rows
3. **Caching**: Cache pagination results for frequently accessed pages
4. **Connection Pooling**: Use database connection pools for better performance

### Memory Usage

- **Sample Data**: Currently generates 100 users in memory
- **Production**: Connect to database for real data
- **Large Datasets**: Consider cursor-based pagination for very large datasets

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Go template system for clean HTML separation
- Bootstrap-inspired CSS for modern styling
- SQL LIMIT/OFFSET standard for pagination

## 📞 Support

If you have any questions or need help with the implementation:

1. Check the debug interface at `/debug`
2. Review the code comments in `main.go`
3. Open an issue on GitHub

---

**Happy Paginating! 🎉** 