package fachadaGeneroServices

import (
	"log"

	generoDTO "servidor.local/grpc-servidorCancion/dominio/genero/dto"
	"servidor.local/grpc-servidorCancion/dominio/genero/repositorio"
	pb "servidor.local/grpc-servidorCancion/serviciosCancion"
)

type GeneroServices struct{}

// ListarGeneros obtiene todos los géneros y devuelve la respuesta lista para gRPC
func ListarGeneros() *pb.RespuestaGenerosDTO {
	log.Printf("Fachada: Listando todos los géneros")
	generos := repositorio.BuscarTodosLosGeneros()

	// Convertir géneros del modelo a Protocol Buffer
	var pbGeneros []*pb.Genero
	for _, g := range generos {
		pbGeneros = append(pbGeneros, generoDTO.ToPbGenero(g))
	}

	return &pb.RespuestaGenerosDTO{
		Mensaje:    "Generos listados exitosamente",
		Codigo:     200,
		ObjGeneros: pbGeneros,
	}
}

// BuscarGenero busca un genero por ID y devuelve la respuesta lista para gRPC
func BuscarGenero(id int32) *pb.RespuestaGeneroDTO {
	log.Printf("Fachada: Buscando genero con ID=%d", id)
	respuesta := repositorio.BuscarGenero(id)

	if respuesta.Codigo == 200 {
		return &pb.RespuestaGeneroDTO{
			Mensaje: "Genero encontrado",
			Codigo:  200,
			Genero:  generoDTO.ToPbGenero(respuesta.ObjGenero),
		}
	}

	return &pb.RespuestaGeneroDTO{
		Mensaje: "Genero no encontrado",
		Codigo:  404,
	}
}
