go build main.go
START=$(date +%s)
INDEX=0
# parallel --gnu --jobs 0 ./main ::: puzzles/big/formatted/* #> ${p//formatted/solved}
# for p in puzzles/[^12]*/formatted/*
for p in puzzles/[^12]*/formatted/*
do
	# echo "$p"
	./main $p & #> ${p//formatted/solved}
	INDEX=$[$INDEX + 1]
done
wait
END=$(date +%s)
DIFF=$(( $END - $START ))
echo "Solved $INDEX puzzles in $DIFF seconds."
rm main
