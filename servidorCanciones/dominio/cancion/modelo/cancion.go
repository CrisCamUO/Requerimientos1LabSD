package modelo

import (
	"servidor.local/grpc-servidorCancion/dominio/genero/modelo"
)

type Cancion struct {
	Id             int32
	Titulo         string
	Autor          string
	AnioLanzamiento int32
	Duracion       string
	Genero         modelo.Genero
}
