# Matrix Ticker

Control an RPI attached matrix with a variety of animations.

View form details
implement nhl scroll
implement nhl/nfl box score
nfl logos
component padding
view definition persistance
frame buffer
decouple display from frame generation



# Local flags for rpi matrix
export CGO_CFLAGS="-I/home/andrew/Lab/rpi-rgb-led-matrix/include -DSHOW_REFRESH_RATE"; export CGO_LDFLAGS="-lrgbmatrix -L/home/andrew/Lab/rpi-rgb-led-matrix/lib -lstdc++ -lm"
