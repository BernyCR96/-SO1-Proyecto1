package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type str_mem struct {
	MemoriaT   int
	MemoriaL   int
	Porcentaje int
}

type Proc struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	State     string  `json:"state"`
	Porram    float64 `json:"ram"`
	CantHijos int     `json:"total"`
	Hijos     string  `json:"children"`
}

type ListaProc []Proc

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/cpu", PorcentajeCPU).Methods("GET")
	router.HandleFunc("/ram", getraminfo).Methods("GET")
	router.HandleFunc("/principal", getprocesosinfo).Methods("GET")
	router.HandleFunc("/kill/{id}", TerminarProceso).Methods("GET")
	http.ListenAndServe(":8081", router)
}

func getraminfo(w http.ResponseWriter, r *http.Request) {
	fileread, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}

	str := string(fileread)
	splitfile := strings.Split(string(str), "\n")

	ramTotal := strings.Replace((splitfile[0])[10:24], " ", "", -1)
	ramFree := strings.Replace((splitfile[1])[10:24], " ", "", -1)

	valTotal, err := strconv.Atoi(ramTotal)
	valFree, err1 := strconv.Atoi(ramFree)

	if err == nil && err1 == nil {
		ramTotalMB := valTotal / 1024
		ramFreeMB := valFree / 1024
		ramconsMB := (valTotal - valFree) / 1024
		rampor := ((ramTotalMB - ramFreeMB) * 100) / ramTotalMB

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(struct {
			Dato1 int `json:"total"`
			Dato2 int `json:"libre"`
			Dato3 int `json:"porcentaje"`
		}{Dato1: ramTotalMB, Dato2: ramconsMB, Dato3: rampor})

	} else {
		return
	}

}

func PorcentajeCPU(w http.ResponseWriter, r *http.Request) {

	idle0, total0 := getcpuinfo()
	time.Sleep(500 * time.Millisecond)
	idle1, total1 := getcpuinfo()
	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Dato1 float64 `json:"porcentaje"`
	}{Dato1: cpuUsage})
}

func getcpuinfo() (idle, total uint64) {

	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i],10,64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val
				if i == 4 {
					idle = val
				}
			}
			return
		}
	}
	return
}

func getprocesosinfo(w http.ResponseWriter, r *http.Request) {

	carpetas, err := ioutil.ReadDir("/proc")

	if err != nil {
		log.Fatal(err)
	}

	NoEjecutados := 0
	NoSuspendidos := 0
	NoDetenidos := 0
	NoZombis := 0

	var procesos = ListaProc{}

	for _, carpeta := range carpetas {
		if carpeta.IsDir() {
			r, _ := regexp.Compile("[0-9]+")
			if !r.MatchString(carpeta.Name()) {
				continue
			}

			var addProc Proc
			addProc.ID = carpeta.Name()

			stat, err := ioutil.ReadFile("/proc/" + addProc.ID + "/stat")
			if err != nil {
				log.Fatal(err)
			}
			statm, err := ioutil.ReadFile(("/proc/") + addProc.ID + "/statm")
			if err != nil {
				log.Fatal(err)
			}

			hijos, err := ioutil.ReadFile("/proc/" + addProc.ID + "/task/" + addProc.ID + "/children")
			if err != nil {
				log.Fatal(err)
			}

			contenidoStatm := strings.Split(string(statm), " ")
			contenido := strings.Split(string(stat), " ")

			ram, err := strconv.ParseFloat(contenidoStatm[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			ram = ram * 4
			ram = (ram * 100) / 8080072

			addProc.Name = strings.Replace(contenido[1], "(", "", -1)
			addProc.Name = strings.Replace(addProc.Name, ")", "", -1)

			addProc.State = contenido[2]
			addProc.Porram = ram
			addProc.Hijos = string(hijos)

			TotalHijos := strings.Split(string(hijos), " ")
			addProc.CantHijos = len(TotalHijos) - 1
			procesos = append(procesos, addProc)

			if contenido[2] == "R" {
				NoEjecutados++
			} else if contenido[2] == "T" {
				NoDetenidos++
			} else if contenido[2] == "S" {
				NoSuspendidos++
			} else if contenido[2] == "Z" {
				NoZombis++
			}

		}
	}
	cantproc := NoEjecutados + NoSuspendidos + NoDetenidos + NoZombis

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(struct {
		Dato1 int       `json:"ejecucion"`
		Dato2 int       `json:"suspendidos"`
		Dato3 int       `json:"detenidos"`
		Dato4 int       `json:"zombie"`
		Dato5 int       `json:"total"`
		Dato6 ListaProc `json:"procesos"`
	}{Dato1: NoEjecutados, Dato2: NoDetenidos, Dato3: NoSuspendidos, Dato4: NoZombis, Dato5: cantproc, Dato6: procesos})

}

func TerminarProceso(w http.ResponseWriter, r *http.Request) {
	procID := mux.Vars(r)["id"] //obtengo la variable
	out, err := exec.Command("kill", procID).Output()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Proceso Terminado") //imprimo mensaje en consola
	output := string(out[:])
	fmt.Println(w, output)
}
