package mime

type mime map[string][2]string

var mimes mime = map[string][2]string{

	"image/jpg":	                {"jpeg","image"},
	"image/png":	                {"png","image"},
	"image/gif":	                {"gif","image"},
	"image/webp":	                {"webp","image"},
	"image/x-canon-cr2":	        {"cr2","image"},
	"image/tiff":	                {"tif","image"},
	"image/bmp":	                {"bmp","image"},
	"image/heif":	                {"heif","image"},
	"image/vnd.ms-photo":	        {"jxr","image"},
	"image/vnd.adobe.photoshop":	{"psd","image"},
	"image/x-icon":	                {"ico","image"},
	"image/vnd.dwg":	            {"dwg","image"},


	"video/mp4":				{"mp4","video"},
	"video/x-m4v":				{"m4v","video"},
	"video/x-matroska":			{"mkv","video"},
	"video/webm":				{"webm","video"},
	"video/quicktime":			{"mov","video"},
	"video/x-msvideo":			{"avi","video"},
	"video/x-ms-wmv":			{"wmv","video"},
	"video/mpeg":				{"mpg","video"},
	"video/x-flv":				{"flv","video"},
	"video/3gpp":				{"3gp","video"},

	"audio/midi":			{"mid","audio"},
	"audio/mpeg":			{"mp3","audio"},
	"audio/m4a":			{"m4a","audio"},
	"audio/ogg":			{"ogg","audio"},
	"audio/x-flac":			{"flac","audio"},
	"audio/x-wav":			{"wav","audio"},
	"audio/amr":			{"amr","audio"},
	"audio/aac":			{"aac","audio"},

	"application/epub+zip":						{"epub","archive"},
	"application/zip":							{"zip","archive"},
	"application/x-tar":						{"tar","archive"},
	"application/x-rar-compressed":				{"rar","archive"},
	"application/gzip":							{"gz","archive"},
	"application/x-bzip2":						{"bz2","archive"},
	"application/x-7z-compressed":				{"7z","archive"},
	"application/x-xz":							{"xz","archive"},
	"application/pdf":							{"pdf","archive"},
	"application/x-msdownload":					{"exe","archive"},
	"application/x-shockwave-flash":			{"swf","archive"},
	"application/rtf":							{"rtf","archive"},
	"application/x-iso9660-image":				{"iso","archive"},
	"application/postscript":					{"ps","archive"},
	"application/x-sqlite3":					{"sqlite","archive"},
	"application/x-nintendo-nes-rom":			{"nes","archive"},
	"application/x-google-chrome-extension":	{"crx","archive"},
	"application/vnd.ms-cab-compressed":		{"cab","archive"},
	"application/x-deb":						{"deb","archive"},
	"application/x-unix-archive":				{"ar","archive"},
	"application/x-compress":					{"Z","archive"},
	"application/x-lzip":						{"lz","archive"},
	"application/x-rpm":						{"rpm","archive"},
	"application/x-executable":					{"elf","archive"},
	"application/dicom":						{"dcm","archive"},

	"application/application/msword":															{"doc","documents"},
	"application/application/vnd.openxmlformats-officedocument.wordprocessingml.document":		{"docx","documents"},
	"application/application/vnd.ms-excel":														{"xls","documents"},
	"application/application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":			{"xlsx","documents"},
	"application/application/vnd.ms-powerpoint":												{"ppt","documents"},
	"application/application/vnd.openxmlformats-officedocument.presentationml.presentation":	{"pptx","documents"},
	"application/vnd.apple.pages": 																{"pages","documents"},
	"application/vnd.apple.numbers":															{"numbers","documents"},
	"application/octet-stream":																	{"key","documents"},

}

//没找到 返回unknow
func GetFileTypeByMIME(mime string) [2]string  {
	theType := [2]string{"unknow","unknow"}
	if v,ok := mimes[mime]; ok {
		theType = v
	}

	return theType
}
