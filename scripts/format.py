import os, sys, math

def convert(s, n):
	return '\n'.join([s[i:i+n] for i in xrange(0, len(s), n)])

for d in os.listdir('./puzzles'):
	dirname = './puzzles/{0}'.format(d)
	filename = './puzzles/{0}/all_puzzles.txt'.format(d)
	with open(filename) as f:
		puzzles = [p.strip() for p in f.readlines()]
		n = int(math.sqrt(len(puzzles[0])))
		for p, i in zip(puzzles,range(1,len(puzzles))):
			with open('{0}/formatted/formatted{1}.txt'.format(dirname,i), 'w') as f2:
				f2.write(convert(p,n))
