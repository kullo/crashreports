# Crash report upload server

## Building
    make

## Running
    ./crashreports

## Testing
To run the unit tests:

    go test

Use curl to do a real upload:

    curl --include --form prod=foobar --form upload_file_minidump=@/etc/issue http://localhost:8080/upload

