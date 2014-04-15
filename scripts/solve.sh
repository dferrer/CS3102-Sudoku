go build solver.go
for p in puzzles/easy1/formatted/*.txt
do
	./solver $p > ${p//formatted/solved}
done
rm solver