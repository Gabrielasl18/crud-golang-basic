package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // Driver Mysql para Go
)

type Names struct {
	Id    int
	Name  string
	Email string
}

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

// abre a conexão com o banco de dados
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "gabriela"
	dbPass := "gabi123"
	dbName := "meubanco"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Conexão com o banco de dados estabelecida com sucesso")
	return db
}

func Index(w http.ResponseWriter, r *http.Request) {
	// abre a conexão com o banco de dados utilizando a função dbConn()
	db := dbConn()
	log.Println("Index route accessed.")

	// realiza a consulta com banco de dados e trata erros
	selDB, err := db.Query("SELECT * FROM names ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	// monta a struct para ser utilizada no template
	n := Names{}

	// monta um array para guardar os valores da struct
	res := []Names{}

	// realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		var id int
		var name, email string

		// faz o scan do SELECT
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		// envia os resultados para a struct
		n.Id = id
		n.Name = name
		n.Email = email

		// junta a struct com Array
		res = append(res, n)
	}

	// Log para verificar os dados recuperados
	log.Println("Data retrieved from the database:", res)

	// logs adicionais para depuração
	log.Println("Executing template...")
	err = tmpl.ExecuteTemplate(w, "Index.html", res)
	if err != nil {
		log.Println("Error executing template:", err)
	}

	// fecha a conexão
	defer db.Close()
}

// Função Show exibe apenas um resultado
func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	// Pega o ID do parametro da URL
	nId := r.URL.Query().Get("id")

	// Usa o ID para fazer a consulta e tratar erros
	selDB, err := db.Query("SELECT * FROM names WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	// Monta a strcut para ser utilizada no template
	n := Names{}

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		// Armazena os valores em variaveis
		var id int
		var name, email string

		// Faz o Scan do SELECT
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		n.Id = id
		n.Name = name
		n.Email = email
	}

	// Mostra o template
	tmpl.ExecuteTemplate(w, "Show.html", n)

	// Fecha a conexão
	defer db.Close()
}

// Função New apenas exibe o formulário para inserir novos dados
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New.html", nil)
}

// Função Edit, edita os dados
func Edit(w http.ResponseWriter, r *http.Request) {
	// Abre a conexão com banco de dados
	db := dbConn()

	// Pega o ID do parametro da URL
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM names WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	// Monta a struct para ser utilizada no template
	n := Names{}

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		//Armazena os valores em variaveis
		var id int
		var name, email string

		// Faz o Scan do SELECT
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		n.Id = id
		n.Name = name
		n.Email = email
	}

	// Mostra o template com formulário preenchido para edição
	tmpl.ExecuteTemplate(w, "Edit.html", n)

	// Fecha a conexão com o banco de dados
	defer db.Close()
}

// insere valores no banco de dados
func Insert(w http.ResponseWriter, r *http.Request) {

	//abre a conexão com banco de dados usando a função: dbConn()
	db := dbConn()

	//verifica o METHOD do formulário passado
	if r.Method == "POST" {
		//pega os campos do formulário
		name := r.FormValue("name")
		email := r.FormValue("email")

		//prepara a SQL e verifica errors
		insForm, err := db.Prepare("INSERT INTO names(name, email) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}

		//insere valores do formulario com a SQL tratada e verifica errors
		insForm.Exec(name, email)

		//exibe um log com os valores digitados no formulario
		log.Println("INSERT: Name: " + name + "| E-mail: " + email)
	}

	//encerra a conexão do dbConn()
	defer db.Close()

	//retorna a HOME
	http.Redirect(w, r, "/", 301)
}

// atualiza valores no banco de dados
func Update(w http.ResponseWriter, r *http.Request) {
	//abre a conexão com o banco de dados usando a função: dbConn()
	db := dbConn()

	//verifica o METHOD do formulario passado
	if r.Method == "POST" {

		//pega os campos do formulário
		name := r.FormValue("name")
		email := r.FormValue("email")
		id := r.FormValue("uid")

		//prepara a SQL e verifica errors
		insForm, err := db.Prepare("UPDATE names SET name=?, email=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}

		//insere valores do formulario com a SQL tratada e verifica erros
		insForm.Exec(name, email, id)

		//exibe um log com os valores digitados no formulario
		log.Println("UPDATE: Name" + name + " |E-mail" + email)
	}

	//encerra a conexão do dbConn()
	defer db.Close()

	//retorna a HOME
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	//abre a conexão com banco de dados usando a função: dbConn()
	db := dbConn()

	nId := r.URL.Query().Get("id")

	//prepara a SQL e verifica errors
	delForm, err := db.Prepare("DELETE FROM names WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	//insere valores do from com a SQL tratada e verifica errors
	delForm.Exec(nId)

	//exibe um log com os valores digitados no form
	log.Println("DELETE")

	//encerra a conexão com o banco de dados
	defer db.Close()

	//retorna a HOME
	http.Redirect(w, r, "/", 301)
}

func main() {
	//exibe a mensagem que o servidor foi iniciado
	log.Println("Server started on: http://localhost:9000")

	//gerencia as URLs
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)

	//ações
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	//inicia o servidor na porta 9000
	http.ListenAndServe(":9000", nil)
}
