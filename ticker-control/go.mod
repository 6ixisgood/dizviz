module github.com/sixisgoood/matrix-ticker

go 1.19

require (
	github.com/fogleman/gg v1.3.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/sixisgoood/go-rpi-rgb-led-matrix v0.0.0-20180401002551-b26063b3169a
	gopkg.in/yaml.v2 v2.4.0
)

require (
	dmitri.shuralyov.com/gpu/mtl v0.0.0-20221208032759-85de2813cf6b // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20200222043503-6f7a984d4dc4 // indirect
	github.com/jezek/xgb v1.0.0 // indirect
	github.com/mcuadros/go-rpi-rgb-led-matrix v0.0.0-20180401002551-b26063b3169a // indirect
	golang.org/x/exp/shiny v0.0.0-20230210204819-062eb4c674ab // indirect
	golang.org/x/image v0.5.0 // indirect
	golang.org/x/mobile v0.0.0-20221110043201-43a038452099 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

// replace syntax
// replace github.com/sixisgoood/go-rpi-rgb-led-matrix => <CLONED_DIR>
replace github.com/sixisgoood/go-rpi-rgb-led-matrix => /home/andrew/Lab/go-rpi-rgb-led-matrix
