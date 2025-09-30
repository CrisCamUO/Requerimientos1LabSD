package repositorio

import (
	dto "servidor.local/grpc-servidorCancion/dominio/genero/dto"
	genero "servidor.local/grpc-servidorCancion/dominio/genero/modelo"
)

var VectorGeneros = []genero.Genero{
	{Id: 1, Nombre: "Salsa"},
	{Id: 2, Nombre: "Cumbia"},
	{Id: 3, Nombre: "Rock"},
}

// BuscarGenero busca un género por ID
func BuscarGenero(id int32) dto.RespuestaDTO {
	var respuesta dto.RespuestaDTO
	for i := 0; i < len(VectorGeneros); i++ {
		if VectorGeneros[i].Id == id {
			respuesta.ObjGenero = VectorGeneros[i]
			respuesta.Codigo = 200
			respuesta.Mensaje = "Genero encontrado correctamente"
			return respuesta
		}
	}
	respuesta.Codigo = 404
	respuesta.Mensaje = "El genero no fue encontrado"
	return respuesta
}

// BuscarTodosLosGeneros devuelve todos los géneros
func BuscarTodosLosGeneros() []genero.Genero {
	return VectorGeneros
}
