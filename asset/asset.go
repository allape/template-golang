package asset

import _ "embed"

const MIME = "image/png"

//go:embed favicon.png
var Favicon []byte

//go:embed noimage.png
var NoImage []byte

//go:embed dameman.png
var DameMan []byte

//go:embed damewoman.png
var DameWoman []byte
