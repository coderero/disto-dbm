package types

// This block of code defines a set of constants for different types of media formats. Each constant
// represents a specific media type and is assigned a string value that corresponds to the MIME type of
// that media format. These constants can be used in Go code to specify the media type of data being
// sent or received over the network, such as in HTTP headers.
const (
	Application_json       string = "application/json"
	Application_xml               = "application/xml"
	Application_xhtml             = "application/xhtml+xml"
	Application_atom              = "application/atom+xml"
	Application_pdf               = "application/pdf"
	Application_msword            = "application/msword"
	Application_octet             = "application/octet-stream"
	Application_x_www_form        = "application/x-www-form-urlencoded"
	Application_multipart         = "multipart/form-data"
	Application_css               = "text/css"
	Application_csv               = "text/csv"
	Application_html              = "text/html"
	Application_javascript        = "text/javascript"
	Application_plain             = "text/plain"
	Application_xml2              = "text/xml"
	Application_zip               = "application/zip"
	Application_gzip              = "application/gzip"
	Application_tar               = "application/tar"
	Application_rar               = "application/rar"
	Application_7z                = "application/x-7z-compressed"
	Application_rtf               = "application/rtf"
	Application_jar               = "application/java-archive"
	Application_swf               = "application/x-shockwave-flash"
	Application_mpeg              = "audio/mpeg"
	Application_webm              = "audio/webm"
	Application_ogg               = "audio/ogg"
	Application_wav               = "audio/wav"
	Application_mp4               = "video/mp4"
	Application_avi               = "video/x-msvideo"
	Application_mpeg2             = "video/mpeg"
	Application_ogg2              = "video/ogg"
	Application_quicktime         = "video/quicktime"
	Application_wmv               = "video/x-ms-wmv"
	Application_flv               = "video/x-flv"
	Application_3gp               = "video/3gpp"
)
