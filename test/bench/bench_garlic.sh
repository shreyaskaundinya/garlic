#!/bin/bash
# parameters
files=1000
warm=2

# cleanup
cleanup() {
    echo "cleaning up"
    # rm -rf /tmp/bench
}
trap cleanup EXIT

# check if hyperfine is installed
if ! command -v hyperfine &>/dev/null; then
    echo "hyperfine is not installed. Please install hyperfine to continue."
    exit 1
fi

# cloning candidates
echo ""
echo "clone SSGs"
echo ""
git clone --depth=1 https://github.com/shreyaskaundinya/garlic /tmp/bench/garlic

echo ""
echo "build SSGs"
echo ""

cd /tmp/bench/garlic && go build -o garlic && mv garlic /tmp/bench/garlic/garlic


# create the content folder
mkdir -p /tmp/bench/garlic/src/content

ls -la /tmp/bench/garlic/src/


# clean content/* dirs
echo ""
echo "Cleaning content directories"
echo ""
rm -rf /tmp/bench/garlic/src/content/*

# create multiple copies of the test file
echo ""
echo "Spawning $files different markdown files..."
for ((i = 0; i < files; i++)); do
    cp /tmp/bench/garlic/test/bench/bench.md "/tmp/bench/garlic/src/content/test_$i.md"
done

# run hyperfine
echo ""
echo "running benchmark: $files md files and $warm warmup runs"
echo ""

# "cd /tmp/bench/sapling/benchmark && ./../sapling run" \

hyperfine -p 'sync' -w $warm \
  "cd /tmp/bench/garlic && ./garlic --src-folder /tmp/bench/garlic/src --dest-folder /tmp/bench/garlic/dest --seed-files"
echo ""


echo "dest"
ls -la /tmp/bench/garlic/dest/