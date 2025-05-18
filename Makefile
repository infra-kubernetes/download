SEALOS_VERSION=4.3.8
REPO_SEALOS=cuisongliu/sealos
download-amd64:
	echo "download sealos"
	wget -q https://github.com/${REPO_SEALOS}/releases/download/v${SEALOS_VERSION}/sealos_${SEALOS_VERSION}_linux_amd64.tar.gz
	tar -zxf sealos_${SEALOS_VERSION}_linux_amd64.tar.gz sealos
	rm -rf sealos_${SEALOS_VERSION}_linux_amd64.tar.gz
	echo "copy sealos to files"
	tar czf files/sealos_${SEALOS_VERSION}_linux_amd64.tar.gz sealos
	rm -rf  files/sealos
download-arm64:
	echo "download sealos"
	wget -q https://github.com/${REPO_SEALOS}/releases/download/v${SEALOS_VERSION}/sealos_${SEALOS_VERSION}_linux_arm64.tar.gz
	tar -zxf sealos_${SEALOS_VERSION}_linux_arm64.tar.gz sealos
	rm -rf sealos_${SEALOS_VERSION}_linux_arm64.tar.gz
	echo "copy sealos to files"
	tar czf files/sealos_${SEALOS_VERSION}_linux_arm64.tar.gz sealos
	rm -rf  files/sealos
init:
	rm -rf files
	mkdir -p files
	echo "clean old sealos"
	rm -rf other/sealos*
	make download-amd64
	make download-arm64