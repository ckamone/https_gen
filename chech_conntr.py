TIME_WAIT = 30 # sysctl net parameter
PORT_RANGE = 65535
IP_COUNT = 5

cps = 10000
test_dur = 300
conntracks = 0

for i in range(test_dur):
  conntracks += cps
  assert IP_COUNT * PORT_RANGE > conntracks, 'not enough clients'
  if i > TIME_WAIT:
    conntracks-=cps

print(conntracks)