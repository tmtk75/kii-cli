XC_ARCH=386 amd64
XC_OS=linux darwin windows

dl:
	curl -v -OL http://dl.bintray.com/tmtk75/generic/kii-cli_darwin_amd64.gz

release: pkg/kii-cli_linux_amd64.gz
	curl -T $< \
		-utmtk75:$(API_KEY) \
		https://api.bintray.com/content/tmtk75/generic/kii-cli/$(version)/kii-cli_darwin_amd64.gz


pkg/kii-cli_amd64.zip pkg/kii-cli_darwin_amd64.gz pkg/kii-cli_linux_amd64.gz:
	gzip pkg/*_386
	gzip pkg/*_amd64
	for e in 386 amd64; do \
		mv pkg/kii-cli_windows_$$e.exe pkg/kii-cli_$$e.exe; \
		zip pkg/kii-cli_$$e.zip pkg/kii-cli_$$e.exe; \
	done

build: clean
	gox \
	  -os="$(XC_OS)" \
	  -arch="$(XC_ARCH)" \
	  -output "pkg/{{.Dir}}_{{.OS}}_{{.Arch}}"

clean:
	rm -rf pkg

setup:
	go get -u github.com/mitchellh/gox
	gox -build-toolchain
