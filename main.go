package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
	"log"
	_ "github.com/go-sql-driver/mysql"
)
// Categoria es una entrada de contenido.
type Categoria struct {
	Id     int
	CategoriaId int
	NombreCategoria  string
	TipoCategoria   string
}
/ Get busca un categoria por ID. El bool es falso si no lo encontramos.
func Get(id int) (Categoria, bool) {
	categorias := getCategorias(id)
	if len(categorias) == 0 {
		// Slice vacío; no se encontró el categoria.
		return Categoria{}, false
	}
	return categorias[0], true
}

// List devuelve un slice de todos los categorias.
func List() []Categoria {
	return getCategorias(-1)
}

// New guarda un categoria nuevo.
func New(p Categoria) []Categoria {
	return newCategoria(p)
}

// Put guarda un categoria existente.
func Put(p Categoria) {
	putCategoria(p)
}

// Del borra un categoria.
func Del(id int) {
	delCategoria(id)
}

// db es la base de datos global
var db *sql.DB

// Prepared statements
type stmtConfig struct {
	stmt *sql.Stmt
	q    string
}

var prepStmts = map[string]*stmtConfig{
	"get":    {q: "select * from categoria where id = ?;"},
	"list":   {q: "select * from categoria;"},
	"insert": {q: "insert into categoria (id, categoriaId, nombreCategoria, tipoCategoria) values (?, ?, ?, ?);"},
	"update": {q: "update categoria set categoriaId = ?, nombreCategoria = ?, tipoCategoria = ? where id = ?;"},
	"delete": {q: "delete from categoria where id = ?;"},
}




func main() {
	// Open database connection
	//db, err := sql.Open("mysql", "root:password@/busqueda-db")
	rand.Seed(time.Now().UnixNano())

	// Info para la DB.
	const (
		driver       = "mysql"
		dsn          = "busqueda-db"
		categoriaTableSQL  = `create table if not exists post(
			id int primary key not null,
			CategoriaId int not null,
			NombreCategoria text not null,
			TipoCategoria text not null
		);`
	)

	// Abrimos la base de datos
	var err error
	db, err = sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("categoria: error opening DB: %v", err)
	}

	// Creamos la tabla para los post, si no existe
	_, err = db.Exec(categoriaTableSQL )
	if err != nil {
		return fmt.Errorf("post: error creating post table: %v", err)
	}

	// Preparamos los "prepared statements" para get, list, new, put y del.
	for verb, sc := range prepStmts {
		sc.stmt, err = db.Prepare(sc.q)
		if err != nil {
			return fmt.Errorf("categoria: error preparing %s statement: %v", verb, err)
		}
	}

	return nil
}
