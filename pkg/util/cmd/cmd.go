package cmd

func NewWriteCmd(filename, filecontent string) string {
	cmd := "cat > " + filename + " <<EOF\n" + filecontent + "\nEOF"
	return cmd
}

func NewReadCmd(filename string) string {
	cmd := "cat " + filename
	return cmd
}

func NewSystemStartCmd(serviceName string) string {
	return "systemctl daemon-reload && systemctl start " + serviceName
}

// when wrong, return err "Process exited with status 1"
func NewMkdirAllCmd(dir string) string {
	return "mkdir -p " + dir
}

func NewTarXCmd(tgzName, dir string) string {
	return "tar -zxvf " + tgzName + " -C " + dir
}

func NewCurlCmd(targetSite, filename string) string {
	return "curl -L " + targetSite+filename + " -o /tmp/" + filename
}