go build main.go
# bash scripts/solve.sh "easy1" "easy"
# for dir in hard[3]
# do 
# bash scripts/solve.sh "hard1" "hard"
# bash scripts/solve.sh "hard2" "hard"
# bash scripts/solve.sh "hard3" "hard"
# bash scripts/solve.sh "hardest1" "hardest"
# bash scripts/solve.sh "hardest2" "hardest"
# bash scripts/solve.sh "hardest3" "hardest"
bash scripts/solve.sh "big1" "assorted"
# done
# START=$(date +%s)
# INDEX=0
# parallel --gnu --jobs 0 ./main ::: puzzles/big/formatted/* #> ${p//formatted/solved}
# for p in puzzles/[^12]*/formatted/*
# for p in puzzles/easy1/formatted/*
# do
# 	# echo "$p"
# 	./main $p & #> ${p//formatted/solved}
# 	INDEX=$[$INDEX + 1]
# done
# wait
# END=$(date +%s)
# DIFF=$(( $END - $START ))
# echo "Solved $INDEX puzzles in $DIFF seconds."
rm main
