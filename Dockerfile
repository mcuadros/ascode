FROM mcuadros/ascode:latest

LABEL MAINTAINER="Máximo Cuadros <mcuadros@gmail.com>"
LABEL "com.github.actions.description"="converts starlark files to HCL"
LABEL "com.github.actions.name"="ascode-action"
LABEL "com.github.actions.color"="blue"

COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
