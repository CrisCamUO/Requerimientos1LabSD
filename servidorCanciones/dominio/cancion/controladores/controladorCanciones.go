package controladores

import (
	"context"
	"log"

	fachadacancionesservices "servidor.local/grpc-servidorCancion/dominio/cancion/fachadaCancionesServices"
	fachadaGeneros "servidor.local/grpc-servidorCancion/dominio/genero/fachadaGeneroServices"
	pb "servidor.local/grpc-servidorCancion/serviciosCancion"
)

type ControladorCanciones struct {
	pb.UnimplementedServiciosCancionesServer
}

func NewControladorCanciones() *ControladorCanciones {
	return &ControladorCanciones{}
}

func (c *ControladorCanciones) ListarGeneros(ctx context.Context, req *pb.Vacio) (*pb.ListaGeneros, error) {
	log.Printf("Llamando al metodo remoto: ListarGeneros")
	return fachadacancionesservices.ListarGeneros()
}

func (c *ControladorCanciones) ListarCancionesPorGenero(ctx context.Context, req *pb.IdGenero) (*pb.ListaCanciones, error) {
	log.Printf("Llamando al metodo remoto: ListarCancionesPorGenero con id=%d", req.Id)
	return fachadacancionesservices.ListarCancionesPorGenero(req.Id)
}

func (c *ControladorCanciones) BuscarGenero(ctx context.Context, req *pb.PeticionGeneroID) (*pb.RespuestaGeneroDTO, error) {
	log.Printf("Llamando al método remoto: BuscarGenero con ID=%d", req.Id)
	// Usar la fachada de géneros para esta operación específica
	return fachadaGeneros.BuscarGenero(req.Id), nil
}

func (c *ControladorCanciones) BuscarCancion(ctx context.Context, req *pb.PeticionCancionDTO) (*pb.RespuestaCancionDTO, error) {
	log.Printf("Llamando al metodo remoto: BuscarCancion con titulo='%s'", req.Titulo)
	return fachadacancionesservices.BuscarCancion(req.Titulo)
}

func (c *ControladorCanciones) ObtenerDetalleCancion(ctx context.Context, req *pb.IdCancion) (*pb.DetalleCancion, error) {
	log.Printf("Llamando al metodo remoto: ObtenerDetalleCancion con id=%d", req.Id)
	return fachadacancionesservices.ObtenerDetalleCancion(req.Id)
}
