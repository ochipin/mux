package basemux

// ContentList : Content-Type 一覧情報を返却する
func ContentList() map[string]string {
	return map[string]string{
		".bz":    "application/x-bzip",
		".bz2":   "application/x-bzip2",
		".css":   "text/css",
		".csv":   "text/csv",
		".doc":   "application/msword",
		".docx":  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".gif":   "image/gif",
		".htm":   "text/html",
		".html":  "text/html",
		".ico":   "image/x-icon",
		".jar":   "application/java-archive",
		".jpg":   "image/jpeg",
		".jpeg":  "image/jpeg",
		".js":    "application/javascript",
		".json":  "application/json",
		".odp":   "application/vnd.oasis.opendocument.presentation",
		".ods":   "application/vnd.oasis.opendocument.spreadsheet",
		".odt":   "application/vnd.oasis.opendocument.text",
		".otf":   "font/otf",
		".png":   "image/png",
		".pdf":   "application/pdf",
		".ppt":   "application/vnd.ms-powerpoint",
		".pptx":  "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".rar":   "application/x-rar-compressed",
		".rtf":   "application/rtf",
		".sh":    "application/x-sh",
		".svg":   "image/svg+xml",
		".tar":   "application/x-tar",
		".tif":   "image/tiff",
		".tiff":  "image/tiff",
		".ts":    "application/typescript",
		".ttf":   "font/ttf",
		".txt":   "text/plain",
		".text":  "text/plain",
		".log":   "text/plain",
		".woff":  "font/woff",
		".woff2": "font/woff2",
		".xhtml": "application/xhtml+xml",
		".xls":   "application/vnd.ms-excel",
		".xlsx":  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".xml":   "application/xml",
		".zip":   "application/zip",
		".7z":    "application/x-7z-compressed",
	}
}