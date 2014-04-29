START=$(date +%s)
INDEX=0
for p in puzzles/$1/formatted/*
do
	./main $p > ${p//formatted/solved} &
	INDEX=$[$INDEX + 1]
done
wait
END=$(date +%s)
DIFF=$(( $END - $START ))
echo "Solved $INDEX $2 puzzles in $DIFF seconds."
