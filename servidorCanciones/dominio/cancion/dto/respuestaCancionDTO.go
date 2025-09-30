package dto

import (
	cancionModelo "servidor.local/grpc-servidorCancion/dominio/cancion/modelo"
	pb "servidor.local/grpc-servidorCancion/serviciosCancion"
)

type RespuestaDTO struct {
	ObjCancion cancionModelo.Cancion
	Codigo     int32
	Mensaje    string
}

// ToPbCancion convierte un modelo interno de Cancion a un mensaje Protobuf
// para ser enviado via gRPC al cliente. Realiza el mapeo entre las estructuras
// del dominio y las estructuras de transporte.

func ToPbCancion(c cancionModelo.Cancion) *pb.Cancion {
	return &pb.Cancion{
		Id:              c.Id,
		Titulo:          c.Titulo,
		Artista:         c.Autor,
		AnioLanzamiento: c.AnioLanzamiento,
		Duracion:        c.Duracion,
		ObjGenero: &pb.Genero{
			Id:     c.Genero.Id,
			Nombre: c.Genero.Nombre,
		},
	}
}
