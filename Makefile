update:
	go build -o terraform-provider-cicd
	mkdir -p ~/.terraform.d/plugins/local/cicd/0.1.0/darwin_amd64/
	mv terraform-provider-cicd ~/.terraform.d/plugins/local/cicd/0.1.0/darwin_amd64/
