go build solver.go
START=$(date +%s)
INDEX=0
for p in puzzles/*/formatted/*
do
	./solver $p > ${p//formatted/solved}
	INDEX=$[$INDEX + 1]
done
END=$(date +%s)
DIFF=$(( $END - $START ))
echo "Solved $INDEX puzzles in $DIFF seconds."
rm solver
