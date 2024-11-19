package helpers

func ConvertVelocityToVolume(velocity uint8) uint32 {
	return uint32(velocity) * 65535 / 127
}
