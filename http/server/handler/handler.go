package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"net/http"
	"os"

	"github.com/Ivan-Bolotov/game-life/internal/service"
)

// создадим новый тип для добавления middleware к обработчикам
type Decorator func(http.Handler) http.Handler
// объект для хранения состояния игры
type LifeStates struct {
	service.LifeService
}

func New(ctx context.Context,
	lifeService service.LifeService,
) (http.Handler, error) {
	serveMux := http.NewServeMux()

	lifeState := LifeStates{
		LifeService: lifeService,
	}

	serveMux.HandleFunc("/nextstate", lifeState.nextState)
	serveMux.HandleFunc("/setstate", lifeState.setState)

	return serveMux, nil
}
// функция добавления middleware
func Decorate(next http.Handler, ds ...Decorator) http.Handler {
	decorated := next
	for d := len(ds) - 1; d >= 0; d-- {
		decorated = ds[d](decorated)
	}

	return decorated
}

// получение очередного состояния игры
func (ls *LifeStates) nextState(w http.ResponseWriter, r *http.Request) {
	worldState := ls.LifeService.NewState()

	err := json.NewEncoder(w).Encode(worldState.Cells)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ls *LifeStates) setState(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		data := map[string]int{}
		decoder := json.NewDecoder(r.Body)
    	err := decoder.Decode(&data)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			file, err := os.Create("state.cfg")
			if err != nil {
				fmt.Println(err.Error())
			} else {
				file.Write([]byte(strconv.Itoa(data["fill"]) + "%"))
			}
		}
	}
}
