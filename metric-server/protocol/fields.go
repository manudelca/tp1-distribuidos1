package protocol

type FieldType uint8

const (
	METRICID FieldType = iota
	VALUE
	AGGREGATION
	AGGREGATIONWINDOWSSECS
	FROM
	TO
)

func parseMetricId(message []byte, i int) (string, int, error) {
	// Es leer el len. Parsearlo (error si no va)
	// Y tomar esos valores de los bytes -> Como manejo el error si me voy de mambo?
	// Tengo que chequear el len antes con un if jeje
	// Si me cuestiono estas cosas es porque oficialmente estoy cansado (? Mañana sigo
	// Me habre pasado con la implementacion de este protocolo maybe?
	return "", 0, nil
}
