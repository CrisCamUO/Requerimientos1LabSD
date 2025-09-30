package repositorio

import (
	dto "servidor.local/grpc-servidorCancion/dominio/cancion/dto"
	canciones "servidor.local/grpc-servidorCancion/dominio/cancion/modelo"
	repoGenero "servidor.local/grpc-servidorCancion/dominio/genero/repositorio"
)

var VectorCanciones = make([]canciones.Cancion, 4)

// CargarCanciones inicializa el vector con datos de ejemplo
func CargarCanciones() {
	var objCancion1, objCancion2, objCancion3 canciones.Cancion

	objCancion1.Id = 1
	objCancion1.Titulo = "La Vida"
	objCancion1.Autor = "Carlos Vives"
	objCancion1.Duracion = "3:45"
	objCancion1.AnioLanzamiento = 1998
	objCancion1.Genero = repoGenero.VectorGeneros[0] // Asignar el género Salsa

	objCancion2.Id = 2
	objCancion2.Titulo = "La Bicicleta"
	objCancion2.Autor = "Shakira"
	objCancion2.Duracion = "3:38"
	objCancion2.AnioLanzamiento = 2016
	objCancion2.Genero = repoGenero.VectorGeneros[1] // Asignar el género Cumbia

	objCancion3.Id = 3
	objCancion3.Titulo = "Ojos Así"
	objCancion3.Autor = "Shakira"
	objCancion3.Duracion = "4:12"
	objCancion3.AnioLanzamiento = 2000
	objCancion3.Genero = repoGenero.VectorGeneros[2] // Asignar el género Rock

	VectorCanciones[0] = objCancion1
	VectorCanciones[1] = objCancion2
	VectorCanciones[2] = objCancion3
}

// BuscarCancion busca una canción por título
func BuscarCancion(titulo string) dto.RespuestaDTO {
	var respuesta dto.RespuestaDTO
	for i := 0; i < len(VectorCanciones); i++ {
		if VectorCanciones[i].Titulo == titulo {
			respuesta.ObjCancion = VectorCanciones[i]
			respuesta.Codigo = 200
			respuesta.Mensaje = "Cancion encontrada correctamente"
			return respuesta
		}
	}
	respuesta.Codigo = 404
	respuesta.Mensaje = "La cancion no se encontro"
	return respuesta
}
