DELAY=125
convert -loop 0 -delay $DELAY s1.png -delay $DELAY s2.png -delay $DELAY s3.png -delay 300 s4.png -delay 500 s5.png -layers Optimize final.gif
