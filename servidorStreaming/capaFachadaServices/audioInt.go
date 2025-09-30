package capafachadaservices

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func StreamAudioFile(idCancion int32, funcionParaEnviarFragmento func([]byte) error) error {
	log.Printf("Canción solicitada: %d", idCancion)
	file, err := os.Open("canciones/" + strconv.FormatInt(int64(idCancion), 10) + ".mp3")
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 64*1024) // 64 KB se envian por fragmento
	fragmento := 0

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			log.Println("Cancion enviada completamente desde la fachada")
			break
		}

		if err != nil {
			return fmt.Errorf("error leyendo el archivo: %w", err)
		}
		fragmento++
		log.Printf("Fragmento #%d leido (%d bytes) y enviando", fragmento, n)
		// Ejecutamos la función para enviar el fragmento al cliente
		err = funcionParaEnviarFragmento(buffer[:n])
		if err != nil {
			return fmt.Errorf("error enviando fragmento #%d: %w", fragmento, err)
		}
	}
	return nil
}
