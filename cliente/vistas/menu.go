package vistas

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	pbStream "servidor.local/grpc-servidor/serviciosStreaming"
	pbSong "servidor.local/grpc-servidorCancion/serviciosCancion"

	util "cliente.local/grpc-cliente/utilidades"
)

var reader = bufio.NewReader(os.Stdin)

// MostrarMenuPrincipal - Punto de entrada principal del menú
func MostrarMenuPrincipal(clienteCanciones pbSong.ServiciosCancionesClient, clienteStreaming pbStream.AudioServiceClient, ctx context.Context) {
	for {
		opcion := mostrarMenuPrincipalYObtenerOpcion()

		switch opcion {
		case 1:
			explorarGeneros(clienteCanciones, clienteStreaming, ctx)
		case 2:
			fmt.Println("\n🎵 ¡Gracias por usar nuestro reproductor de música! ¡Hasta luego! 🎵")
			return
		default:
			fmt.Println("\n❌ Opción no válida. Por favor, seleccione una opción del menú.")
		}
	}
}

// mostrarMenuPrincipalYObtenerOpcion - Muestra el menú principal y obtiene la opción del usuario
func mostrarMenuPrincipalYObtenerOpcion() int {
	for {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("🎵 REPRODUCTOR DE MÚSICA - MENÚ PRINCIPAL 🎵")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("1. 🎸 Explorar géneros musicales")
		fmt.Println("2. 🚪 Salir")
		fmt.Print("\n📝 Seleccione una opción (1-2): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("❌ Error leyendo entrada. Intente nuevamente.")
			continue
		}

		input = strings.TrimSpace(input)
		opcion, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("❌ Por favor, ingrese un número válido.")
			continue
		}

		if opcion >= 1 && opcion <= 2 {
			return opcion
		}

		fmt.Println("❌ Opción fuera de rango. Seleccione 1 o 2.")
	}
}

// explorarGeneros - Maneja la exploración de géneros musicales
func explorarGeneros(clienteCanciones pbSong.ServiciosCancionesClient, clienteStreaming pbStream.AudioServiceClient, ctx context.Context) {
	fmt.Println("\n📡 Obteniendo lista de géneros disponibles...")

	respuestaGeneros, err := clienteCanciones.ListarGeneros(ctx, &pbSong.Vacio{})
	if err != nil {
		fmt.Printf("❌ Error obteniendo géneros: %v\n", err)
		presionarEnterParaContinuar()
		return
	}

	if len(respuestaGeneros.Generos) == 0 {
		fmt.Println("😔 No hay géneros disponibles en este momento.")
		presionarEnterParaContinuar()
		return
	}

	for {
		idGenero := mostrarGenerosYSeleccionar(respuestaGeneros.Generos)
		if idGenero == -1 { // Usuario eligió volver
			return
		}

		genero := buscarGeneroPorId(respuestaGeneros.Generos, idGenero)
		if genero == nil {
			continue // El bucle pedirá otra opción
		}

		explorarCancionesPorGenero(clienteCanciones, clienteStreaming, ctx, genero)
	}
}

// mostrarGenerosYSeleccionar - Muestra la lista de géneros y permite seleccionar uno
func mostrarGenerosYSeleccionar(generos []*pbSong.Genero) int32 {
	for {
		fmt.Println("\n" + strings.Repeat("=", 40))
		fmt.Println("🎶 GÉNEROS MUSICALES DISPONIBLES")
		fmt.Println(strings.Repeat("=", 40))

		for _, g := range generos {
			fmt.Printf("🎵 %d. %s\n", g.Id, g.Nombre)
		}
		fmt.Printf("🔙 0. Volver al menú principal\n")
		fmt.Print("\n📝 Seleccione un género: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("❌ Error leyendo entrada. Intente nuevamente.")
			continue
		}

		input = strings.TrimSpace(input)
		if input == "0" {
			return -1 // Señal para volver
		}

		idGenero, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("❌ Por favor, ingrese un número válido.")
			continue
		}

		return int32(idGenero)
	}
}

// buscarGeneroPorId - Busca un género por su ID
func buscarGeneroPorId(generos []*pbSong.Genero, id int32) *pbSong.Genero {
	for _, g := range generos {
		if g.Id == id {
			return g
		}
	}
	fmt.Printf("❌ Género con ID %d no encontrado. Intente nuevamente.\n", id)
	return nil
}

// explorarCancionesPorGenero - Explora las canciones de un género específico
func explorarCancionesPorGenero(clienteCanciones pbSong.ServiciosCancionesClient, clienteStreaming pbStream.AudioServiceClient, ctx context.Context, genero *pbSong.Genero) {
	fmt.Printf("\n📡 Buscando canciones del género '%s'...\n", genero.Nombre)

	respuestaCanciones, err := clienteCanciones.ListarCancionesPorGenero(ctx, &pbSong.IdGenero{Id: genero.Id})
	if err != nil {
		fmt.Printf("❌ Error obteniendo canciones: %v\n", err)
		presionarEnterParaContinuar()
		return
	}

	if len(respuestaCanciones.Canciones) == 0 {
		fmt.Printf("😔 No se encontraron canciones para el género '%s'.\n", genero.Nombre)
		presionarEnterParaContinuar()
		return
	}

	for {
		mostrarCancionesDelGenero(respuestaCanciones.Canciones, genero.Nombre)

		titulo := solicitarTituloCancion()
		if titulo == "" { // Usuario eligió volver
			return
		}

		buscarYReproducirCancion(clienteCanciones, clienteStreaming, ctx, titulo)
	}
}

// mostrarCancionesDelGenero - Muestra las canciones disponibles de un género
func mostrarCancionesDelGenero(canciones []*pbSong.Cancion, nombreGenero string) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("🎵 CANCIONES DEL GÉNERO: %s\n", strings.ToUpper(nombreGenero))
	fmt.Println(strings.Repeat("=", 50))

	for i, c := range canciones {
		fmt.Printf("🎶 %d. %s - %s\n", i+1, c.Titulo, c.Artista)
	}
	fmt.Println("\n💡 Para reproducir una canción, escriba el título exacto.")
}

// solicitarTituloCancion - Solicita al usuario el título de la canción a reproducir
func solicitarTituloCancion() string {
	for {
		fmt.Print("\n📝 Ingrese el título de la canción (o 'volver' para regresar): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("❌ Error leyendo entrada. Intente nuevamente.")
			continue
		}

		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "volver" {
			return ""
		}

		if input == "" {
			fmt.Println("❌ El título no puede estar vacío. Intente nuevamente.")
			continue
		}

		return input
	}
}

// buscarYReproducirCancion - Busca una canción y ofrece reproducirla
func buscarYReproducirCancion(clienteCanciones pbSong.ServiciosCancionesClient, clienteStreaming pbStream.AudioServiceClient, ctx context.Context, titulo string) {
	fmt.Printf("\n🔍 Buscando la canción '%s'...\n", titulo)

	respuestaCancion, err := clienteCanciones.BuscarCancion(ctx, &pbSong.PeticionCancionDTO{Titulo: titulo})
	if err != nil {
		fmt.Printf("❌ Error buscando la canción: %v\n", err)
		presionarEnterParaContinuar()
		return
	}

	if respuestaCancion.Codigo != 200 {
		fmt.Printf("😔 La canción '%s' no fue encontrada.\n", titulo)
		fmt.Println("💡 Verifique que el título esté escrito exactamente como aparece en la lista.")
		presionarEnterParaContinuar()
		return
	}

	mostrarDetallesCancion(respuestaCancion.ObjCancion)

	if confirmarReproduccion() {
		reproducirCancion(clienteStreaming, ctx, respuestaCancion.ObjCancion)
	}
}

// mostrarDetallesCancion - Muestra los detalles de una canción
func mostrarDetallesCancion(cancion *pbSong.Cancion) {
	fmt.Println("\n" + strings.Repeat("=", 45))
	fmt.Println("🎵 DETALLES DE LA CANCIÓN")
	fmt.Println(strings.Repeat("=", 45))
	fmt.Printf("🎶 Título: %s\n", cancion.Titulo)
	fmt.Printf("🎤 Artista: %s\n", cancion.Artista)
	fmt.Printf("📅 Año: %d\n", cancion.AnioLanzamiento)
	fmt.Printf("⏱️  Duración: %s\n", cancion.Duracion)
	fmt.Printf("🎸 Género: %s\n", cancion.ObjGenero.Nombre)
	fmt.Println(strings.Repeat("=", 45))
}

// confirmarReproduccion - Pregunta al usuario si desea reproducir la canción
func confirmarReproduccion() bool {
	for {
		fmt.Print("\n🎵 ¿Desea reproducir esta canción? (s/n): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("❌ Error leyendo entrada. Intente nuevamente.")
			continue
		}

		input = strings.ToLower(strings.TrimSpace(input))

		switch input {
		case "s", "si", "sí", "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("❌ Por favor, responda 's' para sí o 'n' para no.")
		}
	}
}

// reproducirCancion - Reproduce una canción usando streaming
func reproducirCancion(clienteStreaming pbStream.AudioServiceClient, ctx context.Context, cancion *pbSong.Cancion) {
	fmt.Printf("\n🎵 Iniciando reproducción de '%s'...\n", cancion.Titulo)

	stream, err := clienteStreaming.EnviarCancionMedianteStream(ctx, &pbStream.PeticionDTO{
		Id:      cancion.Id,
		Formato: "mp3",
	})
	if err != nil {
		fmt.Printf("❌ Error iniciando streaming: %v\n", err)
		presionarEnterParaContinuar()
		return
	}

	fmt.Println("🔊 Reproduciendo canción en vivo...")
	fmt.Println("⏸️  Presione Ctrl+C para detener la reproducción")

	reader, writer := io.Pipe()
	canalSincronizacion := make(chan struct{})

	// Goroutine para recibir y escribir los fragmentos en el pipe
	go util.DecodificarReproducir(reader, canalSincronizacion)
	util.RecibirCancion(stream, writer, canalSincronizacion)

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("\n✅ Reproducción finalizada.")
			break
		}
		if err != nil {
			fmt.Printf("\n❌ Error durante la reproducción: %v\n", err)
			break
		}
	}

	presionarEnterParaContinuar()
}

// presionarEnterParaContinuar - Pausa la ejecución hasta que el usuario presione Enter
func presionarEnterParaContinuar() {
	fmt.Print("\n📥 Presione Enter para continuar...")
	reader.ReadString('\n')
}
