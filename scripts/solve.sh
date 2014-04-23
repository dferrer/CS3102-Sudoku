go build main.go
START=$(date +%s)
INDEX=0
for p in puzzles/1*/formatted/*
# for p in puzzles/big/formatted/*
do
	echo "$p"
	parallel ./main ::: $p > ${p//formatted/solved}
	INDEX=$[$INDEX + 1]
done
wait
END=$(date +%s)
DIFF=$(( $END - $START ))
echo "Solved $INDEX puzzles in $DIFF seconds."
rm main
