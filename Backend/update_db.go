package main
import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "strings"
    _ "github.com/microsoft/go-mssqldb"
)
func main() {
    db, err := sql.Open("sqlserver", "server=localhost\\SQLEXPRESS;trusted_connection=yes;encrypt=disable;")
    if err != nil { panic(err) }
    defer db.Close()

    content, err := ioutil.ReadFile("database_fixed.sql")
    if err != nil { panic(err) }
    
    statements := strings.Split(string(content), "\nGO")
    for _, stmt := range statements {
        stmt = strings.TrimSpace(stmt)
        if stmt == "" { continue }
        _, err := db.Exec(stmt)
        if err != nil {
            fmt.Println("Error executing statement:", err)
        }
    }
    fmt.Println("Database updated successfully!")
}
