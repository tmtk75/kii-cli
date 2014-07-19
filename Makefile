XC_ARCH=386 amd64
XC_OS=linux darwin windows
version=0.0.5

bitray_dl:
	curl -v -OL http://dl.bintray.com/tmtk75/generic/$(version)_kii-cli_darwin_amd64.gz

bitray_release: pkg/kii-cli_linux_amd64.gz
	curl -T $< \
		-utmtk75:$(API_KEY) \
		https://api.bintray.com/content/tmtk75/generic/kii-cli/v1/$(version)_kii-cli_darwin_amd64.gz

compress: pkg/kii-cli_amd64.zip pkg/kii-cli_darwin_amd64.gz pkg/kii-cli_linux_amd64.gz

pkg/kii-cli_amd64.zip pkg/kii-cli_darwin_amd64.gz pkg/kii-cli_linux_amd64.gz:
	gzip -k pkg/*_386
	gzip -k pkg/*_amd64
	for e in 386 amd64; do \
		mv pkg/kii-cli_windows_$$e.exe pkg/kii-cli_$$e.exe; \
		zip kii-cli_$$e.zip pkg/kii-cli_$$e.exe; \
		mv kii-cli_$$e.zip pkg; \
	done

build: clean
	gox \
	  -os="$(XC_OS)" \
	  -arch="$(XC_ARCH)" \
	  -output "pkg/{{.Dir}}_{{.OS}}_{{.Arch}}"

clean:
	rm -f pkg/*.gz pkg/*.zip

distclean:
	rm -rf pkg

setup:
	go get -u github.com/mitchellh/gox
	gox -build-toolchain
