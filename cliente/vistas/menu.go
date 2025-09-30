package vistas

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	pbStream "servidor.local/grpc-servidor/serviciosStreaming"
	pbSong "servidor.local/grpc-servidorCancion/serviciosCancion"

	util "cliente.local/grpc-cliente/utilidades"
)

var reader = bufio.NewReader(os.Stdin)

// MostrarMenuPrincipal - Punto de entrada principal del menu
func MostrarMenuPrincipal(clienteCanciones pbSong.ServiciosCancionesClient, clienteStreaming pbStream.AudioServiceClient, ctx context.Context) {
	for {
		opcion := mostrarMenuPrincipalYObtenerOpcion()

		switch opcion {
		case 1:
			explorarGeneros(clienteCanciones, clienteStreaming, ctx)
		case 2:
			fmt.Println("\n🎵 ¡Gracias por usar nuestro reproductor de musica! ¡Hasta luego! 🎵")
			return
		default:
			fmt.Println("\n Opcion no válida. Por favor, seleccione una opcion del menu.")
		}
	}
}

// mostrarMenuPrincipalYObtenerOpcion - Muestra el menú principal y obtiene la opción del usuario
func mostrarMenuPrincipalYObtenerOpcion() int {
	for {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("🎵 REPRODUCTOR DE MUSICA - MENU PRINCIPAL 🎵")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("1. 🎸 Explorar géneros musicales")
		fmt.Println("2. 🚪 Salir")
		fmt.Print("\n📝 Seleccione una opcion (1-2): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("❌ Error leyendo entrada. Intente nuevamente.")
			continue
		}

		input = strings.TrimSpace(input)
		opcion, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("❌ Por favor, ingrese un numero valido.")
			continue
		}

		if opcion >= 1 && opcion <= 2 {
			return opcion
		}

		fmt.Println("❌ Opcion fuera de rango. Seleccione 1 o 2.")
	}
}

// explorarGeneros - Maneja la exploración de géneros musicales
func explorarGeneros(clienteCanciones pbSong.ServiciosCancionesClient, clienteStreaming pbStream.AudioServiceClient, ctx context.Context) {
	fmt.Println("\n📡 Obteniendo lista de géneros disponibles...")

	respuestaGeneros, err := clienteCanciones.ListarGeneros(ctx, &pbSong.Vacio{})
	if err != nil {
		fmt.Printf("❌ Error obteniendo generos: %v\n", err)
		presionarEnterParaContinuar()
		return
	}

	if len(respuestaGeneros.Generos) == 0 {
		fmt.Println("😔 No hay generos disponibles en este momento.")
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
			continue
		}

		if explorarCancionesPorGenero(clienteCanciones, clienteStreaming, ctx, genero) {
			return 
		}
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

// buscarGeneroPorId - Busca un genero por su ID
func buscarGeneroPorId(generos []*pbSong.Genero, id int32) *pbSong.Genero {
	for _, g := range generos {
		if g.Id == id {
			return g
		}
	}
	fmt.Printf("❌ Genero con ID %d no encontrado. Intente nuevamente.\n", id)
	return nil
}

// explorarCancionesPorGenero - Explora las canciones de un género específico
func explorarCancionesPorGenero(clienteCanciones pbSong.ServiciosCancionesClient, clienteStreaming pbStream.AudioServiceClient, ctx context.Context, genero *pbSong.Genero) bool {
	fmt.Printf("\n📡 Buscando canciones del genero '%s'...\n", genero.Nombre)

	respuestaCanciones, err := clienteCanciones.ListarCancionesPorGenero(ctx, &pbSong.IdGenero{Id: genero.Id})
	if err != nil {
		fmt.Printf("❌ Error obteniendo canciones: %v\n", err)
		presionarEnterParaContinuar()
		return false
	}

	if len(respuestaCanciones.Canciones) == 0 {
		fmt.Printf("😔 No se encontraron canciones para el genero '%s'.\n", genero.Nombre)
		presionarEnterParaContinuar()
		return false
	}

	for {
		mostrarCancionesDelGenero(respuestaCanciones.Canciones, genero.Nombre)

		titulo := solicitarTituloCancion()
		if titulo == "" { 
			return false
		}

		if buscarYReproducirCancion(clienteCanciones, clienteStreaming, ctx, titulo) {
			return true 
		}
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
			fmt.Println("❌ El tItulo no puede estar vacío. Intente nuevamente.")
			continue
		}

		return input
	}
}

// buscarYReproducirCancion - Busca una canción y ofrece reproducirla
func buscarYReproducirCancion(clienteCanciones pbSong.ServiciosCancionesClient, clienteStreaming pbStream.AudioServiceClient, ctx context.Context, titulo string) bool {
	fmt.Printf("\n🔍 Buscando la canción '%s'...\n", titulo)

	respuestaCancion, err := clienteCanciones.BuscarCancion(ctx, &pbSong.PeticionCancionDTO{Titulo: titulo})
	if err != nil {
		fmt.Printf("❌ Error buscando la canción: %v\n", err)
		presionarEnterParaContinuar()
		return false
	}

	if respuestaCancion.Codigo != 200 {
		fmt.Printf("😔 La canción '%s' no fue encontrada.\n", titulo)
		fmt.Println("💡 Verifique que el título esté escrito exactamente como aparece en la lista.")
		presionarEnterParaContinuar()
		return false
	}

	mostrarDetallesCancion(respuestaCancion.ObjCancion)

	if confirmarReproduccion() {
		return reproducirCancion(clienteStreaming, ctx, respuestaCancion.ObjCancion)
	}
	return false
}

// mostrarDetallesCancion - Muestra los detalles de una canción
func mostrarDetallesCancion(cancion *pbSong.Cancion) {
	fmt.Println("\n" + strings.Repeat("=", 45))
	fmt.Println("🎵 DETALLES DE LA CANCION")
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

// reproducirCancion - Reproduce una canción usando streaming, con opción de detener con '1'
func reproducirCancion(clienteStreaming pbStream.AudioServiceClient, ctx context.Context, cancion *pbSong.Cancion) bool {

	logOriginal := log.Writer()

	fmt.Printf("\n🎵 Iniciando reproduccion de '%s'...\n", cancion.Titulo)

	ctxCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := clienteStreaming.EnviarCancionMedianteStream(ctxCancel, &pbStream.PeticionDTO{
		Id:      cancion.Id,
		Formato: "mp3",
	})

	if err != nil {
		fmt.Printf("❌ Error iniciando streaming: %v\n", err)
		presionarEnterParaContinuar()
		return false
	}

	fmt.Println("🔊 Reproduciendo cancion en vivo...")
	fmt.Println("⏸️  Escriba '1' y presione Enter en cualquier momento para detener y volver al menú.")

	audioReader, audioWriter := io.Pipe()
	canalSincronizacion := make(chan struct{})
	interrupcion := make(chan bool, 1)

	// Goroutine 1: ReproducciOn de audio
	go util.DecodificarReproducir(audioReader, canalSincronizacion)

	// Goroutine 2 : Escuchar teclado
	go func() {
		stdinReader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print(">>> ")
			input, err := stdinReader.ReadString('\n')
			if err != nil {
				return
			}
			if strings.TrimSpace(input) == "1" {
				interrupcion <- true
				return
			}
		}
	}()

	// Goroutine 3 : Recibir datos del stream
	go func() {
		util.RecibirCancion(stream, audioWriter, canalSincronizacion)
	}()

	select {
		case <-interrupcion:
			fmt.Println("\n⏹️  Reproduccion detenida por el usuario.")
			cancel()
			audioReader.Close()
			audioWriter.Close()

			return true

		case <-canalSincronizacion:
			log.SetOutput(logOriginal)
			fmt.Println("\n Reproduccion finalizada.")
			presionarEnterParaContinuar()

			return false
	}
}

// presionarEnterParaContinuar - Pausa la ejecucion hasta que el usuario presione Enter
func presionarEnterParaContinuar() {
	fmt.Print("\n📥 Presione Enter para continuar...")
	reader.ReadString('\n')
}
