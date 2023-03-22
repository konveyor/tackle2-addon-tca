module github.com/divsan93/tta_tca

go 1.16

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20181213151034-8d9ed539ba31

replace k8s.io/api => k8s.io/api v0.0.0-20181213150558-05914d821849

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20181213153335-0fe22c71c476

require (
	github.com/konveyor/tackle2-addon v0.0.0-20230118205230-1aba7c8edc09
	github.com/konveyor/tackle2-hub v0.0.0-20230214175307-54c9422e54df
	github.com/pelletier/go-toml/v2 v2.0.5 // indirect
)
