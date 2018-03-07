package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// StatusType status
type StatusType struct {
	Status string `json:"status,omitempty"`
}

// Info sys info
type Info struct {
	Servico  string `json:"servico,omitempty"`
	Versao   string `json:"versao,omitempty"`
	Ambiente string `json:"ambiente,omitempty"`
	Delay    string `json:"delay,omitempty"`
}

// Processo info
type Processo struct {
	ID        string `json:"id,omitempty"`
	Nome      string `json:"nome,omitempty"`
	Cpfcnpj   string `json:"cpfcnpj,omitempty"`
	Descricao string `json:"descricao,omitempty"`
	Decisao   string `json:"decisao,omitempty"`
}

// Metrica metricas do servico
type Metrica struct {
	Consultas int `json:"consultas"`
	Erro      int `json:"erro"`
	Sucesso   int `json:"sucesso"`
	Atrasadas int `json:"atrasadas"`
}

var processos []Processo // processo db
var ver = 1              // service version
var delay = 0            // secs to delay service
var metricas = Metrica{}

// GetInfo info about the service
func GetInfo(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "v%d", ver)
	json.NewEncoder(w).Encode(Info{Servico: "processo", Versao: "v1", Ambiente: "golang", Delay: "0"})
}

// GetHealth service health
func GetHealth(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, StatusType{status: "UP"})
	json.NewEncoder(w).Encode(StatusType{Status: "UP"})
}

// GetMetric service metrics
func GetMetric(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, StatusType{status: "UP"})
	json.NewEncoder(w).Encode(metricas)
}

// GetProcessos get all process
func GetProcessos(w http.ResponseWriter, r *http.Request) {
	metricas.Consultas++
	json.NewEncoder(w).Encode(processos)
}

// GetProcesso obtem um processo
func GetProcesso(w http.ResponseWriter, r *http.Request) {
	metricas.Consultas++
	params := mux.Vars(r)
	for _, item := range processos {
		if item.ID == params["id"] {
			//fmt.Println(item)
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Processo{})
}

// LoadDB load database
func LoadDB() {
	processos = append(processos, Processo{ID: "1", Nome: "ACME", Cpfcnpj: "xyz", Descricao: "Gritou com a mãe", Decisao: "culpado"})
	processos = append(processos, Processo{ID: "4", Nome: "Gale Corp", Cpfcnpj: "xyz", Descricao: "Nao lavou a louça", Decisao: "absorvido"})
}

// main function TODO add more endpoints to mimic a component
func main() {
	router := mux.NewRouter()

	LoadDB() // init db
	// TODO add more process
	router.HandleFunc("/processo/{id}", GetProcesso).Methods("GET")
	router.HandleFunc("/processos", GetProcessos).Methods("GET")

	router.HandleFunc("/info", GetInfo).Methods("GET")
	router.HandleFunc("/health", GetHealth).Methods("GET")
	// TODO /setdelay/:segs
	// TODO /seterror/:quando (sempre|par|impar)
	// TODO /metric
	router.HandleFunc("/metrics", GetMetric).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
