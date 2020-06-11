// main.go 
package main 

import ( 
    "fmt" 
    "log" 
    "net/http"
    "encoding/json"
    "strconv"
    "flag"
)
// IndexHandler allows us to handle the request to the path '/' 
// and return "hello world" as a response to the client. 
IndexHandler func (w http.ResponseWriter, r * http.Request) { 
    fmt.Fprint (w, "hello world") 
}
func main () { 

    //  Instance of http.DefaultServerMux
     mux: = http.NewServeMux ()
    //  Path  to  handle
     mux.HandleFunc ("/", IndexHandler)
    //Esta conectada a la funcion NotesHandle 
     mux.HandleFunc(“/notes”, NotesHandler)
    //  server  listening  on  the  port  8080
     http.ListenAndServe ( "8080", mux) 


     // flag para realizar la creación de las tablas en la base
    // de datos.
    migrate := flag.Bool(
        "migrate", false, "Crea las tablas en la base de datos",
    )
    // Parseando todas las flags
    flag.Parse()
    if *migrate {
        if err := MakeMigrations(); err != nil {
            log.Fatal(err)
        }
    }
}



/ GetNotesHandler nos permite manejar las peticiones a la ruta
// ‘/notes’ con el método GET.
func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
    // Puntero a una estructura de tipo Note.
    n := new(Note)
    // Solicitando todas las notas en la base de datos.
    notes, err := n.GetAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    // Convirtiendo el slice de notas a formato JSON,
    // retorna un []byte y un error.
    j, err := json.Marshal(notes)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // Escribiendo el código de respuesta.
    w.WriteHeader(http.StatusOK)
    // Estableciendo el tipo de contenido del cuerpo de la
    // respuesta.
    w.Header().Set(“Content-Type”, “application/json”)
    // Escribiendo la respuesta, es decir nuestro slice de notas
    // en formato JSON.
    w.Write(j)
}


// UpdateNotesHandler nos permite manejar las peticiones a la ruta
// ‘/notes’ con el método UPDATE.
func UpdateNotesHandler(w http.ResponseWriter, r *http.Request) {
    var note Note
err := json.NewDecoder(r.Body).Decode(&note)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // Actualizamos la nota correspondiente.
    err = note.Update()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}


// DeleteNotesHandler nos permite manejar las peticiones a la ruta
// ‘/notes’ con el método DELETE.
func DeleteNotesHandler(w http.ResponseWriter, r *http.Request) {
    // obtenemos el valor pasado en la url como query
    // correspondiente a id, del tipo /notes?id=3.
    idStr := r.URL.Query().Get("id")
    // Verificamos que no esté vacío.
    if idStr == "" {
         http.Error(w, "Query id es requerido",
             http.StatusBadRequest)
         return
    }
    // Convertimos el valor obtenido del query a un int, de ser
    // posible.
    id, err := strconv.Atoi(idStr)
    if err != nil {
         http.Error(w, “Query id debe ser un número”,
             http.StatusBadRequest)
         return
    }
    var note Note
    // Borramos la nota con el id correspondiente.
    err = note.Delete(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

// NotesHandler nos permite manejar la petición a la ruta ‘/tasklist’ // y pasa el control a la función correspondiente según el método
// de la petición.
func NotesHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
        case http.MethodGet:
            GetNotesHandler(w, r)
        case http.MethodPost:
            CreateNotesHandler(w, r)
        case http.MethodPut:
            UpdateNotesHandler(w, r)
        case http.MethodDelete:
            DeleteNotesHandler(w, r)
        default:
            // Caso por defecto en caso de que se realice una
            // petición con un método diferente a los esperados.
            http.Error(w, "Metodo no permitido",
                http.StatusBadRequest)
            return
    }
}