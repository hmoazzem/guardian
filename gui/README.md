# gui - guardian UI

Install dependencies

```sh
cd gui
npm install
```

Generate protobuf stubs
```sh
mkdir -p grpc/generated
protoc -I=../proto ../proto/metrics.proto \
  --js_out=import_style=commonjs:grpc/generated \
  --grpc-web_out=import_style=commonjs,mode=grpcwebtext:grpc/generated
```

Compile the grpc client JavaScript Code
```sh
npx webpack
```

Build the project lib
```sh
npx ng build ng-essential
```

Start dev serve

```sh
npx ng serve
```

## experimental stuff; mostly non-functional

- proto stubs typescript

```sh
protoc -I=../proto ../proto/metrics.proto \
  --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
  --js_out=import_style=commonjs:grpc/generated \
  --grpc-web_out=import_style=typescript,mode=grpcwebtext:grpc/generated \
  --ts_out=service=grpc-web,mode=grpc-js:grpc/generated
```

ts client errors

- pglite

```sh
cp node_modules/@electric-sql/pglite/dist/postgres.wasm public/
```
still trynna get helloworld working in pglite