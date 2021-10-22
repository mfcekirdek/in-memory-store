set -x

arr=(repository service handler)

rm -rf mocks/* || true

for i in ${arr[@]}
do
  $GOPATH/bin/mockgen -destination=./mocks/mock_store_${i}.go -source=./pkg/${i}/store_${i}.go -package=mocks
  echo $?
done