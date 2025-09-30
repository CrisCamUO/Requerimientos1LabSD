package dto

import (
	generoModelo "servidor.local/grpc-servidorCancion/dominio/genero/modelo"
	pb "servidor.local/grpc-servidorCancion/serviciosCancion"
)

type RespuestaDTO struct {
	ObjGenero generoModelo.Genero
	Codigo    int32
	Mensaje   string
}

func ToPbGenero(g generoModelo.Genero) *pb.Genero {
	return &pb.Genero{
		Id:     g.Id,
		Nombre: g.Nombre,
	}
}
