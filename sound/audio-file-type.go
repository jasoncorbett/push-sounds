package sound

const (
	OggVorbisFile AudioFileType = iota
	WavFile
	Mp3File
	FlacFile
	UnknownAudioFile
)

func (t AudioFileType) Name() string {
	switch t {
	case OggVorbisFile:
		return "ogg/vorbis"
	case WavFile:
		return "wav"
	case Mp3File:
		return "mp3"
	case FlacFile:
		return "flac"
	default:
		return "Unknown Audio File"
	}
}

func (t AudioFileType) Extension() string {
	switch t {
	case OggVorbisFile:
		return ".ogg"
	case WavFile:
		return ".wav"
	case Mp3File:
		return ".mp3"
	case FlacFile:
		return ".flac"
	default:
		return "UNKNOWN AUDIO FILE"
	}
}

func KnownAudioFileTypes() []AudioFileType {
	return []AudioFileType{
		OggVorbisFile,
		WavFile,
		Mp3File,
		FlacFile,
	}
}
