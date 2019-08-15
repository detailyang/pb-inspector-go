<p align="center">
  <b>
    <span style="font-size:larger;">pb-inspector-go</span>
  </b>
  <br />
   <a href="https://travis-ci.org/detailyang/pb-inspector-go"><img src="https://travis-ci.org/detailyang/pb-inspector-go.svg?branch=master" /></a>
   <a href="https://ci.appveyor.com/project/detailyang/pb-inspector-go"><img src="https://ci.appveyor.com/api/projects/status/r4w4w09rwc4rpfwj?svg=true" /></a>
   <br />
   <b>"pb-inspector-go" inpsects the protobuf binary file to debug with or without schema</b>

   <img src="fixtures/carbon.png" />
   <blockquote>many thanks to <a href="https://github.com/emicklei">emicklei</a> who writes the protobuf go parser</blockquote>
</p>

# Install

* from github:
    > go get github.com/detailyang/protobuf-insepctor-go/cmd/pb-inspector

* from source:
    > make build

# Usage

## Without schema

````bash
echo 08ffff01100840f7d438 | xxd -r -p | pb-inspector -
````

## With schema

````bash
pb-inspector --file-type hex  --pb-file proto/test/v1/test.proto  fixtures/test1.hex "test.v1" "Test"
````
