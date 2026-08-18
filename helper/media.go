package helper

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type CommandOutput []byte

type CodecType string

const (
	Video CodecType = "video"
	//Audio CodecType = "audio"
)

type FFProbeStream struct {
	Index         int       `json:"index"`
	CodecName     string    `json:"codec_name"`
	CodecLongName string    `json:"codec_long_name"`
	Profile       string    `json:"profile"`
	CodecType     CodecType `json:"codec_type"`
	CodecTagStr   string    `json:"codec_tag_string"`
	CodecTag      string    `json:"codec_tag"`
	NbFrames      string    `json:"nb_frames"`
	Width         int       `json:"width"`
	Height        int       `json:"height"`
}

type FFProbeFormat struct {
	Filename       string `json:"filename"`
	NbStreams      int    `json:"nb_streams"`
	NbPrograms     int    `json:"nb_programs"`
	NbStreamGroups int    `json:"nb_stream_groups"`
	FormatName     string `json:"format_name"`
	FormatLongName string `json:"format_long_name"`
	StartTime      string `json:"start_time"`
	Duration       string `json:"duration"`
	Size           string `json:"size"`
	BitRate        string `json:"bit_rate"`
	ProbeScore     int    `json:"probe_score"`
	Tags           any    `json:"tags"`
}

type FFProbeJson struct {
	Streams []FFProbeStream `json:"streams"`
	Format  FFProbeFormat   `json:"format"`
}

func (f *FFProbeJson) NBFrames(ct CodecType) (uint64, error) {
	for _, stream := range f.Streams {
		if stream.CodecType == ct {
			nbFrames, err := strconv.ParseInt(stream.NbFrames, 10, 64)
			return uint64(nbFrames), err
		}
	}
	return 0, nil
}

func (f *FFProbeJson) Duration() (time.Duration, error) {
	duration, err := strconv.ParseFloat(f.Format.Duration, 64)
	if err != nil {
		return 0, err
	}
	return time.Millisecond * time.Duration(duration*1000), nil
}

func (f *FFProbeJson) Size() image.Point {
	for _, stream := range f.Streams {
		if stream.CodecType == Video {
			return image.Point{
				X: stream.Width,
				Y: stream.Height,
			}
		}
	}
	return image.Point{}
}

func FFProbe(file string) (*FFProbeJson, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return nil, err
	} else if stat.IsDir() {
		return nil, fmt.Errorf("ffprobe: %s is a directory", file)
	}

	cmd := exec.Command(
		"ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		file,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffprobe: %w: %s", err, output)
	}

	var ffprobe FFProbeJson
	err = json.Unmarshal(output, &ffprobe)
	if err != nil {
		return nil, err
	}

	return &ffprobe, nil
}

func FFMpegScale(dst, src string, scale float64) (CommandOutput, error) {
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-hide_banner",
		"-i",
		src,
		"-map_metadata", "-1",
		"-vf",
		fmt.Sprintf("scale=iw*%f:ih*%f", scale, scale),
		dst,
	)
	return cmd.CombinedOutput()
}

func FFMpegVideoSampleImage(dst, src string, scale float64, tile image.Point) (CommandOutput, error) {
	ffprobe, err := FFProbe(src)
	if err != nil {
		return nil, err
	}

	duration, err := ffprobe.Duration()
	if err != nil {
		return nil, err
	}

	size := ffprobe.Size()

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-hide_banner",
		"-i",
		src,
		"-map_metadata", "-1",
		"-fps_mode",
		"vfr",
		"-vf",
		fmt.Sprintf(
			"select='isnan(prev_selected_t)+gte(t-prev_selected_t\\,%.02f)',scale=%.02f:%.02f,tile=%dx%d",
			duration.Seconds()/float64(tile.X*tile.Y),
			float64(size.X)*scale,
			float64(size.Y)*scale,
			tile.X,
			tile.Y,
		),
		"-frames:v",
		"1",
		dst,
	)
	return cmd.CombinedOutput()
}

func ExifToolPreview(dst, src string) error {
	cmd := exec.Command(
		"exiftool",
		"-b",
		"-PreviewImage",
		src,
	)
	bs, err := cmd.Output()
	if err != nil {
		return err
	}

	file, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	_, err = file.Write(bs)
	if err != nil {
		return err
	}

	return nil
}
