XC_ARCH=386 amd64
XC_OS=linux darwin windows
version=`./kii-cli -v | sed 's/kii-cli version //g'`

install:
	go install ./cmd/kii-cli

## GITHUB_TOKEN is needed
release: ./kii-cli
	rm -f pkg/*.exe pkg/*_amd64 pkg/*_386*
	ghr -u tmtk75 v$(version) pkg

./kii-cli:
	go build cmd/kii-cli/kii-cli.go

hash:
	shasum -a1 pkg/*_amd64.{gz,zip}

compress: pkg/kii-cli_win_amd64.zip pkg/kii-cli_darwin_amd64.gz pkg/kii-cli_linux_amd64.gz

pkg/kii-cli_win_amd64.zip pkg/kii-cli_darwin_amd64.gz pkg/kii-cli_linux_amd64.gz:
	gzip -k pkg/*darwin* pkg/*linux*
	for e in 386 amd64; do \
	  mv pkg/kii-cli_windows_$$e kii-cli_windows_$$e.exe; \
	  zip kii-cli_windows_$$e.zip kii-cli_windows_$$e.exe; \
	done
	mv kii-cli_windows_* pkg

build: clean
	for arch in $(XC_ARCH); do \
	  for os in $(XC_OS); do \
	    echo $$arch $$os; \
	    GOARCH=$$arch GOOS=$$os go build -o pkg/kii-cli_$${os}_$$arch ./cmd/kii-cli; \
	  done; \
	done

clean:
	rm -f pkg/*.gz pkg/*.zip

distclean:
	rm -rf kii-cli pkg
