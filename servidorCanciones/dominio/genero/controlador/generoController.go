package controlador

import (
	"context"
	"log"

	//Importamos las fachadas correspondientes
	fachadaCancionesServices "servidor.local/grpc-servidorCancion/dominio/cancion/fachadaCancionesServices"
	fachada "servidor.local/grpc-servidorCancion/dominio/genero/fachadaGeneroServices"
	pb "servidor.local/grpc-servidorCancion/serviciosCancion"
)

type GeneroController struct {
	pb.UnimplementedServiciosCancionesServer
}

func (c *GeneroController) ListarGeneros(ctx context.Context, req *pb.Vacio) (*pb.ListaGeneros, error) {
	log.Printf("Llamando al metodo remoto: ListarGeneros")
	// Usar la fachada de canciones que ya maneja la lógica
	return fachadaCancionesServices.ListarGeneros()
}

func (c *GeneroController) BuscarGenero(ctx context.Context, req *pb.PeticionGeneroID) (*pb.RespuestaGeneroDTO, error) {
	log.Printf("Llamando al método remoto: BuscarGenero con ID=%d", req.Id)
	// La fachada maneja todo el mapeo y la lógica de negocio
	return fachada.BuscarGenero(req.Id), nil
}

func (c *GeneroController) ListarCancionesPorGenero(ctx context.Context, req *pb.IdGenero) (*pb.ListaCanciones, error) {
	log.Printf("Llamando al método remoto: ListarCancionesPorGenero con ID=%d", req.Id)
	// Usar la fachada de canciones que maneja esta lógica
	return fachadaCancionesServices.ListarCancionesPorGenero(req.Id)
}
