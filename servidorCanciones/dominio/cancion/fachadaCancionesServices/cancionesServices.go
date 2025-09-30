package fachadacancionesservices

import (
	"log"

	dto "servidor.local/grpc-servidorCancion/dominio/cancion/dto"
	repositorio "servidor.local/grpc-servidorCancion/dominio/cancion/repositorio"
	fachadagen "servidor.local/grpc-servidorCancion/dominio/genero/fachadaGeneroServices"
	pb "servidor.local/grpc-servidorCancion/serviciosCancion"
)

// Inicializa los datos al iniciar el servicio
func init() {
	repositorio.CargarCanciones()
}

func ListarGeneros() (*pb.ListaGeneros, error) {
	// Usa la función de la fachada de géneros
	respuestaGeneros := fachadagen.ListarGeneros()

	lista := &pb.ListaGeneros{}
	lista.Generos = respuestaGeneros.ObjGeneros

	return lista, nil
}

// ListarCancionesPorGenero devuelve las canciones de un género específico
func ListarCancionesPorGenero(idGenero int32) (*pb.ListaCanciones, error) {
	log.Printf("Ejecutando fachada: ListarCancionesPorGenero con id=%d", idGenero)

	lista := &pb.ListaCanciones{}
	for _, c := range repositorio.VectorCanciones {
		if int32(c.Genero.Id) == idGenero {
			lista.Canciones = append(lista.Canciones, dto.ToPbCancion(c))
		}
	}

	return lista, nil
}

// BuscarCancion busca una canción por título
func BuscarCancion(titulo string) (*pb.RespuestaCancionDTO, error) {
	log.Printf("Ejecutando fachada: BuscarCancion con titulo='%s'", titulo)

	res := repositorio.BuscarCancion(titulo)

	return &pb.RespuestaCancionDTO{
		Codigo:  res.Codigo,
		Mensaje: res.Mensaje,
		ObjCancion: func() *pb.Cancion {
			if res.Codigo == 200 {
				return dto.ToPbCancion(res.ObjCancion)
			}
			return nil
		}(),
	}, nil
}

// ObtenerDetalleCancion devuelve los detalles de una canción por ID
func ObtenerDetalleCancion(id int32) (*pb.DetalleCancion, error) {
	for _, c := range repositorio.VectorCanciones {
		if int32(c.Id) == id {
			return &pb.DetalleCancion{
				Cancion: dto.ToPbCancion(c), // mapea todos los campos básicos, incluyendo ObjGenero
			}, nil
		}
	}
	return &pb.DetalleCancion{}, nil
}
