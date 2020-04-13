package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//main
func main() {
	a := App{} 
    // You need to set your Username and Password here
    a.Initialize("root", "password", "busqueda-db")

    a.Run(":55000")

}

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/categorias", a.getCategorias).Methods("GET")
	a.Router.HandleFunc("/categoria", a.createCategoria).Methods("POST")
	a.Router.HandleFunc("/categoria/{id:[0-9]+}", a.getCategoria).Methods("GET")
	a.Router.HandleFunc("/categoria/{id:[0-9]+}", a.updateCategoria).Methods("PUT")
	a.Router.HandleFunc("/categoria/{id:[0-9]+}", a.deleteCategoria).Methods("DELETE")
}

func (a *App) getCategorias(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getCategorias(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createCategoria(w http.ResponseWriter, r *http.Request) {
	var u categoria
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := u.createCategoria(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) getCategoria(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid categoria ID")
		return
	}

	u := categoria{ID: id}
	if err := u.getCategoria(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Categoria not found press f")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) updateCategoria(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid categoria ID")
		return
	}

	var u categoria
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	u.ID = id

	if err := u.updateCategoria(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteCategoria(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Categoria ID")
		return
	}

	u := categoria{ID: id}
	if err := u.deleteCategoria(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


//model
type categoria struct {
	ID   int    `json:"id"`
	NombreCategoria  string `json:"nombrecategoria"`
	TipoCategoria   string `json:"tipocategoria"`
}

func (u *categoria) getCategoria(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT NombreCategoria, TipoCategoria FROM categorias WHERE id=%d", u.ID)
	return db.QueryRow(statement).Scan(&u.NombreCategoria, &u.TipoCategoria)
}

func (u *categoria) updateCategoria(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE categorias SET NombreCategoria='%s', TipoCategoria='%s' WHERE id=%d", u.NombreCategoria, u.TipoCategoria, u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *categoria) deleteCategoria(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM categorias WHERE id=%d", u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *categoria) createCategoria(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO categorias(NombreCategoria, TipoCategoria) VALUES('%s', '%s')", u.NombreCategoria, u.TipoCategoria)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func getCategorias(db *sql.DB, start, count int) ([]categoria, error) {
	statement := fmt.Sprintf("SELECT id, NombreCategoria, TipoCategoria FROM categorias LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categorias := []categoria{}

	for rows.Next() {
		var u categoria
		if err := rows.Scan(&u.ID, &u.NombreCategoria, &u.TipoCategoria); err != nil {
			return nil, err
		}
		categorias = append(categorias, u)
	}

	return categorias, nil
}
