package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
	"time"
	"strconv"

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
	Delay    int `json:"delay,omitempty"`
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
var metricas = Metrica{}
var info = Info{Servico: "processo", Versao: "v1", Ambiente: "golang", Delay: 0}

func sleep(t int) {
	time.Sleep(time.Duration(t) * time.Second)
}

// GetInfo info about the service
func GetInfo(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "v%d", ver)
	json.NewEncoder(w).Encode(info)
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
	if info.Delay > 0 {
		metricas.Atrasadas++		
		sleep(info.Delay)
	}
	json.NewEncoder(w).Encode(processos)
}

// GetProcesso obtem um processo
func GetProcesso(w http.ResponseWriter, r *http.Request) {
	metricas.Consultas++
	params := mux.Vars(r)
	if info.Delay > 0 {
		metricas.Atrasadas++
		sleep(info.Delay)
	}
	for _, item := range processos {
		if item.ID == params["id"] {
			//fmt.Println(item)
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Processo{})
}

// SetDelay muda delay
func SetDelay(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	delay, _ := strconv.Atoi(params["segs"]) // TODO check conversion errors
	info.Delay = delay
}
// LoadDB load database
func LoadDB() {
	processos = append(processos, Processo{ID:"69372",Nome:"Alberto Azevedo",Cpfcnpj:"323.105.417-61",Descricao:"Gritou com a mae",Decisao:"culpado"})
	processos = append(processos, Processo{ID:"45532",Nome:"Thiago Santos",Cpfcnpj:"633.206.383-19",Descricao:"Nao lavou a louca",Decisao:"absorvido"})
	processos = append(processos, Processo{ID:"95902",Nome:"Ígor Moraes",Cpfcnpj:"065.226.778-57",Descricao:"Cupiditate sed repellendus numquam.",Decisao: "absorvido"})
	processos = append(processos, Processo{ID:"69532",Nome:"Breno Franco",Cpfcnpj:"554.748.026-68",Descricao:"Aliquam pariatur ea eum repellat non dicta.",Decisao:"culpado"})
	processos = append(processos, Processo{ID:"57724",Nome:"Suélen Barros",Cpfcnpj:"420.431.211-03",Descricao:"Maxime quo facere atque dolorum ea magni omnis voluptatem modi.",Decisao:"culpado"})
	processos = append(processos, Processo{ID:"42560",Nome:"Vicente Silva",Cpfcnpj:"137.021.861-31",Descricao:"Quae fugiat voluptatem illum vero quosunt sunt possimus voluptatum.",Decisao:"absorvido"})
	processos = append(processos, Processo{ID:"21825",Nome:"Rafael Nogueira",Cpfcnpj:"376.428.485-44",Descricao:"Tenetur veritatis doloribus autemamet.",Decisao:"absorvido"})
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
	router.HandleFunc("/setdelay/{segs}", SetDelay).Methods("GET")
	// TODO /seterror/:quando (sempre|par|impar)
	// TODO /metric
	router.HandleFunc("/metrics", GetMetric).Methods("GET")
	fmt.Println("Iniciando servidor...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
