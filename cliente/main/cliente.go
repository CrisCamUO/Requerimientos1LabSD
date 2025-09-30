package main

import (
	"context"
	"log"
	"time"

	pbStream "servidor.local/grpc-servidor/serviciosStreaming"    // CORRECTO
	pbSong "servidor.local/grpc-servidorCancion/serviciosCancion" // CORRECTO

	menu "cliente.local/grpc-cliente/vistas"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Conexión al Servidor de Canciones (puerto 50051)
	connSong, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor de canciones: %v", err)
	}
	defer connSong.Close()

	// Conexión al Servidor de Streaming (puerto 50052)
	connStream, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor de streaming: %v", err)
	}
	defer connStream.Close()

	// Crear clientes gRPC
	clienteCanciones := pbSong.NewServiciosCancionesClient(connSong)
	clienteStreaming := pbStream.NewAudioServiceClient(connStream)

	// Contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Mostrar menu principal
	menu.MostrarMenuPrincipal(clienteCanciones, clienteStreaming, ctx)
}
