gen:
	protoc protos/calculator.proto --js_out=import_style=commonjs,binary:./ui/src/ --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./ui/src/ --go-grpc_out=. --go_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
gen-ts:
	protoc -I=protos calculator.proto --js_out=import_style=commonjs,binary:. --grpc-web_out=import_style=typescript,mode=grpcweb:.

gen-js:
	protoc -I=protos calculator.proto --js_out=import_style=commonjs:. --grpc-web_out=import_style=commonjs,mode=grpcwebtext:.
	