XC_ARCH=386 amd64
XC_OS=linux darwin windows
version=`./kii-cli -v | sed 's/kii-cli version //g'`

install:
	go install

## GITHUB_TOKEN is needed
release: ./kii-cli
	rm -f pkg/*.exe pkg/*_amd64 pkg/*_386*
	ghr -u tmtk75 v$(version) pkg

./kii-cli:
	go build

hash:
	shasum -a1 pkg/*_amd64.{gz,zip}

compress: pkg/kii-cli_win_amd64.zip pkg/kii-cli_darwin_amd64.gz pkg/kii-cli_linux_amd64.gz

pkg/kii-cli_win_amd64.zip pkg/kii-cli_darwin_amd64.gz pkg/kii-cli_linux_amd64.gz:
	gzip -k pkg/*_386
	gzip -k pkg/*_amd64
	for e in 386 amd64; do \
		mv pkg/kii-cli_windows_$$e.exe pkg/kii-cli_$$e.exe; \
		zip kii-cli_win_$$e.zip pkg/kii-cli_$$e.exe; \
		mv kii-cli_win_$$e.zip pkg; \
	done

build: clean
	gox \
	  -os="$(XC_OS)" \
	  -arch="$(XC_ARCH)" \
	  -output "pkg/{{.Dir}}_{{.OS}}_{{.Arch}}" \
	  ./cmd/kii-cli

clean:
	rm -f pkg/*.gz pkg/*.zip

distclean:
	rm -rf kii-cli pkg

setup:
	go get -u github.com/mitchellh/gox
	gox -build-toolchain
