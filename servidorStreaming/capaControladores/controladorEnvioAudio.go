package capacontroladores

import (
	capafachadaservices "servidor.local/grpc-servidor/capaFachadaServices"
	pb "servidor.local/grpc-servidor/serviciosStreaming"
)

type ControladorServidor struct {
	pb.UnimplementedAudioServiceServer
}

// Implementación del procedimiento remoto que recibe el título de una canción y envia el archivo de audio en fragmentos mediante un stream.
func (s *ControladorServidor) EnviarCancionMedianteStream(req *pb.PeticionDTO, stream pb.AudioService_EnviarCancionMedianteStreamServer) error {
	return capafachadaservices.StreamAudioFile(
		req.Id,
		func(data []byte) error {
			return stream.Send(&pb.FragmentoCancion{Data: data})
		})
}
