package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
	"log"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)
//mux
type Person struct {
	ID string `json:"id,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName string `json:"lastname,omitempty"`
	Address *Address `json:"address,omitempty"`
  }
  
  type Address struct {
	City string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
  }
  
  var people []Person
  
  // EndPoints
  func GetPersonEndpoint(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	for _, item := range people {
	  if item.ID == params["id"] {
		json.NewEncoder(w).Encode(item)
		return
	  }
	}
	json.NewEncoder(w).Encode(&Person{})
  }
  
  func GetPeopleEndpoint(w http.ResponseWriter, req *http.Request){
	json.NewEncoder(w).Encode(people)
  }
  
  func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	var person Person
	_ = json.NewDecoder(req.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
  
  }
  
  func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range people {
	  if item.ID == params["id"] {
		people = append(people[:index], people[index + 1:]...)
		break
	  }
	}
	json.NewEncoder(w).Encode(people)
  }
// Categoria es una entrada de contenido.
type Categoria struct {
	Id     int
	CategoriaId int
	NombreCategoria  string
	TipoCategoria   string
}
//publico
// Get busca un categoria por ID. El bool es falso si no lo encontramos.
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

//privado 

// getCategorias busca un categoria con id o listado de todos si id es -1.
func getCategorias(id int) []Categoria {
	res := []Categoria{}
	if id != -1 {
		var p Categoria
		// Obtenemos y ejecutamos el get prepared statement.
		get := prepStmts["get"].stmt
		err := get.QueryRow(id).Scan(&p.Id, &p.CategoriaId , &p.NombreCategoria, &p.TipoCategoria )
		if err != nil {
			if err != sql.ErrNoRows {
				log.Printf("categoria: error getting categoria. Id: %d, err: %v\n", id, err)
			}
		} else {
			res = append(res, p)
		}
		return res
	}

	// Obtenemos y ejecutamos el list prepared statement.
	list := prepStmts["list"].stmt
	rows, err := list.Query()
	if err != nil {
		log.Printf("categoria: error getting categorias. err: %v\n", err)
	}
	defer rows.Close()

	// Procesamos los rows.
	for rows.Next() {
		var p Categoria
		if err := rows.Scan(&p.Id, &p.CategoriaId , &p.NombreCategoria, &p.TipoCategoria ); err != nil {
			log.Printf("categoria: error scanning row: %v\n", err)
			continue
		}
		res = append(res, p)
	}
	// Verificamos si hubo error procesando los rows.
	if err := rows.Err(); err != nil {
		log.Printf("categoria: error reading rows: %v\n", err)
	}

	return res
}

// newCategoria inserta un categoria en la DB.
func newCategoria(p Categoria) []Categoria {
	// Generamos ID único para el nuevo categoria.
	p.Id = rand.Intn(1000)
	for {
		l := getCategorias(p.Id)
		if len(l) == 0 {
			break
		}
		p.Id = rand.Intn(1000)
	}

	// Obtenemos y ejecutamos insert prepared statement.
	insert := prepStmts["insert"].stmt
	_, err := insert.Exec(p.Id, p.CategoriaId , p.NombreCategoria, p.TipoCategoria )
	if err != nil {
		log.Printf("categoria: error inserting categoria %d into DB: %v\n", p.Id, err)
	}
	return []Categoria{p}
}

// putCategoria actualiza un categoria en la DB.
func putCategoria(p Categoria) {
	// Obtenemos y ejecutamos update prepared statement.
	update := prepStmts["update"].stmt
	_, err := update.Exec(p.CategoriaId , p.NombreCategoria, p.TipoCategoria , p.Id)
	if err != nil {
		log.Printf("categoria: error updating categoria %d into DB: %v\n", p.Id, err)
	}
}

// delCategoria borra un categoria de la DB.
func delCategoria(id int) {
	// Obtenemos y ejecutamos delete prepared statement.
	del := prepStmts["delete"].stmt
	_, err := del.Exec(id)
	if err != nil {
		log.Printf("categoria: error deleting categoria %d into DB: %v\n", id, err)
	}
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
	//db, err = sql.Open(driver, dsn)
	db, err := sql.Open("mysql", "root:password@/busqueda-db")
        fmt.Println("hello world",db)

	if err != nil {
		panic(err.Error()) 
	}

	// Creamos la tabla para loas categorias, si no existe
	_, err = db.Exec(categoriaTableSQL )
	if err != nil {
		panic(err.Error()) 
	}
	//definimos como escuchar
	 router := mux.NewRouter()
  
	// data ejemplo
	busqueda = append(busqueda, Busqueda{ID: "1", FirstName:"Ryan", LastName:"Ray", Address: &Address{City:"Dubling", State:"California"}})
	busqueda = append(busqueda, Busqueda{ID: "2", FirstName:"Maria", LastName:"Ray"})

	// endpoints
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":1859", router))

}
