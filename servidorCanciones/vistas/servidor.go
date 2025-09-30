package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	capaControladoresCancion "servidor.local/grpc-servidorCancion/dominio/cancion/controladores"
	pb "servidor.local/grpc-servidorCancion/serviciosCancion"
)

func main() {

	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterServiciosCancionesServer(grpcServer, &capaControladoresCancion.ControladorCanciones{})
	fmt.Printf("Servidor de Canciones escuchando en %s...\n", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
